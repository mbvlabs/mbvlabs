// Package queue provides queue-specific resources.
package queue

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"mbvlabs/internal/storage"
	"mbvlabs/queue/jobs"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivertype"
	"github.com/robfig/cron/v3"

	"go.uber.org/fx"
)

const diaryReminderRecipient = "reminder@mbvlabs.com"

type Processor struct {
	Client *river.Client[*sql.Tx]
}

func (p Processor) Shutdown(ctx context.Context) error {
	return p.Client.Stop(ctx)
}

func (p Processor) Start(ctx context.Context) error {
	return p.Client.Start(ctx)
}

func (p Processor) Stop(ctx context.Context) error {
	return p.Client.Stop(ctx)
}

func NewProcessor(
	db storage.Pool,
	workers *river.Workers,
) (Processor, error) {
	periodicJobs, err := diaryReminderPeriodicJobs()
	if err != nil {
		return Processor{}, err
	}

	riverClient, err := river.NewClient(riverdatabasesql.New(db.Conn()), &river.Config{
		PeriodicJobs: periodicJobs,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Logger:  slog.Default(),
		Workers: workers,
	})
	if err != nil {
		return Processor{}, err
	}

	return Processor{riverClient}, nil
}

func diaryReminderPeriodicJobs() ([]*river.PeriodicJob, error) {
	morning, err := diaryReminderPeriodicJob(
		"morning",
		"CRON_TZ=Europe/Madrid 0 9 * * *",
	)
	if err != nil {
		return nil, err
	}

	evening, err := diaryReminderPeriodicJob(
		"evening",
		"CRON_TZ=Europe/Madrid 0 17 * * *",
	)
	if err != nil {
		return nil, err
	}

	return []*river.PeriodicJob{morning, evening}, nil
}

func diaryReminderPeriodicJob(period string, cronSpec string) (*river.PeriodicJob, error) {
	schedule, err := cron.ParseStandard(cronSpec)
	if err != nil {
		return nil, fmt.Errorf("parse %s diary reminder cron: %w", period, err)
	}

	return river.NewPeriodicJob(
		schedule,
		func() (river.JobArgs, *river.InsertOpts) {
			return jobs.SendDiaryReminderArgs{Period: period}, nil
		},
		&river.PeriodicJobOpts{ID: "diary_reminder_" + period},
	), nil
}

type InsertOnly struct {
	client *river.Client[*sql.Tx]
}

// Insert implements storage.InsertQueue.
func (i *InsertOnly) Insert(
	ctx context.Context,
	args river.JobArgs,
	opts *river.InsertOpts,
) (*rivertype.JobInsertResult, error) {
	return i.client.Insert(ctx, args, opts)
}

// InsertMany implements storage.InsertQueue.
func (i *InsertOnly) InsertMany(
	ctx context.Context,
	params []river.InsertManyParams,
) ([]*rivertype.JobInsertResult, error) {
	return i.client.InsertMany(ctx, params)
}

// InsertManyFast implements storage.InsertQueue.
func (i *InsertOnly) InsertManyFast(
	ctx context.Context,
	params []river.InsertManyParams,
) (int, error) {
	return i.client.InsertManyFast(ctx, params)
}

// InsertManyFastTx implements storage.InsertQueue.
func (i *InsertOnly) InsertManyFastTx(
	ctx context.Context,
	tx *sql.Tx,
	params []river.InsertManyParams,
) (int, error) {
	return i.client.InsertManyFastTx(ctx, tx, params)
}

// InsertManyTx implements storage.InsertQueue.
func (i *InsertOnly) InsertManyTx(
	ctx context.Context,
	tx *sql.Tx,
	params []river.InsertManyParams,
) ([]*rivertype.JobInsertResult, error) {
	return i.client.InsertManyTx(ctx, tx, params)
}

// InsertTx implements storage.InsertQueue.
func (i *InsertOnly) InsertTx(
	ctx context.Context,
	tx *sql.Tx,
	args river.JobArgs,
	opts *river.InsertOpts,
) (*rivertype.JobInsertResult, error) {
	return i.client.InsertTx(ctx, tx, args, opts)
}

var _ storage.InsertQueue = (*InsertOnly)(nil)

func NewInsertOnly(db storage.Pool, workers *river.Workers) (InsertOnly, error) {
	riverClient, err := river.NewClient(riverdatabasesql.New(db.Conn()), &river.Config{
		Workers: workers,
	})
	if err != nil {
		return InsertOnly{}, err
	}

	return InsertOnly{riverClient}, nil
}

var Module = fx.Module(
	"queue",
	fx.Provide(
		NewInsertOnly,
		NewProcessor,
	),
)
