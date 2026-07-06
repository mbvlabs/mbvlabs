package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mbvlabs/internal/storage"
	"mbvlabs/models"
)

const maxFeaturedProjects = 3

var (
	ErrFeaturedProjectLimit = errors.New("no more than 3 projects can be featured")
)

type Projects struct {
	db storage.Pool
}

func NewProjects(db storage.Pool) Projects {
	return Projects{db: db}
}

func (p Projects) CreateProject(
	ctx context.Context,
	data models.CreateProjectData,
) (models.ProjectEntity, error) {
	if !data.IsFeatured {
		return models.Project.Create(ctx, p.db.Executor(), data)
	}

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return models.ProjectEntity{}, fmt.Errorf("begin create project transaction: %v", err)
	}
	defer tx.Rollback()

	count, err := models.Project.CountFeatured(ctx, tx, 0)
	if err != nil {
		return models.ProjectEntity{}, err
	}
	if count >= maxFeaturedProjects {
		return models.ProjectEntity{}, ErrFeaturedProjectLimit
	}

	project, err := models.Project.Create(ctx, tx, data)
	if err != nil {
		return models.ProjectEntity{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.ProjectEntity{}, fmt.Errorf("commit create project transaction: %v", err)
	}

	return project, nil
}

func (p Projects) UpdateProject(
	ctx context.Context,
	data models.UpdateProjectData,
) (models.ProjectEntity, error) {
	if !data.IsFeatured {
		return models.Project.Update(ctx, p.db.Executor(), data)
	}

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return models.ProjectEntity{}, fmt.Errorf("begin update project transaction: %v", err)
	}
	defer tx.Rollback()

	count, err := models.Project.CountFeatured(ctx, tx, data.ID)
	if err != nil {
		return models.ProjectEntity{}, err
	}
	if count >= maxFeaturedProjects {
		return models.ProjectEntity{}, ErrFeaturedProjectLimit
	}

	project, err := models.Project.Update(ctx, tx, data)
	if err != nil {
		return models.ProjectEntity{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.ProjectEntity{}, fmt.Errorf("commit update project transaction: %v", err)
	}

	return project, nil
}
