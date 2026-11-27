package main

import (
	"context"
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
	"time"
	"uuid"

	"mbvlabs/domains"
	"mbvlabs/migrations"
	"mbvlabs/queue/jobs"

	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/storage/storagetest"
	"github.com/riverqueue/river"
)

func TestReleaseRequiresMigrationURL(t *testing.T) {
	t.Setenv("RELEASE_COMMAND", "1")
	t.Setenv("DATABASE_MIGRATION_URL", "")
	t.Setenv("DB_KIND", "invalid")

	if _, err := openMigrationDatabase(t.Context()); err == nil ||
		err.Error() != "DATABASE_MIGRATION_URL is required for the Fly release command" {
		t.Fatal("release must reject a missing migration URL before using DB_* settings")
	}
}

func TestInvalidMigrationURLDoesNotLeakCredentials(t *testing.T) {
	const password = "synthetic-migration-password"
	t.Setenv("DATABASE_MIGRATION_URL", "postgres://user:"+password+"@localhost:invalid/db")

	if _, err := openMigrationDatabase(t.Context()); err == nil ||
		!strings.Contains(err.Error(), "DATABASE_MIGRATION_URL") ||
		strings.Contains(err.Error(), password) {
		t.Fatal("expected a migration URL error without credentials")
	}
}

func migrationTestCluster(t *testing.T) *storagetest.Cluster {
	t.Helper()
	cluster, err := storagetest.NewCluster(t.Context())
	if err != nil {
		t.Skipf("start disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() {
		if err := cluster.Close(context.Background()); err != nil {
			t.Errorf("close disposable PostgreSQL: %v", err)
		}
	})
	return cluster
}

func migrationsBefore(t *testing.T, cutoff string) fstest.MapFS {
	t.Helper()
	earlier := fstest.MapFS{}
	entries, err := fs.ReadDir(migrations.Migrations, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".sql") || name >= cutoff {
			continue
		}
		data, err := fs.ReadFile(migrations.Migrations, name)
		if err != nil {
			t.Fatal(err)
		}
		earlier[name] = &fstest.MapFile{Data: data}
	}
	return earlier
}

func TestReleaseMigratesV1UsingMigrationURL(t *testing.T) {
	db, cfg := migrationTestCluster(t).NewTestDBWithConfig(t, migrationsBefore(t, "00016_"), ".")
	migrationURL, err := cfg.DatabaseURL()
	if err != nil {
		t.Fatal("build disposable PostgreSQL URL")
	}
	t.Setenv("DATABASE_MIGRATION_URL", migrationURL)
	t.Setenv("RELEASE_COMMAND", "1")
	t.Setenv("DB_KIND", "invalid")

	users := domains.NewUsers(db)
	var absentUser uuid.UUID
	if _, err := users.Find(t.Context(), absentUser); err == nil || !strings.Contains(err.Error(), "session_version") {
		t.Fatalf("expected V1 to lack session_version, got %v", err)
	}
	if err := run(); err != nil {
		t.Fatalf("migrate V1 database: %v", err)
	}
	if _, err := users.Find(t.Context(), absentUser); !errors.Is(err, domains.ErrNotFound) {
		t.Fatalf("expected V2 user query to succeed with no match: %v", err)
	}
	if err := run(); err != nil {
		t.Fatalf("rerun migrations: %v", err)
	}
}

type retiredReminderArgs struct{}

func (retiredReminderArgs) Kind() string { return "send_diary_reminder" }

func TestReleaseCancelsOnlyQueuedDiaryReminders(t *testing.T) {
	db, cfg := migrationTestCluster(t).NewTestDBWithConfig(t, migrationsBefore(t, "00018_"), ".")
	ctx := t.Context()
	insert, err := storage.NewQueueInsert(db, storage.DefaultQueueConfig())
	if err != nil {
		t.Fatal(err)
	}

	enqueue := func(args river.JobArgs, opts *river.InsertOpts) int64 {
		t.Helper()
		result, err := insert.Insert(ctx, args, opts)
		if err != nil {
			t.Fatal(err)
		}
		return result.Job.ID
	}
	scheduledAt := time.Now().Add(time.Hour)
	available := enqueue(retiredReminderArgs{}, nil)
	scheduled := enqueue(retiredReminderArgs{}, &river.InsertOpts{ScheduledAt: scheduledAt})
	unrelated := enqueue(jobs.DeleteExpiredTokensArgs{}, &river.InsertOpts{ScheduledAt: scheduledAt})
	historical := enqueue(retiredReminderArgs{}, nil)
	cancelled, err := insert.Client().JobCancel(ctx, historical)
	if err != nil {
		t.Fatal(err)
	}

	migrationURL, err := cfg.DatabaseURL()
	if err != nil {
		t.Fatal("build disposable PostgreSQL URL")
	}
	t.Setenv("DATABASE_MIGRATION_URL", migrationURL)
	t.Setenv("RELEASE_COMMAND", "1")
	if err := run(); err != nil {
		t.Fatalf("cancel queued reminders: %v", err)
	}

	queues := domains.NewRiverQueues(db, insert.Client(), nil)
	for _, id := range []int64{available, scheduled} {
		job, err := queues.Find(ctx, id)
		if err != nil || job.State != "cancelled" || job.FinalizedAt == nil {
			t.Fatalf("reminder %d: state=%q, finalized=%v, err=%v", id, job.State, job.FinalizedAt, err)
		}
	}
	job, err := queues.Find(ctx, unrelated)
	if err != nil || job.State != "scheduled" || job.FinalizedAt != nil {
		t.Fatalf("unrelated scheduled job changed: %+v, %v", job, err)
	}
	job, err = queues.Find(ctx, historical)
	if err != nil || job.State != "cancelled" || job.FinalizedAt == nil || !job.FinalizedAt.Equal(*cancelled.FinalizedAt) {
		t.Fatalf("historical cancelled job changed: %+v, %v", job, err)
	}
}
