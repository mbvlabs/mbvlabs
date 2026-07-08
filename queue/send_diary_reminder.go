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
	sender email.TransactionalSender
}

func NewSendDiaryReminderWorker(sender email.TransactionalSender) *SendDiaryReminderWorker {
	return &SendDiaryReminderWorker{sender: sender}
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

	err = email.SendTransactional(ctx, data, w.sender)
	if err != nil {
		if !email.IsRetryable(err) {
			return river.JobCancel(err)
		}
		return err
	}

	return nil
}

func diaryReminderURL(period string) string {
	return routes.AdminDiaryEntryToday.FullURL(
		config.BaseURL,
		routing.QueryParam("focus", period),
	)
}
