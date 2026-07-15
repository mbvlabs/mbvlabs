package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"mbvlabs/internal/storage"
	"mbvlabs/internal/validation"
	"mbvlabs/models"
	"mbvlabs/queue"
	"mbvlabs/queue/jobs"

	"github.com/riverqueue/river"
)

type BlogPosts struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
}

func NewBlogPosts(db storage.Pool, insertOnly queue.InsertOnly) BlogPosts {
	return BlogPosts{db: db, insertOnly: insertOnly}
}

func (bp BlogPosts) Create(
	ctx context.Context,
	data models.CreateBlogPostData,
	scheduledAt *time.Time,
) (models.BlogPostEntity, error) {
	if scheduledAt == nil {
		return models.BlogPost.Create(ctx, bp.db.Executor(), data)
	}

	publishAt := scheduledAt.UTC()
	if !publishAt.After(time.Now()) {
		builder := validation.NewBuilder()
		builder.AddField("scheduledAt", "future", "must be in the future")
		return models.BlogPostEntity{}, errors.Join(models.ErrDomainValidation, builder.Err())
	}
	data.Status = models.Draft.String()

	tx, err := bp.db.BeginTx(ctx, nil)
	if err != nil {
		return models.BlogPostEntity{}, fmt.Errorf("begin scheduled blog post transaction: %v", err)
	}
	defer tx.Rollback()

	article, err := models.BlogPost.Create(ctx, tx, data)
	if err != nil {
		return models.BlogPostEntity{}, err
	}
	publicationReady := article
	publicationReady.Status = models.Published.String()
	publicationReady.PublishedAt = sql.NullTime{Time: publishAt, Valid: true}
	if err := validation.Validate(publicationReady); err != nil {
		return models.BlogPostEntity{}, errors.Join(models.ErrDomainValidation, err)
	}

	result, err := bp.insertOnly.InsertTx(ctx, tx.Tx, jobs.PublishBlogPostArgs{
		BlogPostID: article.ID,
		PublishAt:  publishAt,
	}, &river.InsertOpts{ScheduledAt: publishAt})
	if err != nil {
		return models.BlogPostEntity{}, fmt.Errorf("schedule blog post publication: %v", err)
	}
	if result == nil || result.Job == nil {
		return models.BlogPostEntity{}, errors.New("schedule blog post publication: missing inserted job")
	}

	if err := article.SetPublicationSchedule(models.BlogPostPublicationSchedule{
		JobID:       result.Job.ID,
		ScheduledAt: publishAt,
	}); err != nil {
		return models.BlogPostEntity{}, err
	}
	article, err = models.BlogPost.UpdatePublicationSchedule(ctx, tx, article)
	if err != nil {
		return models.BlogPostEntity{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.BlogPostEntity{}, fmt.Errorf("commit scheduled blog post transaction: %v", err)
	}
	return article, nil
}
