package domains

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/mail"
	"strings"
	"time"

	"mbvlabs/domains/internal/queries"

	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type ProjectInquiries struct{ queries *queries.Queries }

func NewProjectInquiries(db storage.Connection) *ProjectInquiries {
	return &ProjectInquiries{queries: queries.New(db)}
}

func (p *ProjectInquiries) Latest(ctx context.Context, limit int) ([]ProjectInquiryEntity, error) {
	items, err := p.queries.LatestProjectInquiries[ProjectInquiryEntity](ctx, int64(limit))
	if err == nil && items == nil {
		items = []ProjectInquiryEntity{}
	}
	return items, err
}

type ProjectInquiryEntity struct {
	ID          int32           `andurel:"id"`
	CreatedAt   time.Time       `andurel:"created_at"`
	UpdatedAt   time.Time       `andurel:"updated_at"`
	Name        string          `andurel:"name"`
	Email       string          `andurel:"email"`
	Company     sql.NullString  `andurel:"company"`
	Role        sql.NullString  `andurel:"role"`
	ProjectType sql.NullString  `andurel:"project_type"`
	Timeline    sql.NullString  `andurel:"timeline"`
	Message     string          `andurel:"message"`
	Source      sql.NullString  `andurel:"source"`
	Status      string          `andurel:"status"`
	Metadata    json.RawMessage `andurel:"metadata"`
}

func (e ProjectInquiryEntity) Validate() error {
	b := validation.NewBuilder()
	b.Required("name", e.Name)
	b.LenBetween("name", e.Name, 2, 255)
	b.Required("email", e.Email)
	b.MaxLen("email", e.Email, 255)
	b.MaxLen("company", e.Company, 255)
	b.MaxLen("role", e.Role, 255)
	b.OneOf("projectType", e.ProjectType, "web-app", "website", "automation", "integration", "consulting", "open-source", "modernization")
	b.MaxLen("timeline", e.Timeline, 100)
	b.Required("message", e.Message)
	b.LenBetween("message", e.Message, 10, 5000)
	b.MaxLen("source", e.Source, 255)
	b.Required("status", e.Status)
	b.OneOf("status", e.Status, "new", "contacted", "qualified", "proposal", "won", "lost")
	b.Required("metadata", e.Metadata)
	if e.Email != "" {
		if _, err := mail.ParseAddress(e.Email); err != nil {
			b.Add("email", "email", "must be a valid email address")
		}
	}
	return b.Err()
}

type CreateProjectInquiryData struct {
	Name        string
	Email       string
	Company     sql.NullString
	Role        sql.NullString
	ProjectType sql.NullString
	Timeline    sql.NullString
	Message     string
	Source      sql.NullString
	Status      string
	Metadata    json.RawMessage
}

type UpdateProjectInquiryData struct {
	ID          int32
	UpdatedAt   time.Time
	Name        string
	Email       string
	Company     sql.NullString
	Role        sql.NullString
	ProjectType sql.NullString
	Timeline    sql.NullString
	Message     string
	Source      sql.NullString
	Status      string
	Metadata    json.RawMessage
}

type PaginatedProjectInquiries struct {
	ProjectInquiries []ProjectInquiryEntity
	TotalCount       int64
	Page             int64
	PageSize         int64
	TotalPages       int64
}

func (p *ProjectInquiries) Find(ctx context.Context, id int32) (ProjectInquiryEntity, error) {
	return p.queries.GetProjectInquiry[ProjectInquiryEntity](ctx, id)
}

func (p *ProjectInquiries) Create(ctx context.Context, data CreateProjectInquiryData) (ProjectInquiryEntity, error) {
	e := ProjectInquiryEntity{Name: data.Name, Email: data.Email, Company: data.Company, Role: data.Role, ProjectType: data.ProjectType, Timeline: data.Timeline, Message: data.Message, Source: data.Source, Status: data.Status, Metadata: data.Metadata}
	if err := e.Validate(); err != nil {
		return ProjectInquiryEntity{}, errors.Join(ErrDomainValidation, err)
	}
	now := time.Now()
	return p.queries.CreateProjectInquiry[ProjectInquiryEntity](ctx, queries.CreateProjectInquiryParams{
		CreatedAt: now, Name: e.Name, Email: e.Email, Company: contentText(e.Company), Role: contentText(e.Role),
		ProjectType: contentText(e.ProjectType), Timeline: contentText(e.Timeline), Message: e.Message,
		Source: contentText(e.Source), Status: e.Status, Metadata: e.Metadata,
	})
}

// Upsert's data has no ID; like the legacy implementation it inserts a new row.
func (p *ProjectInquiries) Upsert(ctx context.Context, data CreateProjectInquiryData) (ProjectInquiryEntity, error) {
	return p.Create(ctx, data)
}

func (p *ProjectInquiries) Update(ctx context.Context, data UpdateProjectInquiryData) (ProjectInquiryEntity, error) {
	e := ProjectInquiryEntity{ID: data.ID, Name: data.Name, Email: data.Email, Company: data.Company, Role: data.Role, ProjectType: data.ProjectType, Timeline: data.Timeline, Message: data.Message, Source: data.Source, Status: data.Status, Metadata: data.Metadata}
	if err := e.Validate(); err != nil {
		return ProjectInquiryEntity{}, errors.Join(ErrDomainValidation, err)
	}
	return p.queries.UpdateProjectInquiry[ProjectInquiryEntity](ctx, queries.UpdateProjectInquiryParams{
		ID: e.ID, UpdatedAt: time.Now(), Name: e.Name, Email: e.Email, Company: contentText(e.Company), Role: contentText(e.Role),
		ProjectType: contentText(e.ProjectType), Timeline: contentText(e.Timeline), Message: e.Message,
		Source: contentText(e.Source), Status: e.Status, Metadata: e.Metadata,
	})
}

func (p *ProjectInquiries) Destroy(ctx context.Context, id int32) error {
	return p.queries.DeleteProjectInquiry(ctx, id)
}

func (p *ProjectInquiries) All(ctx context.Context) ([]ProjectInquiryEntity, error) {
	items, err := p.queries.ListProjectInquiries[ProjectInquiryEntity](ctx)
	if err == nil && items == nil {
		items = []ProjectInquiryEntity{}
	}
	return items, err
}

func (p *ProjectInquiries) Paginate(ctx context.Context, page, pageSize int64) (PaginatedProjectInquiries, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	total, err := p.queries.CountProjectInquiries(ctx)
	if err != nil {
		return PaginatedProjectInquiries{}, err
	}
	entries, err := p.queries.PaginateProjectInquiries[ProjectInquiryEntity](ctx, queries.PaginateProjectInquiriesParams{Column1: pageSize, Column2: (page - 1) * pageSize})
	if err != nil {
		return PaginatedProjectInquiries{}, err
	}
	if entries == nil {
		entries = []ProjectInquiryEntity{}
	}
	return PaginatedProjectInquiries{entries, total, page, pageSize, (total + pageSize - 1) / pageSize}, nil
}

func projectInquiryEmailDomain(email string) (string, error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}
	at := strings.LastIndexByte(address.Address, '@')
	if at < 0 {
		return "", errors.New("email address has no domain")
	}
	domain := strings.ToLower(strings.TrimSuffix(address.Address[at+1:], "."))
	if domain == "" {
		return "", errors.New("email address has no domain")
	}
	return domain, nil
}

func (p *ProjectInquiries) IsDomainBanned(ctx context.Context, email string) (bool, error) {
	domain, err := projectInquiryEmailDomain(email)
	if err != nil {
		return false, nil
	}
	return p.queries.IsProjectInquiryDomainBanned(ctx, domain)
}

func (p *ProjectInquiries) BanDomain(ctx context.Context, id int32) (string, error) {
	inquiry, err := p.Find(ctx, id)
	if err != nil {
		return "", err
	}
	domain, err := projectInquiryEmailDomain(inquiry.Email)
	if err != nil {
		return "", err
	}
	err = p.queries.BanProjectInquiryDomain(ctx, queries.BanProjectInquiryDomainParams{Domain: domain, CreatedAt: time.Now()})
	return domain, err
}

func (p *ProjectInquiries) UnbanDomain(ctx context.Context, id int32) (string, error) {
	inquiry, err := p.Find(ctx, id)
	if err != nil {
		return "", err
	}
	domain, err := projectInquiryEmailDomain(inquiry.Email)
	if err != nil {
		return "", err
	}
	err = p.queries.UnbanProjectInquiryDomain(ctx, domain)
	return domain, err
}
