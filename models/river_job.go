package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mbvlabs/internal/storage"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

const (
	RiverJobStateAvailable = "available"
	RiverJobStateCancelled = "cancelled"
	RiverJobStateCompleted = "completed"
	RiverJobStateDiscarded = "discarded"
	RiverJobStatePending   = "pending"
	RiverJobStateRetryable = "retryable"
	RiverJobStateRunning   = "running"
	RiverJobStateScheduled = "scheduled"
)

var (
	ErrRiverJobInvalidTransition = errors.New("river job cannot be moved to that state")
	riverJobStates               = map[string]struct{}{
		RiverJobStateAvailable: {},
		RiverJobStateCancelled: {},
		RiverJobStateCompleted: {},
		RiverJobStateDiscarded: {},
		RiverJobStatePending:   {},
		RiverJobStateRetryable: {},
		RiverJobStateRunning:   {},
		RiverJobStateScheduled: {},
	}
)

type RiverJobEntity struct {
	bun.BaseModel `bun:"table:river_job,alias:river_job"`
	ID            int64           `bun:"id,pk,autoincrement"`
	State         string          `bun:"state,type:river_job_state"`
	Attempt       int16           `bun:"attempt"`
	MaxAttempts   int16           `bun:"max_attempts"`
	AttemptedAt   sql.NullTime    `bun:"attempted_at"`
	CreatedAt     time.Time       `bun:"created_at"`
	FinalizedAt   sql.NullTime    `bun:"finalized_at"`
	ScheduledAt   time.Time       `bun:"scheduled_at"`
	Priority      int16           `bun:"priority"`
	Args          json.RawMessage `bun:"args,type:jsonb"`
	AttemptedBy   []string        `bun:"attempted_by,array"`
	Errors        []string        `bun:"errors,array"`
	Kind          string          `bun:"kind"`
	Metadata      json.RawMessage `bun:"metadata,type:jsonb"`
	Queue         string          `bun:"queue"`
	Tags          []string        `bun:"tags,array"`
	UniqueKey     []byte          `bun:"unique_key"`
	UniqueStates  sql.NullString  `bun:"unique_states"`
}

type RiverJobFilters struct {
	State string
	Queue string
	Kind  string
}

type PaginatedRiverJobs struct {
	RiverJobs  []RiverJobEntity
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

type RiverJobStateCount struct {
	State string
	Count int64
}

func IsRiverJobState(state string) bool {
	_, ok := riverJobStates[state]
	return ok
}

func (rj riverJob) Find(ctx context.Context, db storage.Executor, id int64) (RiverJobEntity, error) {
	var entity RiverJobEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RiverJobEntity{}, ErrNotFound
		}
		return RiverJobEntity{}, err
	}

	return entity, nil
}

func (rj riverJob) Recent(ctx context.Context, db storage.Executor, limit int) ([]RiverJobEntity, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	var entities []RiverJobEntity
	if err := db.NewSelect().
		Model(&entities).
		Order("created_at DESC").
		Order("id DESC").
		Limit(limit).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

func (rj riverJob) CountByState(ctx context.Context, db storage.Executor) ([]RiverJobStateCount, error) {
	counts := make([]RiverJobStateCount, 0)
	if err := db.NewRaw(`
		SELECT state::text AS state, count(*) AS count
		FROM river_job
		GROUP BY state
		ORDER BY state
	`).Scan(ctx, &counts); err != nil {
		return nil, err
	}

	return counts, nil
}

func (rj riverJob) Paginate(
	ctx context.Context,
	db storage.Executor,
	filters RiverJobFilters,
	page int64,
	pageSize int64,
) (PaginatedRiverJobs, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
	}
	if pageSize > 100 {
		pageSize = 100
	}

	base := db.NewSelect().Model((*RiverJobEntity)(nil))
	applyRiverJobFilters(base, filters)

	totalCount, err := base.Count(ctx)
	if err != nil {
		return PaginatedRiverJobs{}, err
	}

	entities := make([]RiverJobEntity, 0, int(pageSize))
	query := db.NewSelect().
		Model(&entities).
		Order("created_at DESC").
		Order("id DESC").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize))
	applyRiverJobFilters(query, filters)

	if err := query.Scan(ctx); err != nil {
		return PaginatedRiverJobs{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedRiverJobs{
		RiverJobs:  entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (rj riverJob) Cancel(ctx context.Context, db storage.Executor, id int64) error {
	return rj.moveToFinalState(ctx, db, id, RiverJobStateCancelled)
}

func (rj riverJob) Discard(ctx context.Context, db storage.Executor, id int64) error {
	return rj.moveToFinalState(ctx, db, id, RiverJobStateDiscarded)
}

func (rj riverJob) Retry(ctx context.Context, db storage.Executor, id int64) error {
	var queue string
	if err := db.NewRaw(`
		UPDATE river_job
		SET state = ?, finalized_at = NULL, scheduled_at = now()
		WHERE id = ?
			AND state IN ('cancelled', 'discarded', 'retryable', 'scheduled')
		RETURNING queue
	`, RiverJobStateAvailable, id).Scan(ctx, &queue); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		if _, err := rj.Find(ctx, db, id); err != nil {
			return err
		}
		return ErrRiverJobInvalidTransition
	}

	_, err := db.NewRaw(
		"SELECT pg_notify('river_insert', json_build_object('queue', ?)::text)",
		queue,
	).Exec(ctx)
	return err
}

func (rj riverJob) moveToFinalState(
	ctx context.Context,
	db storage.Executor,
	id int64,
	state string,
) error {
	result, err := db.NewUpdate().
		Model((*RiverJobEntity)(nil)).
		Set("state = ?", state).
		Set("finalized_at = now()").
		Where("id = ?", id).
		Where("state NOT IN ('cancelled', 'completed', 'discarded')").
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}

	if _, err := rj.Find(ctx, db, id); err != nil {
		return err
	}
	return ErrRiverJobInvalidTransition
}

func applyRiverJobFilters(query *bun.SelectQuery, filters RiverJobFilters) {
	if IsRiverJobState(filters.State) {
		query.Where("state = ?", filters.State)
	}

	if queue := strings.TrimSpace(filters.Queue); queue != "" {
		query.Where("queue = ?", queue)
	}

	if kind := strings.TrimSpace(filters.Kind); kind != "" {
		query.Where("kind ILIKE ?", fmt.Sprintf("%%%s%%", kind))
	}
}
