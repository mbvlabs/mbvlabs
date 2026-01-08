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

type Categorys struct {
	db storage.Pool
}

func NewCategorys(db storage.Pool) Categorys {
	return Categorys{db}
}

func (r Categorys) Index(c echo.Context) error {
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

	categoriesList, err := models.PaginateCategorys(
		c.Request().Context(),
		r.db.Conn(),
		page,
		perPage,
	)
	if err != nil {
		return render(c, views.InternalError())
	}

	return render(c, views.CategoryIndex(categoriesList.Categorys))
}

func (r Categorys) Show(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	category, err := models.FindCategory(c.Request().Context(), r.db.Conn(), categoryID)
	if err != nil {
		return render(c, views.NotFound())
	}

	return render(c, views.CategoryShow(category))
}

func (r Categorys) New(c echo.Context) error {
	return render(c, views.CategoryNew())
}

type CreateCategoryFormPayload struct {
	Name        string `form:"name"`
	Slug        string `form:"slug"`
	Description string `form:"description"`
	Color       string `form:"color"`
}

func (r Categorys) Create(c echo.Context) error {
	var payload CreateCategoryFormPayload
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not parse CreateCategoryFormPayload",
			"error",
			err,
		)

		return render(c, views.NotFound())
	}

	data := models.CreateCategoryData{
		Name:        payload.Name,
		Slug:        payload.Slug,
		Description: payload.Description,
		Color:       payload.Color,
	}

	category, err := models.CreateCategory(
		c.Request().Context(),
		r.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to create category: %v", err)); flashErr != nil {
			return flashErr
		}
		return c.Redirect(http.StatusSeeOther, routes.CategoryNew.URL())
	}

	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "Category created successfully"); flashErr != nil {
		return render(c, views.InternalError())
	}

	return c.Redirect(http.StatusSeeOther, routes.CategoryShow.URL(category.ID))
}

func (r Categorys) Edit(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	category, err := models.FindCategory(c.Request().Context(), r.db.Conn(), categoryID)
	if err != nil {
		return render(c, views.NotFound())
	}

	return render(c, views.CategoryEdit(category))
}

type UpdateCategoryFormPayload struct {
	Name        string `form:"name"`
	Slug        string `form:"slug"`
	Description string `form:"description"`
	Color       string `form:"color"`
}

func (r Categorys) Update(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	var payload UpdateCategoryFormPayload
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not parse UpdateCategoryFormPayload",
			"error",
			err,
		)

		return render(c, views.NotFound())
	}

	data := models.UpdateCategoryData{
		ID:          categoryID,
		Name:        payload.Name,
		Slug:        payload.Slug,
		Description: payload.Description,
		Color:       payload.Color,
	}

	category, err := models.UpdateCategory(
		c.Request().Context(),
		r.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to update category: %v", err)); flashErr != nil {
			return render(c, views.InternalError())
		}
		return c.Redirect(
			http.StatusSeeOther,
			routes.CategoryEdit.URL(categoryID),
		)
	}

	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "Category updated successfully"); flashErr != nil {
		return render(c, views.InternalError())
	}

	return c.Redirect(http.StatusSeeOther, routes.CategoryShow.URL(category.ID))
}

func (r Categorys) Destroy(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return render(c, views.BadRequest())
	}

	err = models.DestroyCategory(c.Request().Context(), r.db.Conn(), categoryID)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to delete category: %v", err)); flashErr != nil {
			return render(c, views.InternalError())
		}
		return c.Redirect(http.StatusSeeOther, routes.CategoryIndex.URL())
	}

	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "Category destroyed successfully"); flashErr != nil {
		return render(c, views.InternalError())
	}

	return c.Redirect(http.StatusSeeOther, routes.CategoryIndex.URL())
}
