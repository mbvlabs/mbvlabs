package domains_test

import (
	"errors"
	"testing"

	"mbvlabs/domains"
	"mbvlabs/migrations"

	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/riverqueue/river"
)

type queueTestArgs struct{}

func (queueTestArgs) Kind() string { return "test_kind" }

func TestRiverQueueDashboardAndActions(t *testing.T) {
	db := testCluster.NewTestDB(t, migrations.Migrations, ".")
	insert, err := storage.NewQueueInsert(db, storage.DefaultQueueConfig())
	if err != nil {
		t.Fatal(err)
	}
	q := domains.NewRiverQueues(db, insert.Client(), map[string]int64{"default": 7})
	ctx := t.Context()
	inserted, err := insert.Insert(ctx, queueTestArgs{}, &river.InsertOpts{Queue: "default", MaxAttempts: 3})
	if err != nil {
		t.Fatal(err)
	}
	id := inserted.Job.ID
	jobs, page, err := q.Jobs(ctx, domains.RiverFilters{Kind: "TEST"}, 1, 10)
	if err != nil || len(jobs) != 1 || jobs[0].ID != id || page.TotalCount != 1 || page.TotalPages != 1 {
		t.Fatalf("jobs=%+v page=%+v err=%v", jobs, page, err)
	}
	stats, err := q.Stats(ctx)
	if err != nil || len(stats) != 1 || stats[0].AvailableCount != 1 || stats[0].MaxWorkers != 7 {
		t.Fatalf("stats=%+v err=%v", stats, err)
	}
	if err := q.SetPaused(ctx, "default", true); err != nil {
		t.Fatal(err)
	}
	stats, err = q.Stats(ctx)
	if err != nil || !stats[0].IsPaused {
		t.Fatalf("paused stats=%+v err=%v", stats, err)
	}
	if err := q.SetPaused(ctx, "default", false); err != nil {
		t.Fatal(err)
	}
	if err := q.Action(ctx, id, "cancel"); err != nil {
		t.Fatal(err)
	}
	if err := q.Action(ctx, id, "cancel"); !errors.Is(err, domains.ErrRiverJobInvalidTransition) {
		t.Fatalf("cancel again: %v", err)
	}
	if err := q.Action(ctx, id, "retry"); err != nil {
		t.Fatal(err)
	}
	if err := q.Action(ctx, id, "discard"); err != nil {
		t.Fatal(err)
	}
	job, err := q.Find(ctx, id)
	if err != nil || job.State != "discarded" {
		t.Fatalf("job=%+v err=%v", job, err)
	}
	if _, err := q.Find(ctx, id+1000); !errors.Is(err, domains.ErrNotFound) {
		t.Fatalf("missing job: %v", err)
	}
}
