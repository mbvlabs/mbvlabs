package main

import (
	"context"
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
	"uuid"

	"mbvlabs/domains"
	"mbvlabs/migrations"

	"github.com/mbvlabs/andurel/pkg/storage/storagetest"
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

func TestReleaseMigratesV1UsingMigrationURL(t *testing.T) {
	cluster, err := storagetest.NewCluster(t.Context())
	if err != nil {
		t.Skipf("start disposable PostgreSQL: %v", err)
	}
	t.Cleanup(func() {
		if err := cluster.Close(context.Background()); err != nil {
			t.Errorf("close disposable PostgreSQL: %v", err)
		}
	})

	v1 := fstest.MapFS{}
	entries, err := fs.ReadDir(migrations.Migrations, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".sql") || name >= "00016_" {
			continue
		}
		data, err := fs.ReadFile(migrations.Migrations, name)
		if err != nil {
			t.Fatal(err)
		}
		v1[name] = &fstest.MapFile{Data: data}
	}
	db, cfg := cluster.NewTestDBWithConfig(t, v1, ".")
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
