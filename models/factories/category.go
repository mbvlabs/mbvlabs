package factories

import (
	"context"
	"fmt"
	"time"

	"mbvlabs/internal/storage"
	"mbvlabs/models"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

// CategoryFactory wraps models.Category for testing
type CategoryFactory struct {
	models.Category // Embedded
}

type CategoryOption func(*CategoryFactory)

// BuildCategory creates an in-memory Category with default test values
func BuildCategory(opts ...CategoryOption) models.Category {
	f := &CategoryFactory{
		Category: models.Category{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        faker.Word(),
			Slug:        faker.Word(),
			Description: faker.Word(),
			Color:       faker.Word(),
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.Category
}

// CreateCategory creates and persists a Category to the database
func CreateCategory(ctx context.Context, exec storage.Executor, opts ...CategoryOption) (models.Category, error) {
	// Build with defaults and required FKs
	built := BuildCategory(opts...)

	// Prepare creation data
	data := models.CreateCategoryData{
		Name:        built.Name,
		Slug:        built.Slug,
		Description: built.Description,
		Color:       built.Color,
	}

	// Use model's Create function
	category, err := models.CreateCategory(ctx, exec, data)
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

// CreateCategorys creates multiple Category records at once
func CreateCategorys(ctx context.Context, exec storage.Executor, count int, opts ...CategoryOption) ([]models.Category, error) {
	categorys := make([]models.Category, 0, count)

	for i := 0; i < count; i++ {
		category, err := CreateCategory(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create category %d: %w", i+1, err)
		}
		categorys = append(categorys, category)
	}

	return categorys, nil
}

// Option functions

// WithCategoriesName sets the Name field
func WithCategoriesName(value string) CategoryOption {
	return func(f *CategoryFactory) {
		f.Category.Name = value
	}
}

// WithCategoriesSlug sets the Slug field
func WithCategoriesSlug(value string) CategoryOption {
	return func(f *CategoryFactory) {
		f.Category.Slug = value
	}
}

// WithCategoriesDescription sets the Description field
func WithCategoriesDescription(value string) CategoryOption {
	return func(f *CategoryFactory) {
		f.Category.Description = value
	}
}

// WithCategoriesColor sets the Color field
func WithCategoriesColor(value string) CategoryOption {
	return func(f *CategoryFactory) {
		f.Category.Color = value
	}
}
