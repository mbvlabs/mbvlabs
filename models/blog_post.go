package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mbvlabs/internal/storage"
	"mbvlabs/internal/validation"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/uptrace/bun"
)

type BlogPostEntity struct {
	bun.BaseModel `bun:"table:blog_posts,alias:blog_posts"`
	ID            int32           `bun:"id,pk,autoincrement"`
	CreatedAt     time.Time       `bun:"created_at"`
	UpdatedAt     time.Time       `bun:"updated_at"`
	Title         string          `bun:"title"`
	Slug          string          `bun:"slug"`
	Excerpt       sql.NullString  `bun:"excerpt"`
	Body          string          `bun:"body"`
	Status        string          `bun:"status"`
	CoverImageUrl sql.NullString  `bun:"cover_image_url"`
	Tags          json.RawMessage `bun:"tags,type:jsonb"`
	PublishedAt   sql.NullTime    `bun:"published_at"`
}

func (bp BlogPostEntity) Validate() error {
	builder := validation.NewBuilder()

	builder.Required("title", bp.Title)
	builder.LenBetween("title", bp.Title, 10, 255)
	builder.Required("slug", bp.Slug)
	builder.MaxLen("slug", bp.Slug, 255)
	builder.OneOf("status", bp.Status, "draft", "published", "archived")
	builder.MinLen("excerpt", bp.Excerpt, 10)
	builder.MinLen("body", bp.Body, 10)
	builder.URL("coverImageUrl", bp.CoverImageUrl)

	var tags []string
	tagsErr := json.Unmarshal(bp.Tags, &tags)
	if tagsErr == nil {
		builder.NoBlankItems("tags", tags)
	}

	if strings.TrimSpace(bp.Slug) != "" && slug.Make(bp.Slug) != bp.Slug {
		builder.AddField("slug", "slug", "must be a valid slug")
	}

	if bp.Status == "published" {
		builder.Required("excerpt", bp.Excerpt)
		builder.Required("body", bp.Body)
		builder.Required("coverImageUrl", bp.CoverImageUrl)
		builder.Required("publishedAt", bp.PublishedAt)

		if tagsErr != nil || len(tags) == 0 {
			builder.AddField("tags", "required", "is required")
		}
	}

	return builder.Err()
}

func (bp blogPost) Find(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (BlogPostEntity, error) {
	var entity BlogPostEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return BlogPostEntity{}, ErrNotFound
		}
		return BlogPostEntity{}, err
	}

	return entity, nil
}

func (bp blogPost) FindBySlug(
	ctx context.Context,
	db storage.Executor,
	slug string,
) (BlogPostEntity, error) {
	var entity BlogPostEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("slug = ?", slug).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return BlogPostEntity{}, ErrNotFound
		}
		return BlogPostEntity{}, fmt.Errorf("find blog post by slug: %v", err)
	}

	return entity, nil
}

type CreateBlogPostData struct {
	Title         string
	Slug          string
	Excerpt       string
	Body          string
	Status        string
	CoverImageUrl string
	Tags          []string
}

func (bp blogPost) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateBlogPostData,
) (BlogPostEntity, error) {
	tagsJSON, err := json.Marshal(data.Tags)
	if err != nil {
		return BlogPostEntity{}, fmt.Errorf("marshal tags: %v", err)
	}

	now := time.Now()
	entity := BlogPostEntity{
		CreatedAt: now,
		UpdatedAt: now,
		Title:     data.Title,
		Slug:      data.Slug,
		Excerpt:   sql.NullString{String: data.Excerpt, Valid: data.Excerpt != ""},
		Body:      data.Body,
		Status:    data.Status,
		CoverImageUrl: sql.NullString{
			String: data.CoverImageUrl,
			Valid:  data.CoverImageUrl != "",
		},
		Tags: json.RawMessage(tagsJSON),
	}

	if entity.Status == "published" && !entity.PublishedAt.Valid {
		entity.PublishedAt = sql.NullTime{Time: now, Valid: true}
	}

	if err := validation.Validate(entity); err != nil {
		return BlogPostEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return BlogPostEntity{}, err
	}

	return entity, nil
}

type UpdateBlogPostData struct {
	ID            int32
	Title         string
	Slug          string
	Excerpt       string
	Body          string
	Status        string
	CoverImageUrl string
	Tags          []string
}

func (bp blogPost) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateBlogPostData,
) (BlogPostEntity, error) {
	existing, err := bp.Find(ctx, db, data.ID)
	if err != nil {
		return BlogPostEntity{}, err
	}

	tagsJSON, err := json.Marshal(data.Tags)
	if err != nil {
		return BlogPostEntity{}, fmt.Errorf("marshal tags: %v", err)
	}

	now := time.Now()
	publishedAt := existing.PublishedAt
	if data.Status == "published" && !publishedAt.Valid {
		publishedAt = sql.NullTime{Time: now, Valid: true}
	}

	entity := BlogPostEntity{
		ID:        data.ID,
		UpdatedAt: now,
		Title:     data.Title,
		Slug:      data.Slug,
		Excerpt:   sql.NullString{String: data.Excerpt, Valid: data.Excerpt != ""},
		Body:      data.Body,
		Status:    data.Status,
		CoverImageUrl: sql.NullString{
			String: data.CoverImageUrl,
			Valid:  data.CoverImageUrl != "",
		},
		Tags:        json.RawMessage(tagsJSON),
		PublishedAt: publishedAt,
	}

	if err := validation.Validate(entity); err != nil {
		return BlogPostEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("title").
		Column("slug").
		Column("excerpt").
		Column("body").
		Column("status").
		Column("cover_image_url").
		Column("tags").
		Column("published_at").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return BlogPostEntity{}, ErrNotFound
		}
		return BlogPostEntity{}, err
	}

	return entity, nil
}

func (bp blogPost) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*BlogPostEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (bp blogPost) All(ctx context.Context, db storage.Executor) ([]BlogPostEntity, error) {
	var entities []BlogPostEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

func (bp blogPost) AllPublished(ctx context.Context, db storage.Executor) ([]BlogPostEntity, error) {
	var entities []BlogPostEntity
	if err := db.NewSelect().
		Model(&entities).
		Where("status = ?", "published").
		Order("published_at DESC").
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list published blog posts: %v", err)
	}

	return entities, nil
}

func (bp blogPost) PaginatePublished(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedBlogPosts, error) {
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
		Model(&BlogPostEntity{}).
		Where("status = ?", "published").
		Count(ctx)
	if err != nil {
		return PaginatedBlogPosts{}, err
	}

	entities := make([]BlogPostEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Where("status = ?", "published").
		Order("published_at DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedBlogPosts{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedBlogPosts{
		BlogPosts:  entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

type PaginatedBlogPosts struct {
	BlogPosts  []BlogPostEntity
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func (bp blogPost) Paginate(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedBlogPosts, error) {
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
		Model(&BlogPostEntity{}).Count(ctx)
	if err != nil {
		return PaginatedBlogPosts{}, err
	}

	entities := make([]BlogPostEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedBlogPosts{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedBlogPosts{
		BlogPosts:  entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (bp blogPost) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateBlogPostData,
) (BlogPostEntity, error) {
	tagsJSON, err := json.Marshal(data.Tags)
	if err != nil {
		return BlogPostEntity{}, fmt.Errorf("marshal tags: %v", err)
	}

	now := time.Now()
	entity := BlogPostEntity{
		CreatedAt: now,
		UpdatedAt: now,
		Title:     data.Title,
		Slug:      data.Slug,
		Excerpt:   sql.NullString{String: data.Excerpt, Valid: data.Excerpt != ""},
		Body:      data.Body,
		Status:    data.Status,
		CoverImageUrl: sql.NullString{
			String: data.CoverImageUrl,
			Valid:  data.CoverImageUrl != "",
		},
		Tags: json.RawMessage(tagsJSON),
	}

	if entity.Status == "published" && !entity.PublishedAt.Valid {
		entity.PublishedAt = sql.NullTime{Time: now, Valid: true}
	}

	if err := validation.Validate(entity); err != nil {
		return BlogPostEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (id) DO UPDATE").
		Set("title = excluded.title").
		Set("slug = excluded.slug").
		Set("excerpt = excluded.excerpt").
		Set("body = excluded.body").
		Set("status = excluded.status").
		Set("cover_image_url = excluded.cover_image_url").
		Set("tags = excluded.tags").
		Set("published_at = excluded.published_at").
		Returning("*").
		Scan(ctx); err != nil {
		return BlogPostEntity{}, err
	}

	return entity, nil
}
