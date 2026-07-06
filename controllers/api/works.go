package api

import (
	"errors"
	"fmt"
	"log/slog"
	"mbvlabs/internal/validation"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/routes"
	"mbvlabs/services"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v5"
)

type Works struct {
	workSvc services.Works
}

func NewWorks(works services.Works) Works {
	return Works{workSvc: works}
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

type errorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
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

		return etx.JSON(
			http.StatusBadRequest,
			errorResponse{Error: "invalid JSON request body"},
		)
	}

	validationErrors := map[string]string{}
	startedAt := parseDateField(validationErrors, "startedAt", payload.StartedAt)
	completedAt := parseDateField(validationErrors, "completedAt", payload.CompletedAt)
	publishedAt := parseDateField(validationErrors, "publishedAt", payload.PublishedAt)

	if len(validationErrors) > 0 {
		return etx.JSON(
			http.StatusBadRequest,
			errorResponse{Error: "validation failed", Fields: validationErrors},
		)
	}

	status := strings.TrimSpace(payload.Status)
	if status == "" {
		status = models.Draft.String()
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
		StartedAt:      startedAt,
		CompletedAt:    completedAt,
		Status:         models.StatusEnum(status),
		PublishedAt:    publishedAt,
		IsFeatured:     payload.IsFeatured,
	}

	work, err := w.workSvc.CreateWork(
		etx.Request().Context(),
		data,
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to create work",
			"error",
			err,
		)

		if errs, ok := validation.As(err); ok {
			return etx.JSON(
				http.StatusBadRequest,
				errorResponse{Error: "validation failed", Fields: errs.ToMap()},
			)
		}

		if errors.Is(err, services.ErrFeaturedWorkLimit) {
			return etx.JSON(
				http.StatusBadRequest,
				errorResponse{
					Error:  "validation failed",
					Fields: map[string]string{"isFeatured": "No more than 3 works can be featured."},
				},
			)
		}

		if isUniqueViolation(err, "works_slug_key") {
			return etx.JSON(
				http.StatusConflict,
				errorResponse{
					Error:  "work slug already exists",
					Fields: map[string]string{"slug": "must be unique"},
				},
			)
		}

		return etx.JSON(
			http.StatusInternalServerError,
			errorResponse{Error: fmt.Sprintf("failed to create work: %v", err)},
		)
	}

	return etx.JSON(http.StatusCreated, work)
}

func parseDateField(
	validationErrors map[string]string,
	field string,
	value string,
) time.Time {
	if strings.TrimSpace(value) == "" {
		return time.Time{}
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		validationErrors[field] = "must be a valid date in YYYY-MM-DD format"
		return time.Time{}
	}

	return parsed
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}
