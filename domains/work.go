package domains

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"mbvlabs/domains/internal/queries"

	"github.com/gosimple/slug"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type Works struct{ queries *queries.Queries }

func NewWorks(db storage.Connection) *Works { return &Works{queries: queries.New(db)} }

func (w *Works) InTx(tx storage.Tx) *Works { return &Works{queries: queries.New(storage.Unwrap(tx))} }

type WorkEntity struct {
	ID             int32          `andurel:"id"`
	CreatedAt      time.Time      `andurel:"created_at"`
	UpdatedAt      time.Time      `andurel:"updated_at"`
	Title          string         `andurel:"title"`
	Slug           string         `andurel:"slug"`
	ClientName     sql.NullString `andurel:"client_name"`
	ClientIndustry sql.NullString `andurel:"client_industry"`
	ClientURL      sql.NullString `andurel:"client_url"`
	ClientLogoURL  sql.NullString `andurel:"client_logo_url"`
	Summary        string         `andurel:"summary"`
	CoverImageURL  sql.NullString `andurel:"cover_image_url"`
	Specialisms    []string       `andurel:"specialisms"`
	Platforms      []string       `andurel:"platforms"`
	Technologies   []string       `andurel:"technologies"`
	Challenge      string         `andurel:"challenge"`
	Approach       string         `andurel:"approach"`
	Deliverables   string         `andurel:"deliverables"`
	Outcome        string         `andurel:"outcome"`
	Content        string         `andurel:"content"`
	StartedAt      sql.NullTime   `andurel:"started_at"`
	CompletedAt    sql.NullTime   `andurel:"completed_at"`
	Status         StatusEnum     `andurel:"status"`
	PublishedAt    sql.NullTime   `andurel:"published_at"`
	IsFeatured     bool           `andurel:"is_featured"`
}

func (e WorkEntity) Validate() error {
	b := validation.NewBuilder()
	b.Required("title", e.Title)
	b.LenBetween("title", e.Title, 10, 255)
	b.Required("slug", e.Slug)
	b.MaxLen("slug", e.Slug, 255)
	b.OneOf("status", e.Status, Draft.String(), Published.String())
	b.MinLen("summary", e.Summary, 10)
	b.URL("clientUrl", e.ClientURL)
	b.URL("clientLogoUrl", e.ClientLogoURL)
	b.URL("coverImageUrl", e.CoverImageURL)
	b.NoBlankItems("specialisms", e.Specialisms)
	b.NoBlankItems("platforms", e.Platforms)
	b.NoBlankItems("technologies", e.Technologies)
	b.TimeBeforeOrEqual("startedAt", e.StartedAt, "completedAt", e.CompletedAt)
	if strings.TrimSpace(e.Slug) != "" && slug.Make(e.Slug) != e.Slug {
		b.AddField("slug", "slug", "must be a valid slug")
	}
	if e.Status == Published {
		b.Required("clientName", e.ClientName)
		b.Required("clientIndustry", e.ClientIndustry)
		b.Required("clientUrl", e.ClientURL)
		b.Required("clientLogoUrl", e.ClientLogoURL)
		b.Required("summary", e.Summary)
		b.Required("coverImageUrl", e.CoverImageURL)
		b.Required("specialisms", e.Specialisms)
		b.Required("platforms", e.Platforms)
		b.Required("technologies", e.Technologies)
		b.Required("challenge", e.Challenge)
		b.Required("approach", e.Approach)
		b.Required("deliverables", e.Deliverables)
		b.Required("outcome", e.Outcome)
		b.Required("content", e.Content)
		b.Required("startedAt", e.StartedAt)
		b.Required("completedAt", e.CompletedAt)
		b.Required("publishedAt", e.PublishedAt)
	}
	return b.Err()
}

type CreateWorkData struct {
	Title, Slug, ClientName, ClientIndustry, ClientURL, ClientLogoURL string
	Summary, Content, CoverImageURL                                   string
	Specialisms, Platforms, Technologies                              []string
	Challenge, Approach, Deliverables, Outcome                        string
	TestimonialQuote, TestimonialName, TestimonialRole                string
	StartedAt, CompletedAt                                            time.Time
	Status                                                            StatusEnum
	PublishedAt                                                       time.Time
	IsFeatured                                                        bool
}
type UpdateWorkData struct {
	ID                                                                int32
	Title, Slug, ClientName, ClientIndustry, ClientURL, ClientLogoURL string
	Summary, Content, CoverImageURL                                   string
	Specialisms, Platforms, Technologies                              []string
	Challenge, Approach, Deliverables, Outcome                        string
	TestimonialQuote, TestimonialName, TestimonialRole                string
	StartedAt, CompletedAt                                            time.Time
	Status                                                            StatusEnum
	PublishedAt                                                       time.Time
	IsFeatured                                                        bool
}
type PaginatedWorks struct {
	Works                                  []WorkEntity
	TotalCount, Page, PageSize, TotalPages int64
}

func workString(s string) sql.NullString { return sql.NullString{String: s, Valid: s != ""} }
func workTime(t time.Time) sql.NullTime  { return sql.NullTime{Time: t, Valid: !t.IsZero()} }
func workFromData(d CreateWorkData) WorkEntity {
	now := time.Now()
	e := WorkEntity{CreatedAt: now, UpdatedAt: now, Title: d.Title, Slug: d.Slug,
		ClientName: workString(d.ClientName), ClientIndustry: workString(d.ClientIndustry), ClientURL: workString(d.ClientURL),
		ClientLogoURL: workString(d.ClientLogoURL), Summary: d.Summary, CoverImageURL: workString(d.CoverImageURL),
		Specialisms: d.Specialisms, Platforms: d.Platforms, Technologies: d.Technologies,
		Challenge: d.Challenge, Approach: d.Approach, Deliverables: d.Deliverables, Outcome: d.Outcome, Content: d.Content,
		StartedAt: workTime(d.StartedAt), CompletedAt: workTime(d.CompletedAt), Status: d.Status, PublishedAt: workTime(d.PublishedAt), IsFeatured: d.IsFeatured}
	if e.Status == Published && !e.PublishedAt.Valid {
		e.PublishedAt = workTime(now)
	}
	return e
}

func contentRows[T any](rows []T, err error) ([]T, error) {
	if err == nil && rows == nil {
		rows = []T{}
	}
	return rows, err
}

func workParams(e WorkEntity) queries.CreateWorkParams {
	return queries.CreateWorkParams{CreatedAt: e.UpdatedAt, Title: e.Title, Slug: e.Slug,
		ClientName: contentText(e.ClientName), ClientIndustry: contentText(e.ClientIndustry), ClientUrl: contentText(e.ClientURL),
		ClientLogoUrl: contentText(e.ClientLogoURL), Summary: e.Summary, CoverImageUrl: contentText(e.CoverImageURL),
		Specialisms: e.Specialisms, Platforms: e.Platforms, Technologies: e.Technologies,
		Challenge: e.Challenge, Approach: e.Approach, Deliverables: e.Deliverables, Outcome: e.Outcome, Content: e.Content,
		StartedAt: contentDate(e.StartedAt), CompletedAt: contentDate(e.CompletedAt), Status: e.Status.String(),
		PublishedAt: contentTimestamp(e.PublishedAt), IsFeatured: e.IsFeatured}
}
func (w *Works) Find(ctx context.Context, id int32) (WorkEntity, error) {
	return w.queries.GetWork[WorkEntity](ctx, id)
}
func (w *Works) Latest(ctx context.Context, limit int) ([]WorkEntity, error) {
	return contentRows(w.queries.LatestWorks[WorkEntity](ctx, int64(limit)))
}
func (w *Works) FindBySlug(ctx context.Context, slug string) (WorkEntity, error) {
	return w.queries.GetWorkBySlug[WorkEntity](ctx, slug)
}
func (w *Works) Create(ctx context.Context, d CreateWorkData) (WorkEntity, error) {
	e := workFromData(d)
	if err := e.Validate(); err != nil {
		return WorkEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return w.queries.CreateWork[WorkEntity](ctx, workParams(e))
}
func (w *Works) Upsert(ctx context.Context, d CreateWorkData) (WorkEntity, error) {
	e := workFromData(d)
	if err := e.Validate(); err != nil {
		return WorkEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return w.queries.UpsertWork[WorkEntity](ctx, queries.UpsertWorkParams(workParams(e)))
}
func (w *Works) Update(ctx context.Context, d UpdateWorkData) (WorkEntity, error) {
	old, err := w.Find(ctx, d.ID)
	if err != nil {
		return WorkEntity{}, err
	}
	e := workFromData(CreateWorkData{Title: d.Title, Slug: d.Slug, ClientName: d.ClientName, ClientIndustry: d.ClientIndustry, ClientURL: d.ClientURL, ClientLogoURL: d.ClientLogoURL, Summary: d.Summary, Content: d.Content, CoverImageURL: d.CoverImageURL, Specialisms: d.Specialisms, Platforms: d.Platforms, Technologies: d.Technologies, Challenge: d.Challenge, Approach: d.Approach, Deliverables: d.Deliverables, Outcome: d.Outcome, StartedAt: d.StartedAt, CompletedAt: d.CompletedAt, Status: d.Status, IsFeatured: d.IsFeatured})
	e.PublishedAt = old.PublishedAt
	if e.Status == Published && !e.PublishedAt.Valid {
		e.PublishedAt = workTime(time.Now())
	}
	if err := e.Validate(); err != nil {
		return WorkEntity{}, errors.Join(ErrDomainValidation, err)
	}
	a := workParams(e)
	return w.queries.UpdateWork[WorkEntity](ctx, queries.UpdateWorkParams{
		ID: d.ID, UpdatedAt: a.CreatedAt, Title: a.Title, Slug: a.Slug, ClientName: a.ClientName,
		ClientIndustry: a.ClientIndustry, ClientUrl: a.ClientUrl, ClientLogoUrl: a.ClientLogoUrl, Summary: a.Summary,
		CoverImageUrl: a.CoverImageUrl, Specialisms: a.Specialisms, Platforms: a.Platforms, Technologies: a.Technologies,
		Challenge: a.Challenge, Approach: a.Approach, Deliverables: a.Deliverables, Outcome: a.Outcome, Content: a.Content,
		StartedAt: a.StartedAt, CompletedAt: a.CompletedAt, Status: a.Status, PublishedAt: a.PublishedAt, IsFeatured: a.IsFeatured,
	})
}

// EditEditorial must be called on a transaction-bound store.
func (w *Works) EditEditorial(ctx context.Context, slug string, edit func(*WorkEntity)) (WorkEntity, error) {
	e, err := w.queries.GetWorkBySlugForUpdate[WorkEntity](ctx, slug)
	if err != nil {
		return WorkEntity{}, err
	}
	edit(&e)
	if err := e.Validate(); err != nil {
		return WorkEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return w.queries.EditWork[WorkEntity](ctx, queries.EditWorkParams{
		ID: e.ID, UpdatedAt: time.Now(), Title: e.Title, Summary: e.Summary, Content: e.Content,
		ClientUrl: contentText(e.ClientURL), ClientLogoUrl: contentText(e.ClientLogoURL), CoverImageUrl: contentText(e.CoverImageURL),
		Challenge: e.Challenge, Approach: e.Approach, Deliverables: e.Deliverables, Outcome: e.Outcome,
	})
}
func (w *Works) LockFeatured(ctx context.Context) error      { return w.queries.LockFeaturedWorks(ctx) }
func (w *Works) Destroy(ctx context.Context, id int32) error { return w.queries.DeleteWork(ctx, id) }
func (w *Works) All(ctx context.Context) ([]WorkEntity, error) {
	return contentRows(w.queries.ListWorks[WorkEntity](ctx))
}
func (w *Works) AllPublished(ctx context.Context) ([]WorkEntity, error) {
	return contentRows(w.queries.ListPublishedWorks[WorkEntity](ctx))
}
func (w *Works) FeaturedPublished(ctx context.Context, limit int) ([]WorkEntity, error) {
	if limit < 1 || limit > 3 {
		limit = 3
	}
	return contentRows(w.queries.ListFeaturedWorks[WorkEntity](ctx, int32(limit)))
}
func (w *Works) CountFeatured(ctx context.Context, exceptID int32) (int64, error) {
	return w.queries.CountFeaturedWorks(ctx, exceptID)
}
func (w *Works) paginate(ctx context.Context, page, size int64, published bool) (PaginatedWorks, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	var total int64
	var entries []WorkEntity
	var err error
	if published {
		total, err = w.queries.CountPublishedWorks(ctx)
		if err == nil {
			entries, err = w.queries.PagePublishedWorks[WorkEntity](ctx, queries.PagePublishedWorksParams{Column1: size, Column2: (page - 1) * size})
		}
	} else {
		total, err = w.queries.CountWorks(ctx)
		if err == nil {
			entries, err = w.queries.PageWorks[WorkEntity](ctx, queries.PageWorksParams{Column1: size, Column2: (page - 1) * size})
		}
	}
	if err != nil {
		return PaginatedWorks{}, err
	}
	if entries == nil {
		entries = []WorkEntity{}
	}
	return PaginatedWorks{entries, total, page, size, (total + size - 1) / size}, nil
}
func (w *Works) Paginate(ctx context.Context, page, size int64) (PaginatedWorks, error) {
	return w.paginate(ctx, page, size, false)
}
func (w *Works) PaginatePublished(ctx context.Context, page, size int64) (PaginatedWorks, error) {
	return w.paginate(ctx, page, size, true)
}
