package queue

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"

	"mbvlabs/config"
	"mbvlabs/email"
	"mbvlabs/internal/routing"
	"mbvlabs/queue/jobs"
	"mbvlabs/router/routes"
)

type SendDiaryReminderWorker struct {
	river.WorkerDefaults[jobs.SendDiaryReminderArgs]
	insertOnly InsertOnly
}

func NewSendDiaryReminderWorker(insertOnly InsertOnly) *SendDiaryReminderWorker {
	return &SendDiaryReminderWorker{insertOnly: insertOnly}
}

func (w *SendDiaryReminderWorker) Register(workers *river.Workers) error {
	return river.AddWorkerSafely(workers, w)
}

func (w *SendDiaryReminderWorker) Work(ctx context.Context, job *river.Job[jobs.SendDiaryReminderArgs]) error {
	period := job.Args.Period
	if period != "morning" && period != "evening" {
		return river.JobCancel(fmt.Errorf("invalid diary reminder period %q", period))
	}

	data, err := email.NewDiaryReminderTransactionalData(
		period,
		diaryReminderURL(period),
		diaryReminderRecipient,
		config.DefaultSenderSignature,
	)
	if err != nil {
		return err
	}

	if _, err := w.insertOnly.Insert(ctx, jobs.SendTransactionalEmailArgs{Data: data}, nil); err != nil {
		return fmt.Errorf("queue diary reminder email: %w", err)
	}

	return nil
}

func diaryReminderURL(period string) string {
	return routes.AdminDiaryEntryToday.FullURL(
		config.BaseURL,
		routing.QueryParam("focus", period),
	)
}
