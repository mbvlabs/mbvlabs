package domains

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type DiaryEntries struct{ db storage.Connection }

func NewDiaryEntries(db storage.Connection) *DiaryEntries { return &DiaryEntries{db: db} }

type DiaryEntryEntity struct {
	ID                               int32
	CreatedAt, UpdatedAt, EntryDate  time.Time
	MorningThoughts, EveningThoughts sql.NullString
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

func scanDiary(row interface{ Scan(...any) error }) (DiaryEntryEntity, error) {
	var e DiaryEntryEntity
	err := row.Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.EntryDate, &e.MorningThoughts, &e.EveningThoughts)
	return storage.Row(e, err)
}

const diaryColumns = "id, created_at, updated_at, entry_date, morning_thoughts, evening_thoughts"

func (d *DiaryEntries) Find(ctx context.Context, id int32) (DiaryEntryEntity, error) {
	return scanDiary(d.db.QueryRow(ctx, "SELECT "+diaryColumns+" FROM diary_entries WHERE id=$1", id))
}
func (d *DiaryEntries) FindByDate(ctx context.Context, date time.Time) (DiaryEntryEntity, error) {
	return scanDiary(d.db.QueryRow(ctx, "SELECT "+diaryColumns+" FROM diary_entries WHERE entry_date=$1", normalizeDate(date)))
}
func (d *DiaryEntries) write(ctx context.Context, data CreateDiaryEntryData, conflict string) (DiaryEntryEntity, error) {
	date := normalizeDate(data.EntryDate)
	if err := (DiaryEntryEntity{EntryDate: date}).Validate(); err != nil {
		return DiaryEntryEntity{}, errors.Join(ErrDomainValidation, err)
	}
	now := time.Now()
	return scanDiary(d.db.QueryRow(ctx, "INSERT INTO diary_entries (created_at, updated_at, entry_date, morning_thoughts, evening_thoughts) VALUES ($1,$1,$2,$3,$4) "+conflict+" RETURNING "+diaryColumns,
		now, date, sql.NullString{String: data.MorningThoughts, Valid: data.MorningThoughts != ""}, sql.NullString{String: data.EveningThoughts, Valid: data.EveningThoughts != ""}))
}
func (d *DiaryEntries) Create(ctx context.Context, data CreateDiaryEntryData) (DiaryEntryEntity, error) {
	return d.write(ctx, data, "")
}
func (d *DiaryEntries) Upsert(ctx context.Context, data CreateDiaryEntryData) (DiaryEntryEntity, error) {
	return d.write(ctx, data, "ON CONFLICT (entry_date) DO UPDATE SET updated_at=excluded.updated_at, morning_thoughts=excluded.morning_thoughts, evening_thoughts=excluded.evening_thoughts")
}
func (d *DiaryEntries) FindOrCreateForDate(ctx context.Context, date time.Time) (DiaryEntryEntity, error) {
	return d.write(ctx, CreateDiaryEntryData{EntryDate: date}, "ON CONFLICT (entry_date) DO UPDATE SET entry_date=excluded.entry_date")
}
func (d *DiaryEntries) Update(ctx context.Context, data UpdateDiaryEntryData) (DiaryEntryEntity, error) {
	date := normalizeDate(data.EntryDate)
	if err := (DiaryEntryEntity{EntryDate: date}).Validate(); err != nil {
		return DiaryEntryEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return scanDiary(d.db.QueryRow(ctx, "UPDATE diary_entries SET updated_at=$2, entry_date=$3, morning_thoughts=$4, evening_thoughts=$5 WHERE id=$1 RETURNING "+diaryColumns,
		data.ID, time.Now(), date, sql.NullString{String: data.MorningThoughts, Valid: data.MorningThoughts != ""}, sql.NullString{String: data.EveningThoughts, Valid: data.EveningThoughts != ""}))
}
func (d *DiaryEntries) Destroy(ctx context.Context, id int32) error {
	_, err := d.db.Exec(ctx, "DELETE FROM diary_entries WHERE id=$1", id)
	return err
}
func (d *DiaryEntries) list(ctx context.Context, query string, args ...any) ([]DiaryEntryEntity, error) {
	rows, err := d.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]DiaryEntryEntity, 0)
	for rows.Next() {
		e, err := scanDiary(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (d *DiaryEntries) All(ctx context.Context) ([]DiaryEntryEntity, error) {
	return d.list(ctx, "SELECT "+diaryColumns+" FROM diary_entries")
}
func (d *DiaryEntries) FindBetweenDates(ctx context.Context, startDate, endDate time.Time) ([]DiaryEntryEntity, error) {
	return d.list(ctx, "SELECT "+diaryColumns+" FROM diary_entries WHERE entry_date >= $1 AND entry_date <= $2 ORDER BY entry_date ASC", normalizeDate(startDate), normalizeDate(endDate))
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
	var total int64
	if err := d.db.QueryRow(ctx, "SELECT count(*) FROM diary_entries").Scan(&total); err != nil {
		return PaginatedDiaryEntries{}, err
	}
	entries, err := d.list(ctx, "SELECT "+diaryColumns+" FROM diary_entries ORDER BY entry_date DESC LIMIT $1 OFFSET $2", pageSize, (page-1)*pageSize)
	if err != nil {
		return PaginatedDiaryEntries{}, err
	}
	return PaginatedDiaryEntries{entries, total, page, pageSize, (total + pageSize - 1) / pageSize}, nil
}
