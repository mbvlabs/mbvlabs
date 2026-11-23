package services

import (
	"context"
	"errors"

	"mbvlabs/domains"

	"github.com/mbvlabs/andurel/pkg/storage"
)

var (
	ErrFeaturedWorkLimit    = errors.New("no more than 3 works can be featured")
	ErrFeaturedProjectLimit = errors.New("no more than 3 projects can be featured")
)

// Works and Projects enforce the featured limit with a transaction that
// serializes each featured set's count and write.
type Works struct {
	works      *domains.Works
	transactor storage.Transactor
}

func NewWorks(works *domains.Works, transactor storage.Transactor) Works {
	return Works{works, transactor}
}
func (s Works) CreateWork(ctx context.Context, data domains.CreateWorkData) (domains.WorkEntity, error) {
	if !data.IsFeatured {
		return s.works.Create(ctx, data)
	}
	var item domains.WorkEntity
	err := storage.AtomicOps(ctx, s.transactor, func(ctx context.Context, tx storage.Tx) error {
		works := s.works.InTx(tx)
		if err := works.LockFeatured(ctx); err != nil {
			return err
		}
		count, err := works.CountFeatured(ctx, 0)
		if err != nil {
			return err
		}
		if count >= 3 {
			return ErrFeaturedWorkLimit
		}
		item, err = works.Create(ctx, data)
		return err
	})
	return item, err
}
func (s Works) UpdateWork(ctx context.Context, data domains.UpdateWorkData) (domains.WorkEntity, error) {
	if !data.IsFeatured {
		return s.works.Update(ctx, data)
	}
	var item domains.WorkEntity
	err := storage.AtomicOps(ctx, s.transactor, func(ctx context.Context, tx storage.Tx) error {
		works := s.works.InTx(tx)
		if err := works.LockFeatured(ctx); err != nil {
			return err
		}
		count, err := works.CountFeatured(ctx, data.ID)
		if err != nil {
			return err
		}
		if count >= 3 {
			return ErrFeaturedWorkLimit
		}
		item, err = works.Update(ctx, data)
		return err
	})
	return item, err
}

type Projects struct {
	projects   *domains.Projects
	transactor storage.Transactor
}

func NewProjects(projects *domains.Projects, transactor storage.Transactor) Projects {
	return Projects{projects, transactor}
}
func (s Projects) CreateProject(ctx context.Context, data domains.CreateProjectData) (domains.ProjectEntity, error) {
	if !data.IsFeatured {
		return s.projects.Create(ctx, data)
	}
	var item domains.ProjectEntity
	err := storage.AtomicOps(ctx, s.transactor, func(ctx context.Context, tx storage.Tx) error {
		projects := s.projects.InTx(tx)
		if err := projects.LockFeatured(ctx); err != nil {
			return err
		}
		count, err := projects.CountFeatured(ctx, 0)
		if err != nil {
			return err
		}
		if count >= 3 {
			return ErrFeaturedProjectLimit
		}
		item, err = projects.Create(ctx, data)
		return err
	})
	return item, err
}
func (s Projects) UpdateProject(ctx context.Context, data domains.UpdateProjectData) (domains.ProjectEntity, error) {
	if !data.IsFeatured {
		return s.projects.Update(ctx, data)
	}
	var item domains.ProjectEntity
	err := storage.AtomicOps(ctx, s.transactor, func(ctx context.Context, tx storage.Tx) error {
		projects := s.projects.InTx(tx)
		if err := projects.LockFeatured(ctx); err != nil {
			return err
		}
		count, err := projects.CountFeatured(ctx, data.ID)
		if err != nil {
			return err
		}
		if count >= 3 {
			return ErrFeaturedProjectLimit
		}
		item, err = projects.Update(ctx, data)
		return err
	})
	return item, err
}
