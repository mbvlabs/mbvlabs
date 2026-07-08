package factories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"mbvlabs/internal/storage"
	"mbvlabs/models"

	"github.com/go-faker/faker/v4"
)

// RiverJobFactory wraps models.RiverJobEntity for testing
type RiverJobFactory struct {
	models.RiverJobEntity
}

type RiverJobOption func(*RiverJobFactory)

// BuildRiverJob creates an in-memory RiverJob with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateRiverJob.
func BuildRiverJob(opts ...RiverJobOption) models.RiverJobEntity {
	f := &RiverJobFactory{
		RiverJobEntity: models.RiverJobEntity{
			State:        models.RiverJobStateAvailable,
			Attempt:      0,
			MaxAttempts:  25,
			AttemptedAt:  sql.NullTime{},
			FinalizedAt:  sql.NullTime{},
			ScheduledAt:  time.Now(),
			Priority:     1,
			Args:         json.RawMessage(`{}`),
			AttemptedBy:  []string{},
			Errors:       []string{},
			Kind:         faker.Word(),
			Metadata:     json.RawMessage(`{}`),
			Queue:        "default",
			Tags:         []string{},
			UniqueKey:    []byte{},
			UniqueStates: sql.NullString{},
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.RiverJobEntity
}

// CreateRiverJob creates and persists a RiverJob to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateRiverJob(ctx context.Context, exec storage.Executor, opts ...RiverJobOption) (models.RiverJobEntity, error) {
	built := BuildRiverJob(opts...)

	entity := models.RiverJobEntity{
		CreatedAt:    time.Now(),
		State:        built.State,
		Attempt:      built.Attempt,
		MaxAttempts:  built.MaxAttempts,
		AttemptedAt:  built.AttemptedAt,
		FinalizedAt:  built.FinalizedAt,
		ScheduledAt:  built.ScheduledAt,
		Priority:     built.Priority,
		Args:         built.Args,
		AttemptedBy:  built.AttemptedBy,
		Errors:       built.Errors,
		Kind:         built.Kind,
		Metadata:     built.Metadata,
		Queue:        built.Queue,
		Tags:         built.Tags,
		UniqueKey:    built.UniqueKey,
		UniqueStates: built.UniqueStates,
	}

	if err := exec.NewInsert().Model(&entity).Returning("*").Scan(ctx); err != nil {
		return models.RiverJobEntity{}, err
	}

	return entity, nil
}

// CreateRiverJobs creates multiple RiverJob records at once
func CreateRiverJobs(ctx context.Context, exec storage.Executor, count int, opts ...RiverJobOption) ([]models.RiverJobEntity, error) {
	riverjobs := make([]models.RiverJobEntity, 0, count)

	for i := 0; i < count; i++ {
		entity, err := CreateRiverJob(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create riverjob %d: %w", i+1, err)
		}
		riverjobs = append(riverjobs, entity)
	}

	return riverjobs, nil
}

// Option functions

// WithRiverJobState sets the State field
func WithRiverJobState(value string) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.State = value
	}
}

// WithRiverJobAttempt sets the Attempt field
func WithRiverJobAttempt(value int16) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.Attempt = value
	}
}

// WithRiverJobMaxAttempts sets the MaxAttempts field
func WithRiverJobMaxAttempts(value int16) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.MaxAttempts = value
	}
}

// WithRiverJobAttemptedAt sets the AttemptedAt field
func WithRiverJobAttemptedAt(value sql.NullTime) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.AttemptedAt = value
	}
}

// WithRiverJobFinalizedAt sets the FinalizedAt field
func WithRiverJobFinalizedAt(value sql.NullTime) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.FinalizedAt = value
	}
}

// WithRiverJobScheduledAt sets the ScheduledAt field
func WithRiverJobScheduledAt(value time.Time) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.ScheduledAt = value
	}
}

// WithRiverJobPriority sets the Priority field
func WithRiverJobPriority(value int16) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.Priority = value
	}
}

// WithRiverJobArgs sets the Args field
func WithRiverJobArgs(value json.RawMessage) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.Args = value
	}
}

// WithRiverJobAttemptedBy sets the AttemptedBy field
func WithRiverJobAttemptedBy(value []string) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.AttemptedBy = value
	}
}

// WithRiverJobErrors sets the Errors field
func WithRiverJobErrors(value []string) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.Errors = value
	}
}

// WithRiverJobKind sets the Kind field
func WithRiverJobKind(value string) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.Kind = value
	}
}

// WithRiverJobMetadata sets the Metadata field
func WithRiverJobMetadata(value json.RawMessage) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.Metadata = value
	}
}

// WithRiverJobQueue sets the Queue field
func WithRiverJobQueue(value string) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.Queue = value
	}
}

// WithRiverJobTags sets the Tags field
func WithRiverJobTags(value []string) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.Tags = value
	}
}

// WithRiverJobUniqueKey sets the UniqueKey field
func WithRiverJobUniqueKey(value []byte) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.UniqueKey = value
	}
}

// WithRiverJobUniqueStates sets the UniqueStates field
func WithRiverJobUniqueStates(value sql.NullString) RiverJobOption {
	return func(f *RiverJobFactory) {
		f.RiverJobEntity.UniqueStates = value
	}
}
