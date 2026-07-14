package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"mbvlabs/internal/hypermedia"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/cookies"
	"mbvlabs/router/middleware"
	"mbvlabs/router/routes"
	"mbvlabs/views"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

type ProjectInquiries struct {
	db storage.Pool
}

func NewProjectInquiries(db storage.Pool) ProjectInquiries {
	return ProjectInquiries{db}
}

func (pi ProjectInquiries) Index(etx *echo.Context) error {
	return hypermedia.RenderPage(etx, views.ProjectInquiryIndex{}.Page())
}

type createProjectInquiryPayload struct {
	Name        string `form:"name" json:"name"`
	Email       string `form:"email" json:"email"`
	Company     string `form:"company" json:"company"`
	Role        string `form:"role" json:"role"`
	ProjectType string `form:"projectType" json:"projectType"`
	Timeline    string `form:"timeline" json:"timeline"`
	Message     string `form:"message" json:"message"`
}

func (pi ProjectInquiries) Create(etx *echo.Context) error {
	var payload createProjectInquiryPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse project inquiry form payload",
			"error",
			err,
		)
		return hypermedia.RenderPage(etx, views.BadRequest())
	}

	metadata, err := json.Marshal(map[string]string{
		"referer":     etx.Request().Referer(),
		"user_agent":  etx.Request().UserAgent(),
		"remote_addr": etx.RealIP(),
	})
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not marshal project inquiry metadata",
			"error",
			err,
		)
		return hypermedia.RenderPage(etx, views.InternalError())
	}

	_, err = models.ProjectInquiry.Create(
		etx.Request().Context(),
		pi.db.Executor(),
		models.CreateProjectInquiryData{
			Name:        strings.TrimSpace(payload.Name),
			Email:       strings.TrimSpace(payload.Email),
			Company:     optionalString(payload.Company),
			Role:        optionalString(payload.Role),
			ProjectType: optionalString(payload.ProjectType),
			Timeline:    optionalString(payload.Timeline),
			Message:     strings.TrimSpace(payload.Message),
			Source:      sql.NullString{String: "contact", Valid: true},
			Status:      "new",
			Metadata:    json.RawMessage(metadata),
		},
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to create project inquiry",
			"error",
			err,
		)

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			"Please check the form and try again.",
		); flashErr != nil {
			return hypermedia.RenderPage(etx, views.InternalError())
		}

		return etx.Redirect(http.StatusSeeOther, routes.ProjectInquiryIndex.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Thanks. I'll get back to you soon.",
	); flashErr != nil {
		return hypermedia.RenderPage(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.ProjectInquiryIndex.URL())
}

func optionalString(value string) sql.NullString {
	value = strings.TrimSpace(value)
	return sql.NullString{String: value, Valid: value != ""}
}

func (pi ProjectInquiries) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.ProjectInquiryIndex.Path(),
		Name:        routes.ProjectInquiryIndex.Name(),
		Handler:     pi.Index,
		Middlewares: []echo.MiddlewareFunc{middleware.MarkdownNegotiation},
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.ProjectInquiryCreate.Path(),
		Name:    routes.ProjectInquiryCreate.Name(),
		Handler: pi.Create,
		Middlewares: []echo.MiddlewareFunc{
			middleware.IPRateLimiter(3, routes.ProjectInquiryIndex),
		},
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.LegacyProjectInquiryIndex.Path(),
		Name:    routes.LegacyProjectInquiryIndex.Name(),
		Handler: pi.RedirectToIndex,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.LegacyProjectInquiryCreate.Path(),
		Name:    routes.LegacyProjectInquiryCreate.Name(),
		Handler: pi.Create,
		Middlewares: []echo.MiddlewareFunc{
			middleware.IPRateLimiter(3, routes.ProjectInquiryIndex),
		},
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (pi ProjectInquiries) RedirectToIndex(etx *echo.Context) error {
	return etx.Redirect(http.StatusMovedPermanently, routes.ProjectInquiryIndex.URL())
}
