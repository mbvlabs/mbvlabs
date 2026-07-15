package queue

import (
	"context"

	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/queue/jobs"

	"github.com/riverqueue/river"
)

type PublishBlogPostWorker struct {
	river.WorkerDefaults[jobs.PublishBlogPostArgs]
	db                     storage.Pool
	invalidateSitemapCache func()
}

func NewPublishBlogPostWorker(db storage.Pool, invalidateSitemapCache func()) *PublishBlogPostWorker {
	if invalidateSitemapCache == nil {
		invalidateSitemapCache = func() {}
	}
	return &PublishBlogPostWorker{db: db, invalidateSitemapCache: invalidateSitemapCache}
}

func (w *PublishBlogPostWorker) Register(workers *river.Workers) error {
	return river.AddWorkerSafely(workers, w)
}

func (w *PublishBlogPostWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.PublishBlogPostArgs],
) error {
	published, err := models.BlogPost.PublishIfDraft(
		ctx,
		w.db.Executor(),
		job.Args.BlogPostID,
		job.Args.PublishAt,
	)
	if err != nil {
		return err
	}
	if published {
		w.invalidateSitemapCache()
	}
	return nil
}
