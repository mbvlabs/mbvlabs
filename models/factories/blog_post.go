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

// BlogPostFactory wraps models.BlogPostEntity for testing
type BlogPostFactory struct {
	models.BlogPostEntity
}

type BlogPostOption func(*BlogPostFactory)

// BuildBlogPost creates an in-memory BlogPost with default test values.
// Auto-managed fields (ID, timestamps) are left at zero and set by CreateBlogPost.
func BuildBlogPost(opts ...BlogPostOption) models.BlogPostEntity {
	f := &BlogPostFactory{
		BlogPostEntity: models.BlogPostEntity{
			Title:         faker.Word(),
			Slug:          faker.Word(),
			Excerpt:       sql.NullString{String: faker.Word(), Valid: true},
			Body:          faker.Word(),
			Status:        faker.Word(),
			CoverImageUrl: sql.NullString{String: faker.Word(), Valid: true},
			Tags:          json.RawMessage{},
			PublishedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		},
	}

	for _, opt := range opts {
		opt(f)
	}

	return f.BlogPostEntity
}

// CreateBlogPost creates and persists a BlogPost to the database.
// It returns the entity populated with all DB-assigned values via RETURNING *.
func CreateBlogPost(
	ctx context.Context,
	exec storage.Executor,
	opts ...BlogPostOption,
) (models.BlogPostEntity, error) {
	built := BuildBlogPost(opts...)

	entity := models.BlogPostEntity{
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Title:         built.Title,
		Slug:          built.Slug,
		Excerpt:       built.Excerpt,
		Body:          built.Body,
		Status:        built.Status,
		CoverImageUrl: built.CoverImageUrl,
		Tags:          built.Tags,
		PublishedAt:   built.PublishedAt,
	}

	if err := exec.NewInsert().Model(&entity).Returning("*").Scan(ctx); err != nil {
		return models.BlogPostEntity{}, err
	}

	return entity, nil
}

// CreateBlogPosts creates multiple BlogPost records at once
func CreateBlogPosts(
	ctx context.Context,
	exec storage.Executor,
	count int,
	opts ...BlogPostOption,
) ([]models.BlogPostEntity, error) {
	blogposts := make([]models.BlogPostEntity, 0, count)

	for i := range count {
		entity, err := CreateBlogPost(ctx, exec, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create blogpost %d: %w", i+1, err)
		}
		blogposts = append(blogposts, entity)
	}

	return blogposts, nil
}

// Option functions

// WithBlogPostsTitle sets the Title field
func WithBlogPostsTitle(value string) BlogPostOption {
	return func(f *BlogPostFactory) {
		f.BlogPostEntity.Title = value
	}
}

// WithBlogPostsSlug sets the Slug field
func WithBlogPostsSlug(value string) BlogPostOption {
	return func(f *BlogPostFactory) {
		f.BlogPostEntity.Slug = value
	}
}

// WithBlogPostsExcerpt sets the Excerpt field
func WithBlogPostsExcerpt(value sql.NullString) BlogPostOption {
	return func(f *BlogPostFactory) {
		f.BlogPostEntity.Excerpt = value
	}
}

// WithBlogPostsBody sets the Body field
func WithBlogPostsBody(value string) BlogPostOption {
	return func(f *BlogPostFactory) {
		f.BlogPostEntity.Body = value
	}
}

// WithBlogPostsStatus sets the Status field
func WithBlogPostsStatus(value string) BlogPostOption {
	return func(f *BlogPostFactory) {
		f.BlogPostEntity.Status = value
	}
}

// WithBlogPostsCoverImageUrl sets the CoverImageUrl field
func WithBlogPostsCoverImageUrl(value sql.NullString) BlogPostOption {
	return func(f *BlogPostFactory) {
		f.BlogPostEntity.CoverImageUrl = value
	}
}

// WithBlogPostsTags sets the Tags field
func WithBlogPostsTags(value json.RawMessage) BlogPostOption {
	return func(f *BlogPostFactory) {
		f.BlogPostEntity.Tags = value
	}
}

// WithBlogPostsPublishedAt sets the PublishedAt field
func WithBlogPostsPublishedAt(value sql.NullTime) BlogPostOption {
	return func(f *BlogPostFactory) {
		f.BlogPostEntity.PublishedAt = value
	}
}
