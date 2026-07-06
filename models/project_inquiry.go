package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"mbvlabs/internal/storage"
	"mbvlabs/internal/validation"
	"net/mail"
	"time"

	"github.com/uptrace/bun"
)

type ProjectInquiryEntity struct {
	bun.BaseModel `bun:"table:project_inquiries,alias:project_inquiries"`
	ID            int32           `bun:"id,pk,autoincrement"`
	CreatedAt     time.Time       `bun:"created_at"`
	UpdatedAt     time.Time       `bun:"updated_at"`
	Name          string          `bun:"name"`
	Email         string          `bun:"email"`
	Company       sql.NullString  `bun:"company"`
	Role          sql.NullString  `bun:"role"`
	ProjectType   sql.NullString  `bun:"project_type"`
	Timeline      sql.NullString  `bun:"timeline"`
	Message       string          `bun:"message"`
	Source        sql.NullString  `bun:"source"`
	Status        string          `bun:"status"`
	Metadata      json.RawMessage `bun:"metadata,type:jsonb"`
}

func (pie ProjectInquiryEntity) Validate() error {
	builder := validation.NewBuilder()

	builder.Required("name", pie.Name)
	builder.LenBetween("name", pie.Name, 2, 255)
	builder.Required("email", pie.Email)
	builder.MaxLen("email", pie.Email, 255)
	builder.MaxLen("company", pie.Company, 255)
	builder.MaxLen("role", pie.Role, 255)
	builder.OneOf("projectType", pie.ProjectType, "web-app", "website", "automation", "integration", "consulting", "open-source", "modernization")
	builder.MaxLen("timeline", pie.Timeline, 100)
	builder.Required("message", pie.Message)
	builder.LenBetween("message", pie.Message, 10, 5000)
	builder.MaxLen("source", pie.Source, 255)
	builder.Required("status", pie.Status)
	builder.OneOf("status", pie.Status, "new", "contacted", "qualified", "proposal", "won", "lost")
	builder.Required("metadata", pie.Metadata)

	if pie.Email != "" {
		if _, err := mail.ParseAddress(pie.Email); err != nil {
			builder.AddField("email", "email", "must be a valid email address")
		}
	}

	return builder.Err()
}

func (pi projectInquiry) Find(
	ctx context.Context,
	db storage.Executor,
	id int32,
) (ProjectInquiryEntity, error) {
	var entity ProjectInquiryEntity
	if err := db.NewSelect().
		Model(&entity).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProjectInquiryEntity{}, ErrNotFound
		}
		return ProjectInquiryEntity{}, err
	}

	return entity, nil
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

func (pi projectInquiry) Create(
	ctx context.Context,
	db storage.Executor,
	data CreateProjectInquiryData,
) (ProjectInquiryEntity, error) {
	entity := ProjectInquiryEntity{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        data.Name,
		Email:       data.Email,
		Company:     data.Company,
		Role:        data.Role,
		ProjectType: data.ProjectType,
		Timeline:    data.Timeline,
		Message:     data.Message,
		Source:      data.Source,
		Status:      data.Status,
		Metadata:    data.Metadata,
	}

	if err := validation.Validate(entity); err != nil {
		return ProjectInquiryEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if _, err := db.NewInsert().Model(&entity).Exec(ctx); err != nil {
		return ProjectInquiryEntity{}, err
	}

	return entity, nil
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

func (pi projectInquiry) Update(
	ctx context.Context,
	db storage.Executor,
	data UpdateProjectInquiryData,
) (ProjectInquiryEntity, error) {
	entity := ProjectInquiryEntity{
		ID:          data.ID,
		UpdatedAt:   time.Now(),
		Name:        data.Name,
		Email:       data.Email,
		Company:     data.Company,
		Role:        data.Role,
		ProjectType: data.ProjectType,
		Timeline:    data.Timeline,
		Message:     data.Message,
		Source:      data.Source,
		Status:      data.Status,
		Metadata:    data.Metadata,
	}

	if err := validation.Validate(entity); err != nil {
		return ProjectInquiryEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewUpdate().
		Model(&entity).
		Column("updated_at").
		Column("name").
		Column("email").
		Column("company").
		Column("role").
		Column("project_type").
		Column("timeline").
		Column("message").
		Column("source").
		Column("status").
		Column("metadata").
		WherePK().
		Returning("*").
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProjectInquiryEntity{}, ErrNotFound
		}
		return ProjectInquiryEntity{}, err
	}

	return entity, nil
}

func (pi projectInquiry) Destroy(ctx context.Context, db storage.Executor, id int32) error {
	_, err := db.NewDelete().
		Model((*ProjectInquiryEntity)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	return err
}

func (pi projectInquiry) All(
	ctx context.Context,
	db storage.Executor,
) ([]ProjectInquiryEntity, error) {
	var entities []ProjectInquiryEntity
	if err := db.NewSelect().
		Model(&entities).
		Scan(ctx); err != nil {
		return nil, err
	}

	return entities, nil
}

type PaginatedProjectInquiries struct {
	ProjectInquiries []ProjectInquiryEntity
	TotalCount       int64
	Page             int64
	PageSize         int64
	TotalPages       int64
}

func (pi projectInquiry) Paginate(
	ctx context.Context,
	db storage.Executor,
	page, pageSize int64,
) (PaginatedProjectInquiries, error) {
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
		Model(&ProjectInquiryEntity{}).Count(ctx)
	if err != nil {
		return PaginatedProjectInquiries{}, err
	}

	entities := make([]ProjectInquiryEntity, 0, int(pageSize))
	if err := db.NewSelect().
		Model(&entities).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(ctx); err != nil {
		return PaginatedProjectInquiries{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedProjectInquiries{
		ProjectInquiries: entities,
		TotalCount:       int64(totalCount),
		Page:             page,
		PageSize:         pageSize,
		TotalPages:       totalPages,
	}, nil
}

func (pi projectInquiry) Upsert(
	ctx context.Context,
	db storage.Executor,
	data CreateProjectInquiryData,
) (ProjectInquiryEntity, error) {
	entity := ProjectInquiryEntity{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        data.Name,
		Email:       data.Email,
		Company:     data.Company,
		Role:        data.Role,
		ProjectType: data.ProjectType,
		Timeline:    data.Timeline,
		Message:     data.Message,
		Source:      data.Source,
		Status:      data.Status,
		Metadata:    data.Metadata,
	}

	if err := validation.Validate(entity); err != nil {
		return ProjectInquiryEntity{}, errors.Join(ErrDomainValidation, err)
	}

	if err := db.NewInsert().
		Model(&entity).
		On("CONFLICT (id) DO UPDATE").
		Set("name = excluded.name").
		Set("email = excluded.email").
		Set("company = excluded.company").
		Set("role = excluded.role").
		Set("project_type = excluded.project_type").
		Set("timeline = excluded.timeline").
		Set("message = excluded.message").
		Set("source = excluded.source").
		Set("status = excluded.status").
		Set("metadata = excluded.metadata").
		Returning("*").
		Scan(ctx); err != nil {
		return ProjectInquiryEntity{}, err
	}

	return entity, nil
}
