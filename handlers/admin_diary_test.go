package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mbvlabs/domains"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
)

type diaryStub struct {
	entry      domains.DiaryEntryEntity
	page, size int64
}

func (s *diaryStub) Find(context.Context, int32) (domains.DiaryEntryEntity, error) {
	return s.entry, nil
}
func (s *diaryStub) FindOrCreateForDate(_ context.Context, date time.Time) (domains.DiaryEntryEntity, error) {
	s.entry.EntryDate = date
	return s.entry, nil
}
func (s *diaryStub) Paginate(_ context.Context, page, size int64) (domains.PaginatedDiaryEntries, error) {
	s.page, s.size = page, size
	return domains.PaginatedDiaryEntries{DiaryEntries: []domains.DiaryEntryEntity{s.entry}}, nil
}
func (s *diaryStub) Create(context.Context, domains.CreateDiaryEntryData) (domains.DiaryEntryEntity, error) {
	return s.entry, nil
}
func (s *diaryStub) Update(context.Context, domains.UpdateDiaryEntryData) (domains.DiaryEntryEntity, error) {
	return s.entry, nil
}
func (s *diaryStub) Destroy(context.Context, int32) error { return nil }

func TestAdminDiaryPagesAndMiddleware(t *testing.T) {
	renderer, err := inertia.NewRenderer(inertia.RendererConfig{
		ContainerID: "app", BuildPathURL: "/assets/*", EntryPoint: "app.ts",
		ViteDevURL: "http://localhost:5173", SSRURL: "http://localhost:13714",
		SSRTimeout: time.Second, SSRMaxResponseBytes: 1 << 20,
	}, inertia.WithRootFunc(func(io.Writer, inertia.RootData) error { return nil }))
	if err != nil {
		t.Fatal(err)
	}
	store := &diaryStub{entry: domains.DiaryEntryEntity{ID: 4, EntryDate: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC), MorningThoughts: domains.OptionalString{String: "hello", Valid: true}}}
	rt := anduin.NewRouter()
	var authCalls int
	auth := func(next anduin.Handler) anduin.Handler {
		return func(ctx context.Context) error { authCalls++; return next(ctx) }
	}
	(AdminDiaryHandler{entries: store, renderer: renderer}).Register(rt, auth)
	for _, tc := range []struct{ path, component string }{
		{"/admin/diary?page=2&per_page=10", "Admin/DiaryEntry/Index"},
		{"/admin/diary/today?focus=evening", "Admin/DiaryEntry/Edit"},
	} {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		req.Header.Set("X-Inertia", "true")
		w := httptest.NewRecorder()
		rt.Handler().ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("%s: status %d: %s", tc.path, w.Code, w.Body.String())
		}
		var page struct {
			Component string                     `json:"component"`
			Props     map[string]json.RawMessage `json:"props"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
			t.Fatal(err)
		}
		if page.Component != tc.component {
			t.Fatalf("component = %q, want %q", page.Component, tc.component)
		}
		if page.Component == "Admin/DiaryEntry/Index" && (store.page != 2 || store.size != 10 || len(page.Props["items"]) == 0) {
			t.Fatalf("pagination or items incorrect: %#v", page.Props)
		}
		if page.Component == "Admin/DiaryEntry/Edit" && string(page.Props["focus"]) != `"evening"` {
			t.Fatalf("focus = %s", page.Props["focus"])
		}
	}
	if authCalls != 2 {
		t.Fatalf("auth middleware calls = %d", authCalls)
	}
}
