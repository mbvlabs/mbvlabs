package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mbvlabs/internal/storage"
	"mbvlabs/internal/validation"
	"time"

	"github.com/uptrace/bun"
)

type DiaryEntryEntity struct {
	bun.BaseModel   `bun:"table:diary_entries,alias:diary_entries"`
	ID              int32          `bun:"id,pk,autoincrement"`
	CreatedAt       time.Time      `bun:"created_at"`
	UpdatedAt       time.Time      `bun:"updated_at"`
	EntryDate       time.Time      `bun:"entry_date,type:date"`
	MorningThoughts sql.NullString `bun:"morning_thoughts"`
	EveningThoughts sql.NullString `bun:"evening_thoughts"`
}

func (e DiaryEntryEntity) Validate() error {
	builder := validation.NewBuilder()
	builder.Required("entryDate", e.EntryDate)
	return builder.Err()
}

func (de diaryEntry) Find(ctx context.Context, db storage.Executor, id int32) (DiaryEntryEntity, error) {
	var entity DiaryEntryEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DiaryEntryEntity{}, ErrNotFound
		}
		return DiaryEntryEntity{}, err
	}

	return entity, nil
}

func (de diaryEntry) FindByDate(
	ctx context.Context,
	db storage.Executor,
	entryDate time.Time,
) (DiaryEntryEntity, error) {
	var entity DiaryEntryEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("entry_date = ?", normalizeDate(entryDate)).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DiaryEntryEntity{}, ErrNotFound
		}
		return DiaryEntryEntity{}, fmt.Errorf("find diary entry by date: %v", err)
	}

	return entity, nil
}

type CreateDiaryEntryData struct {
	EntryDate       time.Time
	MorningThoughts string
	EveningThoughts string
}

func (de diaryEntry) Create(ctx context.Context, db storage.Executor, data CreateDiaryEntryData) (DiaryEntryEntity, error) {
	now := time.Now()
	entity := DiaryEntryEntity{
		CreatedAt: now,
		UpdatedAt: now,
		EntryDate: normalizeDate(data.EntryDate),
		MorningThoughts: sql.NullString{
			String: data.MorningThoughts,
			Valid:  data.MorningThoughts != "",
		},
		EveningThoughts: sql.NullString{
			String: data.EveningThoughts,
			Valid:  data.EveningThoughts != "",
		},
	}

	if err := validation.Validate(entity); err != nil {
		return DiaryEntryEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return DiaryEntryEntity{}, err
	}

	return entity, nil
}

func (de diaryEntry) FindOrCreateForDate(
	ctx context.Context,
	db storage.Executor,
	entryDate time.Time,
) (DiaryEntryEntity, error) {
	now := time.Now()
	entity := DiaryEntryEntity{
		CreatedAt: now,
		UpdatedAt: now,
		EntryDate: normalizeDate(entryDate),
	}

	if err := validation.Validate(entity); err != nil {
		return DiaryEntryEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (entry_date) DO UPDATE").
		Set("entry_date = excluded.entry_date").
		Returning("*").
		Scan(ctx); err != nil {
		return DiaryEntryEntity{}, fmt.Errorf("find or create diary entry for date: %v", err)
	}

	return entity, nil
}

type UpdateDiaryEntryData struct {
	ID              int32
	EntryDate       time.Time
	MorningThoughts string
	EveningThoughts string
}

func (de diaryEntry) Update(ctx context.Context, db storage.Executor, data UpdateDiaryEntryData) (DiaryEntryEntity, error) {
	entity := DiaryEntryEntity{
		ID:        data.ID,
		UpdatedAt: time.Now(),
		EntryDate: normalizeDate(data.EntryDate),
		MorningThoughts: sql.NullString{
			String: data.MorningThoughts,
			Valid:  data.MorningThoughts != "",
		},
		EveningThoughts: sql.NullString{
			String: data.EveningThoughts,
			Valid:  data.EveningThoughts != "",
		},
	}

	if err := validation.Validate(entity); err != nil {
		return DiaryEntryEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("entry_date").
		Column("morning_thoughts").
		Column("evening_thoughts").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DiaryEntryEntity{}, ErrNotFound
		}
		return DiaryEntryEntity{}, err
	}

	return entity, nil
}

func (de diaryEntry) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*DiaryEntryEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (de diaryEntry) All(ctx context.Context, db storage.Executor) ([]DiaryEntryEntity, error) {
	var entities []DiaryEntryEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

type PaginatedDiaryEntries struct {
	DiaryEntries []DiaryEntryEntity
	TotalCount   int64
	Page         int64
	PageSize     int64
	TotalPages   int64
}

func (de diaryEntry) Paginate(ctx context.Context, db storage.Executor, page, pageSize int64) (PaginatedDiaryEntries, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	totalCount, err := db.NewSelect().
		Model(&DiaryEntryEntity{}).Count(ctx)
	if err != nil {
		return PaginatedDiaryEntries{}, err
	}

	entities := make([]DiaryEntryEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Order("entry_date DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedDiaryEntries{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedDiaryEntries{
		DiaryEntries: entities,
		TotalCount:   int64(totalCount),
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}, nil
}

func (de diaryEntry) Upsert(ctx context.Context, db storage.Executor, data CreateDiaryEntryData) (DiaryEntryEntity, error) {
	now := time.Now()
	entity := DiaryEntryEntity{
		CreatedAt: now,
		UpdatedAt: now,
		EntryDate: normalizeDate(data.EntryDate),
		MorningThoughts: sql.NullString{
			String: data.MorningThoughts,
			Valid:  data.MorningThoughts != "",
		},
		EveningThoughts: sql.NullString{
			String: data.EveningThoughts,
			Valid:  data.EveningThoughts != "",
		},
	}

	if err := validation.Validate(entity); err != nil {
		return DiaryEntryEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (entry_date) DO UPDATE").
		Set("updated_at = excluded.updated_at").
		Set("entry_date = excluded.entry_date").
		Set("morning_thoughts = excluded.morning_thoughts").
		Set("evening_thoughts = excluded.evening_thoughts").
		Returning("*").
		Scan(ctx); err != nil {
		return DiaryEntryEntity{}, err
	}

	return entity, nil
}

func normalizeDate(value time.Time) time.Time {
	if value.IsZero() {
		return time.Time{}
	}

	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}
