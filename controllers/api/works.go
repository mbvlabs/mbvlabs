package api

import (
	"errors"
	"fmt"
	"log/slog"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/routes"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type Works struct {
	db storage.Pool
}

func NewWorks(db storage.Pool) Works {
	return Works{db}
}

func (w Works) RegisterRoutes(r *router.Router) error {
	errs := []error{}
	var err error

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.ApiWorkCreate.Path(),
		Name:    routes.ApiWorkCreate.Name(),
		Handler: w.Create,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type CreateWorkPayload struct {
	Title          string   `json:"title"`
	Slug           string   `json:"slug"`
	ClientName     string   `json:"clientName"`
	ClientIndustry string   `json:"clientIndustry"`
	ClientUrl      string   `json:"clientUrl"`
	ClientLogoUrl  string   `json:"clientLogoUrl"`
	Summary        string   `json:"summary"`
	CoverImageUrl  string   `json:"coverImageUrl"`
	Specialisms    []string `json:"specialisms"`
	Platforms      []string `json:"platforms"`
	Technologies   []string `json:"technologies"`
	Challenge      string   `json:"challenge"`
	Approach       string   `json:"approach"`
	Deliverables   string   `json:"deliverables"`
	Outcome        string   `json:"outcome"`
	Content        string   `json:"content"`
	StartedAt      string   `json:"startedAt"`
	CompletedAt    string   `json:"completedAt"`
	Status         string   `json:"status"`
	PublishedAt    string   `json:"publishedAt"`
	IsFeatured     bool     `json:"isFeatured"`
}

func (w Works) Create(etx *echo.Context) error {
	var payload CreateWorkPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateWorkPayload",
			"error",
			err,
		)

		return etx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	data := models.CreateWorkData{
		Title:          payload.Title,
		Slug:           payload.Slug,
		ClientName:     payload.ClientName,
		ClientIndustry: payload.ClientIndustry,
		ClientURL:      payload.ClientUrl,
		ClientLogoURL:  payload.ClientLogoUrl,
		Summary:        payload.Summary,
		CoverImageURL:  payload.CoverImageUrl,
		Specialisms:    payload.Specialisms,
		Platforms:      payload.Platforms,
		Technologies:   payload.Technologies,
		Challenge:      payload.Challenge,
		Approach:       payload.Approach,
		Deliverables:   payload.Deliverables,
		Outcome:        payload.Outcome,
		Content:        payload.Content,
		StartedAt: func() time.Time {
			if t, err := time.Parse("2006-01-02", payload.StartedAt); err == nil {
				return t
			}

			return time.Time{}
		}(),
		CompletedAt: func() time.Time {
			if t, err := time.Parse("2006-01-02", payload.CompletedAt); err == nil {
				return t
			}

			return time.Time{}
		}(),
		Status: models.StatusEnum(payload.Status),
		PublishedAt: func() time.Time {
			if t, err := time.Parse("2006-01-02", payload.PublishedAt); err == nil {
				return t
			}

			return time.Time{}
		}(),
		IsFeatured: payload.IsFeatured,
	}

	work, err := models.Work.Create(
		etx.Request().Context(),
		w.db.Executor(),
		data,
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to create work",
			"error",
			err,
		)

		return etx.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": fmt.Sprintf("failed to create work: %v", err)},
		)
	}

	return etx.JSON(http.StatusCreated, work)
}
