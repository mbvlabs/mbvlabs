package queue

import (
	"time"

	"github.com/riverqueue/river"
	"github.com/robfig/cron/v3"

	"mbvlabs/queue/jobs"
)

// PeriodicJobs returns the recurring jobs the worker process schedules.
func PeriodicJobs() []*river.PeriodicJob {
	return []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.DeleteExpiredTokensArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
		diaryReminderJob("morning", "CRON_TZ=Europe/Madrid 0 9 * * 1-5"),
		diaryReminderJob("evening", "CRON_TZ=Europe/Madrid 0 17 * * 1-5"),
	}
}

func diaryReminderJob(period, spec string) *river.PeriodicJob {
	schedule, err := cron.ParseStandard(spec)
	if err != nil {
		panic(err) // programmer error: these schedules are fixed above
	}
	return river.NewPeriodicJob(
		schedule,
		func() (river.JobArgs, *river.InsertOpts) {
			return jobs.SendDiaryReminderArgs{Period: period}, nil
		},
		&river.PeriodicJobOpts{ID: "diary_reminder_" + period},
	)
}
