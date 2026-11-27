package domains

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/riverqueue/river"
)

var ErrRiverJobInvalidTransition = errors.New("river job cannot be moved to that state")

type RiverJob struct {
	ID                    int64
	State                 string
	Attempt, MaxAttempts  int16
	AttemptedAt           *time.Time
	CreatedAt             time.Time
	FinalizedAt           *time.Time
	ScheduledAt           time.Time
	Priority              int16
	Args                  string
	AttemptedBy, Errors   []string
	Kind, Metadata, Queue string
	Tags                  []string
}

type RiverQueueStats struct {
	Name                                                                                                                                   string
	CreatedAt, UpdatedAt, PausedAt                                                                                                         *time.Time
	IsPaused                                                                                                                               bool
	AvailableCount, CancelledCount, CompletedCount, DiscardedCount, PendingCount, RetryableCount, RunningCount, ScheduledCount, TotalCount int64
	ActiveClients, MaxWorkers, NumJobsRunning, NumJobsCompleted                                                                            int64
}

type RiverStateCount struct {
	State string
	Count int64
}
type RiverFilters struct{ State, Queue, Kind string }
type RiverPagination struct{ Page, PageSize, TotalCount, TotalPages int64 }

type RiverQueues struct {
	db         storage.Connection
	client     *river.Client[pgx.Tx]
	maxWorkers map[string]int64
}

func NewRiverQueues(db storage.Connection, client *river.Client[pgx.Tx], maxWorkers map[string]int64) *RiverQueues {
	return &RiverQueues{db: db, client: client, maxWorkers: maxWorkers}
}

const riverJobColumns = `id, state::text, attempt, max_attempts, attempted_at, created_at, finalized_at, scheduled_at, priority, args::text, attempted_by, ARRAY(SELECT e::text FROM unnest(errors) e), kind, metadata::text, queue, tags`

func scanRiverJob(row interface{ Scan(...any) error }) (RiverJob, error) {
	var j RiverJob
	err := row.Scan(&j.ID, &j.State, &j.Attempt, &j.MaxAttempts, &j.AttemptedAt, &j.CreatedAt, &j.FinalizedAt, &j.ScheduledAt, &j.Priority, &j.Args, &j.AttemptedBy, &j.Errors, &j.Kind, &j.Metadata, &j.Queue, &j.Tags)
	if j.AttemptedBy == nil {
		j.AttemptedBy = []string{}
	}
	if j.Errors == nil {
		j.Errors = []string{}
	}
	if j.Tags == nil {
		j.Tags = []string{}
	}
	return storage.Row(j, err)
}

func (q *RiverQueues) Find(ctx context.Context, id int64) (RiverJob, error) {
	return scanRiverJob(q.db.QueryRow(ctx, "SELECT "+riverJobColumns+" FROM river_job WHERE id=$1", id))
}

func (q *RiverQueues) Jobs(ctx context.Context, f RiverFilters, page, size int64) ([]RiverJob, RiverPagination, error) {
	where := ` WHERE ($1::text = '' OR state::text = $1) AND ($2::text = '' OR queue = $2) AND ($3::text = '' OR kind ILIKE '%' || $3 || '%')`
	args := []any{f.State, strings.TrimSpace(f.Queue), strings.TrimSpace(f.Kind)}
	var total int64
	if err := q.db.QueryRow(ctx, "SELECT count(*) FROM river_job"+where, args...).Scan(&total); err != nil {
		return nil, RiverPagination{}, err
	}
	rows, err := q.db.Query(ctx, "SELECT "+riverJobColumns+" FROM river_job"+where+" ORDER BY created_at DESC, id DESC LIMIT $4 OFFSET $5", append(args, size, (page-1)*size)...)
	if err != nil {
		return nil, RiverPagination{}, err
	}
	defer rows.Close()
	items := []RiverJob{}
	for rows.Next() {
		j, err := scanRiverJob(rows)
		if err != nil {
			return nil, RiverPagination{}, err
		}
		items = append(items, j)
	}
	return items, RiverPagination{page, size, total, (total + size - 1) / size}, rows.Err()
}

func (q *RiverQueues) StateCounts(ctx context.Context) ([]RiverStateCount, error) {
	rows, err := q.db.Query(ctx, "SELECT state::text, count(*) FROM river_job GROUP BY state ORDER BY state")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RiverStateCount{}
	for rows.Next() {
		var s RiverStateCount
		if err := rows.Scan(&s.State, &s.Count); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}

func (q *RiverQueues) Stats(ctx context.Context) ([]RiverQueueStats, error) {
	rows, err := q.db.Query(ctx, `WITH names AS (SELECT name FROM river_queue UNION SELECT queue FROM river_job),
	counts AS (SELECT queue,
	count(*) FILTER (WHERE state='available') available, count(*) FILTER (WHERE state='cancelled') cancelled,
	count(*) FILTER (WHERE state='completed') completed, count(*) FILTER (WHERE state='discarded') discarded,
	count(*) FILTER (WHERE state='pending') pending, count(*) FILTER (WHERE state='retryable') retryable,
	count(*) FILTER (WHERE state='running') running, count(*) FILTER (WHERE state='scheduled') scheduled,
	count(*) total FROM river_job GROUP BY queue)
	SELECT names.name, q.created_at, q.updated_at, q.paused_at, coalesce(c.available,0), coalesce(c.cancelled,0),
	coalesce(c.completed,0), coalesce(c.discarded,0), coalesce(c.pending,0), coalesce(c.retryable,0),
	coalesce(c.running,0), coalesce(c.scheduled,0), coalesce(c.total,0)
	FROM names LEFT JOIN river_queue q ON q.name=names.name LEFT JOIN counts c ON c.queue=names.name ORDER BY names.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RiverQueueStats{}
	for rows.Next() {
		var s RiverQueueStats
		if err := rows.Scan(&s.Name, &s.CreatedAt, &s.UpdatedAt, &s.PausedAt, &s.AvailableCount, &s.CancelledCount, &s.CompletedCount, &s.DiscardedCount, &s.PendingCount, &s.RetryableCount, &s.RunningCount, &s.ScheduledCount, &s.TotalCount); err != nil {
			return nil, err
		}
		s.IsPaused = s.PausedAt != nil
		s.MaxWorkers = q.maxWorkers[s.Name]
		items = append(items, s)
	}
	return items, rows.Err()
}

func (q *RiverQueues) SetPaused(ctx context.Context, name string, paused bool) error {
	name = strings.TrimSpace(name)
	if name == "" || name == "*" {
		return errors.New("queue name is required")
	}
	// River's control API updates existing queues and sends the control notification.
	if _, err := q.db.Exec(ctx, `INSERT INTO river_queue (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, name); err != nil {
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
		tag, e := q.db.Exec(ctx, `UPDATE river_job SET state='discarded', finalized_at=now() WHERE id=$1 AND state NOT IN ('cancelled','completed','discarded')`, id)
		if e != nil {
			return e
		}
		if tag.RowsAffected() == 0 {
			return ErrRiverJobInvalidTransition
		}
	default:
		return errors.New("unknown job action")
	}
	return err
}
