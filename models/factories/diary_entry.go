package factories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"mbvlabs/internal/storage"
	"mbvlabs/models"

	"github.com/go-faker/faker/v4"
)

// DiaryEntryFactory wraps models.DiaryEntryEntity for testing
type DiaryEntryFactory struct {
	models.DiaryEntryEntity
}

type DiaryEntryOption func(*DiaryEntryFactory)

// BuildDiaryEntry creates an in-memory DiaryEntry with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateDiaryEntry.
func BuildDiaryEntry(opts ...DiaryEntryOption) models.DiaryEntryEntity {
	f := &DiaryEntryFactory{
		DiaryEntryEntity: models.DiaryEntryEntity{
			EntryDate:       time.Now(),
			MorningThoughts: sql.NullString{String: faker.Word(), Valid: true},
			EveningThoughts: sql.NullString{String: faker.Word(), Valid: true},
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.DiaryEntryEntity
}

// CreateDiaryEntry creates and persists a DiaryEntry to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateDiaryEntry(ctx context.Context, exec storage.Executor, opts ...DiaryEntryOption) (models.DiaryEntryEntity, error) {
	built := BuildDiaryEntry(opts...)

	entity := models.DiaryEntryEntity{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		EntryDate:       built.EntryDate,
		MorningThoughts: built.MorningThoughts,
		EveningThoughts: built.EveningThoughts,
	}

	if err := exec.NewInsert().Model(&entity).Returning("*").Scan(ctx); err != nil {
		return models.DiaryEntryEntity{}, err
	}

	return entity, nil
}

// CreateDiaryEntrys creates multiple DiaryEntry records at once
func CreateDiaryEntrys(ctx context.Context, exec storage.Executor, count int, opts ...DiaryEntryOption) ([]models.DiaryEntryEntity, error) {
	diaryentrys := make([]models.DiaryEntryEntity, 0, count)

	for i := 0; i < count; i++ {
		entryOpts := append(
			[]DiaryEntryOption{WithDiaryEntriesEntryDate(time.Now().AddDate(0, 0, i))},
			opts...,
		)
		entity, err := CreateDiaryEntry(ctx, exec, entryOpts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create diaryentry %d: %w", i+1, err)
		}
		diaryentrys = append(diaryentrys, entity)
	}

	return diaryentrys, nil
}

// Option functions

// WithDiaryEntriesEntryDate sets the EntryDate field
func WithDiaryEntriesEntryDate(value time.Time) DiaryEntryOption {
	return func(f *DiaryEntryFactory) {
		f.DiaryEntryEntity.EntryDate = value
	}
}

// WithDiaryEntriesMorningThoughts sets the MorningThoughts field
func WithDiaryEntriesMorningThoughts(value sql.NullString) DiaryEntryOption {
	return func(f *DiaryEntryFactory) {
		f.DiaryEntryEntity.MorningThoughts = value
	}
}

// WithDiaryEntriesEveningThoughts sets the EveningThoughts field
func WithDiaryEntriesEveningThoughts(value sql.NullString) DiaryEntryOption {
	return func(f *DiaryEntryFactory) {
		f.DiaryEntryEntity.EveningThoughts = value
	}
}
