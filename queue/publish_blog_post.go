package queue

import (
	"context"
	"time"

	"mbvlabs/queue/jobs"

	"github.com/riverqueue/river"
)

type BlogPublisher interface {
	PublishIfDraft(context.Context, int32, time.Time) (bool, error)
}

type publishBlogPostWorker struct {
	river.WorkerDefaults[jobs.PublishBlogPostArgs]
	posts BlogPublisher
}

func (w *publishBlogPostWorker) Work(ctx context.Context, job *river.Job[jobs.PublishBlogPostArgs]) error {
	_, err := w.posts.PublishIfDraft(ctx, job.Args.BlogPostID, job.Args.PublishAt)
	if err != nil {
		return failJob(ctx, job.ID, job.Kind, err)
	}
	return nil
}
