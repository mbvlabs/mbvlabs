package factories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"mbvlabs/internal/storage"
	"mbvlabs/models"

	"github.com/go-faker/faker/v4"
)

// WorkFactory wraps models.WorkEntity for testing
type WorkFactory struct {
	models.WorkEntity
}

type WorkOption func(*WorkFactory)

// BuildWork creates an in-memory Work with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateWork.
func BuildWork(opts ...WorkOption) models.WorkEntity {
	f := &WorkFactory{
		WorkEntity: models.WorkEntity{
			Title:          faker.Word(),
			Slug:           faker.Word(),
			ClientName:     sql.NullString{String: faker.Word(), Valid: true},
			ClientIndustry: sql.NullString{String: faker.Word(), Valid: true},
			ClientURL:      sql.NullString{String: faker.Word(), Valid: true},
			ClientLogoURL:  sql.NullString{String: faker.Word(), Valid: true},
			Summary:        faker.Word(),
			CoverImageURL:  sql.NullString{String: faker.Word(), Valid: true},
			Specialisms:    []string{},
			Platforms:      []string{},
			Technologies:   []string{},
			Challenge:      faker.Word(),
			Approach:       faker.Word(),
			Deliverables:   faker.Word(),
			Outcome:        faker.Word(),
			Content:        faker.Word(),
			StartedAt:      sql.NullTime{Time: time.Now(), Valid: true},
			CompletedAt:    sql.NullTime{Time: time.Now(), Valid: true},
			Status:         models.Draft,
			PublishedAt:    sql.NullTime{Time: time.Now(), Valid: true},
			IsFeatured:     randomBool(),
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.WorkEntity
}

// CreateWork creates and persists a Work to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateWork(
	ctx context.Context,
	exec storage.Executor,
	opts ...WorkOption,
) (models.WorkEntity, error) {
	built := BuildWork(opts...)

	entity := models.WorkEntity{
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Title:          built.Title,
		Slug:           built.Slug,
		ClientName:     built.ClientName,
		ClientIndustry: built.ClientIndustry,
		ClientURL:      built.ClientURL,
		ClientLogoURL:  built.ClientLogoURL,
		Summary:        built.Summary,
		CoverImageURL:  built.CoverImageURL,
		Specialisms:    built.Specialisms,
		Platforms:      built.Platforms,
		Technologies:   built.Technologies,
		Challenge:      built.Challenge,
		Approach:       built.Approach,
		Deliverables:   built.Deliverables,
		Outcome:        built.Outcome,
		Content:        built.Content,
		StartedAt:      built.StartedAt,
		CompletedAt:    built.CompletedAt,
		Status:         built.Status,
		PublishedAt:    built.PublishedAt,
		IsFeatured:     built.IsFeatured,
	}

	if err := exec.NewInsert().Model(&entity).Returning("*").Scan(ctx); err != nil {
		return models.WorkEntity{}, err
	}

	return entity, nil
}

// CreateWorks creates multiple Work records at once
func CreateWorks(
	ctx context.Context,
	exec storage.Executor,
	count int,
	opts ...WorkOption,
) ([]models.WorkEntity, error) {
	works := make([]models.WorkEntity, 0, count)

	for i := range count {
		entity, err := CreateWork(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create work %d: %w", i+1, err)
		}
		works = append(works, entity)
	}

	return works, nil
}

// Option functions

// WithWorksTitle sets the Title field
func WithWorksTitle(value string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Title = value
	}
}

// WithWorksSlug sets the Slug field
func WithWorksSlug(value string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Slug = value
	}
}

// WithWorksClientName sets the ClientName field
func WithWorksClientName(value sql.NullString) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.ClientName = value
	}
}

// WithWorksClientIndustry sets the ClientIndustry field
func WithWorksClientIndustry(value sql.NullString) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.ClientIndustry = value
	}
}

// WithWorksClientUrl sets the ClientUrl field
func WithWorksClientUrl(value sql.NullString) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.ClientURL = value
	}
}

// WithWorksClientLogoUrl sets the ClientLogoUrl field
func WithWorksClientLogoUrl(value sql.NullString) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.ClientLogoURL = value
	}
}

// WithWorksSummary sets the Summary field
func WithWorksSummary(value string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Summary = value
	}
}

// WithWorksCoverImageUrl sets the CoverImageUrl field
func WithWorksCoverImageUrl(value sql.NullString) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.CoverImageURL = value
	}
}

// WithWorksSpecialisms sets the Specialisms field
func WithWorksSpecialisms(value []string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Specialisms = value
	}
}

// WithWorksPlatforms sets the Platforms field
func WithWorksPlatforms(value []string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Platforms = value
	}
}

// WithWorksTechnologies sets the Technologies field
func WithWorksTechnologies(value []string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Technologies = value
	}
}

// WithWorksChallenge sets the Challenge field
func WithWorksChallenge(value string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Challenge = value
	}
}

// WithWorksApproach sets the Approach field
func WithWorksApproach(value string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Approach = value
	}
}

// WithWorksDeliverables sets the Deliverables field
func WithWorksDeliverables(value string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Deliverables = value
	}
}

// WithWorksOutcome sets the Outcome field
func WithWorksOutcome(value string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Outcome = value
	}
}

// WithWorksContent sets the Content field
func WithWorksContent(value string) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Content = value
	}
}

// WithWorksStartedAt sets the StartedAt field
func WithWorksStartedAt(value sql.NullTime) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.StartedAt = value
	}
}

// WithWorksCompletedAt sets the CompletedAt field
func WithWorksCompletedAt(value sql.NullTime) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.CompletedAt = value
	}
}

// WithWorksStatus sets the Status field
func WithWorksStatus(value models.StatusEnum) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.Status = value
	}
}

// WithWorksPublishedAt sets the PublishedAt field
func WithWorksPublishedAt(value sql.NullTime) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.PublishedAt = value
	}
}

// WithWorksIsFeatured sets the IsFeatured field
func WithWorksIsFeatured(value bool) WorkOption {
	return func(f *WorkFactory) {
		f.WorkEntity.IsFeatured = value
	}
}
