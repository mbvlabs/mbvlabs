package admin

import (
	"errors"
	"fmt"
	"log/slog"
	"mbvlabs/internal/inertia"
	"mbvlabs/internal/storage"
	"mbvlabs/internal/validation"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/cookies"
	"mbvlabs/router/routes"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v5"
)

type BlogPosts struct {
	db storage.Pool
}

func NewBlogPosts(db storage.Pool) BlogPosts {
	return BlogPosts{db: db}
}

type BlogPostData struct {
	ID            int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Slug          string
	Excerpt       string
	Body          string
	Status        string
	CoverImageUrl string
	Tags          string
	PublishedAt   string
}

func newBlogPostData(entity models.BlogPostEntity) BlogPostData {
	return BlogPostData{
		ID:            entity.ID,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Title:         entity.Title,
		Slug:          entity.Slug,
		Excerpt:       entity.Excerpt.String,
		Body:          entity.Body,
		Status:        entity.Status,
		CoverImageUrl: entity.CoverImageUrl.String,
		Tags:          jsonArrayCSV(entity.Tags),
		PublishedAt:   adminDateString(entity.PublishedAt),
	}
}

func newBlogPostDataList(entities []models.BlogPostEntity) []BlogPostData {
	items := make([]BlogPostData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newBlogPostData(entity))
	}

	return items
}

func (bp BlogPosts) Index(etx *echo.Context) error {
	page := int64(1)
	if p := etx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	perPage := int64(25)
	if pp := etx.QueryParam("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 &&
			parsed <= 100 {
			perPage = int64(parsed)
		}
	}

	blogPostsList, err := models.BlogPost.Paginate(
		etx.Request().Context(),
		bp.db.Executor(),
		page,
		perPage,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/BlogPost/Index", inertia.Props{
		"items": newBlogPostDataList(blogPostsList.BlogPosts),
	})
}

func (bp BlogPosts) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	blogPostID := int32(parsed)

	blogPost, err := models.BlogPost.Find(etx.Request().Context(), bp.db.Executor(), blogPostID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/BlogPost/Show", inertia.Props{
		"item": newBlogPostData(blogPost),
	})
}

func (bp BlogPosts) New(etx *echo.Context) error {
	return inertia.Page(etx, "Admin/BlogPost/Create", inertia.Props{})
}

type CreateBlogPostFormPayload struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Excerpt       string `json:"excerpt"`
	Body          string `json:"body"`
	Status        string `json:"status"`
	CoverImageUrl string `json:"coverImageUrl"`
	Tags          string `json:"tags"`
}

func (bp BlogPosts) Create(etx *echo.Context) error {
	var payload CreateBlogPostFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateBlogPostFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	slugSource := payload.Slug
	if slugSource == "" {
		slugSource = payload.Title
	}

	data := models.CreateBlogPostData{
		Title:         payload.Title,
		Slug:          slug.Make(slugSource),
		Excerpt:       payload.Excerpt,
		Body:          payload.Body,
		Status:        payload.Status,
		CoverImageUrl: payload.CoverImageUrl,
		Tags:          strings.FieldsFunc(payload.Tags, func(r rune) bool { return r == ',' }),
	}

	blogPost, err := models.BlogPost.Create(
		etx.Request().Context(),
		bp.db.Executor(),
		data,
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to create blog post",
			"error",
			err,
			"title",
			payload.Title,
		)
		if errs, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/BlogPost/Create",
				inertia.Props{},
				inertia.WithValidationErrors(errs.ToMap()),
			)
		}

		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create blog post: %v", err),
		); flashErr != nil {
			return flashErr
		}

		return inertia.Redirect(etx, routes.AdminBlogPostNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Blog post created successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminBlogPostShow.URL(blogPost.ID))
}

func (bp BlogPosts) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	blogPostID := int32(parsed)

	blogPost, err := models.BlogPost.Find(etx.Request().Context(), bp.db.Executor(), blogPostID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/BlogPost/Edit", inertia.Props{
		"item": newBlogPostData(blogPost),
	})
}

type UpdateBlogPostFormPayload struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Excerpt       string `json:"excerpt"`
	Body          string `json:"body"`
	Status        string `json:"status"`
	CoverImageUrl string `json:"coverImageUrl"`
	Tags          string `json:"tags"`
}

func blogPostDataFromUpdatePayload(id int32, payload UpdateBlogPostFormPayload) BlogPostData {
	return BlogPostData{
		ID:            id,
		Title:         payload.Title,
		Slug:          payload.Slug,
		Excerpt:       payload.Excerpt,
		Body:          payload.Body,
		Status:        payload.Status,
		CoverImageUrl: payload.CoverImageUrl,
		Tags:          payload.Tags,
	}
}

func (bp BlogPosts) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	blogPostID := int32(parsed)

	var payload UpdateBlogPostFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateBlogPostFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	slugSource := payload.Slug
	if slugSource == "" {
		slugSource = payload.Title
	}

	data := models.UpdateBlogPostData{
		ID:            blogPostID,
		Title:         payload.Title,
		Slug:          slug.Make(slugSource),
		Excerpt:       payload.Excerpt,
		Body:          payload.Body,
		Status:        payload.Status,
		CoverImageUrl: payload.CoverImageUrl,
		Tags:          strings.FieldsFunc(payload.Tags, func(r rune) bool { return r == ',' }),
	}

	blogPost, err := models.BlogPost.Update(
		etx.Request().Context(),
		bp.db.Executor(),
		data,
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to update blog post",
			"id",
			blogPostID,
			"error",
			err,
		)
		if errs, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/BlogPost/Edit",
				inertia.Props{"item": blogPostDataFromUpdatePayload(blogPostID, payload)},
				inertia.WithValidationErrors(errs.ToMap()),
			)
		}

		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update blog post: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminBlogPostEdit.URL(blogPostID),
		)
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Blog post updated successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminBlogPostShow.URL(blogPost.ID))
}

func (bp BlogPosts) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	blogPostID := int32(parsed)

	err = models.BlogPost.Destroy(etx.Request().Context(), bp.db.Executor(), blogPostID)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete blogPost: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminBlogPostIndex.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"BlogPost destroyed successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminBlogPostIndex.URL())
}

func (bp BlogPosts) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminBlogPostIndex.Path(),
		Name:        routes.AdminBlogPostIndex.Name(),
		Handler:     bp.Index,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminBlogPostShow.Path(),
		Name:        routes.AdminBlogPostShow.Name(),
		Handler:     bp.Show,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminBlogPostNew.Path(),
		Name:        routes.AdminBlogPostNew.Name(),
		Handler:     bp.New,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodPost,
		Path:        routes.AdminBlogPostCreate.Path(),
		Name:        routes.AdminBlogPostCreate.Name(),
		Handler:     bp.Create,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminBlogPostEdit.Path(),
		Name:        routes.AdminBlogPostEdit.Name(),
		Handler:     bp.Edit,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodPut,
		Path:        routes.AdminBlogPostUpdate.Path(),
		Name:        routes.AdminBlogPostUpdate.Name(),
		Handler:     bp.Update,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodDelete,
		Path:        routes.AdminBlogPostDestroy.Path(),
		Name:        routes.AdminBlogPostDestroy.Name(),
		Handler:     bp.Destroy,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
