package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mbvlabs/internal/storage"
	"mbvlabs/models"
)

const maxFeaturedWorks = 3

var (
	ErrFeaturedWorkLimit = errors.New("no more than 3 works can be featured")
)

type Works struct {
	db storage.Pool
}

func NewWorks(db storage.Pool) Works {
	return Works{db: db}
}

func (w Works) CreateWork(
	ctx context.Context,
	data models.CreateWorkData,
) (models.WorkEntity, error) {
	if !data.IsFeatured {
		return models.Work.Create(ctx, w.db.Executor(), data)
	}

	tx, err := w.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return models.WorkEntity{}, fmt.Errorf("begin create work transaction: %v", err)
	}
	defer tx.Rollback()

	count, err := models.Work.CountFeatured(ctx, tx, 0)
	if err != nil {
		return models.WorkEntity{}, err
	}
	if count >= maxFeaturedWorks {
		return models.WorkEntity{}, ErrFeaturedWorkLimit
	}

	work, err := models.Work.Create(ctx, tx, data)
	if err != nil {
		return models.WorkEntity{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.WorkEntity{}, fmt.Errorf("commit create work transaction: %v", err)
	}

	return work, nil
}

func (w Works) UpdateWork(
	ctx context.Context,
	data models.UpdateWorkData,
) (models.WorkEntity, error) {
	if !data.IsFeatured {
		return models.Work.Update(ctx, w.db.Executor(), data)
	}

	tx, err := w.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return models.WorkEntity{}, fmt.Errorf("begin update work transaction: %v", err)
	}
	defer tx.Rollback()

	count, err := models.Work.CountFeatured(ctx, tx, data.ID)
	if err != nil {
		return models.WorkEntity{}, err
	}
	if count >= maxFeaturedWorks {
		return models.WorkEntity{}, ErrFeaturedWorkLimit
	}

	work, err := models.Work.Update(ctx, tx, data)
	if err != nil {
		return models.WorkEntity{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.WorkEntity{}, fmt.Errorf("commit update work transaction: %v", err)
	}

	return work, nil
}
