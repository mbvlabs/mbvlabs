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
	"time"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v5"
)

type Works struct {
	db      storage.Pool
	workSvc services.Works
}

func NewWorks(db storage.Pool, works services.Works) Works {
	return Works{db: db, workSvc: works}
}

type WorkData struct {
	ID               int32
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Title            string
	Slug             string
	ClientName       string
	ClientIndustry   string
	ClientUrl        string
	ClientLogoUrl    string
	Summary          string
	CoverImageUrl    string
	Specialisms      []string
	Platforms        []string
	Technologies     []string
	Challenge        string
	Approach         string
	Deliverables     string
	Outcome          string
	Content          string
	TestimonialQuote string
	TestimonialName  string
	TestimonialRole  string
	StartedAt        string
	CompletedAt      string
	Status           string
	PublishedAt      string
	IsFeatured       bool
}

func newWorkData(entity models.WorkEntity) WorkData {
	return WorkData{
		ID:             entity.ID,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		Title:          entity.Title,
		Slug:           entity.Slug,
		ClientName:     entity.ClientName.String,
		ClientIndustry: entity.ClientIndustry.String,
		ClientUrl:      entity.ClientURL.String,
		ClientLogoUrl:  entity.ClientLogoURL.String,
		Summary:        entity.Summary,
		CoverImageUrl:  entity.CoverImageURL.String,
		Specialisms:    entity.Specialisms,
		Platforms:      entity.Platforms,
		Technologies:   entity.Technologies,
		Challenge:      entity.Challenge,
		Approach:       entity.Approach,
		Deliverables:   entity.Deliverables,
		Outcome:        entity.Outcome,
		Content:        entity.Content,
		StartedAt: func() string {
			if !entity.StartedAt.Valid {
				return ""
			}
			return entity.StartedAt.Time.Format("2006-01-02")
		}(),
		CompletedAt: func() string {
			if !entity.CompletedAt.Valid {
				return ""
			}
			return entity.CompletedAt.Time.Format("2006-01-02")
		}(),
		Status: entity.Status.String(),
		PublishedAt: func() string {
			if !entity.PublishedAt.Valid {
				return ""
			}
			return entity.PublishedAt.Time.Format("2006-01-02")
		}(),
		IsFeatured: entity.IsFeatured,
	}
}

func newWorkDataList(entities []models.WorkEntity) []WorkData {
	items := make([]WorkData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newWorkData(entity))
	}

	return items
}

func (cs Works) Index(etx *echo.Context) error {
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

	slog.InfoContext(etx.Request().Context(), "listing works", "page", page, "perPage", perPage)

	worksList, err := models.Work.Paginate(
		etx.Request().Context(),
		cs.db.Executor(),
		page,
		perPage,
	)
	if err != nil {
		slog.ErrorContext(etx.Request().Context(), "failed to paginate works", "error", err)
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Work/Index", inertia.Props{
		"items": newWorkDataList(worksList.Works),
	})
}

func (cs Works) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		slog.ErrorContext(etx.Request().Context(), "invalid work id", "error", err)
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	workID := int32(parsed)

	work, err := models.Work.Find(etx.Request().Context(), cs.db.Executor(), workID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to find work",
			"id",
			workID,
			"error",
			err,
		)
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Work/Show", inertia.Props{
		"item": newWorkData(work),
	})
}

func (cs Works) New(etx *echo.Context) error {
	slog.InfoContext(etx.Request().Context(), "showing new work form")
	return inertia.Page(etx, "Admin/Work/Create", inertia.Props{})
}

type CreateWorkFormPayload struct {
	Title            string   `json:"title"`
	Slug             string   `json:"slug"`
	ClientName       string   `json:"clientName"`
	ClientIndustry   string   `json:"clientIndustry"`
	ClientURL        string   `json:"clientUrl"`
	ClientLogoURL    string   `json:"clientLogoUrl"`
	Summary          string   `json:"summary"`
	CoverImageURL    string   `json:"coverImageUrl"`
	Specialisms      []string `json:"specialisms"`
	Platforms        []string `json:"platforms"`
	Technologies     []string `json:"technologies"`
	Challenge        string   `json:"challenge"`
	Approach         string   `json:"approach"`
	Deliverables     string   `json:"deliverables"`
	Outcome          string   `json:"outcome"`
	Content          string   `json:"content"`
	TestimonialQuote string   `json:"testimonialQuote"`
	TestimonialName  string   `json:"testimonialName"`
	TestimonialRole  string   `json:"testimonialRole"`
	StartedAt        string   `json:"startedAt"`
	CompletedAt      string   `json:"completedAt"`
	Status           string   `json:"status"`
	IsFeatured       bool     `json:"isFeatured"`
}

func (cs Works) Create(etx *echo.Context) error {
	var payload CreateWorkFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateWorkFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	validationErrors := map[string]string{}

	slugSource := payload.Slug
	if slugSource == "" {
		slugSource = payload.Title
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

	var completedAt time.Time
	if payload.CompletedAt != "" {
		t, err := time.Parse("2006-01-02", payload.CompletedAt)
		if err != nil {
			validationErrors["completedAt"] = "must be a valid date"
		} else {
			completedAt = t
		}
	}

	if len(validationErrors) > 0 {
		return inertia.Page(
			etx,
			"Admin/Work/Create",
			inertia.Props{},
			inertia.WithValidationErrors(validationErrors),
		)
	}

	data := models.CreateWorkData{
		Title:            payload.Title,
		Slug:             slug.Make(slugSource),
		ClientName:       payload.ClientName,
		ClientIndustry:   payload.ClientIndustry,
		ClientURL:        payload.ClientURL,
		ClientLogoURL:    payload.ClientLogoURL,
		Summary:          payload.Summary,
		CoverImageURL:    payload.CoverImageURL,
		Specialisms:      payload.Specialisms,
		Platforms:        payload.Platforms,
		Technologies:     payload.Technologies,
		Challenge:        payload.Challenge,
		Approach:         payload.Approach,
		Deliverables:     payload.Deliverables,
		Outcome:          payload.Outcome,
		Content:          payload.Content,
		TestimonialQuote: payload.TestimonialQuote,
		TestimonialName:  payload.TestimonialName,
		TestimonialRole:  payload.TestimonialRole,
		StartedAt:        startedAt,
		CompletedAt:      completedAt,
		Status:           models.StatusEnum(payload.Status),
		IsFeatured:       payload.IsFeatured,
	}

	work, err := cs.workSvc.CreateWork(
		etx.Request().Context(),
		data,
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to create work",
			"error",
			err,
			"title",
			payload.Title,
		)
		if errs, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/Work/Create",
				inertia.Props{},
				inertia.WithValidationErrors(errs.ToMap()),
			)
		}

		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}

		if errors.Is(err, services.ErrFeaturedWorkLimit) {
			return inertia.Page(
				etx,
				"Admin/Work/Create",
				inertia.Props{},
				inertia.WithValidationErrors(map[string]string{
					"isFeatured": "No more than 3 works can be featured.",
				}),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create work: %v", err),
		); flashErr != nil {
			return flashErr
		}

		return inertia.Redirect(etx, routes.AdminWorkNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Work created successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Redirect(etx, routes.AdminWorkShow.URL(work.ID))
}

func (cs Works) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		slog.ErrorContext(etx.Request().Context(), "invalid work id", "error", err)
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	workID := int32(parsed)

	work, err := models.Work.Find(etx.Request().Context(), cs.db.Executor(), workID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to find work for edit",
			"id",
			workID,
			"error",
			err,
		)
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Work/Edit", inertia.Props{
		"item": newWorkData(work),
	})
}

type UpdateWorkFormPayload struct {
	Title            string   `json:"title"`
	Slug             string   `json:"slug"`
	ClientName       string   `json:"clientName"`
	ClientIndustry   string   `json:"clientIndustry"`
	ClientUrl        string   `json:"clientUrl"`
	ClientLogoUrl    string   `json:"clientLogoUrl"`
	Summary          string   `json:"summary"`
	CoverImageUrl    string   `json:"coverImageUrl"`
	Specialisms      []string `json:"specialisms"`
	Platforms        []string `json:"platforms"`
	Technologies     []string `json:"technologies"`
	Challenge        string   `json:"challenge"`
	Approach         string   `json:"approach"`
	Deliverables     string   `json:"deliverables"`
	Outcome          string   `json:"outcome"`
	Content          string   `json:"content"`
	TestimonialQuote string   `json:"testimonialQuote"`
	TestimonialName  string   `json:"testimonialName"`
	TestimonialRole  string   `json:"testimonialRole"`
	StartedAt        string   `json:"startedAt"`
	CompletedAt      string   `json:"completedAt"`
	Status           string   `json:"status"`
	IsFeatured       bool     `json:"isFeatured"`
}

func workDataFromUpdatePayload(id int32, payload UpdateWorkFormPayload) WorkData {
	return WorkData{
		ID:               id,
		Title:            payload.Title,
		Slug:             payload.Slug,
		ClientName:       payload.ClientName,
		ClientIndustry:   payload.ClientIndustry,
		ClientUrl:        payload.ClientUrl,
		ClientLogoUrl:    payload.ClientLogoUrl,
		Summary:          payload.Summary,
		CoverImageUrl:    payload.CoverImageUrl,
		Specialisms:      payload.Specialisms,
		Platforms:        payload.Platforms,
		Technologies:     payload.Technologies,
		Challenge:        payload.Challenge,
		Approach:         payload.Approach,
		Deliverables:     payload.Deliverables,
		Outcome:          payload.Outcome,
		Content:          payload.Content,
		TestimonialQuote: payload.TestimonialQuote,
		TestimonialName:  payload.TestimonialName,
		TestimonialRole:  payload.TestimonialRole,
		StartedAt:        payload.StartedAt,
		CompletedAt:      payload.CompletedAt,
		Status:           payload.Status,
		IsFeatured:       payload.IsFeatured,
	}
}

func (cs Works) Update(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		slog.ErrorContext(etx.Request().Context(), "invalid work id", "error", err)
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	workID := int32(parsed)

	var payload UpdateWorkFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateWorkFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
	}

	slugSource := payload.Slug
	if slugSource == "" {
		slugSource = payload.Title
	}

	validationErrors := map[string]string{}

	var startedAt time.Time
	if payload.StartedAt != "" {
		t, err := time.Parse("2006-01-02", payload.StartedAt)
		if err != nil {
			validationErrors["startedAt"] = "must be a valid date"
		} else {
			startedAt = t
		}
	}

	var completedAt time.Time
	if payload.CompletedAt != "" {
		t, err := time.Parse("2006-01-02", payload.CompletedAt)
		if err != nil {
			validationErrors["completedAt"] = "must be a valid date"
		} else {
			completedAt = t
		}
	}

	if len(validationErrors) > 0 {
		return inertia.Page(
			etx,
			"Admin/Work/Edit",
			inertia.Props{"item": workDataFromUpdatePayload(workID, payload)},
			inertia.WithValidationErrors(validationErrors),
		)
	}

	data := models.UpdateWorkData{
		ID:               workID,
		Title:            payload.Title,
		Slug:             slug.Make(slugSource),
		ClientName:       payload.ClientName,
		ClientIndustry:   payload.ClientIndustry,
		ClientURL:        payload.ClientUrl,
		ClientLogoURL:    payload.ClientLogoUrl,
		Summary:          payload.Summary,
		CoverImageURL:    payload.CoverImageUrl,
		Specialisms:      payload.Specialisms,
		Platforms:        payload.Platforms,
		Technologies:     payload.Technologies,
		Challenge:        payload.Challenge,
		Approach:         payload.Approach,
		Deliverables:     payload.Deliverables,
		Outcome:          payload.Outcome,
		Content:          payload.Content,
		TestimonialQuote: payload.TestimonialQuote,
		TestimonialName:  payload.TestimonialName,
		TestimonialRole:  payload.TestimonialRole,
		StartedAt:        startedAt,
		CompletedAt:      completedAt,
		Status:           models.StatusEnum(payload.Status),
		IsFeatured:       payload.IsFeatured,
	}

	work, err := cs.workSvc.UpdateWork(
		etx.Request().Context(),
		data,
	)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to update work",
			"id",
			workID,
			"error",
			err,
		)
		if errs, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/Work/Edit",
				inertia.Props{"item": workDataFromUpdatePayload(workID, payload)},
				inertia.WithValidationErrors(errs.ToMap()),
			)
		}

		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}

		if errors.Is(err, services.ErrFeaturedWorkLimit) {
			return inertia.Page(
				etx,
				"Admin/Work/Edit",
				inertia.Props{"item": workDataFromUpdatePayload(workID, payload)},
				inertia.WithValidationErrors(map[string]string{
					"isFeatured": "No more than 3 works can be featured.",
				}),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update work: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.AdminWorkEdit.URL(workID),
		)
	}

	slog.InfoContext(etx.Request().Context(), "work updated", "id", work.ID, "title", work.Title)
	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Work updated successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminWorkShow.URL(work.ID))
}

func (cs Works) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		slog.ErrorContext(etx.Request().Context(), "invalid work id", "error", err)
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}
	workID := int32(parsed)

	err = models.Work.Destroy(etx.Request().Context(), cs.db.Executor(), workID)
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to destroy work",
			"id",
			workID,
			"error",
			err,
		)
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete work: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminWorkIndex.URL())
	}

	slog.InfoContext(etx.Request().Context(), "work destroyed", "id", workID)
	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Work destroyed successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminWorkIndex.URL())
}

func (cs Works) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminWorkIndex.Path(),
		Name:    routes.AdminWorkIndex.Name(),
		Handler: cs.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminWorkShow.Path(),
		Name:    routes.AdminWorkShow.Name(),
		Handler: cs.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminWorkNew.Path(),
		Name:    routes.AdminWorkNew.Name(),
		Handler: cs.New,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.AdminWorkCreate.Path(),
		Name:    routes.AdminWorkCreate.Name(),
		Handler: cs.Create,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.AdminWorkEdit.Path(),
		Name:    routes.AdminWorkEdit.Name(),
		Handler: cs.Edit,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPut,
		Path:    routes.AdminWorkUpdate.Path(),
		Name:    routes.AdminWorkUpdate.Name(),
		Handler: cs.Update,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodDelete,
		Path:    routes.AdminWorkDestroy.Path(),
		Name:    routes.AdminWorkDestroy.Name(),
		Handler: cs.Destroy,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
