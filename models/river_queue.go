package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"mbvlabs/internal/storage"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

var ErrRiverQueueNameRequired = errors.New("river queue name is required")

type RiverQueueEntity struct {
	bun.BaseModel `bun:"table:river_queue,alias:river_queue"`
	Name          string          `bun:"name,pk"`
	CreatedAt     time.Time       `bun:"created_at"`
	Metadata      json.RawMessage `bun:"metadata,type:jsonb"`
	PausedAt      sql.NullTime    `bun:"paused_at"`
	UpdatedAt     time.Time       `bun:"updated_at"`
}

type RiverQueueStats struct {
	Name             string
	CreatedAt        sql.NullTime
	UpdatedAt        sql.NullTime
	PausedAt         sql.NullTime
	AvailableCount   int64
	CancelledCount   int64
	CompletedCount   int64
	DiscardedCount   int64
	PendingCount     int64
	RetryableCount   int64
	RunningCount     int64
	ScheduledCount   int64
	TotalCount       int64
	ActiveClients    int64
	MaxWorkers       int64
	NumJobsRunning   int64
	NumJobsCompleted int64
}

type RiverQueueDashboard struct {
	Queues      []RiverQueueStats
	StateCounts []RiverJobStateCount
	RecentJobs  []RiverJobEntity
}

func (rq riverQueue) Dashboard(ctx context.Context, db storage.Executor) (RiverQueueDashboard, error) {
	queues, err := rq.WithStats(ctx, db)
	if err != nil {
		return RiverQueueDashboard{}, err
	}

	stateCounts, err := RiverJob.CountByState(ctx, db)
	if err != nil {
		return RiverQueueDashboard{}, err
	}

	recentJobs, err := RiverJob.Recent(ctx, db, 10)
	if err != nil {
		return RiverQueueDashboard{}, err
	}

	return RiverQueueDashboard{
		Queues:      queues,
		StateCounts: stateCounts,
		RecentJobs:  recentJobs,
	}, nil
}

func (rq riverQueue) WithStats(ctx context.Context, db storage.Executor) ([]RiverQueueStats, error) {
	queues := make([]RiverQueueStats, 0)
	if err := db.NewRaw(`
		WITH queue_names AS (
			SELECT name FROM river_queue
			UNION
			SELECT DISTINCT queue AS name FROM river_job
			UNION
			SELECT DISTINCT name FROM river_client_queue
		),
		job_counts AS (
			SELECT
				queue,
				count(*) FILTER (WHERE state = 'available') AS available_count,
				count(*) FILTER (WHERE state = 'cancelled') AS cancelled_count,
				count(*) FILTER (WHERE state = 'completed') AS completed_count,
				count(*) FILTER (WHERE state = 'discarded') AS discarded_count,
				count(*) FILTER (WHERE state = 'pending') AS pending_count,
				count(*) FILTER (WHERE state = 'retryable') AS retryable_count,
				count(*) FILTER (WHERE state = 'running') AS running_count,
				count(*) FILTER (WHERE state = 'scheduled') AS scheduled_count,
				count(*) AS total_count
			FROM river_job
			GROUP BY queue
		),
		client_counts AS (
			SELECT
				name,
				count(*) AS active_clients,
				coalesce(sum(max_workers), 0) AS max_workers,
				coalesce(sum(num_jobs_running), 0) AS num_jobs_running,
				coalesce(sum(num_jobs_completed), 0) AS num_jobs_completed
			FROM river_client_queue
			GROUP BY name
		)
		SELECT
			queue_names.name,
			river_queue.created_at,
			river_queue.updated_at,
			river_queue.paused_at,
			coalesce(job_counts.available_count, 0) AS available_count,
			coalesce(job_counts.cancelled_count, 0) AS cancelled_count,
			coalesce(job_counts.completed_count, 0) AS completed_count,
			coalesce(job_counts.discarded_count, 0) AS discarded_count,
			coalesce(job_counts.pending_count, 0) AS pending_count,
			coalesce(job_counts.retryable_count, 0) AS retryable_count,
			coalesce(job_counts.running_count, 0) AS running_count,
			coalesce(job_counts.scheduled_count, 0) AS scheduled_count,
			coalesce(job_counts.total_count, 0) AS total_count,
			coalesce(client_counts.active_clients, 0) AS active_clients,
			coalesce(client_counts.max_workers, 0) AS max_workers,
			coalesce(client_counts.num_jobs_running, 0) AS num_jobs_running,
			coalesce(client_counts.num_jobs_completed, 0) AS num_jobs_completed
		FROM queue_names
		LEFT JOIN river_queue ON river_queue.name = queue_names.name
		LEFT JOIN job_counts ON job_counts.queue = queue_names.name
		LEFT JOIN client_counts ON client_counts.name = queue_names.name
		ORDER BY queue_names.name
	`).Scan(ctx, &queues); err != nil {
		return nil, err
	}

	return queues, nil
}

func (rq riverQueue) Pause(ctx context.Context, db storage.Executor, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrRiverQueueNameRequired
	}

	now := time.Now()
	entity := RiverQueueEntity{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  json.RawMessage(`{}`),
		PausedAt:  sql.NullTime{Time: now, Valid: true},
	}

	_, err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (name) DO UPDATE").
		Set("updated_at = excluded.updated_at").
		Set("paused_at = excluded.paused_at").
		Exec(ctx)

	return err
}

func (rq riverQueue) Resume(ctx context.Context, db storage.Executor, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrRiverQueueNameRequired
	}

	now := time.Now()
	entity := RiverQueueEntity{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  json.RawMessage(`{}`),
	}

	_, err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (name) DO UPDATE").
		Set("updated_at = excluded.updated_at").
		Set("paused_at = NULL").
		Exec(ctx)

	return err
}
