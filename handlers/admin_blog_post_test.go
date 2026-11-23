package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"mbvlabs/domains"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type adminBlogStub struct{ entry domains.BlogPostEntity }

type adminBlogCreatorStub struct {
	data domains.CreateBlogPostData
	at   *time.Time
}

func (s *adminBlogCreatorStub) Create(_ context.Context, data domains.CreateBlogPostData, at *time.Time) (domains.BlogPostEntity, error) {
	s.data, s.at = data, at
	return domains.BlogPostEntity{}, validation.Field("body", "required", "is required")
}

func (s *adminBlogStub) Find(context.Context, int32) (domains.BlogPostEntity, error) {
	return s.entry, nil
}
func (s *adminBlogStub) Paginate(context.Context, int64, int64) (domains.PaginatedBlogPosts, error) {
	return domains.PaginatedBlogPosts{BlogPosts: []domains.BlogPostEntity{s.entry}}, nil
}
func (s *adminBlogStub) Create(context.Context, domains.CreateBlogPostData) (domains.BlogPostEntity, error) {
	return s.entry, nil
}
func (s *adminBlogStub) Update(context.Context, domains.UpdateBlogPostData) (domains.BlogPostEntity, error) {
	return s.entry, nil
}
func (s *adminBlogStub) Destroy(context.Context, int32) error { return nil }

func adminTestRenderer(t *testing.T) *inertia.Renderer {
	t.Helper()
	renderer, err := inertia.NewRenderer(inertia.RendererConfig{ContainerID: "app", BuildPathURL: "/assets/*", EntryPoint: "app.ts", ViteDevURL: "http://localhost:5173", SSRURL: "http://localhost:13714", SSRTimeout: time.Second, SSRMaxResponseBytes: 1 << 20}, inertia.WithRootFunc(func(io.Writer, inertia.RootData) error { return nil }))
	if err != nil {
		t.Fatal(err)
	}
	return renderer
}

func TestAdminBlogPostShowAndScheduleValidation(t *testing.T) {
	entry := domains.BlogPostEntity{ID: 7, Title: "Scheduled post", Tags: json.RawMessage(`["one","two"]`), PublicationSchedule: json.RawMessage(`{"job_id":13,"scheduled_at":"2030-01-02T12:00:00Z"}`)}
	rt := anduin.NewRouter()
	var authCalls int
	auth := func(next anduin.Handler) anduin.Handler {
		return func(ctx context.Context) error { authCalls++; return next(ctx) }
	}
	(AdminBlogPostHandler{posts: &adminBlogStub{entry: entry}, renderer: adminTestRenderer(t)}).Register(rt, auth)
	req := httptest.NewRequest(http.MethodGet, "/admin/blog-posts/7", nil)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	rt.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var page struct {
		Component string `json:"component"`
		Props     struct {
			Item struct {
				Tags                string
				PublicationSchedule struct {
					JobID int64 `json:"job_id"`
				}
			} `json:"item"`
		} `json:"props"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if page.Component != "Admin/BlogPost/Show" || page.Props.Item.Tags != "one, two" || page.Props.Item.PublicationSchedule.JobID != 13 {
		t.Fatalf("unexpected page: %+v", page)
	}

	req = httptest.NewRequest(http.MethodPost, "/admin/blog-posts", strings.NewReader(`{"status":"scheduled","scheduledAt":"invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Inertia", "true")
	w = httptest.NewRecorder()
	rt.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "scheduledAt") {
		t.Fatalf("schedule validation status %d: %s", w.Code, w.Body.String())
	}
	if authCalls != 2 {
		t.Fatalf("middleware calls = %d", authCalls)
	}
}

func TestAdminBlogPostPassesScheduleToCreator(t *testing.T) {
	creator := &adminBlogCreatorStub{}
	rt := anduin.NewRouter()
	(AdminBlogPostHandler{posts: &adminBlogStub{}, creator: creator, renderer: adminTestRenderer(t)}).Register(rt)
	req := httptest.NewRequest(http.MethodPost, "/admin/blog-posts", strings.NewReader(`{"title":"My Article","status":"scheduled","scheduledAt":"2099-01-02T12:00:00Z","tags":"one,two"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	rt.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "body") {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	if creator.at == nil || creator.at.Format(time.RFC3339) != "2099-01-02T12:00:00Z" || creator.data.Slug != "my-article" || len(creator.data.Tags) != 2 {
		t.Fatalf("creator: %+v", creator)
	}
}
