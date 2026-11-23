package queue

import (
	"context"
	"fmt"

	"mbvlabs/emails"
	"mbvlabs/queue/jobs"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/postmand"
	"github.com/mbvlabs/andurel/pkg/routing"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/riverqueue/river"
)

const diaryReminderRecipient = "reminder@mbvlabs.com"

type sendDiaryReminderWorker struct {
	river.WorkerDefaults[jobs.SendDiaryReminderArgs]
	sender postmand.TransactionalSender
	from   string
	hosts  *routing.Hosts
}

func (w *sendDiaryReminderWorker) Work(ctx context.Context, job *river.Job[jobs.SendDiaryReminderArgs]) error {
	var subject string
	switch job.Args.Period {
	case "morning":
		subject = "Morning Diary Reminder"
	case "evening":
		subject = "Evening Diary Reminder"
	default:
		return river.JobCancel(fmt.Errorf("invalid diary reminder period %q", job.Args.Period))
	}
	email := emails.DiaryReminder{
		Title:    subject,
		Period:   job.Args.Period,
		DiaryURL: routes.AdminDiaryEntryToday.FullURL(w.hosts, routing.QueryParam("focus", job.Args.Period)),
	}
	data, err := postmand.Compose(email, diaryReminderRecipient, w.from, subject)
	if err != nil {
		return failJob(ctx, job.ID, job.Kind, err)
	}
	data.Metadata = map[string]string{"kind": "diary_reminder", "period": job.Args.Period}
	return storage.SendTransactional(ctx, data, w.sender)
}
