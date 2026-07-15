package admin

import (
	"errors"
	"log/slog"
	"mbvlabs/internal/inertia"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/queue"
	"mbvlabs/router"
	"mbvlabs/router/routes"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Admin struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
}

func NewAdmin(
	db storage.Pool,
	insertOnly queue.InsertOnly,
) Admin {
	return Admin{db, insertOnly}
}

func (a Admin) RegisterRoutes(r *router.Router) error {
	errs := []error{}

	_, err := r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminHomePage.Path(),
		Name:        routes.AdminHomePage.Name(),
		Handler:     a.Home,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodHead,
		Path:        routes.AdminHomePage.Path(),
		Name:        routes.AdminHomePage.Name() + ".head",
		Handler:     a.Home,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (a Admin) Home(etx *echo.Context) error {

	var latestInquiries []models.ProjectInquiryEntity
	if err := a.db.Executor().NewSelect().
		Model(&latestInquiries).
		Order("created_at DESC").
		Limit(5).
		Scan(etx.Request().Context()); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to fetch latest project inquiries",
			"error",
			err,
		)
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	var latestWorks []models.WorkEntity
	if err := a.db.Executor().NewSelect().
		Model(&latestWorks).
		Order("created_at DESC").
		Limit(5).
		Scan(etx.Request().Context()); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to fetch latest works",
			"error",
			err,
		)
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	var latestPosts []models.BlogPostEntity
	if err := a.db.Executor().NewSelect().
		Model(&latestPosts).
		Order("created_at DESC").
		Limit(5).
		Scan(etx.Request().Context()); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to fetch latest blog posts",
			"error",
			err,
		)
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	blogPosts, err := newBlogPostDataList(latestPosts)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	var latestProjects []models.ProjectEntity
	if err := a.db.Executor().NewSelect().
		Model(&latestProjects).
		Order("created_at DESC").
		Limit(5).
		Scan(etx.Request().Context()); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to fetch latest projects",
			"error",
			err,
		)
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Home", inertia.Props{
		"appName":          "mbvlabs",
		"projectInquiries": newProjectInquiryDataList(latestInquiries),
		"works":            newWorkDataList(latestWorks),
		"blogPosts":        blogPosts,
		"projects":         newProjectDataList(latestProjects),
	})
}
