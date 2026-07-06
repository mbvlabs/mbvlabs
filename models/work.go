package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mbvlabs/internal/storage"
	"mbvlabs/internal/validation"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/uptrace/bun"
)

type WorkEntity struct {
	bun.BaseModel `bun:"table:works,alias:works"`

	ID             int32          `bun:"id,pk,autoincrement"`
	CreatedAt      time.Time      `bun:"created_at"`
	UpdatedAt      time.Time      `bun:"updated_at"`
	Title          string         `bun:"title"`
	Slug           string         `bun:"slug"`
	ClientName     sql.NullString `bun:"client_name"`
	ClientIndustry sql.NullString `bun:"client_industry"`
	ClientURL      sql.NullString `bun:"client_url"`
	ClientLogoURL  sql.NullString `bun:"client_logo_url"`
	Summary        string         `bun:"summary"`
	CoverImageURL  sql.NullString `bun:"cover_image_url"`
	Specialisms    []string       `bun:"specialisms,array"`
	Platforms      []string       `bun:"platforms,array"`
	Technologies   []string       `bun:"technologies,array"`
	Challenge      string         `bun:"challenge"`
	Approach       string         `bun:"approach"`
	Deliverables   string         `bun:"deliverables"`
	Outcome        string         `bun:"outcome"`
	Content        string         `bun:"content"`
	StartedAt      sql.NullTime   `bun:"started_at"`
	CompletedAt    sql.NullTime   `bun:"completed_at"`
	Status         StatusEnum     `bun:"status,type:text"`
	PublishedAt    sql.NullTime   `bun:"published_at"`
	IsFeatured     bool           `bun:"is_featured"`
}

func (we WorkEntity) Validate() error {
	builder := validation.NewBuilder()

	builder.Required("title", we.Title)
	builder.LenBetween("title", we.Title, 10, 255)
	builder.Required("slug", we.Slug)
	builder.MaxLen("slug", we.Slug, 255)
	builder.OneOf("status", we.Status, Draft.String(), Published.String())
	builder.MinLen("summary", we.Summary, 10)
	builder.URL("clientUrl", we.ClientURL)
	builder.URL("clientLogoUrl", we.ClientLogoURL)
	builder.URL("coverImageUrl", we.CoverImageURL)
	builder.NoBlankItems("specialisms", we.Specialisms)
	builder.NoBlankItems("platforms", we.Platforms)
	builder.NoBlankItems("technologies", we.Technologies)
	builder.TimeBeforeOrEqual("startedAt", we.StartedAt, "completedAt", we.CompletedAt)

	if strings.TrimSpace(we.Slug) != "" && slug.Make(we.Slug) != we.Slug {
		builder.AddField("slug", "slug", "must be a valid slug")
	}

	if we.Status == Published {
		builder.Required("clientName", we.ClientName)
		builder.Required("clientIndustry", we.ClientIndustry)
		builder.Required("clientUrl", we.ClientURL)
		builder.Required("clientLogoUrl", we.ClientLogoURL)
		builder.Required("summary", we.Summary)
		builder.Required("coverImageUrl", we.CoverImageURL)
		builder.Required("specialisms", we.Specialisms)
		builder.Required("platforms", we.Platforms)
		builder.Required("technologies", we.Technologies)
		builder.Required("challenge", we.Challenge)
		builder.Required("approach", we.Approach)
		builder.Required("deliverables", we.Deliverables)
		builder.Required("outcome", we.Outcome)
		builder.Required("content", we.Content)
		builder.Required("startedAt", we.StartedAt)
		builder.Required("completedAt", we.CompletedAt)
		builder.Required("publishedAt", we.PublishedAt)
	}

	return builder.Err()
}

func (cs work) Find(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (WorkEntity, error) {
	var entity WorkEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WorkEntity{}, ErrNotFound
		}
		return WorkEntity{}, fmt.Errorf("find work: %v", err)
	}

	return entity, nil
}

func (cs work) FindBySlug(
	ctx context.Context,
	db storage.Executor,
	slug string,
) (WorkEntity, error) {
	var entity WorkEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("slug = ?", slug).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WorkEntity{}, ErrNotFound
		}
		return WorkEntity{}, fmt.Errorf("find work by slug: %v", err)
	}

	return entity, nil
}

type CreateWorkData struct {
	Title            string
	Slug             string
	ClientName       string
	ClientIndustry   string
	ClientURL        string
	ClientLogoURL    string
	Summary          string
	Content          string
	CoverImageURL    string
	Specialisms      []string
	Platforms        []string
	Technologies     []string
	Challenge        string
	Approach         string
	Deliverables     string
	Outcome          string
	TestimonialQuote string
	TestimonialName  string
	TestimonialRole  string
	StartedAt        time.Time
	CompletedAt      time.Time
	Status           StatusEnum
	PublishedAt      time.Time
	IsFeatured       bool
}

func (cs work) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateWorkData,
) (WorkEntity, error) {
	entity := WorkEntity{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Title:      data.Title,
		Slug:       data.Slug,
		ClientName: sql.NullString{String: data.ClientName, Valid: data.ClientName != ""},
		ClientIndustry: sql.NullString{
			String: data.ClientIndustry,
			Valid:  data.ClientIndustry != "",
		},
		ClientURL: sql.NullString{String: data.ClientURL, Valid: data.ClientURL != ""},
		ClientLogoURL: sql.NullString{
			String: data.ClientLogoURL,
			Valid:  data.ClientLogoURL != "",
		},
		Summary: data.Summary,
		Content: data.Content,
		CoverImageURL: sql.NullString{
			String: data.CoverImageURL,
			Valid:  data.CoverImageURL != "",
		},
		Specialisms:  data.Specialisms,
		Platforms:    data.Platforms,
		Technologies: data.Technologies,
		Challenge:    data.Challenge,
		Approach:     data.Approach,
		Deliverables: data.Deliverables,
		Outcome:      data.Outcome,
		StartedAt:    sql.NullTime{Time: data.StartedAt, Valid: !data.StartedAt.IsZero()},
		CompletedAt:  sql.NullTime{Time: data.CompletedAt, Valid: !data.CompletedAt.IsZero()},
		Status:       data.Status,
		PublishedAt:  sql.NullTime{Time: data.PublishedAt, Valid: !data.PublishedAt.IsZero()},
		IsFeatured:   data.IsFeatured,
	}

	if entity.Status == Published && !entity.PublishedAt.Valid {
		entity.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	if err := validation.Validate(entity); err != nil {
		return WorkEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return WorkEntity{}, fmt.Errorf("create work: %v", err)
	}

	return entity, nil
}

type UpdateWorkData struct {
	ID               int32
	Title            string
	Slug             string
	ClientName       string
	ClientIndustry   string
	ClientURL        string
	ClientLogoURL    string
	Summary          string
	Content          string
	CoverImageURL    string
	Specialisms      []string
	Platforms        []string
	Technologies     []string
	Challenge        string
	Approach         string
	Deliverables     string
	Outcome          string
	TestimonialQuote string
	TestimonialName  string
	TestimonialRole  string
	StartedAt        time.Time
	CompletedAt      time.Time
	Status           StatusEnum
	PublishedAt      time.Time
	IsFeatured       bool
}

func (cs work) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateWorkData,
) (WorkEntity, error) {
	existing, err := cs.Find(ctx, db, data.ID)
	if err != nil {
		return WorkEntity{}, err
	}

	publishedAt := existing.PublishedAt
	if data.Status == Published && !publishedAt.Valid {
		publishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	entity := WorkEntity{
		ID:         data.ID,
		UpdatedAt:  time.Now(),
		Title:      data.Title,
		Slug:       data.Slug,
		ClientName: sql.NullString{String: data.ClientName, Valid: data.ClientName != ""},
		ClientIndustry: sql.NullString{
			String: data.ClientIndustry,
			Valid:  data.ClientIndustry != "",
		},
		ClientURL: sql.NullString{String: data.ClientURL, Valid: data.ClientURL != ""},
		ClientLogoURL: sql.NullString{
			String: data.ClientLogoURL,
			Valid:  data.ClientLogoURL != "",
		},
		Summary: data.Summary,
		Content: data.Content,
		CoverImageURL: sql.NullString{
			String: data.CoverImageURL,
			Valid:  data.CoverImageURL != "",
		},
		Specialisms:  data.Specialisms,
		Platforms:    data.Platforms,
		Technologies: data.Technologies,
		Challenge:    data.Challenge,
		Approach:     data.Approach,
		Deliverables: data.Deliverables,
		Outcome:      data.Outcome,
		StartedAt:    sql.NullTime{Time: data.StartedAt, Valid: !data.StartedAt.IsZero()},
		CompletedAt:  sql.NullTime{Time: data.CompletedAt, Valid: !data.CompletedAt.IsZero()},
		Status:       data.Status,
		PublishedAt:  publishedAt,
		IsFeatured:   data.IsFeatured,
	}

	if err := validation.Validate(entity); err != nil {
		return WorkEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("title").
		Column("slug").
		Column("client_name").
		Column("client_industry").
		Column("client_url").
		Column("client_logo_url").
		Column("summary").
		Column("content").
		Column("cover_image_url").
		Column("specialisms").
		Column("platforms").
		Column("technologies").
		Column("challenge").
		Column("approach").
		Column("deliverables").
		Column("outcome").
		Column("started_at").
		Column("completed_at").
		Column("status").
		Column("published_at").
		Column("is_featured").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WorkEntity{}, ErrNotFound
		}
		return WorkEntity{}, fmt.Errorf("update work: %v", err)
	}

	return entity, nil
}

func (cs work) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*WorkEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("destroy work: %v", err)
	}

	return nil
}

func (cs work) All(ctx context.Context, db storage.Executor) ([]WorkEntity, error) {
	var entities []WorkEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list works: %v", err)
	}

	return entities, nil
}

func (cs work) CountFeatured(
	ctx context.Context,
	db storage.Executor,
	exceptID int32,
) (int64, error) {
	query := db.NewSelect().
		Model((*WorkEntity)(nil)).
		Where("is_featured = TRUE")
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count featured works: %v", err)
	}

	return int64(count), nil
}

func (cs work) PaginatePublished(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedWorks, error) {
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
		Model(&WorkEntity{}).
		Where("status = ?", Published.String()).
		Count(ctx)
	if err != nil {
		return PaginatedWorks{}, fmt.Errorf("count published works: %v", err)
	}

	entities := make([]WorkEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Where("status = ?", Published.String()).
		Order("completed_at DESC NULLS LAST").
		Order("started_at DESC NULLS LAST").
		Order("published_at DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedWorks{}, fmt.Errorf("paginate published works: %v", err)
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedWorks{
		Works:      entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (cs work) AllPublished(ctx context.Context, db storage.Executor) ([]WorkEntity, error) {
	var entities []WorkEntity
	if err := db.NewSelect().
		Model(&entities).
		Where("status = ?", Published.String()).
		Order("completed_at DESC NULLS LAST").
		Order("started_at DESC NULLS LAST").
		Order("published_at DESC").
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list published works: %v", err)
	}

	return entities, nil
}

func (cs work) FeaturedPublished(
	ctx context.Context,
	db storage.Executor,
	limit int,
) ([]WorkEntity, error) {
	if limit < 1 {
		limit = 3
	}
	if limit > 3 {
		limit = 3
	}

	entities := make([]WorkEntity, 0, limit)
	if err := db.NewSelect().
		Model(&entities).
		Where("is_featured = TRUE").
		Where("status = ?", Published.String()).
		Order("completed_at DESC NULLS LAST").
		Order("started_at DESC NULLS LAST").
		Order("published_at DESC").
		Limit(limit).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list featured published works: %v", err)
	}

	return entities, nil
}

type PaginatedWorks struct {
	Works      []WorkEntity
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func (cs work) Paginate(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedWorks, error) {
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
		Model(&WorkEntity{}).Count(ctx)
	if err != nil {
		return PaginatedWorks{}, fmt.Errorf("count works: %v", err)
	}

	entities := make([]WorkEntity, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedWorks{}, fmt.Errorf("paginate works: %v", err)
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedWorks{
		Works:      entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (cs work) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateWorkData,
) (WorkEntity, error) {
	entity := WorkEntity{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Title:      data.Title,
		Slug:       data.Slug,
		ClientName: sql.NullString{String: data.ClientName, Valid: data.ClientName != ""},
		ClientIndustry: sql.NullString{
			String: data.ClientIndustry,
			Valid:  data.ClientIndustry != "",
		},
		ClientURL: sql.NullString{String: data.ClientURL, Valid: data.ClientURL != ""},
		ClientLogoURL: sql.NullString{
			String: data.ClientLogoURL,
			Valid:  data.ClientLogoURL != "",
		},
		Summary: data.Summary,
		Content: data.Content,
		CoverImageURL: sql.NullString{
			String: data.CoverImageURL,
			Valid:  data.CoverImageURL != "",
		},
		Specialisms:  data.Specialisms,
		Platforms:    data.Platforms,
		Technologies: data.Technologies,
		Challenge:    data.Challenge,
		Approach:     data.Approach,
		Deliverables: data.Deliverables,
		Outcome:      data.Outcome,
		StartedAt:    sql.NullTime{Time: data.StartedAt, Valid: !data.StartedAt.IsZero()},
		CompletedAt:  sql.NullTime{Time: data.CompletedAt, Valid: !data.CompletedAt.IsZero()},
		Status:       data.Status,
		PublishedAt:  sql.NullTime{Time: data.PublishedAt, Valid: !data.PublishedAt.IsZero()},
		IsFeatured:   data.IsFeatured,
	}

	if entity.Status == Published && !entity.PublishedAt.Valid {
		entity.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	if err := validation.Validate(entity); err != nil {
		return WorkEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (id) DO UPDATE").
		Set("title = excluded.title").
		Set("slug = excluded.slug").
		Set("client_name = excluded.client_name").
		Set("client_industry = excluded.client_industry").
		Set("client_url = excluded.client_url").
		Set("client_logo_url = excluded.client_logo_url").
		Set("summary = excluded.summary").
		Set("content = excluded.content").
		Set("cover_image_url = excluded.cover_image_url").
		Set("specialisms = excluded.specialisms").
		Set("platforms = excluded.platforms").
		Set("technologies = excluded.technologies").
		Set("challenge = excluded.challenge").
		Set("approach = excluded.approach").
		Set("deliverables = excluded.deliverables").
		Set("outcome = excluded.outcome").
		Set("started_at = excluded.started_at").
		Set("completed_at = excluded.completed_at").
		Set("status = excluded.status").
		Set("published_at = excluded.published_at").
		Set("is_featured = excluded.is_featured").
		Returning("*").
		Scan(ctx); err != nil {
		return WorkEntity{}, fmt.Errorf("upsert work: %v", err)
	}

	return entity, nil
}
