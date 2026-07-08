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
	"time"

	"github.com/labstack/echo/v5"
)

type DiaryEntries struct {
	db storage.Pool
}

func NewDiaryEntries(db storage.Pool) DiaryEntries {
	return DiaryEntries{db}
}

type DiaryEntryData struct {
	ID              int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
	EntryDate       string
	MorningThoughts string
	EveningThoughts string
}

func newDiaryEntryData(entity models.DiaryEntryEntity) DiaryEntryData {
	return DiaryEntryData{
		ID:              entity.ID,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
		EntryDate:       entity.EntryDate.Format("2006-01-02"),
		MorningThoughts: entity.MorningThoughts.String,
		EveningThoughts: entity.EveningThoughts.String,
	}
}

func newDiaryEntryDataList(entities []models.DiaryEntryEntity) []DiaryEntryData {
	items := make([]DiaryEntryData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newDiaryEntryData(entity))
	}

	return items
}

func (de DiaryEntries) Index(etx *echo.Context) error {
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

	diaryEntriesList, err := models.DiaryEntry.Paginate(
		etx.Request().Context(),
		de.db.Executor(),
		page,
		perPage,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/DiaryEntry/Index", inertia.Props{
		"items": newDiaryEntryDataList(diaryEntriesList.DiaryEntries),
	})
}

func (de DiaryEntries) Show(etx *echo.Context) error {
	diaryEntryID, err := parseDiaryEntryID(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	diaryEntry, err := models.DiaryEntry.Find(etx.Request().Context(), de.db.Executor(), diaryEntryID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/DiaryEntry/Show", inertia.Props{
		"item": newDiaryEntryData(diaryEntry),
	})
}

func (de DiaryEntries) New(etx *echo.Context) error {
	return inertia.Page(etx, "Admin/DiaryEntry/Create", inertia.Props{
		"today": todayInMadrid().Format("2006-01-02"),
	})
}

func (de DiaryEntries) Today(etx *echo.Context) error {
	focus := etx.QueryParam("focus")
	if focus != "morning" && focus != "evening" {
		focus = ""
	}

	diaryEntry, err := models.DiaryEntry.FindOrCreateForDate(
		etx.Request().Context(),
		de.db.Executor(),
		todayInMadrid(),
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/DiaryEntry/Edit", inertia.Props{
		"item":  newDiaryEntryData(diaryEntry),
		"focus": focus,
	})
}

type CreateDiaryEntryFormPayload struct {
	EntryDate       string `json:"entryDate"`
	MorningThoughts string `json:"morningThoughts"`
	EveningThoughts string `json:"eveningThoughts"`
}

func (de DiaryEntries) Create(etx *echo.Context) error {
	var payload CreateDiaryEntryFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse CreateDiaryEntryFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	data := models.CreateDiaryEntryData{
		EntryDate:       parseDiaryDate(payload.EntryDate),
		MorningThoughts: payload.MorningThoughts,
		EveningThoughts: payload.EveningThoughts,
	}

	diaryEntry, err := models.DiaryEntry.Create(
		etx.Request().Context(),
		de.db.Executor(),
		data,
	)
	if err != nil {
		if errs, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/DiaryEntry/Create",
				inertia.Props{
					"item":  diaryEntryDataFromCreatePayload(payload),
					"today": todayInMadrid().Format("2006-01-02"),
				},
				inertia.WithValidationErrors(errs.ToMap()),
			)
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to create diary entry: %v", err),
		); flashErr != nil {
			return flashErr
		}
		return inertia.Redirect(etx, routes.AdminDiaryEntryNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Diary entry created successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminDiaryEntryShow.URL(diaryEntry.ID))
}

func (de DiaryEntries) Edit(etx *echo.Context) error {
	diaryEntryID, err := parseDiaryEntryID(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	diaryEntry, err := models.DiaryEntry.Find(etx.Request().Context(), de.db.Executor(), diaryEntryID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/DiaryEntry/Edit", inertia.Props{
		"item":  newDiaryEntryData(diaryEntry),
		"focus": etx.QueryParam("focus"),
	})
}

type UpdateDiaryEntryFormPayload struct {
	EntryDate       string `json:"entryDate"`
	MorningThoughts string `json:"morningThoughts"`
	EveningThoughts string `json:"eveningThoughts"`
}

func diaryEntryDataFromCreatePayload(payload CreateDiaryEntryFormPayload) DiaryEntryData {
	return DiaryEntryData{
		EntryDate:       payload.EntryDate,
		MorningThoughts: payload.MorningThoughts,
		EveningThoughts: payload.EveningThoughts,
	}
}

func diaryEntryDataFromUpdatePayload(id int32, payload UpdateDiaryEntryFormPayload) DiaryEntryData {
	return DiaryEntryData{
		ID:              id,
		EntryDate:       payload.EntryDate,
		MorningThoughts: payload.MorningThoughts,
		EveningThoughts: payload.EveningThoughts,
	}
}

func (de DiaryEntries) Update(etx *echo.Context) error {
	diaryEntryID, err := parseDiaryEntryID(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	var payload UpdateDiaryEntryFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse UpdateDiaryEntryFormPayload",
			"error",
			err,
		)

		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	data := models.UpdateDiaryEntryData{
		ID:              diaryEntryID,
		EntryDate:       parseDiaryDate(payload.EntryDate),
		MorningThoughts: payload.MorningThoughts,
		EveningThoughts: payload.EveningThoughts,
	}

	diaryEntry, err := models.DiaryEntry.Update(
		etx.Request().Context(),
		de.db.Executor(),
		data,
	)
	if err != nil {
		if errs, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Admin/DiaryEntry/Edit",
				inertia.Props{
					"item":  diaryEntryDataFromUpdatePayload(diaryEntryID, payload),
					"focus": etx.QueryParam("focus"),
				},
				inertia.WithValidationErrors(errs.ToMap()),
			)
		}

		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}

		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to update diary entry: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminDiaryEntryEdit.URL(diaryEntryID))
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Diary entry updated successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminDiaryEntryShow.URL(diaryEntry.ID))
}

func (de DiaryEntries) Destroy(etx *echo.Context) error {
	diaryEntryID, err := parseDiaryEntryID(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	err = models.DiaryEntry.Destroy(etx.Request().Context(), de.db.Executor(), diaryEntryID)
	if err != nil {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			fmt.Sprintf("Failed to delete diary entry: %v", err),
		); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(etx, routes.AdminDiaryEntryIndex.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Diary entry deleted successfully",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminDiaryEntryIndex.URL())
}

func (de DiaryEntries) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminDiaryEntryIndex.Path(),
		Name:        routes.AdminDiaryEntryIndex.Name(),
		Handler:     de.Index,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminDiaryEntryToday.Path(),
		Name:        routes.AdminDiaryEntryToday.Name(),
		Handler:     de.Today,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminDiaryEntryShow.Path(),
		Name:        routes.AdminDiaryEntryShow.Name(),
		Handler:     de.Show,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminDiaryEntryNew.Path(),
		Name:        routes.AdminDiaryEntryNew.Name(),
		Handler:     de.New,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodPost,
		Path:        routes.AdminDiaryEntryCreate.Path(),
		Name:        routes.AdminDiaryEntryCreate.Name(),
		Handler:     de.Create,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AdminDiaryEntryEdit.Path(),
		Name:        routes.AdminDiaryEntryEdit.Name(),
		Handler:     de.Edit,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodPut,
		Path:        routes.AdminDiaryEntryUpdate.Path(),
		Name:        routes.AdminDiaryEntryUpdate.Name(),
		Handler:     de.Update,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodDelete,
		Path:        routes.AdminDiaryEntryDestroy.Path(),
		Name:        routes.AdminDiaryEntryDestroy.Name(),
		Handler:     de.Destroy,
		Middlewares: authOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func parseDiaryEntryID(etx *echo.Context) (int32, error) {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(parsed), nil
}

func parseDiaryDate(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func todayInMadrid() time.Time {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		return time.Now()
	}
	return time.Now().In(location)
}
