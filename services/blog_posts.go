package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mbvlabs/domains"
	"mbvlabs/queue/jobs"

	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
	"github.com/riverqueue/river"
)

// BlogPosts coordinates post creation with its publication job.
type BlogPosts struct {
	transactor storage.Transactor
	posts      *domains.BlogPosts
	queue      storage.InsertQueue
}

func NewBlogPosts(transactor storage.Transactor, posts *domains.BlogPosts, queue storage.InsertQueue) BlogPosts {
	return BlogPosts{transactor: transactor, posts: posts, queue: queue}
}

func (b BlogPosts) Create(ctx context.Context, data domains.CreateBlogPostData, scheduledAt *time.Time) (domains.BlogPostEntity, error) {
	if scheduledAt == nil {
		return b.posts.Create(ctx, data)
	}
	at := scheduledAt.UTC()
	if !at.After(time.Now()) {
		return domains.BlogPostEntity{}, validation.Field("scheduledAt", "future", "must be in the future")
	}
	data.Status = "draft"
	var result domains.BlogPostEntity
	err := storage.AtomicOps(ctx, b.transactor, func(ctx context.Context, tx storage.Tx) error {
		posts := b.posts.InTx(tx)
		e, err := posts.Create(ctx, data)
		if err != nil {
			return err
		}
		ready := e
		ready.Status = "published"
		ready.PublishedAt = domains.OptionalTime{Time: at, Valid: true}
		if err := ready.Validate(); err != nil {
			return err
		}
		job, err := b.queue.InTx(tx).Insert(ctx, jobs.PublishBlogPostArgs{BlogPostID: e.ID, PublishAt: at}, &river.InsertOpts{ScheduledAt: at})
		if err != nil {
			return fmt.Errorf("schedule blog post: %w", err)
		}
		if job == nil || job.Job == nil {
			return errors.New("schedule blog post: missing job")
		}
		if err := e.SetPublicationSchedule(domains.BlogPostPublicationSchedule{JobID: job.Job.ID, ScheduledAt: at}); err != nil {
			return err
		}
		result, err = posts.UpdatePublicationSchedule(ctx, e)
		return err
	})
	if err != nil {
		return domains.BlogPostEntity{}, err
	}
	return result, nil
}
