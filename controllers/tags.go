package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router/cookies"
	"mbvlabs/router/routes"
	"mbvlabs/views"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Tags struct {
	db storage.Pool
}

func NewTags(db storage.Pool) Tags {
	return Tags{db}
}

func (r Tags) Index(c echo.Context) error {
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

	tagsList, err := models.PaginateTags(
		c.Request().Context(),
		r.db.Conn(),
		page,
		perPage,
	)
	if err != nil {
		return render(c, views.InternalError())
	}

	return render(c, views.TagIndex(tagsList.Tags))
}

func (r Tags) Show(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	tag, err := models.FindTag(c.Request().Context(), r.db.Conn(), tagID)
	if err != nil {
		return render(c, views.NotFound())
	}

	return render(c, views.TagShow(tag))
}

func (r Tags) New(c echo.Context) error {
	return render(c, views.TagNew())
}

type CreateTagFormPayload struct {
	Name        string `form:"name"`
	Slug        string `form:"slug"`
	Description string `form:"description"`
	Color       string `form:"color"`
}

func (r Tags) Create(c echo.Context) error {
	var payload CreateTagFormPayload
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not parse CreateTagFormPayload",
			"error",
			err,
		)

		return render(c, views.NotFound())
	}

	data := models.CreateTagData{
		Name:        payload.Name,
		Slug:        payload.Slug,
		Description: payload.Description,
		Color:       payload.Color,
	}

	tag, err := models.CreateTag(
		c.Request().Context(),
		r.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to create tag: %v", err)); flashErr != nil {
			return flashErr
		}
		return c.Redirect(http.StatusSeeOther, routes.TagNew.URL())
	}

	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "Tag created successfully"); flashErr != nil {
		return render(c, views.InternalError())
	}

	return c.Redirect(http.StatusSeeOther, routes.TagShow.URL(tag.ID))
}

func (r Tags) Edit(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	tag, err := models.FindTag(c.Request().Context(), r.db.Conn(), tagID)
	if err != nil {
		return render(c, views.NotFound())
	}

	return render(c, views.TagEdit(tag))
}

type UpdateTagFormPayload struct {
	Name        string `form:"name"`
	Slug        string `form:"slug"`
	Description string `form:"description"`
	Color       string `form:"color"`
}

func (r Tags) Update(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	var payload UpdateTagFormPayload
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not parse UpdateTagFormPayload",
			"error",
			err,
		)

		return render(c, views.NotFound())
	}

	data := models.UpdateTagData{
		ID:          tagID,
		Name:        payload.Name,
		Slug:        payload.Slug,
		Description: payload.Description,
		Color:       payload.Color,
	}

	tag, err := models.UpdateTag(
		c.Request().Context(),
		r.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to update tag: %v", err)); flashErr != nil {
			return render(c, views.InternalError())
		}
		return c.Redirect(
			http.StatusSeeOther,
			routes.TagEdit.URL(tagID),
		)
	}

	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "Tag updated successfully"); flashErr != nil {
		return render(c, views.InternalError())
	}

	return c.Redirect(http.StatusSeeOther, routes.TagShow.URL(tag.ID))
}

func (r Tags) Destroy(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	err = models.DestroyTag(c.Request().Context(), r.db.Conn(), tagID)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to delete tag: %v", err)); flashErr != nil {
			return render(c, views.InternalError())
		}
		return c.Redirect(http.StatusSeeOther, routes.TagIndex.URL())
	}

	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "Tag destroyed successfully"); flashErr != nil {
		return render(c, views.InternalError())
	}

	return c.Redirect(http.StatusSeeOther, routes.TagIndex.URL())
}
