package api_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mbvlabs/config"
	"mbvlabs/controllers/api"
	"mbvlabs/database"
	"mbvlabs/models/factories"

	"github.com/labstack/echo/v5"
)

func TestDiaryEntries_PreviousWeekThoughts(t *testing.T) {
	ctx := context.Background()
	db := testCluster.NewTestDB(t, database.Migrations, "migrations")
	weekStart, weekEnd := previousMadridWeek()

	create := func(date time.Time, morning, evening string) {
		t.Helper()
		if _, err := factories.CreateDiaryEntry(ctx, db.Executor(),
			factories.WithDiaryEntriesEntryDate(date),
			factories.WithDiaryEntriesMorningThoughts(sql.NullString{String: morning, Valid: morning != ""}),
			factories.WithDiaryEntriesEveningThoughts(sql.NullString{String: evening, Valid: evening != ""}),
		); err != nil {
			t.Fatalf("create diary entry: %v", err)
		}
	}

	create(weekEnd, "Sunday morning", "Sunday evening")
	create(weekStart, "Monday morning", "Monday evening")
	create(weekStart.AddDate(0, 0, -1), "Too early", "")
	create(weekEnd.AddDate(0, 0, 1), "Too late", "")

	controller := api.NewDiaryEntries(db, config.Config{})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/diary/thoughts/previous-week", nil)
	rec := httptest.NewRecorder()
	if err := controller.PreviousWeekThoughts(e.NewContext(req, rec)); err != nil {
		t.Fatalf("get previous week thoughts: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var response api.CurrentWeekThoughtsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.WeekStart != weekStart.Format("2006-01-02") || response.WeekEnd != weekEnd.Format("2006-01-02") {
		t.Fatalf("week = %s to %s, want %s to %s", response.WeekStart, response.WeekEnd, weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))
	}
	if len(response.Items) != 2 {
		t.Fatalf("item count = %d, want 2: %s", len(response.Items), rec.Body.String())
	}
	if response.Items[0].EntryDate != weekStart.Format("2006-01-02") || response.Items[0].MorningThoughts != "Monday morning" || response.Items[0].EveningThoughts != "Monday evening" {
		t.Errorf("first item = %#v", response.Items[0])
	}
	if response.Items[1].EntryDate != weekEnd.Format("2006-01-02") || response.Items[1].MorningThoughts != "Sunday morning" || response.Items[1].EveningThoughts != "Sunday evening" {
		t.Errorf("second item = %#v", response.Items[1])
	}
}

func TestDiaryEntries_PreviousWeekThoughtsEmpty(t *testing.T) {
	db := testCluster.NewTestDB(t, database.Migrations, "migrations")
	controller := api.NewDiaryEntries(db, config.Config{})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/diary/thoughts/previous-week", nil)
	rec := httptest.NewRecorder()

	if err := controller.PreviousWeekThoughts(e.NewContext(req, rec)); err != nil {
		t.Fatalf("get previous week thoughts: %v", err)
	}

	var response api.CurrentWeekThoughtsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Items == nil || len(response.Items) != 0 {
		t.Fatalf("items = %#v, want empty array", response.Items)
	}
}

func previousMadridWeek() (time.Time, time.Time) {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		location = time.Local
	}
	now := time.Now().In(location)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := now.AddDate(0, 0, -(weekday - 1)).Date()
	currentWeekStart := time.Date(year, month, day, 0, 0, 0, 0, location)
	weekStart := currentWeekStart.AddDate(0, 0, -7)
	return weekStart, weekStart.AddDate(0, 0, 6)
}
