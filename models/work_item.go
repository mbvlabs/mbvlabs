package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"mbvlabs/internal/storage"
	"mbvlabs/models/internal/db"
)

type WorkItem struct {
	ID               uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Title            string
	Slug             string
	ShortDescription string
	Content          string
	Client           string
	Industry         string
	ProjectDate      time.Time
	ProjectDuration  string
	HeroImageUrl     string
	HeroImageAlt     string
	ExternalUrl      string
	IsPublished      bool
	IsFeatured       bool
	DisplayOrder     int32
	MetaTitle        string
	MetaDescription  string
	MetaKeywords     []string
}

func FindWorkItem(
	ctx context.Context,
	exec storage.Executor,
	id uuid.UUID,
) (WorkItem, error) {
	row, err := queries.QueryWorkItemByID(ctx, exec, id)
	if err != nil {
		return WorkItem{}, err
	}

	return rowToWorkItem(row), nil
}

type CreateWorkItemData struct {
	Title            string
	Slug             string
	ShortDescription string
	Content          string
	Client           string
	Industry         string
	ProjectDate      time.Time
	ProjectDuration  string
	HeroImageUrl     string
	HeroImageAlt     string
	ExternalUrl      string
	IsPublished      bool
	IsFeatured       bool
	DisplayOrder     int32
	MetaTitle        string
	MetaDescription  string
	MetaKeywords     []string
}

func CreateWorkItem(
	ctx context.Context,
	exec storage.Executor,
	data CreateWorkItemData,
) (WorkItem, error) {
	if err := validate.Struct(data); err != nil {
		return WorkItem{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.InsertWorkItemParams{
		ID:               uuid.New(),
		Title:            data.Title,
		Slug:             data.Slug,
		ShortDescription: data.ShortDescription,
		Content:          data.Content,
		Client:           pgtype.Text{String: data.Client, Valid: true},
		Industry:         pgtype.Text{String: data.Industry, Valid: true},
		ProjectDate:      pgtype.Date{Time: data.ProjectDate, Valid: true},
		ProjectDuration:  pgtype.Text{String: data.ProjectDuration, Valid: true},
		HeroImageUrl:     pgtype.Text{String: data.HeroImageUrl, Valid: true},
		HeroImageAlt:     pgtype.Text{String: data.HeroImageAlt, Valid: true},
		ExternalUrl:      pgtype.Text{String: data.ExternalUrl, Valid: true},
		IsPublished:      data.IsPublished,
		IsFeatured:       data.IsFeatured,
		DisplayOrder:     data.DisplayOrder,
		MetaTitle:        pgtype.Text{String: data.MetaTitle, Valid: true},
		MetaDescription:  pgtype.Text{String: data.MetaDescription, Valid: true},
		MetaKeywords:     data.MetaKeywords,
	}
	row, err := queries.InsertWorkItem(ctx, exec, params)
	if err != nil {
		return WorkItem{}, err
	}

	return rowToWorkItem(row), nil
}

type UpdateWorkItemData struct {
	ID               uuid.UUID
	UpdatedAt        time.Time
	Title            string
	Slug             string
	ShortDescription string
	Content          string
	Client           string
	Industry         string
	ProjectDate      time.Time
	ProjectDuration  string
	HeroImageUrl     string
	HeroImageAlt     string
	ExternalUrl      string
	IsPublished      bool
	IsFeatured       bool
	DisplayOrder     int32
	MetaTitle        string
	MetaDescription  string
	MetaKeywords     []string
}

func UpdateWorkItem(
	ctx context.Context,
	exec storage.Executor,
	data UpdateWorkItemData,
) (WorkItem, error) {
	if err := validate.Struct(data); err != nil {
		return WorkItem{}, errors.Join(ErrDomainValidation, err)
	}

	_, err := queries.QueryWorkItemByID(ctx, exec, data.ID)
	if err != nil {
		return WorkItem{}, err
	}

	params := db.UpdateWorkItemParams{
		ID:               data.ID,
		Title:            data.Title,
		Slug:             data.Slug,
		ShortDescription: data.ShortDescription,
		Content:          data.Content,
		Client:           pgtype.Text{String: data.Client, Valid: true},
		Industry:         pgtype.Text{String: data.Industry, Valid: true},
		ProjectDate:      pgtype.Date{Time: data.ProjectDate, Valid: true},
		ProjectDuration:  pgtype.Text{String: data.ProjectDuration, Valid: true},
		HeroImageUrl:     pgtype.Text{String: data.HeroImageUrl, Valid: true},
		HeroImageAlt:     pgtype.Text{String: data.HeroImageAlt, Valid: true},
		ExternalUrl:      pgtype.Text{String: data.ExternalUrl, Valid: true},
		IsPublished:      data.IsPublished,
		IsFeatured:       data.IsFeatured,
		DisplayOrder:     data.DisplayOrder,
		MetaTitle:        pgtype.Text{String: data.MetaTitle, Valid: true},
		MetaDescription:  pgtype.Text{String: data.MetaDescription, Valid: true},
		MetaKeywords:     data.MetaKeywords,
	}

	row, err := queries.UpdateWorkItem(ctx, exec, params)
	if err != nil {
		return WorkItem{}, err
	}

	return rowToWorkItem(row), nil
}

func DestroyWorkItem(
	ctx context.Context,
	exec storage.Executor,
	id uuid.UUID,
) error {
	return queries.DeleteWorkItem(ctx, exec, id)
}

func AllWorkItems(
	ctx context.Context,
	exec storage.Executor,
) ([]WorkItem, error) {
	rows, err := queries.QueryAllWorkItems(ctx, exec)
	if err != nil {
		return nil, err
	}

	workitems := make([]WorkItem, len(rows))
	for i, row := range rows {
		workitems[i] = rowToWorkItem(row)
	}

	return workitems, nil
}

type PaginatedWorkItems struct {
	WorkItems  []WorkItem
	TotalCount int64
	Page       int64
	PageSize   int64
	TotalPages int64
}

func PaginateWorkItems(
	ctx context.Context,
	exec storage.Executor,
	page int64,
	pageSize int64,
) (PaginatedWorkItems, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	totalCount, err := queries.CountWorkItems(ctx, exec)
	if err != nil {
		return PaginatedWorkItems{}, err
	}

	rows, err := queries.QueryPaginatedWorkItems(
		ctx,
		exec,
		db.QueryPaginatedWorkItemsParams{
			Limit:  pageSize,
			Offset: offset,
		},
	)
	if err != nil {
		return PaginatedWorkItems{}, err
	}

	workitems := make([]WorkItem, len(rows))
	for i, row := range rows {
		workitems[i] = rowToWorkItem(row)
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return PaginatedWorkItems{
		WorkItems:  workitems,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func UpsertWorkItem(
	ctx context.Context,
	exec storage.Executor,
	data CreateWorkItemData,
) (WorkItem, error) {
	if err := validate.Struct(data); err != nil {
		return WorkItem{}, errors.Join(ErrDomainValidation, err)
	}

	params := db.UpsertWorkItemParams{
		ID:               uuid.New(),
		Title:            data.Title,
		Slug:             data.Slug,
		ShortDescription: data.ShortDescription,
		Content:          data.Content,
		Client:           pgtype.Text{String: data.Client, Valid: true},
		Industry:         pgtype.Text{String: data.Industry, Valid: true},
		ProjectDate:      pgtype.Date{Time: data.ProjectDate, Valid: true},
		ProjectDuration:  pgtype.Text{String: data.ProjectDuration, Valid: true},
		HeroImageUrl:     pgtype.Text{String: data.HeroImageUrl, Valid: true},
		HeroImageAlt:     pgtype.Text{String: data.HeroImageAlt, Valid: true},
		ExternalUrl:      pgtype.Text{String: data.ExternalUrl, Valid: true},
		IsPublished:      data.IsPublished,
		IsFeatured:       data.IsFeatured,
		DisplayOrder:     data.DisplayOrder,
		MetaTitle:        pgtype.Text{String: data.MetaTitle, Valid: true},
		MetaDescription:  pgtype.Text{String: data.MetaDescription, Valid: true},
		MetaKeywords:     data.MetaKeywords,
	}
	row, err := queries.UpsertWorkItem(ctx, exec, params)
	if err != nil {
		return WorkItem{}, err
	}

	return rowToWorkItem(row), nil
}

func rowToWorkItem(row db.WorkItem) WorkItem {
	return WorkItem{
		ID:               row.ID,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
		Title:            row.Title,
		Slug:             row.Slug,
		ShortDescription: row.ShortDescription,
		Content:          row.Content,
		Client:           row.Client.String,
		Industry:         row.Industry.String,
		ProjectDate:      row.ProjectDate.Time,
		ProjectDuration:  row.ProjectDuration.String,
		HeroImageUrl:     row.HeroImageUrl.String,
		HeroImageAlt:     row.HeroImageAlt.String,
		ExternalUrl:      row.ExternalUrl.String,
		IsPublished:      row.IsPublished,
		IsFeatured:       row.IsFeatured,
		DisplayOrder:     row.DisplayOrder,
		MetaTitle:        row.MetaTitle.String,
		MetaDescription:  row.MetaDescription.String,
		MetaKeywords:     row.MetaKeywords,
	}
}
