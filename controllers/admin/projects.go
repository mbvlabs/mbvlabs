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
	"mbvlabs/services"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v5"
)

type Projects struct {
	db  storage.Pool
	svc services.Projects
}

func NewProjects(db storage.Pool, svc services.Projects) Projects {
	return Projects{db: db, svc: svc}
}

type ProjectData struct {
	ID            int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Slug          string
	Tagline       string
	Description   string
	ProjectType   string
	RepositoryUrl string
	LiveUrl       string
	ImageUrl      string
	Technologies  string
	StartedAt     string
	LaunchedAt    string
	PublishedAt   string
	IsFeatured    bool
}

func newProjectData(entity models.ProjectEntity) ProjectData {
	return ProjectData{
		ID:            entity.ID,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Name:          entity.Name,
		Slug:          entity.Slug,
		Tagline:       entity.Tagline,
		Description:   entity.Description.String,
		ProjectType:   entity.ProjectType,
		RepositoryUrl: entity.RepositoryURL.String,
		LiveUrl:       entity.LiveURL.String,
		ImageUrl:      entity.ImageURL.String,
		Technologies:  jsonArrayCSV(entity.Technologies),
		StartedAt:     adminDateString(entity.StartedAt),
		LaunchedAt:    adminDateString(entity.LaunchedAt),
		PublishedAt:   adminDateString(entity.PublishedAt),
		IsFeatured:    entity.IsFeatured,
	}
}

func newProjectDataList(entities []models.ProjectEntity) []ProjectData {
	items := make([]ProjectData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newProjectData(entity))
	}

	return items
}

func (p Projects) Index(etx *echo.Context) error {
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

	projectsList, err := models.Project.Paginate(
		etx.Request().Context(),
		p.db.Executor(),
		page,
		perPage,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Project/Index", inertia.Props{
		"items": newProjectDataList(projectsList.Projects),
	})
}

func (p Projects) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectID := int32(parsed)

	project, err := models.Project.Find(etx.Request().Context(), p.db.Executor(), projectID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Project/Show", inertia.Props{
		"item": newProjectData(project),
	})
}

func (p Projects) New(etx *echo.Context) error {
	return inertia.Page(etx, "Admin/Project/Create", inertia.Props{})
}

type CreateProjectFormPayload struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Tagline       string `json:"tagline"`
	Description   string `json:"description"`
	ProjectType   string `json:"projectType"`
	RepositoryUrl string `json:"repositoryUrl"`
	LiveUrl       string `json:"liveUrl"`
	ImageUrl      string `json:"imageUrl"`
	Technologies  string `json:"technologies"`
	StartedAt     string `json:"startedAt"`
	LaunchedAt    string `json:"launchedAt"`
	PublishedAt   string `json:"publishedAt"`
	IsFeatured    bool   `json:"isFeatured"`
}

func (p Projects) Create(etx *echo.Context) error {
	var payload CreateProjectFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateProjectFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	validationErrors := map[string]string{}

	slugSource := payload.Slug
	if slugSource == "" {
		slugSource = payload.Name
	}

	var startedAt time.Time
	if payload.StartedAt != "" {
		t, err := time.Parse("2006-01-02", payload.StartedAt)
		if err != nil {
			validationErrors["startedAt"] = "must be a valid date"
		} else {
			startedAt = t
		}
	}

	var launchedAt time.Time
	if payload.LaunchedAt != "" {
		t, err := time.Parse("2006-01-02", payload.LaunchedAt)
		if err != nil {
			validationErrors["launchedAt"] = "must be a valid date"
		} else {
			launchedAt = t
		}
	}

	var publishedAt time.Time
	if payload.PublishedAt != "" {
		t, err := time.Parse("2006-01-02", payload.PublishedAt)
		if err != nil {
			validationErrors["publishedAt"] = "must be a valid date"
		} else {
			publishedAt = t
		}
	}

	if len(validationErrors) > 0 {
		return inertia.Page(
			etx,
			"Admin/Project/Create",
			inertia.Props{},
			inertia.WithValidationErrors(validationErrors),
		)
	}

	data := models.CreateProjectData{
		Name:          payload.Name,
		Slug:          slug.Make(slugSource),
		Tagline:       payload.Tagline,
		Description:   payload.Description,
		ProjectType:   payload.ProjectType,
		RepositoryUrl: payload.RepositoryUrl,
		LiveUrl:       payload.LiveUrl,
		ImageUrl:      payload.ImageUrl,
		Technologies: strings.FieldsFunc(
			payload.Technologies,
			func(r rune) bool { return r == ',' },
		),
		StartedAt:   startedAt,
		LaunchedAt:  launchedAt,
		PublishedAt: publishedAt,
		IsFeatured:  payload.IsFeatured,
	}

	project, err := p.svc.CreateProject(
		etx.Request().Context(),
		data,
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to create project",
			"error",
			err,
			"name",
			payload.Name,
		)
		if errs, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/Project/Create",
				inertia.Props{},
				inertia.WithValidationErrors(errs.ToMap()),
			)
		}

		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}

		if errors.Is(err, services.ErrFeaturedProjectLimit) {
			return inertia.Page(
				etx,
				"Admin/Project/Create",
				inertia.Props{},
				inertia.WithValidationErrors(map[string]string{
					"isFeatured": "No more than 3 projects can be featured.",
				}),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create project: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, routes.AdminProjectNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Project created successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminProjectShow.URL(project.ID))
}

func (p Projects) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectID := int32(parsed)

	project, err := models.Project.Find(etx.Request().Context(), p.db.Executor(), projectID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Project/Edit", inertia.Props{
		"item": newProjectData(project),
	})
}

type UpdateProjectFormPayload struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Tagline       string `json:"tagline"`
	Description   string `json:"description"`
	ProjectType   string `json:"projectType"`
	RepositoryUrl string `json:"repositoryUrl"`
	LiveUrl       string `json:"liveUrl"`
	ImageUrl      string `json:"imageUrl"`
	Technologies  string `json:"technologies"`
	StartedAt     string `json:"startedAt"`
	LaunchedAt    string `json:"launchedAt"`
	PublishedAt   string `json:"publishedAt"`
	IsFeatured    bool   `json:"isFeatured"`
}

func projectDataFromUpdatePayload(id int32, payload UpdateProjectFormPayload) ProjectData {
	return ProjectData{
		ID:            id,
		Name:          payload.Name,
		Slug:          payload.Slug,
		Tagline:       payload.Tagline,
		Description:   payload.Description,
		ProjectType:   payload.ProjectType,
		RepositoryUrl: payload.RepositoryUrl,
		LiveUrl:       payload.LiveUrl,
		ImageUrl:      payload.ImageUrl,
		Technologies:  payload.Technologies,
		StartedAt:     payload.StartedAt,
		LaunchedAt:    payload.LaunchedAt,
		PublishedAt:   payload.PublishedAt,
		IsFeatured:    payload.IsFeatured,
	}
}

func (p Projects) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectID := int32(parsed)

	var payload UpdateProjectFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateProjectFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	validationErrors := map[string]string{}

	slugSource := payload.Slug
	if slugSource == "" {
		slugSource = payload.Name
	}

	var startedAt time.Time
	if payload.StartedAt != "" {
		t, err := time.Parse("2006-01-02", payload.StartedAt)
		if err != nil {
			validationErrors["startedAt"] = "must be a valid date"
		} else {
			startedAt = t
		}
	}

	var launchedAt time.Time
	if payload.LaunchedAt != "" {
		t, err := time.Parse("2006-01-02", payload.LaunchedAt)
		if err != nil {
			validationErrors["launchedAt"] = "must be a valid date"
		} else {
			launchedAt = t
		}
	}

	var publishedAt time.Time
	if payload.PublishedAt != "" {
		t, err := time.Parse("2006-01-02", payload.PublishedAt)
		if err != nil {
			validationErrors["publishedAt"] = "must be a valid date"
		} else {
			publishedAt = t
		}
	}

	if len(validationErrors) > 0 {
		return inertia.Page(
			etx,
			"Admin/Project/Edit",
			inertia.Props{"item": projectDataFromUpdatePayload(projectID, payload)},
			inertia.WithValidationErrors(validationErrors),
		)
	}

	data := models.UpdateProjectData{
		ID:            projectID,
		Name:          payload.Name,
		Slug:          slug.Make(slugSource),
		Tagline:       payload.Tagline,
		Description:   payload.Description,
		ProjectType:   payload.ProjectType,
		RepositoryUrl: payload.RepositoryUrl,
		LiveUrl:       payload.LiveUrl,
		ImageUrl:      payload.ImageUrl,
		Technologies: strings.FieldsFunc(
			payload.Technologies,
			func(r rune) bool { return r == ',' },
		),
		StartedAt:   startedAt,
		LaunchedAt:  launchedAt,
		PublishedAt: publishedAt,
		IsFeatured:  payload.IsFeatured,
	}

	project, err := p.svc.UpdateProject(
		etx.Request().Context(),
		data,
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to update project",
			"id",
			projectID,
			"error",
			err,
		)
		if errs, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/Project/Edit",
				inertia.Props{"item": projectDataFromUpdatePayload(projectID, payload)},
				inertia.WithValidationErrors(errs.ToMap()),
			)
		}

		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}

		if errors.Is(err, services.ErrFeaturedProjectLimit) {
			return inertia.Page(
				etx,
				"Admin/Project/Edit",
				inertia.Props{"item": projectDataFromUpdatePayload(projectID, payload)},
				inertia.WithValidationErrors(map[string]string{
					"isFeatured": "No more than 3 projects can be featured.",
				}),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update project: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminProjectEdit.URL(projectID),
		)
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Project updated successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminProjectShow.URL(project.ID))
}

func (p Projects) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	projectID := int32(parsed)

	err = models.Project.Destroy(etx.Request().Context(), p.db.Executor(), projectID)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete project: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminProjectIndex.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Project destroyed successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminProjectIndex.URL())
}

func (p Projects) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminProjectIndex.Path(),
		Name:        routes.AdminProjectIndex.Name(),
		Handler:     p.Index,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminProjectShow.Path(),
		Name:        routes.AdminProjectShow.Name(),
		Handler:     p.Show,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminProjectNew.Path(),
		Name:        routes.AdminProjectNew.Name(),
		Handler:     p.New,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodPost,
		Path:        routes.AdminProjectCreate.Path(),
		Name:        routes.AdminProjectCreate.Name(),
		Handler:     p.Create,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminProjectEdit.Path(),
		Name:        routes.AdminProjectEdit.Name(),
		Handler:     p.Edit,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodPut,
		Path:        routes.AdminProjectUpdate.Path(),
		Name:        routes.AdminProjectUpdate.Name(),
		Handler:     p.Update,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodDelete,
		Path:        routes.AdminProjectDestroy.Path(),
		Name:        routes.AdminProjectDestroy.Name(),
		Handler:     p.Destroy,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
