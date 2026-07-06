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

type ProjectEntity struct {
	bun.BaseModel `bun:"table:projects,alias:projects"`
	ID            int32           `bun:"id,pk,autoincrement"`
	CreatedAt     time.Time       `bun:"created_at"`
	UpdatedAt     time.Time       `bun:"updated_at"`
	Name          string          `bun:"name"`
	Slug          string          `bun:"slug"`
	Tagline       string          `bun:"tagline"`
	Description   sql.NullString  `bun:"description"`
	ProjectType   string          `bun:"project_type"`
	RepositoryURL sql.NullString  `bun:"repository_url"`
	LiveURL       sql.NullString  `bun:"live_url"`
	ImageURL      sql.NullString  `bun:"image_url"`
	Technologies  json.RawMessage `bun:"technologies,type:jsonb"`
	StartedAt     sql.NullTime    `bun:"started_at"`
	LaunchedAt    sql.NullTime    `bun:"launched_at"`
	PublishedAt   sql.NullTime    `bun:"published_at"`
	IsFeatured    bool            `bun:"is_featured"`
}

func (p ProjectEntity) Validate() error {
	builder := validation.NewBuilder()

	builder.Required("name", p.Name)
	builder.LenBetween("name", p.Name, 3, 255)
	builder.Required("slug", p.Slug)
	builder.MaxLen("slug", p.Slug, 255)
	builder.Required("tagline", p.Tagline)
	builder.MinLen("tagline", p.Tagline, 10)
	builder.OneOf(
		"projectType",
		p.ProjectType,
		"web-app",
		"website",
		"automation",
		"integration",
		"consulting",
		"framework",
	)
	builder.URL("repositoryUrl", p.RepositoryURL)
	builder.URL("liveUrl", p.LiveURL)
	builder.URL("imageUrl", p.ImageURL)

	var technologies []string
	tagsErr := json.Unmarshal(p.Technologies, &technologies)
	if tagsErr == nil {
		builder.NoBlankItems("technologies", technologies)
	}

	if p.IsFeatured {
		builder.Required("description", p.Description)
		builder.MinLen("description", p.Description, 10)
	}

	if strings.TrimSpace(p.Slug) != "" && slug.Make(p.Slug) != p.Slug {
		builder.AddField("slug", "slug", "must be a valid slug")
	}

	return builder.Err()
}

func (p project) Find(ctx context.Context, db storage.Executor, id int32) (ProjectEntity, error) {
	var entity ProjectEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProjectEntity{}, ErrNotFound
		}
		return ProjectEntity{}, err
	}

	return entity, nil
}

func (p project) FindBySlug(
	ctx context.Context,
	db storage.Executor,
	slug string,
) (ProjectEntity, error) {
	var entity ProjectEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("slug = ?", slug).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProjectEntity{}, ErrNotFound
		}
		return ProjectEntity{}, fmt.Errorf("find project by slug: %v", err)
	}

	return entity, nil
}

type CreateProjectData struct {
	Name          string
	Slug          string
	Tagline       string
	Description   string
	ProjectType   string
	RepositoryUrl string
	LiveUrl       string
	ImageUrl      string
	Technologies  []string
	StartedAt     time.Time
	LaunchedAt    time.Time
	PublishedAt   time.Time
	IsFeatured    bool
}

func (p project) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateProjectData,
) (ProjectEntity, error) {
	technologiesJSON, err := json.Marshal(data.Technologies)
	if err != nil {
		return ProjectEntity{}, fmt.Errorf("marshal technologies: %v", err)
	}

	entity := ProjectEntity{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      data.Name,
		Slug:      data.Slug,
		Tagline:   data.Tagline,
		Description: sql.NullString{
			String: data.Description,
			Valid:  data.Description != "",
		},
		ProjectType: data.ProjectType,
		RepositoryURL: sql.NullString{
			String: data.RepositoryUrl,
			Valid:  data.RepositoryUrl != "",
		},
		LiveURL: sql.NullString{
			String: data.LiveUrl,
			Valid:  data.LiveUrl != "",
		},
		ImageURL: sql.NullString{
			String: data.ImageUrl,
			Valid:  data.ImageUrl != "",
		},
		Technologies: json.RawMessage(technologiesJSON),
		StartedAt: sql.NullTime{
			Time:  data.StartedAt,
			Valid: !data.StartedAt.IsZero(),
		},
		LaunchedAt: sql.NullTime{
			Time:  data.LaunchedAt,
			Valid: !data.LaunchedAt.IsZero(),
		},
		PublishedAt: sql.NullTime{
			Time:  data.PublishedAt,
			Valid: !data.PublishedAt.IsZero(),
		},
		IsFeatured: data.IsFeatured,
	}

	if err := validation.Validate(entity); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return ProjectEntity{}, err
	}

	return entity, nil
}

type UpdateProjectData struct {
	ID            int32
	Name          string
	Slug          string
	Tagline       string
	Description   string
	ProjectType   string
	RepositoryUrl string
	LiveUrl       string
	ImageUrl      string
	Technologies  []string
	StartedAt     time.Time
	LaunchedAt    time.Time
	PublishedAt   time.Time
	IsFeatured    bool
}

func (p project) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateProjectData,
) (ProjectEntity, error) {
	technologiesJSON, err := json.Marshal(data.Technologies)
	if err != nil {
		return ProjectEntity{}, fmt.Errorf("marshal technologies: %v", err)
	}

	entity := ProjectEntity{
		ID:        data.ID,
		UpdatedAt: time.Now(),
		Name:      data.Name,
		Slug:      data.Slug,
		Tagline:   data.Tagline,
		Description: sql.NullString{
			String: data.Description,
			Valid:  data.Description != "",
		},
		ProjectType: data.ProjectType,
		RepositoryURL: sql.NullString{
			String: data.RepositoryUrl,
			Valid:  data.RepositoryUrl != "",
		},
		LiveURL: sql.NullString{
			String: data.LiveUrl,
			Valid:  data.LiveUrl != "",
		},
		ImageURL: sql.NullString{
			String: data.ImageUrl,
			Valid:  data.ImageUrl != "",
		},
		Technologies: json.RawMessage(technologiesJSON),
		StartedAt: sql.NullTime{
			Time:  data.StartedAt,
			Valid: !data.StartedAt.IsZero(),
		},
		LaunchedAt: sql.NullTime{
			Time:  data.LaunchedAt,
			Valid: !data.LaunchedAt.IsZero(),
		},
		PublishedAt: sql.NullTime{
			Time:  data.PublishedAt,
			Valid: !data.PublishedAt.IsZero(),
		},
		IsFeatured: data.IsFeatured,
	}

	if err := validation.Validate(entity); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("name").
		Column("slug").
		Column("tagline").
		Column("description").
		Column("project_type").
		Column("repository_url").
		Column("live_url").
		Column("image_url").
		Column("technologies").
		Column("started_at").
		Column("launched_at").
		Column("published_at").
		Column("is_featured").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProjectEntity{}, ErrNotFound
		}
		return ProjectEntity{}, err
	}

	return entity, nil
}

func (p project) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*ProjectEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (p project) All(ctx context.Context, db storage.Executor) ([]ProjectEntity, error) {
	var entities []ProjectEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

type PaginatedProjects struct {
	Projects   []ProjectEntity
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func (p project) Paginate(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedProjects, error) {
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
		Model(&ProjectEntity{}).Count(ctx)
	if err != nil {
		return PaginatedProjects{}, err
	}

	entities := make([]ProjectEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedProjects{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedProjects{
		Projects:   entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (p project) PaginatePublished(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedProjects, error) {
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
		Model(&ProjectEntity{}).
		Where("published_at IS NOT NULL").
		Count(ctx)
	if err != nil {
		return PaginatedProjects{}, fmt.Errorf("count published projects: %v", err)
	}

	entities := make([]ProjectEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Where("published_at IS NOT NULL").
		Order("published_at DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedProjects{}, fmt.Errorf("paginate published projects: %v", err)
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedProjects{
		Projects:   entities,
		TotalCount: int64(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (p project) AllPublished(ctx context.Context, db storage.Executor) ([]ProjectEntity, error) {
	var entities []ProjectEntity
	if err := db.NewSelect().
		Model(&entities).
		Where("published_at IS NOT NULL").
		Order("published_at DESC").
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list published projects: %v", err)
	}

	return entities, nil
}

func (p project) FeaturedPublished(
	ctx context.Context,
	db storage.Executor,
	limit int,
) ([]ProjectEntity, error) {
	if limit < 1 {
		limit = 3
	}
	if limit > 3 {
		limit = 3
	}

	entities := make([]ProjectEntity, 0, limit)
	if err := db.NewSelect().
		Model(&entities).
		Where("is_featured = TRUE").
		Where("published_at IS NOT NULL").
		Order("published_at DESC").
		Limit(limit).
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list featured published projects: %v", err)
	}

	return entities, nil
}

func (p project) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateProjectData,
) (ProjectEntity, error) {
	technologiesJSON, err := json.Marshal(data.Technologies)
	if err != nil {
		return ProjectEntity{}, fmt.Errorf("marshal technologies: %v", err)
	}

	entity := ProjectEntity{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      data.Name,
		Slug:      data.Slug,
		Tagline:   data.Tagline,
		Description: sql.NullString{
			String: data.Description,
			Valid:  data.Description != "",
		},
		ProjectType: data.ProjectType,
		RepositoryURL: sql.NullString{
			String: data.RepositoryUrl,
			Valid:  data.RepositoryUrl != "",
		},
		LiveURL: sql.NullString{
			String: data.LiveUrl,
			Valid:  data.LiveUrl != "",
		},
		ImageURL: sql.NullString{
			String: data.ImageUrl,
			Valid:  data.ImageUrl != "",
		},
		Technologies: json.RawMessage(technologiesJSON),
		StartedAt: sql.NullTime{
			Time:  data.StartedAt,
			Valid: !data.StartedAt.IsZero(),
		},
		LaunchedAt: sql.NullTime{
			Time:  data.LaunchedAt,
			Valid: !data.LaunchedAt.IsZero(),
		},
		PublishedAt: sql.NullTime{
			Time:  data.PublishedAt,
			Valid: !data.PublishedAt.IsZero(),
		},
		IsFeatured: data.IsFeatured,
	}

	if err := validation.Validate(entity); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (id) DO UPDATE").
		Set("name = excluded.name").
		Set("slug = excluded.slug").
		Set("tagline = excluded.tagline").
		Set("description = excluded.description").
		Set("project_type = excluded.project_type").
		Set("repository_url = excluded.repository_url").
		Set("live_url = excluded.live_url").
		Set("image_url = excluded.image_url").
		Set("technologies = excluded.technologies").
		Set("started_at = excluded.started_at").
		Set("launched_at = excluded.launched_at").
		Set("published_at = excluded.published_at").
		Set("is_featured = excluded.is_featured").
		Returning("*").
		Scan(ctx); err != nil {
		return ProjectEntity{}, err
	}

	return entity, nil
}

func (p project) CountFeatured(
	ctx context.Context,
	db storage.Executor,
	exceptID int32,
) (int64, error) {
	query := db.NewSelect().
		Model((*ProjectEntity)(nil)).
		Where("is_featured = TRUE")
	if exceptID > 0 {
		query = query.Where("id <> ?", exceptID)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count featured projects: %v", err)
	}

	return int64(count), nil
}
