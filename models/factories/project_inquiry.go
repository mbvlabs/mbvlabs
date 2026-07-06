package factories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"mbvlabs/internal/storage"
	"mbvlabs/models"

	"github.com/go-faker/faker/v4"
)

// ProjectInquiryFactory wraps models.ProjectInquiryEntity for testing
type ProjectInquiryFactory struct {
	models.ProjectInquiryEntity
}

type ProjectInquiryOption func(*ProjectInquiryFactory)

// BuildProjectInquiry creates an in-memory ProjectInquiry with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateProjectInquiry.
func BuildProjectInquiry(opts ...ProjectInquiryOption) models.ProjectInquiryEntity {
	f := &ProjectInquiryFactory{
		ProjectInquiryEntity: models.ProjectInquiryEntity{
			Name:        faker.Word(),
			Email:       faker.Word(),
			Company:     sql.NullString{String: faker.Word(), Valid: true},
			Role:        sql.NullString{String: faker.Word(), Valid: true},
			ProjectType: sql.NullString{String: faker.Word(), Valid: true},
			Timeline:    sql.NullString{String: faker.Word(), Valid: true},
			Message:     faker.Word(),
			Source:      sql.NullString{String: faker.Word(), Valid: true},
			Status:      faker.Word(),
			Metadata:    json.RawMessage{},
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.ProjectInquiryEntity
}

// CreateProjectInquiry creates and persists a ProjectInquiry to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateProjectInquiry(
	ctx context.Context,
	exec storage.Executor,
	opts ...ProjectInquiryOption,
) (models.ProjectInquiryEntity, error) {
	built := BuildProjectInquiry(opts...)

	entity := models.ProjectInquiryEntity{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        built.Name,
		Email:       built.Email,
		Company:     built.Company,
		Role:        built.Role,
		ProjectType: built.ProjectType,
		Timeline:    built.Timeline,
		Message:     built.Message,
		Source:      built.Source,
		Status:      built.Status,
		Metadata:    built.Metadata,
	}

	if err := exec.NewInsert().Model(&entity).Returning("*").Scan(ctx); err != nil {
		return models.ProjectInquiryEntity{}, err
	}

	return entity, nil
}

// CreateProjectInquirys creates multiple ProjectInquiry records at once
func CreateProjectInquirys(
	ctx context.Context,
	exec storage.Executor,
	count int,
	opts ...ProjectInquiryOption,
) ([]models.ProjectInquiryEntity, error) {
	projectinquirys := make([]models.ProjectInquiryEntity, 0, count)

	for i := range count {
		entity, err := CreateProjectInquiry(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create projectinquiry %d: %w", i+1, err)
		}
		projectinquirys = append(projectinquirys, entity)
	}

	return projectinquirys, nil
}

// Option functions

// WithProjectInquiriesName sets the Name field
func WithProjectInquiriesName(value string) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.Name = value
	}
}

// WithProjectInquiriesEmail sets the Email field
func WithProjectInquiriesEmail(value string) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.Email = value
	}
}

// WithProjectInquiriesCompany sets the Company field
func WithProjectInquiriesCompany(value sql.NullString) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.Company = value
	}
}

// WithProjectInquiriesRole sets the Role field
func WithProjectInquiriesRole(value sql.NullString) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.Role = value
	}
}

// WithProjectInquiriesProjectType sets the ProjectType field
func WithProjectInquiriesProjectType(value sql.NullString) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.ProjectType = value
	}
}

// WithProjectInquiriesTimeline sets the Timeline field
func WithProjectInquiriesTimeline(value sql.NullString) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.Timeline = value
	}
}

// WithProjectInquiriesMessage sets the Message field
func WithProjectInquiriesMessage(value string) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.Message = value
	}
}

// WithProjectInquiriesSource sets the Source field
func WithProjectInquiriesSource(value sql.NullString) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.Source = value
	}
}

// WithProjectInquiriesStatus sets the Status field
func WithProjectInquiriesStatus(value string) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.Status = value
	}
}

// WithProjectInquiriesMetadata sets the Metadata field
func WithProjectInquiriesMetadata(value json.RawMessage) ProjectInquiryOption {
	return func(f *ProjectInquiryFactory) {
		f.ProjectInquiryEntity.Metadata = value
	}
}
