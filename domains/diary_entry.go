package domains

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"mbvlabs/domains/internal/queries"

	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type DiaryEntries struct{ queries *queries.Queries }

func NewDiaryEntries(db storage.Connection) *DiaryEntries {
	return &DiaryEntries{queries: queries.New(db)}
}

type DiaryEntryEntity struct {
	ID              int32          `andurel:"id"`
	CreatedAt       time.Time      `andurel:"created_at"`
	UpdatedAt       time.Time      `andurel:"updated_at"`
	EntryDate       time.Time      `andurel:"entry_date"`
	MorningThoughts sql.NullString `andurel:"morning_thoughts"`
	EveningThoughts sql.NullString `andurel:"evening_thoughts"`
}

func (e DiaryEntryEntity) Validate() error {
	b := validation.NewBuilder()
	b.Required("entryDate", e.EntryDate)
	return b.Err()
}

type CreateDiaryEntryData struct {
	EntryDate                        time.Time
	MorningThoughts, EveningThoughts string
}

type UpdateDiaryEntryData struct {
	ID                               int32
	EntryDate                        time.Time
	MorningThoughts, EveningThoughts string
}

type PaginatedDiaryEntries struct {
	DiaryEntries                           []DiaryEntryEntity
	TotalCount, Page, PageSize, TotalPages int64
}

func normalizeDate(value time.Time) time.Time {
	if value.IsZero() {
		return value
	}
	y, m, d := value.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, value.Location())
}

func (d *DiaryEntries) Find(ctx context.Context, id int32) (DiaryEntryEntity, error) {
	return d.queries.GetDiaryEntry[DiaryEntryEntity](ctx, id)
}
func (d *DiaryEntries) FindByDate(ctx context.Context, date time.Time) (DiaryEntryEntity, error) {
	return d.queries.GetDiaryEntryByDate[DiaryEntryEntity](ctx, normalizeDate(date))
}

type diaryWriteMode uint8

const (
	diaryCreate diaryWriteMode = iota
	diaryUpsert
	diaryFindOrCreate
)

func (d *DiaryEntries) write(ctx context.Context, data CreateDiaryEntryData, mode diaryWriteMode) (DiaryEntryEntity, error) {
	date := normalizeDate(data.EntryDate)
	if err := (DiaryEntryEntity{EntryDate: date}).Validate(); err != nil {
		return DiaryEntryEntity{}, errors.Join(ErrDomainValidation, err)
	}
	now := time.Now()
	params := queries.CreateDiaryEntryParams{CreatedAt: now, EntryDate: date, MorningThoughts: contentText(sql.NullString{String: data.MorningThoughts, Valid: data.MorningThoughts != ""}), EveningThoughts: contentText(sql.NullString{String: data.EveningThoughts, Valid: data.EveningThoughts != ""})}
	switch mode {
	case diaryUpsert:
		return d.queries.UpsertDiaryEntry[DiaryEntryEntity](ctx, queries.UpsertDiaryEntryParams(params))
	case diaryFindOrCreate:
		return d.queries.FindOrCreateDiaryEntry[DiaryEntryEntity](ctx, queries.FindOrCreateDiaryEntryParams(params))
	default:
		return d.queries.CreateDiaryEntry[DiaryEntryEntity](ctx, params)
	}
}
func (d *DiaryEntries) Create(ctx context.Context, data CreateDiaryEntryData) (DiaryEntryEntity, error) {
	return d.write(ctx, data, diaryCreate)
}
func (d *DiaryEntries) Upsert(ctx context.Context, data CreateDiaryEntryData) (DiaryEntryEntity, error) {
	return d.write(ctx, data, diaryUpsert)
}
func (d *DiaryEntries) FindOrCreateForDate(ctx context.Context, date time.Time) (DiaryEntryEntity, error) {
	return d.write(ctx, CreateDiaryEntryData{EntryDate: date}, diaryFindOrCreate)
}
func (d *DiaryEntries) Update(ctx context.Context, data UpdateDiaryEntryData) (DiaryEntryEntity, error) {
	date := normalizeDate(data.EntryDate)
	if err := (DiaryEntryEntity{EntryDate: date}).Validate(); err != nil {
		return DiaryEntryEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return d.queries.UpdateDiaryEntry[DiaryEntryEntity](ctx, queries.UpdateDiaryEntryParams{ID: data.ID, UpdatedAt: time.Now(), EntryDate: date, MorningThoughts: contentText(sql.NullString{String: data.MorningThoughts, Valid: data.MorningThoughts != ""}), EveningThoughts: contentText(sql.NullString{String: data.EveningThoughts, Valid: data.EveningThoughts != ""})})
}
func (d *DiaryEntries) Destroy(ctx context.Context, id int32) error {
	return d.queries.DeleteDiaryEntry(ctx, id)
}
func nonNilDiaryEntries(items []DiaryEntryEntity, err error) ([]DiaryEntryEntity, error) {
	if err == nil && items == nil {
		items = []DiaryEntryEntity{}
	}
	return items, err
}
func (d *DiaryEntries) All(ctx context.Context) ([]DiaryEntryEntity, error) {
	return nonNilDiaryEntries(d.queries.ListDiaryEntries[DiaryEntryEntity](ctx))
}
func (d *DiaryEntries) FindBetweenDates(ctx context.Context, startDate, endDate time.Time) ([]DiaryEntryEntity, error) {
	return nonNilDiaryEntries(d.queries.ListDiaryEntriesBetweenDates[DiaryEntryEntity](ctx, queries.ListDiaryEntriesBetweenDatesParams{EntryDate: normalizeDate(startDate), EntryDate_2: normalizeDate(endDate)}))
}
func (d *DiaryEntries) Paginate(ctx context.Context, page, pageSize int64) (PaginatedDiaryEntries, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	total, err := d.queries.CountDiaryEntries(ctx)
	if err != nil {
		return PaginatedDiaryEntries{}, err
	}
	entries, err := nonNilDiaryEntries(d.queries.ListDiaryEntriesNewestFirst[DiaryEntryEntity](ctx, queries.ListDiaryEntriesNewestFirstParams{Column1: pageSize, Column2: (page - 1) * pageSize}))
	if err != nil {
		return PaginatedDiaryEntries{}, err
	}
	return PaginatedDiaryEntries{entries, total, page, pageSize, (total + pageSize - 1) / pageSize}, nil
}
