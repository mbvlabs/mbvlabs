package domains

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type Works struct{ db contentDB }

func NewWorks(db storage.Connection) *Works { return &Works{db: db} }

func (w *Works) InTx(tx storage.Tx) *Works { return &Works{db: storage.Unwrap(tx)} }

type WorkEntity struct {
	ID                                                   int32
	CreatedAt, UpdatedAt                                 time.Time
	Title, Slug                                          string
	ClientName, ClientIndustry, ClientURL, ClientLogoURL sql.NullString
	Summary                                              string
	CoverImageURL                                        sql.NullString
	Specialisms, Platforms, Technologies                 []string
	Challenge, Approach, Deliverables, Outcome, Content  string
	StartedAt, CompletedAt                               sql.NullTime
	Status                                               StatusEnum
	PublishedAt                                          sql.NullTime
	IsFeatured                                           bool
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

const workColumns = "id, created_at, updated_at, title, slug, client_name, client_industry, client_url, client_logo_url, summary, cover_image_url, specialisms, platforms, technologies, challenge, approach, deliverables, outcome, content, started_at, completed_at, status, published_at, is_featured"

func scanWork(row interface{ Scan(...any) error }) (WorkEntity, error) {
	var e WorkEntity
	err := row.Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.Title, &e.Slug, &e.ClientName, &e.ClientIndustry, &e.ClientURL, &e.ClientLogoURL, &e.Summary, &e.CoverImageURL, &e.Specialisms, &e.Platforms, &e.Technologies, &e.Challenge, &e.Approach, &e.Deliverables, &e.Outcome, &e.Content, &e.StartedAt, &e.CompletedAt, &e.Status, &e.PublishedAt, &e.IsFeatured)
	return storage.Row(e, err)
}
func workArgs(e WorkEntity) []any {
	return []any{e.UpdatedAt, e.Title, e.Slug, e.ClientName, e.ClientIndustry, e.ClientURL, e.ClientLogoURL, e.Summary, e.CoverImageURL, e.Specialisms, e.Platforms, e.Technologies, e.Challenge, e.Approach, e.Deliverables, e.Outcome, e.Content, e.StartedAt, e.CompletedAt, e.Status, e.PublishedAt, e.IsFeatured}
}
func (w *Works) Find(ctx context.Context, id int32) (WorkEntity, error) {
	return scanWork(w.db.QueryRow(ctx, "SELECT "+workColumns+" FROM works WHERE id=$1", id))
}
func (w *Works) FindBySlug(ctx context.Context, s string) (WorkEntity, error) {
	return scanWork(w.db.QueryRow(ctx, "SELECT "+workColumns+" FROM works WHERE slug=$1", s))
}
func (w *Works) Create(ctx context.Context, d CreateWorkData) (WorkEntity, error) {
	e := workFromData(d)
	if err := e.Validate(); err != nil {
		return WorkEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return scanWork(w.db.QueryRow(ctx, `INSERT INTO works (created_at,updated_at,title,slug,client_name,client_industry,client_url,client_logo_url,summary,cover_image_url,specialisms,platforms,technologies,challenge,approach,deliverables,outcome,content,started_at,completed_at,status,published_at,is_featured)
	VALUES ($1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22) RETURNING `+workColumns, workArgs(e)...))
}
func (w *Works) Upsert(ctx context.Context, d CreateWorkData) (WorkEntity, error) {
	e := workFromData(d)
	if err := e.Validate(); err != nil {
		return WorkEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return scanWork(w.db.QueryRow(ctx, `INSERT INTO works (created_at,updated_at,title,slug,client_name,client_industry,client_url,client_logo_url,summary,cover_image_url,specialisms,platforms,technologies,challenge,approach,deliverables,outcome,content,started_at,completed_at,status,published_at,is_featured)
	VALUES ($1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
	ON CONFLICT (id) DO UPDATE SET title=excluded.title,slug=excluded.slug,client_name=excluded.client_name,client_industry=excluded.client_industry,client_url=excluded.client_url,client_logo_url=excluded.client_logo_url,summary=excluded.summary,cover_image_url=excluded.cover_image_url,specialisms=excluded.specialisms,platforms=excluded.platforms,technologies=excluded.technologies,challenge=excluded.challenge,approach=excluded.approach,deliverables=excluded.deliverables,outcome=excluded.outcome,content=excluded.content,started_at=excluded.started_at,completed_at=excluded.completed_at,status=excluded.status,published_at=excluded.published_at,is_featured=excluded.is_featured RETURNING `+workColumns, workArgs(e)...))
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
	args := append(workArgs(e), d.ID)
	return scanWork(w.db.QueryRow(ctx, `UPDATE works SET updated_at=$1,title=$2,slug=$3,client_name=$4,client_industry=$5,client_url=$6,client_logo_url=$7,summary=$8,cover_image_url=$9,specialisms=$10,platforms=$11,technologies=$12,challenge=$13,approach=$14,deliverables=$15,outcome=$16,content=$17,started_at=$18,completed_at=$19,status=$20,published_at=$21,is_featured=$22 WHERE id=$23 RETURNING `+workColumns, args...))
}
func (w *Works) Destroy(ctx context.Context, id int32) error {
	_, err := w.db.Exec(ctx, "DELETE FROM works WHERE id=$1", id)
	return err
}
func (w *Works) list(ctx context.Context, suffix string, args ...any) ([]WorkEntity, error) {
	rows, err := w.db.Query(ctx, "SELECT "+workColumns+" FROM works "+suffix, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]WorkEntity, 0)
	for rows.Next() {
		e, err := scanWork(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (w *Works) All(ctx context.Context) ([]WorkEntity, error) { return w.list(ctx, "") }

const workPublishedOrder = " ORDER BY completed_at DESC NULLS LAST, started_at DESC NULLS LAST, published_at DESC"

func (w *Works) AllPublished(ctx context.Context) ([]WorkEntity, error) {
	return w.list(ctx, "WHERE status='published'"+workPublishedOrder)
}
func (w *Works) FeaturedPublished(ctx context.Context, limit int) ([]WorkEntity, error) {
	if limit < 1 || limit > 3 {
		limit = 3
	}
	return w.list(ctx, "WHERE status='published' AND is_featured"+workPublishedOrder+" LIMIT $1", limit)
}
func (w *Works) CountFeatured(ctx context.Context, exceptID int32) (int64, error) {
	var n int64
	err := w.db.QueryRow(ctx, "SELECT count(*) FROM works WHERE is_featured AND ($1::integer <= 0 OR id <> $1)", exceptID).Scan(&n)
	return n, err
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
	where := ""
	order := ""
	if published {
		where = " WHERE status='published'"
		order = workPublishedOrder
	}
	var total int64
	if err := w.db.QueryRow(ctx, "SELECT count(*) FROM works"+where).Scan(&total); err != nil {
		return PaginatedWorks{}, err
	}
	entries, err := w.list(ctx, where+order+" LIMIT $1 OFFSET $2", size, (page-1)*size)
	if err != nil {
		return PaginatedWorks{}, err
	}
	return PaginatedWorks{entries, total, page, size, (total + size - 1) / size}, nil
}
func (w *Works) Paginate(ctx context.Context, page, size int64) (PaginatedWorks, error) {
	return w.paginate(ctx, page, size, false)
}
func (w *Works) PaginatePublished(ctx context.Context, page, size int64) (PaginatedWorks, error) {
	return w.paginate(ctx, page, size, true)
}
