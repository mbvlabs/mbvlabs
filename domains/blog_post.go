package domains

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type BlogPosts struct{ db contentDB }

func NewBlogPosts(db storage.Connection) *BlogPosts { return &BlogPosts{db: db} }

func (bp *BlogPosts) InTx(tx storage.Tx) *BlogPosts { return &BlogPosts{db: storage.Unwrap(tx)} }

type BlogPostEntity struct {
	ID                  int32
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Title               string
	Slug                string
	Excerpt             sql.NullString
	Body                string
	Status              string
	CoverImageUrl       sql.NullString
	Tags                json.RawMessage
	PublishedAt         sql.NullTime
	PublicationSchedule json.RawMessage
}

type BlogPostPublicationSchedule struct {
	JobID       int64     `json:"job_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

func (bp *BlogPostEntity) SetPublicationSchedule(schedule BlogPostPublicationSchedule) error {
	data, err := json.Marshal(schedule)
	if err != nil {
		return fmt.Errorf("marshal blog post publication schedule: %v", err)
	}
	bp.PublicationSchedule = data
	return nil
}

func (bp BlogPostEntity) PublicationScheduleData() (*BlogPostPublicationSchedule, error) {
	data := bytes.TrimSpace(bp.PublicationSchedule)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil, nil
	}
	var schedule BlogPostPublicationSchedule
	if err := json.Unmarshal(data, &schedule); err != nil {
		return nil, fmt.Errorf("unmarshal blog post publication schedule: %v", err)
	}
	return &schedule, nil
}

func (bp BlogPostEntity) Validate() error {
	b := validation.NewBuilder()
	b.Required("title", bp.Title)
	b.LenBetween("title", bp.Title, 10, 255)
	b.Required("slug", bp.Slug)
	b.MaxLen("slug", bp.Slug, 255)
	b.OneOf("status", bp.Status, "draft", "published", "archived")
	b.MinLen("excerpt", bp.Excerpt, 10)
	b.MinLen("body", bp.Body, 10)
	b.URL("coverImageUrl", bp.CoverImageUrl)
	schedule, err := bp.PublicationScheduleData()
	if err != nil {
		b.AddField("publicationSchedule", "json", "must be valid JSON")
	} else if schedule != nil {
		b.Required("publicationSchedule.jobId", schedule.JobID)
		b.Required("publicationSchedule.scheduledAt", schedule.ScheduledAt)
	}
	var tags []string
	tagsErr := json.Unmarshal(bp.Tags, &tags)
	if tagsErr == nil {
		b.NoBlankItems("tags", tags)
	}
	if strings.TrimSpace(bp.Slug) != "" && slug.Make(bp.Slug) != bp.Slug {
		b.AddField("slug", "slug", "must be a valid slug")
	}
	if bp.Status == "published" {
		b.Required("excerpt", bp.Excerpt)
		b.Required("body", bp.Body)
		b.Required("coverImageUrl", bp.CoverImageUrl)
		b.Required("publishedAt", bp.PublishedAt)
		if tagsErr != nil || len(tags) == 0 {
			b.AddField("tags", "required", "is required")
		}
	}
	return b.Err()
}

type CreateBlogPostData struct {
	Title, Slug, Excerpt, Body, Status, CoverImageUrl string
	Tags                                              []string
}

type UpdateBlogPostData struct {
	ID                                                int32
	Title, Slug, Excerpt, Body, Status, CoverImageUrl string
	Tags                                              []string
}

type PaginatedBlogPosts struct {
	BlogPosts                              []BlogPostEntity
	TotalCount, Page, PageSize, TotalPages int64
}

const blogPostColumns = "id, created_at, updated_at, title, slug, excerpt, body, status, cover_image_url, tags, published_at, publication_schedule"

func scanBlogPost(row interface{ Scan(...any) error }) (BlogPostEntity, error) {
	var e BlogPostEntity
	err := row.Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.Title, &e.Slug, &e.Excerpt, &e.Body, &e.Status, &e.CoverImageUrl, &e.Tags, &e.PublishedAt, &e.PublicationSchedule)
	return storage.Row(e, err)
}

func (bp *BlogPosts) Find(ctx context.Context, id int32) (BlogPostEntity, error) {
	return scanBlogPost(bp.db.QueryRow(ctx, "SELECT "+blogPostColumns+" FROM blog_posts WHERE id=$1", id))
}

func (bp *BlogPosts) FindBySlug(ctx context.Context, slug string) (BlogPostEntity, error) {
	return scanBlogPost(bp.db.QueryRow(ctx, "SELECT "+blogPostColumns+" FROM blog_posts WHERE slug=$1", slug))
}

func newBlogPost(data CreateBlogPostData) (BlogPostEntity, error) {
	tags, err := json.Marshal(data.Tags)
	if err != nil {
		return BlogPostEntity{}, fmt.Errorf("marshal tags: %v", err)
	}
	now := time.Now()
	e := BlogPostEntity{CreatedAt: now, UpdatedAt: now, Title: data.Title, Slug: data.Slug,
		Excerpt: sql.NullString{String: data.Excerpt, Valid: data.Excerpt != ""}, Body: data.Body,
		Status: data.Status, CoverImageUrl: sql.NullString{String: data.CoverImageUrl, Valid: data.CoverImageUrl != ""}, Tags: tags}
	if e.Status == "published" {
		e.PublishedAt = sql.NullTime{Time: now, Valid: true}
	}
	return e, nil
}

func validateBlogPost(e BlogPostEntity) error {
	if err := e.Validate(); err != nil {
		return errors.Join(ErrDomainValidation, err)
	}
	return nil
}

func (bp *BlogPosts) Create(ctx context.Context, data CreateBlogPostData) (BlogPostEntity, error) {
	return bp.insert(ctx, data, "")
}

func (bp *BlogPosts) insert(ctx context.Context, data CreateBlogPostData, conflict string) (BlogPostEntity, error) {
	e, err := newBlogPost(data)
	if err != nil {
		return BlogPostEntity{}, err
	}
	if err := validateBlogPost(e); err != nil {
		return BlogPostEntity{}, err
	}
	return scanBlogPost(bp.db.QueryRow(ctx, "INSERT INTO blog_posts (created_at, updated_at, title, slug, excerpt, body, status, cover_image_url, tags, published_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) "+conflict+" RETURNING "+blogPostColumns,
		e.CreatedAt, e.UpdatedAt, e.Title, e.Slug, e.Excerpt, e.Body, e.Status, e.CoverImageUrl, e.Tags, e.PublishedAt))
}

func (bp *BlogPosts) Upsert(ctx context.Context, data CreateBlogPostData) (BlogPostEntity, error) {
	return bp.insert(ctx, data, "ON CONFLICT (id) DO UPDATE SET title=excluded.title, slug=excluded.slug, excerpt=excluded.excerpt, body=excluded.body, status=excluded.status, cover_image_url=excluded.cover_image_url, tags=excluded.tags, published_at=excluded.published_at")
}

func (bp *BlogPosts) Update(ctx context.Context, data UpdateBlogPostData) (BlogPostEntity, error) {
	existing, err := bp.Find(ctx, data.ID)
	if err != nil {
		return BlogPostEntity{}, err
	}
	e, err := newBlogPost(CreateBlogPostData{data.Title, data.Slug, data.Excerpt, data.Body, data.Status, data.CoverImageUrl, data.Tags})
	if err != nil {
		return BlogPostEntity{}, err
	}
	if existing.PublishedAt.Valid {
		e.PublishedAt = existing.PublishedAt
	}
	if data.Status == "draft" {
		e.PublicationSchedule = existing.PublicationSchedule
	}
	if err := validateBlogPost(e); err != nil {
		return BlogPostEntity{}, err
	}
	return scanBlogPost(bp.db.QueryRow(ctx, "UPDATE blog_posts SET updated_at=$2, title=$3, slug=$4, excerpt=$5, body=$6, status=$7, cover_image_url=$8, tags=$9, published_at=$10, publication_schedule=$11 WHERE id=$1 RETURNING "+blogPostColumns,
		data.ID, e.UpdatedAt, e.Title, e.Slug, e.Excerpt, e.Body, e.Status, e.CoverImageUrl, e.Tags, e.PublishedAt, e.PublicationSchedule))
}

func (bp *BlogPosts) UpdatePublicationSchedule(ctx context.Context, e BlogPostEntity) (BlogPostEntity, error) {
	if err := validateBlogPost(e); err != nil {
		return BlogPostEntity{}, err
	}
	return scanBlogPost(bp.db.QueryRow(ctx, "UPDATE blog_posts SET publication_schedule=$2 WHERE id=$1 RETURNING "+blogPostColumns, e.ID, e.PublicationSchedule))
}

func (bp *BlogPosts) PublishIfDraft(ctx context.Context, id int32, publishAt time.Time) (bool, error) {
	e, err := bp.Find(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	schedule, err := e.PublicationScheduleData()
	if err != nil {
		return false, err
	}
	if schedule == nil || e.Status != "draft" || e.PublishedAt.Valid {
		return false, nil
	}
	e.Status = "published"
	e.PublishedAt = sql.NullTime{Time: publishAt, Valid: true}
	if err := validateBlogPost(e); err != nil {
		return false, err
	}
	result, err := bp.db.Exec(ctx, "UPDATE blog_posts SET status='published', published_at=$2, publication_schedule=NULL, updated_at=$3 WHERE id=$1 AND status='draft' AND published_at IS NULL AND publication_schedule=$4", id, publishAt, time.Now(), e.PublicationSchedule)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}

func (bp *BlogPosts) Destroy(ctx context.Context, id int32) error {
	_, err := bp.db.Exec(ctx, "DELETE FROM blog_posts WHERE id=$1", id)
	return err
}

func (bp *BlogPosts) list(ctx context.Context, query string, args ...any) ([]BlogPostEntity, error) {
	rows, err := bp.db.Query(ctx, "SELECT "+blogPostColumns+" FROM blog_posts "+query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entities := make([]BlogPostEntity, 0)
	for rows.Next() {
		e, err := scanBlogPost(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}
	return entities, rows.Err()
}

func (bp *BlogPosts) All(ctx context.Context) ([]BlogPostEntity, error) {
	return bp.list(ctx, "")
}
func (bp *BlogPosts) AllDraftOrPublished(ctx context.Context) ([]BlogPostEntity, error) {
	return bp.list(ctx, "WHERE status IN ('draft', 'published') ORDER BY created_at DESC, id DESC")
}
func (bp *BlogPosts) AllPublished(ctx context.Context) ([]BlogPostEntity, error) {
	return bp.list(ctx, "WHERE status='published' ORDER BY published_at DESC")
}
func (bp *BlogPosts) paginate(ctx context.Context, page, pageSize int64, published bool) (PaginatedBlogPosts, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	where := ""
	if published {
		where = " WHERE status='published'"
	}
	var count int64
	if err := bp.db.QueryRow(ctx, "SELECT count(*) FROM blog_posts"+where).Scan(&count); err != nil {
		return PaginatedBlogPosts{}, err
	}
	order := ""
	if published {
		order = " ORDER BY published_at DESC"
	}
	entities, err := bp.list(ctx, where+order+" LIMIT $1 OFFSET $2", pageSize, (page-1)*pageSize)
	if err != nil {
		return PaginatedBlogPosts{}, err
	}
	return PaginatedBlogPosts{entities, count, page, pageSize, (count + pageSize - 1) / pageSize}, nil
}
func (bp *BlogPosts) Paginate(ctx context.Context, page, pageSize int64) (PaginatedBlogPosts, error) {
	return bp.paginate(ctx, page, pageSize, false)
}
func (bp *BlogPosts) PaginatePublished(ctx context.Context, page, pageSize int64) (PaginatedBlogPosts, error) {
	return bp.paginate(ctx, page, pageSize, true)
}
