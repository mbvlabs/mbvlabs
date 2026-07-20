package api

import (
	"errors"
	"log/slog"
	"mbvlabs/config"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/middleware"
	"mbvlabs/router/routes"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type DiaryEntries struct {
	db  storage.Pool
	cfg config.Config
}

func NewDiaryEntries(db storage.Pool, cfg config.Config) DiaryEntries {
	return DiaryEntries{db: db, cfg: cfg}
}

func (de DiaryEntries) RegisterRoutes(r *router.Router) error {
	errs := []error{}

	for _, route := range []echo.Route{
		{
			Method:  http.MethodGet,
			Path:    routes.ApiDiaryThoughtsCurrentWeek.Path(),
			Name:    routes.ApiDiaryThoughtsCurrentWeek.Name(),
			Handler: de.CurrentWeekThoughts,
		},
		{
			Method:  http.MethodGet,
			Path:    routes.ApiDiaryThoughtsPreviousWeek.Path(),
			Name:    routes.ApiDiaryThoughtsPreviousWeek.Name(),
			Handler: de.PreviousWeekThoughts,
		},
	} {
		route.Middlewares = []echo.MiddlewareFunc{middleware.APIBasicAuth(de.cfg)}
		if _, err := r.AddRoute(route); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

type CurrentWeekThoughtsResponse struct {
	WeekStart string                     `json:"weekStart"`
	WeekEnd   string                     `json:"weekEnd"`
	Items     []CurrentWeekThoughtsEntry `json:"items"`
}

type CurrentWeekThoughtsEntry struct {
	EntryDate       string `json:"entryDate"`
	MorningThoughts string `json:"morningThoughts"`
	EveningThoughts string `json:"eveningThoughts"`
}

func (de DiaryEntries) CurrentWeekThoughts(etx *echo.Context) error {
	weekStart, weekEnd := currentMadridWeek()
	return de.weekThoughts(etx, weekStart, weekEnd, "current")
}

func (de DiaryEntries) PreviousWeekThoughts(etx *echo.Context) error {
	currentWeekStart, _ := currentMadridWeek()
	weekStart := currentWeekStart.AddDate(0, 0, -7)
	return de.weekThoughts(etx, weekStart, weekStart.AddDate(0, 0, 6), "previous")
}

func (de DiaryEntries) weekThoughts(etx *echo.Context, weekStart, weekEnd time.Time, week string) error {
	entries, err := models.DiaryEntry.FindBetweenDates(
		etx.Request().Context(),
		de.db.Executor(),
		weekStart,
		weekEnd,
	)
	if err != nil {
		message := "failed to fetch " + week + " week diary thoughts"
		slog.ErrorContext(etx.Request().Context(), message, "error", err)
		return etx.JSON(http.StatusInternalServerError, errorResponse{Error: message})
	}

	return etx.JSON(http.StatusOK, CurrentWeekThoughtsResponse{
		WeekStart: weekStart.Format("2006-01-02"),
		WeekEnd:   weekEnd.Format("2006-01-02"),
		Items:     currentWeekThoughtEntries(entries),
	})
}

func currentWeekThoughtEntries(entries []models.DiaryEntryEntity) []CurrentWeekThoughtsEntry {
	items := make([]CurrentWeekThoughtsEntry, 0, len(entries))
	for _, entry := range entries {
		items = append(items, CurrentWeekThoughtsEntry{
			EntryDate:       entry.EntryDate.Format("2006-01-02"),
			MorningThoughts: entry.MorningThoughts.String,
			EveningThoughts: entry.EveningThoughts.String,
		})
	}

	return items
}

func currentMadridWeek() (time.Time, time.Time) {
	now := time.Now().In(madridLocation())
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	year, month, day := now.AddDate(0, 0, -(weekday - 1)).Date()
	weekStart := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	return weekStart, weekStart.AddDate(0, 0, 6)
}

func madridLocation() *time.Location {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		return time.Local
	}
	return location
}
