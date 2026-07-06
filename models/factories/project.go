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

// ProjectFactory wraps models.ProjectEntity for testing
type ProjectFactory struct {
	models.ProjectEntity
}

type ProjectOption func(*ProjectFactory)

// BuildProject creates an in-memory Project with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateProject.
func BuildProject(opts ...ProjectOption) models.ProjectEntity {
	f := &ProjectFactory{
		ProjectEntity: models.ProjectEntity{
			Name:          faker.Word(),
			Slug:          faker.Word(),
			Tagline:       faker.Word(),
			Description:   sql.NullString{String: faker.Word(), Valid: true},
			ProjectType:   faker.Word(),
			RepositoryURL: sql.NullString{String: faker.Word(), Valid: true},
			LiveURL:       sql.NullString{String: faker.Word(), Valid: true},
			ImageURL:      sql.NullString{String: faker.Word(), Valid: true},
			Technologies:  json.RawMessage{},
			StartedAt:     sql.NullTime{Time: time.Now(), Valid: true},
			LaunchedAt:    sql.NullTime{Time: time.Now(), Valid: true},
			PublishedAt:   sql.NullTime{Time: time.Now(), Valid: true},
			IsFeatured:    randomBool(),
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.ProjectEntity
}

// CreateProject creates and persists a Project to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateProject(
	ctx context.Context,
	exec storage.Executor,
	opts ...ProjectOption,
) (models.ProjectEntity, error) {
	built := BuildProject(opts...)

	entity := models.ProjectEntity{
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Name:          built.Name,
		Slug:          built.Slug,
		Tagline:       built.Tagline,
		Description:   built.Description,
		ProjectType:   built.ProjectType,
		RepositoryURL: built.RepositoryURL,
		LiveURL:       built.LiveURL,
		ImageURL:      built.ImageURL,
		Technologies:  built.Technologies,
		StartedAt:     built.StartedAt,
		LaunchedAt:    built.LaunchedAt,
		PublishedAt:   built.PublishedAt,
		IsFeatured:    built.IsFeatured,
	}

	if err := exec.NewInsert().Model(&entity).Returning("*").Scan(ctx); err != nil {
		return models.ProjectEntity{}, err
	}

	return entity, nil
}

// CreateProjects creates multiple Project records at once
func CreateProjects(
	ctx context.Context,
	exec storage.Executor,
	count int,
	opts ...ProjectOption,
) ([]models.ProjectEntity, error) {
	projects := make([]models.ProjectEntity, 0, count)

	for i := range count {
		entity, err := CreateProject(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create project %d: %w", i+1, err)
		}
		projects = append(projects, entity)
	}

	return projects, nil
}

// Option functions

// WithProjectsName sets the Name field
func WithProjectsName(value string) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.Name = value
	}
}

// WithProjectsSlug sets the Slug field
func WithProjectsSlug(value string) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.Slug = value
	}
}

// WithProjectsTagline sets the Tagline field
func WithProjectsTagline(value string) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.Tagline = value
	}
}

// WithProjectsDescription sets the Description field
func WithProjectsDescription(value sql.NullString) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.Description = value
	}
}

// WithProjectsProjectType sets the ProjectType field
func WithProjectsProjectType(value string) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.ProjectType = value
	}
}

// WithProjectsRepositoryUrl sets the RepositoryUrl field
func WithProjectsRepositoryUrl(value sql.NullString) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.RepositoryURL = value
	}
}

// WithProjectsLiveUrl sets the LiveUrl field
func WithProjectsLiveUrl(value sql.NullString) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.LiveURL = value
	}
}

// WithProjectsImageUrl sets the ImageUrl field
func WithProjectsImageUrl(value sql.NullString) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.ImageURL = value
	}
}

// WithProjectsTechnologies sets the Technologies field
func WithProjectsTechnologies(value json.RawMessage) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.Technologies = value
	}
}

// WithProjectsStartedAt sets the StartedAt field
func WithProjectsStartedAt(value sql.NullTime) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.StartedAt = value
	}
}

// WithProjectsLaunchedAt sets the LaunchedAt field
func WithProjectsLaunchedAt(value sql.NullTime) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.LaunchedAt = value
	}
}

// WithProjectsPublishedAt sets the PublishedAt field
func WithProjectsPublishedAt(value sql.NullTime) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.PublishedAt = value
	}
}

// WithProjectsIsFeatured sets the IsFeatured field
func WithProjectsIsFeatured(value bool) ProjectOption {
	return func(f *ProjectFactory) {
		f.ProjectEntity.IsFeatured = value
	}
}
