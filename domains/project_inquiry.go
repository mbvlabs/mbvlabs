package domains

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type ProjectInquiries struct{ db storage.Connection }

func NewProjectInquiries(db storage.Connection) *ProjectInquiries {
	return &ProjectInquiries{db: db}
}

type ProjectInquiryEntity struct {
	ID          int32
	CreatedAt   time.Time
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

const inquiryColumns = "id, created_at, updated_at, name, email, company, role, project_type, timeline, message, source, status, metadata"

func scanInquiry(row interface{ Scan(...any) error }) (ProjectInquiryEntity, error) {
	var e ProjectInquiryEntity
	err := row.Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.Name, &e.Email, &e.Company, &e.Role, &e.ProjectType, &e.Timeline, &e.Message, &e.Source, &e.Status, &e.Metadata)
	return storage.Row(e, err)
}

func (p *ProjectInquiries) Find(ctx context.Context, id int32) (ProjectInquiryEntity, error) {
	return scanInquiry(p.db.QueryRow(ctx, "SELECT "+inquiryColumns+" FROM project_inquiries WHERE id=$1", id))
}

func (p *ProjectInquiries) Create(ctx context.Context, data CreateProjectInquiryData) (ProjectInquiryEntity, error) {
	e := ProjectInquiryEntity{Name: data.Name, Email: data.Email, Company: data.Company, Role: data.Role, ProjectType: data.ProjectType, Timeline: data.Timeline, Message: data.Message, Source: data.Source, Status: data.Status, Metadata: data.Metadata}
	if err := e.Validate(); err != nil {
		return ProjectInquiryEntity{}, errors.Join(ErrDomainValidation, err)
	}
	now := time.Now()
	return scanInquiry(p.db.QueryRow(ctx, "INSERT INTO project_inquiries (created_at, updated_at, name, email, company, role, project_type, timeline, message, source, status, metadata) VALUES ($1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "+inquiryColumns,
		now, e.Name, e.Email, e.Company, e.Role, e.ProjectType, e.Timeline, e.Message, e.Source, e.Status, e.Metadata))
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
	return scanInquiry(p.db.QueryRow(ctx, "UPDATE project_inquiries SET updated_at=$2, name=$3, email=$4, company=$5, role=$6, project_type=$7, timeline=$8, message=$9, source=$10, status=$11, metadata=$12 WHERE id=$1 RETURNING "+inquiryColumns,
		e.ID, time.Now(), e.Name, e.Email, e.Company, e.Role, e.ProjectType, e.Timeline, e.Message, e.Source, e.Status, e.Metadata))
}

func (p *ProjectInquiries) Destroy(ctx context.Context, id int32) error {
	_, err := p.db.Exec(ctx, "DELETE FROM project_inquiries WHERE id=$1", id)
	return err
}

func (p *ProjectInquiries) list(ctx context.Context, query string, args ...any) ([]ProjectInquiryEntity, error) {
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]ProjectInquiryEntity, 0)
	for rows.Next() {
		e, err := scanInquiry(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (p *ProjectInquiries) All(ctx context.Context) ([]ProjectInquiryEntity, error) {
	return p.list(ctx, "SELECT "+inquiryColumns+" FROM project_inquiries")
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
	var total int64
	if err := p.db.QueryRow(ctx, "SELECT count(*) FROM project_inquiries").Scan(&total); err != nil {
		return PaginatedProjectInquiries{}, err
	}
	entries, err := p.list(ctx, "SELECT "+inquiryColumns+" FROM project_inquiries ORDER BY id LIMIT $1 OFFSET $2", pageSize, (page-1)*pageSize)
	if err != nil {
		return PaginatedProjectInquiries{}, err
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
	var banned bool
	err = p.db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM project_inquiry_domain_bans WHERE domain=$1)", domain).Scan(&banned)
	return banned, err
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
	_, err = p.db.Exec(ctx, "INSERT INTO project_inquiry_domain_bans (domain, created_at) VALUES ($1,$2) ON CONFLICT (domain) DO NOTHING", domain, time.Now())
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
	_, err = p.db.Exec(ctx, "DELETE FROM project_inquiry_domain_bans WHERE domain=$1", domain)
	return domain, err
}
