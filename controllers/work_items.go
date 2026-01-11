package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router/cookies"
	"mbvlabs/router/routes"
	"mbvlabs/views"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WorkItems struct {
	db storage.Pool
}

func NewWorkItems(db storage.Pool) WorkItems {
	return WorkItems{db}
}

func (r WorkItems) Index(c echo.Context) error {
	page := int64(1)
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	perPage := int64(25)
	if pp := c.QueryParam("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 &&
			parsed <= 100 {
			perPage = int64(parsed)
		}
	}

	workItemsList, err := models.PaginateWorkItems(
		c.Request().Context(),
		r.db.Conn(),
		page,
		perPage,
	)
	if err != nil {
		return render(c, views.InternalError())
	}

	return render(c, views.WorkItemIndex(workItemsList.WorkItems))
}

func (r WorkItems) Show(c echo.Context) error {
	workItemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	workItem, err := models.FindWorkItem(c.Request().Context(), r.db.Conn(), workItemID)
	if err != nil {
		return render(c, views.NotFound())
	}

	return render(c, views.WorkItemShow(workItem))
}

func (r WorkItems) New(c echo.Context) error {
	return render(c, views.WorkItemNew())
}

type CreateWorkItemFormPayload struct {
	Title            string `form:"title"`
	Slug             string `form:"slug"`
	ShortDescription string `form:"short_description"`
	Content          string `form:"content"`
	Client           string `form:"client"`
	Industry         string `form:"industry"`
	ProjectDate      string `form:"project_date"`
	ProjectDuration  string `form:"project_duration"`
	HeroImageUrl     string `form:"hero_image_url"`
	HeroImageAlt     string `form:"hero_image_alt"`
	ExternalUrl      string `form:"external_url"`
	IsPublished      bool   `form:"is_published"`
	IsFeatured       bool   `form:"is_featured"`
	DisplayOrder     int32  `form:"display_order"`
	MetaTitle        string `form:"meta_title"`
	MetaDescription  string `form:"meta_description"`
	MetaKeywords     string `form:"meta_keywords"`
}

func (r WorkItems) Create(c echo.Context) error {
	var payload CreateWorkItemFormPayload
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not parse CreateWorkItemFormPayload",
			"error",
			err,
		)

		return render(c, views.NotFound())
	}

	data := models.CreateWorkItemData{
		Title:            payload.Title,
		Slug:             payload.Slug,
		ShortDescription: payload.ShortDescription,
		Content:          payload.Content,
		Client:           payload.Client,
		Industry:         payload.Industry,
		ProjectDate: func() time.Time {
			if payload.ProjectDate == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.ProjectDate); err == nil {
				return t
			}
			return time.Time{}
		}(),
		ProjectDuration: payload.ProjectDuration,
		HeroImageUrl:    payload.HeroImageUrl,
		HeroImageAlt:    payload.HeroImageAlt,
		ExternalUrl:     payload.ExternalUrl,
		IsPublished:     payload.IsPublished,
		IsFeatured:      payload.IsFeatured,
		DisplayOrder:    payload.DisplayOrder,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		MetaKeywords:    []string{payload.MetaKeywords},
	}

	workItem, err := models.CreateWorkItem(
		c.Request().Context(),
		r.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to create workItem: %v", err)); flashErr != nil {
			return flashErr
		}
		return c.Redirect(http.StatusSeeOther, routes.WorkItemNew.URL())
	}

	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "WorkItem created successfully"); flashErr != nil {
		return render(c, views.InternalError())
	}

	return c.Redirect(http.StatusSeeOther, routes.WorkItemShow.URL(workItem.ID))
}

func (r WorkItems) Edit(c echo.Context) error {
	workItemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	workItem, err := models.FindWorkItem(c.Request().Context(), r.db.Conn(), workItemID)
	if err != nil {
		return render(c, views.NotFound())
	}

	return render(c, views.WorkItemEdit(workItem))
}

type UpdateWorkItemFormPayload struct {
	Title            string `form:"title"`
	Slug             string `form:"slug"`
	ShortDescription string `form:"short_description"`
	Content          string `form:"content"`
	Client           string `form:"client"`
	Industry         string `form:"industry"`
	ProjectDate      string `form:"project_date"`
	ProjectDuration  string `form:"project_duration"`
	HeroImageUrl     string `form:"hero_image_url"`
	HeroImageAlt     string `form:"hero_image_alt"`
	ExternalUrl      string `form:"external_url"`
	IsPublished      bool   `form:"is_published"`
	IsFeatured       bool   `form:"is_featured"`
	DisplayOrder     int32  `form:"display_order"`
	MetaTitle        string `form:"meta_title"`
	MetaDescription  string `form:"meta_description"`
	MetaKeywords     string `form:"meta_keywords"`
}

func (r WorkItems) Update(c echo.Context) error {
	workItemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	var payload UpdateWorkItemFormPayload
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not parse UpdateWorkItemFormPayload",
			"error",
			err,
		)

		return render(c, views.NotFound())
	}

	data := models.UpdateWorkItemData{
		ID:               workItemID,
		Title:            payload.Title,
		Slug:             payload.Slug,
		ShortDescription: payload.ShortDescription,
		Content:          payload.Content,
		Client:           payload.Client,
		Industry:         payload.Industry,
		ProjectDate: func() time.Time {
			if payload.ProjectDate == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.ProjectDate); err == nil {
				return t
			}
			return time.Time{}
		}(),
		ProjectDuration: payload.ProjectDuration,
		HeroImageUrl:    payload.HeroImageUrl,
		HeroImageAlt:    payload.HeroImageAlt,
		ExternalUrl:     payload.ExternalUrl,
		IsPublished:     payload.IsPublished,
		IsFeatured:      payload.IsFeatured,
		DisplayOrder:    payload.DisplayOrder,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		MetaKeywords:    []string{payload.MetaKeywords},
	}

	workItem, err := models.UpdateWorkItem(
		c.Request().Context(),
		r.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to update workItem: %v", err)); flashErr != nil {
			return render(c, views.InternalError())
		}
		return c.Redirect(
			http.StatusSeeOther,
			routes.WorkItemEdit.URL(workItemID),
		)
	}

	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "WorkItem updated successfully"); flashErr != nil {
		return render(c, views.InternalError())
	}

	return c.Redirect(http.StatusSeeOther, routes.WorkItemShow.URL(workItem.ID))
}

func (r WorkItems) Destroy(c echo.Context) error {
	workItemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	err = models.DestroyWorkItem(c.Request().Context(), r.db.Conn(), workItemID)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to delete workItem: %v", err)); flashErr != nil {
			return render(c, views.InternalError())
		}
		return c.Redirect(http.StatusSeeOther, routes.WorkItemIndex.URL())
	}

	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "WorkItem destroyed successfully"); flashErr != nil {
		return render(c, views.InternalError())
	}

	return c.Redirect(http.StatusSeeOther, routes.WorkItemIndex.URL())
}
