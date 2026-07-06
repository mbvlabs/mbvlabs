package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"mbvlabs/internal/inertia"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/cookies"
	"mbvlabs/router/routes"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

type ProjectInquiries struct {
	db storage.Pool
}

func NewProjectInquiries(db storage.Pool) ProjectInquiries {
	return ProjectInquiries{db}
}

type ProjectInquiryData struct {
	ID          int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Email       string
	Company     string
	Role        string
	ProjectType string
	Timeline    string
	Message     string
	Source      string
	Status      string
	Metadata    string
}

func newProjectInquiryData(entity models.ProjectInquiryEntity) ProjectInquiryData {
	return ProjectInquiryData{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Name:        entity.Name,
		Email:       entity.Email,
		Company:     entity.Company.String,
		Role:        entity.Role.String,
		ProjectType: entity.ProjectType.String,
		Timeline:    entity.Timeline.String,
		Message:     entity.Message,
		Source:      entity.Source.String,
		Status:      entity.Status,
		Metadata:    objectJSONString(entity.Metadata),
	}
}

func newProjectInquiryDataList(entities []models.ProjectInquiryEntity) []ProjectInquiryData {
	items := make([]ProjectInquiryData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newProjectInquiryData(entity))
	}

	return items
}

func (pi ProjectInquiries) Index(etx *echo.Context) error {
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

	projectInquiriesList, err := models.ProjectInquiry.Paginate(
		etx.Request().Context(),
		pi.db.Executor(),
		page,
		perPage,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/ProjectInquiry/Index", inertia.Props{
		"items": newProjectInquiryDataList(projectInquiriesList.ProjectInquiries),
	})
}

func (pi ProjectInquiries) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectInquiryID := int32(parsed)

	projectInquiry, err := models.ProjectInquiry.Find(
		etx.Request().Context(),
		pi.db.Executor(),
		projectInquiryID,
	)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/ProjectInquiry/Show", inertia.Props{
		"item": newProjectInquiryData(projectInquiry),
	})
}

func (pi ProjectInquiries) New(etx *echo.Context) error {
	return inertia.Page(etx, "Admin/ProjectInquiry/Create", inertia.Props{})
}

type CreateProjectInquiryFormPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Company     string `json:"company"`
	Role        string `json:"role"`
	ProjectType string `json:"projectType"`
	Timeline    string `json:"timeline"`
	Message     string `json:"message"`
	Source      string `json:"source"`
	Status      string `json:"status"`
	Metadata    string `json:"metadata"`
}

func (pi ProjectInquiries) Create(etx *echo.Context) error {
	var payload CreateProjectInquiryFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateProjectInquiryFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	metadata, err := objectJSON(payload.Metadata)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			err.Error(),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, routes.AdminProjectInquiryNew.URL())
	}

	data := models.CreateProjectInquiryData{
		Name:        payload.Name,
		Email:       payload.Email,
		Company:     sql.NullString{String: payload.Company, Valid: true},
		Role:        sql.NullString{String: payload.Role, Valid: true},
		ProjectType: sql.NullString{String: payload.ProjectType, Valid: true},
		Timeline:    sql.NullString{String: payload.Timeline, Valid: true},
		Message:     payload.Message,
		Source:      sql.NullString{String: payload.Source, Valid: true},
		Status:      payload.Status,
		Metadata:    metadata,
	}

	projectInquiry, err := models.ProjectInquiry.Create(
		etx.Request().Context(),
		pi.db.Executor(),
		data,
	)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create projectInquiry: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, routes.AdminProjectInquiryNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"ProjectInquiry created successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminProjectInquiryShow.URL(projectInquiry.ID))
}

func (pi ProjectInquiries) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectInquiryID := int32(parsed)

	projectInquiry, err := models.ProjectInquiry.Find(
		etx.Request().Context(),
		pi.db.Executor(),
		projectInquiryID,
	)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/ProjectInquiry/Edit", inertia.Props{
		"item": newProjectInquiryData(projectInquiry),
	})
}

type UpdateProjectInquiryFormPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Company     string `json:"company"`
	Role        string `json:"role"`
	ProjectType string `json:"projectType"`
	Timeline    string `json:"timeline"`
	Message     string `json:"message"`
	Source      string `json:"source"`
	Status      string `json:"status"`
	Metadata    string `json:"metadata"`
}

func (pi ProjectInquiries) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectInquiryID := int32(parsed)

	var payload UpdateProjectInquiryFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateProjectInquiryFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	metadata, err := objectJSON(payload.Metadata)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			err.Error(),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminProjectInquiryEdit.URL(projectInquiryID),
		)
	}

	data := models.UpdateProjectInquiryData{
		ID:          projectInquiryID,
		Name:        payload.Name,
		Email:       payload.Email,
		Company:     sql.NullString{String: payload.Company, Valid: true},
		Role:        sql.NullString{String: payload.Role, Valid: true},
		ProjectType: sql.NullString{String: payload.ProjectType, Valid: true},
		Timeline:    sql.NullString{String: payload.Timeline, Valid: true},
		Message:     payload.Message,
		Source:      sql.NullString{String: payload.Source, Valid: true},
		Status:      payload.Status,
		Metadata:    metadata,
	}

	projectInquiry, err := models.ProjectInquiry.Update(
		etx.Request().Context(),
		pi.db.Executor(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update projectInquiry: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminProjectInquiryEdit.URL(projectInquiryID),
		)
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"ProjectInquiry updated successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminProjectInquiryShow.URL(projectInquiry.ID))
}

func (pi ProjectInquiries) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectInquiryID := int32(parsed)

	err = models.ProjectInquiry.Destroy(etx.Request().Context(), pi.db.Executor(), projectInquiryID)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete projectInquiry: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminProjectInquiryIndex.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"ProjectInquiry destroyed successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminProjectInquiryIndex.URL())
}

func (pi ProjectInquiries) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminProjectInquiryIndex.Path(),
		Name:    routes.AdminProjectInquiryIndex.Name(),
		Handler: pi.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminProjectInquiryShow.Path(),
		Name:    routes.AdminProjectInquiryShow.Name(),
		Handler: pi.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminProjectInquiryNew.Path(),
		Name:    routes.AdminProjectInquiryNew.Name(),
		Handler: pi.New,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.AdminProjectInquiryCreate.Path(),
		Name:    routes.AdminProjectInquiryCreate.Name(),
		Handler: pi.Create,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminProjectInquiryEdit.Path(),
		Name:    routes.AdminProjectInquiryEdit.Name(),
		Handler: pi.Edit,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPut,
		Path:    routes.AdminProjectInquiryUpdate.Path(),
		Name:    routes.AdminProjectInquiryUpdate.Name(),
		Handler: pi.Update,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodDelete,
		Path:    routes.AdminProjectInquiryDestroy.Path(),
		Name:    routes.AdminProjectInquiryDestroy.Name(),
		Handler: pi.Destroy,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
