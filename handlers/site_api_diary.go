package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"mbvlabs/domains"
	"mbvlabs/middleware"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
)

type apiDiaryStore interface {
	FindBetweenDates(context.Context, time.Time, time.Time) ([]domains.DiaryEntryEntity, error)
}

type APIDiaryHandler struct {
	entries            apiDiaryStore
	username, password string
}

func NewAPIDiaryHandler(entries apiDiaryStore, username, password string) APIDiaryHandler {
	return APIDiaryHandler{entries: entries, username: username, password: password}
}

func (h APIDiaryHandler) Register(rt *anduin.Router) {
	auth := middleware.APIBasicAuth(h.username, h.password)
	rt.Route(routes.ApiDiaryThoughtsCurrentWeek).Get(h.week(false), auth)
	rt.Route(routes.ApiDiaryThoughtsPreviousWeek).Get(h.week(true), auth)
}

type apiThought struct {
	EntryDate       string `json:"entryDate"`
	MorningThoughts string `json:"morningThoughts"`
	EveningThoughts string `json:"eveningThoughts"`
}

func (h APIDiaryHandler) week(previous bool) anduin.Handler {
	return func(ctx context.Context) error {
		location, err := time.LoadLocation("Europe/Madrid")
		if err != nil {
			location = time.Local
		}
		now := time.Now().In(location)
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		y, m, d := now.AddDate(0, 0, 1-weekday).Date()
		start := time.Date(y, m, d, 0, 0, 0, 0, location)
		name := "current"
		if previous {
			start = start.AddDate(0, 0, -7)
			name = "previous"
		}
		end := start.AddDate(0, 0, 6)
		entries, err := h.entries.FindBetweenDates(ctx, start, end)
		if err != nil {
			message := "failed to fetch " + name + " week diary thoughts"
			slog.ErrorContext(ctx, message, "error", err)
			return apiFailure(ctx, http.StatusInternalServerError, message)
		}
		items := make([]apiThought, 0, len(entries))
		for _, entry := range entries {
			items = append(items, apiThought{entry.EntryDate.Format("2006-01-02"), entry.MorningThoughts.String, entry.EveningThoughts.String})
		}
		return anduin.JSON(ctx, http.StatusOK, struct {
			WeekStart string       `json:"weekStart"`
			WeekEnd   string       `json:"weekEnd"`
			Items     []apiThought `json:"items"`
		}{start.Format("2006-01-02"), end.Format("2006-01-02"), items})
	}
}
