package domains

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type Projects struct{ db contentDB }

func NewProjects(db storage.Connection) *Projects { return &Projects{db: db} }

func (p *Projects) InTx(tx storage.Tx) *Projects { return &Projects{db: storage.Unwrap(tx)} }

type ProjectEntity struct {
	ID                                 int32
	CreatedAt, UpdatedAt               time.Time
	Name, Slug, Tagline                string
	Description                        sql.NullString
	ProjectType                        string
	RepositoryURL, LiveURL, ImageURL   sql.NullString
	Technologies                       json.RawMessage
	StartedAt, LaunchedAt, PublishedAt sql.NullTime
	IsFeatured                         bool
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

const projectColumns = "id, created_at, updated_at, name, slug, tagline, description, project_type, repository_url, live_url, image_url, technologies, started_at, launched_at, published_at, is_featured"

func scanProject(row interface{ Scan(...any) error }) (ProjectEntity, error) {
	var p ProjectEntity
	err := row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Name, &p.Slug, &p.Tagline, &p.Description, &p.ProjectType, &p.RepositoryURL, &p.LiveURL, &p.ImageURL, &p.Technologies, &p.StartedAt, &p.LaunchedAt, &p.PublishedAt, &p.IsFeatured)
	return storage.Row(p, err)
}
func projectArgs(p ProjectEntity) []any {
	return []any{p.UpdatedAt, p.Name, p.Slug, p.Tagline, p.Description, p.ProjectType, p.RepositoryURL, p.LiveURL, p.ImageURL, p.Technologies, p.StartedAt, p.LaunchedAt, p.PublishedAt, p.IsFeatured}
}
func (p *Projects) Find(ctx context.Context, id int32) (ProjectEntity, error) {
	return scanProject(p.db.QueryRow(ctx, "SELECT "+projectColumns+" FROM projects WHERE id=$1", id))
}
func (p *Projects) FindBySlug(ctx context.Context, s string) (ProjectEntity, error) {
	return scanProject(p.db.QueryRow(ctx, "SELECT "+projectColumns+" FROM projects WHERE slug=$1", s))
}
func (p *Projects) Create(ctx context.Context, d CreateProjectData) (ProjectEntity, error) {
	e, err := projectFromData(d)
	if err != nil {
		return ProjectEntity{}, err
	}
	if err = e.Validate(); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return scanProject(p.db.QueryRow(ctx, `INSERT INTO projects (created_at,updated_at,name,slug,tagline,description,project_type,repository_url,live_url,image_url,technologies,started_at,launched_at,published_at,is_featured)
	VALUES ($1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING `+projectColumns, projectArgs(e)...))
}
func (p *Projects) Upsert(ctx context.Context, d CreateProjectData) (ProjectEntity, error) {
	e, err := projectFromData(d)
	if err != nil {
		return ProjectEntity{}, err
	}
	if err = e.Validate(); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return scanProject(p.db.QueryRow(ctx, `INSERT INTO projects (created_at,updated_at,name,slug,tagline,description,project_type,repository_url,live_url,image_url,technologies,started_at,launched_at,published_at,is_featured)
	VALUES ($1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	ON CONFLICT (id) DO UPDATE SET name=excluded.name,slug=excluded.slug,tagline=excluded.tagline,description=excluded.description,project_type=excluded.project_type,repository_url=excluded.repository_url,live_url=excluded.live_url,image_url=excluded.image_url,technologies=excluded.technologies,started_at=excluded.started_at,launched_at=excluded.launched_at,published_at=excluded.published_at,is_featured=excluded.is_featured RETURNING `+projectColumns, projectArgs(e)...))
}
func (p *Projects) Update(ctx context.Context, d UpdateProjectData) (ProjectEntity, error) {
	e, err := projectFromData(CreateProjectData{Name: d.Name, Slug: d.Slug, Tagline: d.Tagline, Description: d.Description, ProjectType: d.ProjectType, RepositoryUrl: d.RepositoryUrl, LiveUrl: d.LiveUrl, ImageUrl: d.ImageUrl, Technologies: d.Technologies, StartedAt: d.StartedAt, LaunchedAt: d.LaunchedAt, PublishedAt: d.PublishedAt, IsFeatured: d.IsFeatured})
	if err != nil {
		return ProjectEntity{}, err
	}
	if err = e.Validate(); err != nil {
		return ProjectEntity{}, errors.Join(ErrDomainValidation, err)
	}
	args := append(projectArgs(e), d.ID)
	return scanProject(p.db.QueryRow(ctx, `UPDATE projects SET updated_at=$1,name=$2,slug=$3,tagline=$4,description=$5,project_type=$6,repository_url=$7,live_url=$8,image_url=$9,technologies=$10,started_at=$11,launched_at=$12,published_at=$13,is_featured=$14 WHERE id=$15 RETURNING `+projectColumns, args...))
}
func (p *Projects) Destroy(ctx context.Context, id int32) error {
	_, err := p.db.Exec(ctx, "DELETE FROM projects WHERE id=$1", id)
	return err
}
func (p *Projects) list(ctx context.Context, suffix string, args ...any) ([]ProjectEntity, error) {
	rows, err := p.db.Query(ctx, "SELECT "+projectColumns+" FROM projects "+suffix, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]ProjectEntity, 0)
	for rows.Next() {
		e, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (p *Projects) All(ctx context.Context) ([]ProjectEntity, error) { return p.list(ctx, "") }
func (p *Projects) AllPublished(ctx context.Context) ([]ProjectEntity, error) {
	return p.list(ctx, "WHERE published_at IS NOT NULL ORDER BY published_at DESC")
}
func (p *Projects) FeaturedPublished(ctx context.Context, limit int) ([]ProjectEntity, error) {
	if limit < 1 || limit > 3 {
		limit = 3
	}
	return p.list(ctx, "WHERE published_at IS NOT NULL AND is_featured ORDER BY published_at DESC LIMIT $1", limit)
}
func (p *Projects) CountFeatured(ctx context.Context, exceptID int32) (int64, error) {
	var n int64
	err := p.db.QueryRow(ctx, "SELECT count(*) FROM projects WHERE is_featured AND ($1::integer <= 0 OR id <> $1)", exceptID).Scan(&n)
	return n, err
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
	where := ""
	order := ""
	if published {
		where = " WHERE published_at IS NOT NULL"
		order = " ORDER BY published_at DESC"
	}
	var total int64
	if err := p.db.QueryRow(ctx, "SELECT count(*) FROM projects"+where).Scan(&total); err != nil {
		return PaginatedProjects{}, err
	}
	entries, err := p.list(ctx, where+order+" LIMIT $1 OFFSET $2", size, (page-1)*size)
	if err != nil {
		return PaginatedProjects{}, err
	}
	return PaginatedProjects{entries, total, page, size, (total + size - 1) / size}, nil
}
func (p *Projects) Paginate(ctx context.Context, page, size int64) (PaginatedProjects, error) {
	return p.paginate(ctx, page, size, false)
}
func (p *Projects) PaginatePublished(ctx context.Context, page, size int64) (PaginatedProjects, error) {
	return p.paginate(ctx, page, size, true)
}
