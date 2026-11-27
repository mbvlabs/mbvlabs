package domains

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"mbvlabs/domains/internal/queries"

	"github.com/gosimple/slug"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type Projects struct{ queries *queries.Queries }

func NewProjects(db storage.Connection) *Projects { return &Projects{queries: queries.New(db)} }

func (p *Projects) InTx(tx storage.Tx) *Projects {
	return &Projects{queries: queries.New(storage.Unwrap(tx))}
}

type ProjectEntity struct {
	ID            int32           `andurel:"id"`
	CreatedAt     time.Time       `andurel:"created_at"`
	UpdatedAt     time.Time       `andurel:"updated_at"`
	Name          string          `andurel:"name"`
	Slug          string          `andurel:"slug"`
	Tagline       string          `andurel:"tagline"`
	Description   sql.NullString  `andurel:"description"`
	ProjectType   string          `andurel:"project_type"`
	RepositoryURL sql.NullString  `andurel:"repository_url"`
	LiveURL       sql.NullString  `andurel:"live_url"`
	ImageURL      sql.NullString  `andurel:"image_url"`
	Technologies  json.RawMessage `andurel:"technologies"`
	StartedAt     sql.NullTime    `andurel:"started_at"`
	LaunchedAt    sql.NullTime    `andurel:"launched_at"`
	PublishedAt   sql.NullTime    `andurel:"published_at"`
	IsFeatured    bool            `andurel:"is_featured"`
}

func (p ProjectEntity) Validate() error {
	b := validation.NewBuilder()
	b.Required("name", p.Name)
	b.LenBetween("name", p.Name, 3, 255)
	b.Required("slug", p.Slug)
	b.MaxLen("slug", p.Slug, 255)
	b.Required("tagline", p.Tagline)
	b.MinLen("tagline", p.Tagline, 10)
	b.OneOf("projectType", p.ProjectType, "web-app", "website", "automation", "integration", "consulting", "framework")
	b.URL("repositoryUrl", p.RepositoryURL)
	b.URL("liveUrl", p.LiveURL)
	b.URL("imageUrl", p.ImageURL)
	var technologies []string
	if json.Unmarshal(p.Technologies, &technologies) == nil {
		b.NoBlankItems("technologies", technologies)
	}
	if p.IsFeatured {
		b.Required("description", p.Description)
		b.MinLen("description", p.Description, 10)
	}
	if strings.TrimSpace(p.Slug) != "" && slug.Make(p.Slug) != p.Slug {
		b.AddField("slug", "slug", "must be a valid slug")
	}
	return b.Err()
}

type CreateProjectData struct {
	Name, Slug, Tagline, Description, ProjectType string
	RepositoryUrl, LiveUrl, ImageUrl              string
	Technologies                                  []string
	StartedAt, LaunchedAt, PublishedAt            time.Time
	IsFeatured                                    bool
}
type UpdateProjectData struct {
	ID                                            int32
	Name, Slug, Tagline, Description, ProjectType string
	RepositoryUrl, LiveUrl, ImageUrl              string
	Technologies                                  []string
	StartedAt, LaunchedAt, PublishedAt            time.Time
	IsFeatured                                    bool
}
type PaginatedProjects struct {
	Projects                               []ProjectEntity
	TotalCount, Page, PageSize, TotalPages int64
}

func projectFromData(d CreateProjectData) (ProjectEntity, error) {
	technologies, err := json.Marshal(d.Technologies)
	if err != nil {
		return ProjectEntity{}, err
	}
	now := time.Now()
	return ProjectEntity{CreatedAt: now, UpdatedAt: now, Name: d.Name, Slug: d.Slug, Tagline: d.Tagline, Description: workString(d.Description), ProjectType: d.ProjectType,
		RepositoryURL: workString(d.RepositoryUrl), LiveURL: workString(d.LiveUrl), ImageURL: workString(d.ImageUrl), Technologies: technologies,
		StartedAt: workTime(d.StartedAt), LaunchedAt: workTime(d.LaunchedAt), PublishedAt: workTime(d.PublishedAt), IsFeatured: d.IsFeatured}, nil
}

func projectParams(e ProjectEntity) queries.CreateProjectParams {
	return queries.CreateProjectParams{CreatedAt: e.UpdatedAt, Name: e.Name, Slug: e.Slug, Tagline: e.Tagline,
		Description: contentText(e.Description), ProjectType: e.ProjectType, RepositoryUrl: contentText(e.RepositoryURL),
		LiveUrl: contentText(e.LiveURL), ImageUrl: contentText(e.ImageURL), Technologies: e.Technologies,
		StartedAt: contentDate(e.StartedAt), LaunchedAt: contentDate(e.LaunchedAt), PublishedAt: contentTimestamp(e.PublishedAt), IsFeatured: e.IsFeatured}
}
func (p *Projects) Find(ctx context.Context, id int32) (ProjectEntity, error) {
	return p.queries.GetProject[ProjectEntity](ctx, id)
}
func (p *Projects) Latest(ctx context.Context, limit int) ([]ProjectEntity, error) {
	return contentRows(p.queries.LatestProjects[ProjectEntity](ctx, int64(limit)))
}
func (p *Projects) FindBySlug(ctx context.Context, slug string) (ProjectEntity, error) {
	return p.queries.GetProjectBySlug[ProjectEntity](ctx, slug)
}
func (p *Projects) Create(ctx context.Context, d CreateProjectData) (ProjectEntity, error) {
	e, err := projectFromData(d)
	if err != nil {
		return ProjectEntity{}, err
	}
	if err = e.Validate(); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return p.queries.CreateProject[ProjectEntity](ctx, projectParams(e))
}
func (p *Projects) Upsert(ctx context.Context, d CreateProjectData) (ProjectEntity, error) {
	e, err := projectFromData(d)
	if err != nil {
		return ProjectEntity{}, err
	}
	if err = e.Validate(); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return p.queries.UpsertProject[ProjectEntity](ctx, queries.UpsertProjectParams(projectParams(e)))
}
func (p *Projects) Update(ctx context.Context, d UpdateProjectData) (ProjectEntity, error) {
	e, err := projectFromData(CreateProjectData{Name: d.Name, Slug: d.Slug, Tagline: d.Tagline, Description: d.Description, ProjectType: d.ProjectType, RepositoryUrl: d.RepositoryUrl, LiveUrl: d.LiveUrl, ImageUrl: d.ImageUrl, Technologies: d.Technologies, StartedAt: d.StartedAt, LaunchedAt: d.LaunchedAt, PublishedAt: d.PublishedAt, IsFeatured: d.IsFeatured})
	if err != nil {
		return ProjectEntity{}, err
	}
	if err = e.Validate(); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}
	a := projectParams(e)
	return p.queries.UpdateProject[ProjectEntity](ctx, queries.UpdateProjectParams{
		ID: d.ID, UpdatedAt: a.CreatedAt, Name: a.Name, Slug: a.Slug, Tagline: a.Tagline,
		Description: a.Description, ProjectType: a.ProjectType, RepositoryUrl: a.RepositoryUrl, LiveUrl: a.LiveUrl,
		ImageUrl: a.ImageUrl, Technologies: a.Technologies, StartedAt: a.StartedAt, LaunchedAt: a.LaunchedAt,
		PublishedAt: a.PublishedAt, IsFeatured: a.IsFeatured,
	})
}

// EditEditorial must be called on a transaction-bound store.
func (p *Projects) EditEditorial(ctx context.Context, slug string, edit func(*ProjectEntity)) (ProjectEntity, error) {
	e, err := p.queries.GetProjectBySlugForUpdate[ProjectEntity](ctx, slug)
	if err != nil {
		return ProjectEntity{}, err
	}
	edit(&e)
	if err := e.Validate(); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return p.queries.EditProject[ProjectEntity](ctx, queries.EditProjectParams{
		ID: e.ID, UpdatedAt: time.Now(), Name: e.Name, Tagline: e.Tagline, Description: contentText(e.Description),
		RepositoryUrl: contentText(e.RepositoryURL), LiveUrl: contentText(e.LiveURL), ImageUrl: contentText(e.ImageURL), Technologies: e.Technologies,
	})
}
func (p *Projects) LockFeatured(ctx context.Context) error {
	return p.queries.LockFeaturedProjects(ctx)
}
func (p *Projects) Destroy(ctx context.Context, id int32) error {
	return p.queries.DeleteProject(ctx, id)
}
func (p *Projects) All(ctx context.Context) ([]ProjectEntity, error) {
	return contentRows(p.queries.ListProjects[ProjectEntity](ctx))
}
func (p *Projects) AllPublished(ctx context.Context) ([]ProjectEntity, error) {
	return contentRows(p.queries.ListPublishedProjects[ProjectEntity](ctx))
}
func (p *Projects) FeaturedPublished(ctx context.Context, limit int) ([]ProjectEntity, error) {
	if limit < 1 || limit > 3 {
		limit = 3
	}
	return contentRows(p.queries.ListFeaturedProjects[ProjectEntity](ctx, int32(limit)))
}
func (p *Projects) CountFeatured(ctx context.Context, exceptID int32) (int64, error) {
	return p.queries.CountFeaturedProjects(ctx, exceptID)
}
func (p *Projects) paginate(ctx context.Context, page, size int64, published bool) (PaginatedProjects, error) {
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
	var entries []ProjectEntity
	var err error
	if published {
		total, err = p.queries.CountPublishedProjects(ctx)
		if err == nil {
			entries, err = p.queries.PagePublishedProjects[ProjectEntity](ctx, queries.PagePublishedProjectsParams{Column1: size, Column2: (page - 1) * size})
		}
	} else {
		total, err = p.queries.CountProjects(ctx)
		if err == nil {
			entries, err = p.queries.PageProjects[ProjectEntity](ctx, queries.PageProjectsParams{Column1: size, Column2: (page - 1) * size})
		}
	}
	if err != nil {
		return PaginatedProjects{}, err
	}
	if entries == nil {
		entries = []ProjectEntity{}
	}
	return PaginatedProjects{entries, total, page, size, (total + size - 1) / size}, nil
}
func (p *Projects) Paginate(ctx context.Context, page, size int64) (PaginatedProjects, error) {
	return p.paginate(ctx, page, size, false)
}
func (p *Projects) PaginatePublished(ctx context.Context, page, size int64) (PaginatedProjects, error) {
	return p.paginate(ctx, page, size, true)
}
