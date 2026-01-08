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

// WorkItemFactory wraps models.WorkItem for testing
type WorkItemFactory struct {
	models.WorkItem // Embedded
}

type WorkItemOption func(*WorkItemFactory)

// BuildWorkItem creates an in-memory WorkItem with default test values
func BuildWorkItem(opts ...WorkItemOption) models.WorkItem {
	f := &WorkItemFactory{
		WorkItem: models.WorkItem{
			ID:               uuid.New(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Title:            faker.Word(),
			Slug:             faker.Word(),
			ShortDescription: faker.Word(),
			Content:          faker.Word(),
			Client:           faker.Word(),
			Industry:         faker.Word(),
			ProjectDate:      time.Time{}, // Optional timestamp - zero by default
			ProjectDuration:  faker.Word(),
			HeroImageUrl:     faker.Word(),
			HeroImageAlt:     faker.Word(),
			ExternalUrl:      faker.Word(),
			IsPublished:      faker.Bool(),
			IsFeatured:       faker.Bool(),
			DisplayOrder:     randomInt(1, 1000, 100),
			MetaTitle:        faker.Word(),
			MetaDescription:  faker.Word(),
			MetaKeywords:     []string{},
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.WorkItem
}

// CreateWorkItem creates and persists a WorkItem to the database
func CreateWorkItem(ctx context.Context, exec storage.Executor, opts ...WorkItemOption) (models.WorkItem, error) {
	// Build with defaults and required FKs
	built := BuildWorkItem(opts...)

	// Prepare creation data
	data := models.CreateWorkItemData{
		Title:            built.Title,
		Slug:             built.Slug,
		ShortDescription: built.ShortDescription,
		Content:          built.Content,
		Client:           built.Client,
		Industry:         built.Industry,
		ProjectDate:      built.ProjectDate,
		ProjectDuration:  built.ProjectDuration,
		HeroImageUrl:     built.HeroImageUrl,
		HeroImageAlt:     built.HeroImageAlt,
		ExternalUrl:      built.ExternalUrl,
		IsPublished:      built.IsPublished,
		IsFeatured:       built.IsFeatured,
		DisplayOrder:     built.DisplayOrder,
		MetaTitle:        built.MetaTitle,
		MetaDescription:  built.MetaDescription,
		MetaKeywords:     built.MetaKeywords,
	}

	// Use model's Create function
	workitem, err := models.CreateWorkItem(ctx, exec, data)
	if err != nil {
		return models.WorkItem{}, err
	}

	return workitem, nil
}

// CreateWorkItems creates multiple WorkItem records at once
func CreateWorkItems(ctx context.Context, exec storage.Executor, count int, opts ...WorkItemOption) ([]models.WorkItem, error) {
	workitems := make([]models.WorkItem, 0, count)

	for i := 0; i < count; i++ {
		workitem, err := CreateWorkItem(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create workitem %d: %w", i+1, err)
		}
		workitems = append(workitems, workitem)
	}

	return workitems, nil
}

// Option functions

// WithWork_itemsTitle sets the Title field
func WithWork_itemsTitle(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.Title = value
	}
}

// WithWork_itemsSlug sets the Slug field
func WithWork_itemsSlug(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.Slug = value
	}
}

// WithWork_itemsShortDescription sets the ShortDescription field
func WithWork_itemsShortDescription(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.ShortDescription = value
	}
}

// WithWork_itemsContent sets the Content field
func WithWork_itemsContent(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.Content = value
	}
}

// WithWork_itemsClient sets the Client field
func WithWork_itemsClient(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.Client = value
	}
}

// WithWork_itemsIndustry sets the Industry field
func WithWork_itemsIndustry(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.Industry = value
	}
}

// WithWork_itemsProjectDate sets the ProjectDate field
func WithWork_itemsProjectDate(value time.Time) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.ProjectDate = value
	}
}

// WithWork_itemsProjectDuration sets the ProjectDuration field
func WithWork_itemsProjectDuration(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.ProjectDuration = value
	}
}

// WithWork_itemsHeroImageUrl sets the HeroImageUrl field
func WithWork_itemsHeroImageUrl(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.HeroImageUrl = value
	}
}

// WithWork_itemsHeroImageAlt sets the HeroImageAlt field
func WithWork_itemsHeroImageAlt(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.HeroImageAlt = value
	}
}

// WithWork_itemsExternalUrl sets the ExternalUrl field
func WithWork_itemsExternalUrl(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.ExternalUrl = value
	}
}

// WithWork_itemsIsPublished sets the IsPublished field
func WithWork_itemsIsPublished(value bool) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.IsPublished = value
	}
}

// WithWork_itemsIsFeatured sets the IsFeatured field
func WithWork_itemsIsFeatured(value bool) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.IsFeatured = value
	}
}

// WithWork_itemsDisplayOrder sets the DisplayOrder field
func WithWork_itemsDisplayOrder(value int32) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.DisplayOrder = value
	}
}

// WithWork_itemsMetaTitle sets the MetaTitle field
func WithWork_itemsMetaTitle(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.MetaTitle = value
	}
}

// WithWork_itemsMetaDescription sets the MetaDescription field
func WithWork_itemsMetaDescription(value string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.MetaDescription = value
	}
}

// WithWork_itemsMetaKeywords sets the MetaKeywords field
func WithWork_itemsMetaKeywords(value []string) WorkItemOption {
	return func(f *WorkItemFactory) {
		f.WorkItem.MetaKeywords = value
	}
}
