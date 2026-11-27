package domains

import (
	"context"
	"errors"
	"strings"
	"time"

	"mbvlabs/domains/internal/queries"

	"github.com/jackc/pgx/v5"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/riverqueue/river"
)

var ErrRiverJobInvalidTransition = errors.New("river job cannot be moved to that state")

type RiverJob struct {
	ID          int64      `andurel:"id"`
	State       string     `andurel:"state"`
	Attempt     int16      `andurel:"attempt"`
	MaxAttempts int16      `andurel:"max_attempts"`
	AttemptedAt *time.Time `andurel:"attempted_at"`
	CreatedAt   time.Time  `andurel:"created_at"`
	FinalizedAt *time.Time `andurel:"finalized_at"`
	ScheduledAt time.Time  `andurel:"scheduled_at"`
	Priority    int16      `andurel:"priority"`
	Args        string     `andurel:"args"`
	AttemptedBy []string   `andurel:"attempted_by"`
	Errors      []string   `andurel:"errors"`
	Kind        string     `andurel:"kind"`
	Metadata    string     `andurel:"metadata"`
	Queue       string     `andurel:"queue"`
	Tags        []string   `andurel:"tags"`
}

type RiverQueueStats struct {
	Name                                                        string     `andurel:"name"`
	CreatedAt                                                   *time.Time `andurel:"created_at"`
	UpdatedAt                                                   *time.Time `andurel:"updated_at"`
	PausedAt                                                    *time.Time `andurel:"paused_at"`
	IsPaused                                                    bool
	AvailableCount                                              int64 `andurel:"available_count"`
	CancelledCount                                              int64 `andurel:"cancelled_count"`
	CompletedCount                                              int64 `andurel:"completed_count"`
	DiscardedCount                                              int64 `andurel:"discarded_count"`
	PendingCount                                                int64 `andurel:"pending_count"`
	RetryableCount                                              int64 `andurel:"retryable_count"`
	RunningCount                                                int64 `andurel:"running_count"`
	ScheduledCount                                              int64 `andurel:"scheduled_count"`
	TotalCount                                                  int64 `andurel:"total_count"`
	ActiveClients, MaxWorkers, NumJobsRunning, NumJobsCompleted int64
}

type RiverStateCount struct {
	State string `andurel:"state"`
	Count int64  `andurel:"count"`
}
type RiverFilters struct{ State, Queue, Kind string }
type RiverPagination struct{ Page, PageSize, TotalCount, TotalPages int64 }

type RiverQueues struct {
	queries    *queries.Queries
	client     *river.Client[pgx.Tx]
	maxWorkers map[string]int64
}

func NewRiverQueues(db storage.Connection, client *river.Client[pgx.Tx], maxWorkers map[string]int64) *RiverQueues {
	return &RiverQueues{queries: queries.New(db), client: client, maxWorkers: maxWorkers}
}

func normalizeRiverJob(j RiverJob) RiverJob {
	if j.AttemptedBy == nil {
		j.AttemptedBy = []string{}
	}
	if j.Errors == nil {
		j.Errors = []string{}
	}
	if j.Tags == nil {
		j.Tags = []string{}
	}
	return j
}

func (q *RiverQueues) Find(ctx context.Context, id int64) (RiverJob, error) {
	j, err := q.queries.GetRiverJob[RiverJob](ctx, id)
	if err != nil {
		return RiverJob{}, err
	}
	return normalizeRiverJob(j), nil
}

func (q *RiverQueues) Jobs(ctx context.Context, f RiverFilters, page, size int64) ([]RiverJob, RiverPagination, error) {
	queue, kind := strings.TrimSpace(f.Queue), strings.TrimSpace(f.Kind)
	total, err := q.queries.CountRiverJobs(ctx, queries.CountRiverJobsParams{Column1: f.State, Column2: queue, Column3: kind})
	if err != nil {
		return nil, RiverPagination{}, err
	}
	items, err := q.queries.ListRiverJobs[RiverJob](ctx, queries.ListRiverJobsParams{Column1: f.State, Column2: queue, Column3: kind, Column4: size, Column5: (page - 1) * size})
	if err != nil {
		return nil, RiverPagination{}, err
	}
	if items == nil {
		items = []RiverJob{}
	}
	for i := range items {
		items[i] = normalizeRiverJob(items[i])
	}
	return items, RiverPagination{page, size, total, (total + size - 1) / size}, nil
}

func (q *RiverQueues) StateCounts(ctx context.Context) ([]RiverStateCount, error) {
	items, err := q.queries.CountRiverJobStates[RiverStateCount](ctx)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []RiverStateCount{}
	}
	return items, nil
}

func (q *RiverQueues) Stats(ctx context.Context) ([]RiverQueueStats, error) {
	items, err := q.queries.ListRiverQueueStats[RiverQueueStats](ctx)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []RiverQueueStats{}
	}
	for i := range items {
		items[i].IsPaused = items[i].PausedAt != nil
		items[i].MaxWorkers = q.maxWorkers[items[i].Name]
	}
	return items, nil
}

func (q *RiverQueues) SetPaused(ctx context.Context, name string, paused bool) error {
	name = strings.TrimSpace(name)
	if name == "" || name == "*" {
		return errors.New("queue name is required")
	}
	// River's control API updates existing queues and sends the control notification.
	if err := q.queries.EnsureRiverQueue(ctx, name); err != nil {
		return err
	}
	if paused {
		return q.client.QueuePause(ctx, name, nil)
	}
	return q.client.QueueResume(ctx, name, nil)
}

func (q *RiverQueues) Action(ctx context.Context, id int64, action string) error {
	job, err := q.Find(ctx, id)
	if err != nil {
		return err
	}
	switch action {
	case "cancel":
		if job.State == "cancelled" || job.State == "completed" || job.State == "discarded" {
			return ErrRiverJobInvalidTransition
		}
		_, err = q.client.JobCancel(ctx, id)
	case "retry":
		if job.State != "cancelled" && job.State != "discarded" && job.State != "retryable" && job.State != "scheduled" {
			return ErrRiverJobInvalidTransition
		}
		_, err = q.client.JobRetry(ctx, id)
	case "discard":
		count, e := q.queries.DiscardRiverJob(ctx, id)
		if e != nil {
			return e
		}
		if count == 0 {
			return ErrRiverJobInvalidTransition
		}
	default:
		return errors.New("unknown job action")
	}
	return err
}
