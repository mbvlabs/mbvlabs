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

type queueStub struct {
	filters    domains.RiverFilters
	page, size int64
}

func (*queueStub) Find(context.Context, int64) (domains.RiverJob, error) {
	return domains.RiverJob{ID: 3}, nil
}
func (s *queueStub) Jobs(_ context.Context, f domains.RiverFilters, page, size int64) ([]domains.RiverJob, domains.RiverPagination, error) {
	s.filters, s.page, s.size = f, page, size
	return []domains.RiverJob{{ID: 3}}, domains.RiverPagination{Page: page, PageSize: size, TotalCount: 1, TotalPages: 1}, nil
}
func (*queueStub) Stats(context.Context) ([]domains.RiverQueueStats, error) {
	return []domains.RiverQueueStats{}, nil
}
func (*queueStub) StateCounts(context.Context) ([]domains.RiverStateCount, error) {
	return []domains.RiverStateCount{}, nil
}
func (*queueStub) SetPaused(context.Context, string, bool) error { return nil }
func (*queueStub) Action(context.Context, int64, string) error   { return nil }

func TestAdminQueuePages(t *testing.T) {
	renderer, err := inertia.NewRenderer(inertia.RendererConfig{
		ContainerID: "app", BuildPathURL: "/assets/*", EntryPoint: "app.ts",
		ViteDevURL: "http://localhost:5173", SSRURL: "http://localhost:13714",
		SSRTimeout: time.Second, SSRMaxResponseBytes: 1 << 20,
	}, inertia.WithRootFunc(func(io.Writer, inertia.RootData) error { return nil }))
	if err != nil {
		t.Fatal(err)
	}
	store := &queueStub{}
	rt := anduin.NewRouter()
	calls := 0
	auth := func(next anduin.Handler) anduin.Handler {
		return func(ctx context.Context) error { calls++; return next(ctx) }
	}
	(AdminQueueHandler{queues: store, renderer: renderer}).Register(rt, auth)
	for _, tc := range []struct{ path, component string }{
		{"/admin/queue", "Admin/Queue/Index"},
		{"/admin/queue/jobs?state=wrong&queue=default&kind=mail&page=2&per_page=5", "Admin/Queue/Jobs"},
		{"/admin/queue/jobs/3", "Admin/Queue/Show"},
	} {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		req.Header.Set("X-Inertia", "true")
		w := httptest.NewRecorder()
		rt.Handler().ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("%s: %d %s", tc.path, w.Code, w.Body.String())
		}
		var page struct {
			Component string                     `json:"component"`
			Props     map[string]json.RawMessage `json:"props"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
			t.Fatal(err)
		}
		if page.Component != tc.component {
			t.Fatalf("component %q", page.Component)
		}
		if tc.component == "Admin/Queue/Jobs" && (store.page != 2 || store.size != 5 || store.filters.State != "" || store.filters.Queue != "default" || store.filters.Kind != "mail" || len(page.Props["pagination"]) == 0) {
			t.Fatalf("store=%+v props=%v", store, page.Props)
		}
	}
	if calls != 3 {
		t.Fatalf("middleware calls=%d", calls)
	}
}
