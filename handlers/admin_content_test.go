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
	"mbvlabs/services"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
)

type workAdminStub struct {
	entity     domains.WorkEntity
	create     domains.CreateWorkData
	update     domains.UpdateWorkData
	deleted    int32
	page, size int64
	err        error
}

func (s *workAdminStub) Paginate(_ context.Context, page, size int64) (domains.PaginatedWorks, error) {
	s.page, s.size = page, size
	return domains.PaginatedWorks{Works: []domains.WorkEntity{s.entity}}, nil
}
func (s *workAdminStub) Find(context.Context, int32) (domains.WorkEntity, error) {
	return s.entity, nil
}
func (s *workAdminStub) Destroy(_ context.Context, id int32) error { s.deleted = id; return nil }
func (s *workAdminStub) CreateWork(_ context.Context, d domains.CreateWorkData) (domains.WorkEntity, error) {
	s.create = d
	return s.entity, s.err
}
func (s *workAdminStub) UpdateWork(_ context.Context, d domains.UpdateWorkData) (domains.WorkEntity, error) {
	s.update = d
	return s.entity, s.err
}

type projectAdminStub struct {
	entity  domains.ProjectEntity
	create  domains.CreateProjectData
	update  domains.UpdateProjectData
	deleted int32
	err     error
}

func (s *projectAdminStub) Paginate(context.Context, int64, int64) (domains.PaginatedProjects, error) {
	return domains.PaginatedProjects{Projects: []domains.ProjectEntity{s.entity}}, nil
}
func (s *projectAdminStub) Find(context.Context, int32) (domains.ProjectEntity, error) {
	return s.entity, nil
}
func (s *projectAdminStub) Destroy(_ context.Context, id int32) error { s.deleted = id; return nil }
func (s *projectAdminStub) CreateProject(_ context.Context, d domains.CreateProjectData) (domains.ProjectEntity, error) {
	s.create = d
	return s.entity, s.err
}
func (s *projectAdminStub) UpdateProject(_ context.Context, d domains.UpdateProjectData) (domains.ProjectEntity, error) {
	s.update = d
	return s.entity, s.err
}

func adminContentTestRouter(t *testing.T, work *workAdminStub, project *projectAdminStub, invalidations *int, authCalls *int) http.Handler {
	t.Helper()
	renderer, err := inertia.NewRenderer(inertia.RendererConfig{
		ContainerID: "app", BuildPathURL: "/assets/*", EntryPoint: "app.ts",
		ViteDevURL: "http://localhost:5173", SSRURL: "http://localhost:13714",
		SSRTimeout: time.Second, SSRMaxResponseBytes: 1 << 20,
	}, inertia.WithRootFunc(func(io.Writer, inertia.RootData) error { return nil }), inertia.WithVersion("test"))
	if err != nil {
		t.Fatal(err)
	}
	keys := kiks.Keys{Authentication: []byte(strings.Repeat("a", 32)), Encryption: []byte(strings.Repeat("b", 32))}
	store, err := kiks.NewCookieStore(keys)
	if err != nil {
		t.Fatal(err)
	}
	jar, err := kiks.NewJar(keys, store, kiks.NewSession[*struct{}]("test_session"))
	if err != nil {
		t.Fatal(err)
	}
	rt := anduin.NewRouter()
	rt.Use(jar.Middleware(), renderer.Middleware())
	auth := func(next anduin.Handler) anduin.Handler {
		return func(ctx context.Context) error { *authCalls++; return next(ctx) }
	}
	invalidate := func() { *invalidations++ }
	(AdminWorkHandler{work, work, renderer, invalidate}).Register(rt, auth)
	(AdminProjectHandler{project, project, renderer, invalidate}).Register(rt, auth)
	return rt.Handler()
}

func adminContentRequest(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("X-Inertia", "true")
	req.Header.Set("X-Inertia-Version", "test")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestAdminWorkAndProjectRoutes(t *testing.T) {
	work := &workAdminStub{entity: domains.WorkEntity{ID: 7, Title: "A Work", ClientURL: domains.OptionalString{String: "https://example.com", Valid: true}, Status: domains.Draft}}
	project := &projectAdminStub{entity: domains.ProjectEntity{ID: 8, Name: "Project", Technologies: json.RawMessage(`["Go","TypeScript"]`)}}
	var invalidations, authCalls int
	h := adminContentTestRouter(t, work, project, &invalidations, &authCalls)
	for _, tc := range []struct{ path, component, expected string }{
		{"/admin/works?page=2&per_page=10", "Admin/Work/Index", `"ClientUrl":"https://example.com"`},
		{"/admin/works/7", "Admin/Work/Show", `"Title":"A Work"`},
		{"/admin/works/new", "Admin/Work/Create", ""},
		{"/admin/works/7/edit", "Admin/Work/Edit", `"Title":"A Work"`},
		{"/admin/projects", "Admin/Project/Index", `"Technologies":"Go, TypeScript"`},
		{"/admin/projects/8", "Admin/Project/Show", `"Name":"Project"`},
		{"/admin/projects/new", "Admin/Project/Create", ""},
		{"/admin/projects/8/edit", "Admin/Project/Edit", `"Name":"Project"`},
	} {
		w := adminContentRequest(t, h, "GET", tc.path, "")
		if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"component":"`+tc.component+`"`) || !strings.Contains(w.Body.String(), tc.expected) {
			t.Fatalf("%s: %d %s", tc.path, w.Code, w.Body.String())
		}
	}
	if work.page != 2 || work.size != 10 {
		t.Fatalf("pagination = %d/%d", work.page, work.size)
	}
	for _, tc := range []struct{ method, path, body, location string }{
		{"POST", "/admin/works", `{"title":"A New Work","status":"draft","startedAt":"2026-01-02"}`, "/admin/works/7"},
		{"PUT", "/admin/works/7", `{"title":"Updated Work","slug":"Explicit Slug","status":"draft"}`, "/admin/works/7"},
		{"DELETE", "/admin/works/7", "", "/admin/works"},
		{"POST", "/admin/projects", `{"name":"New Project","technologies":"Go, TypeScript","publishedAt":"2026-02-03"}`, "/admin/projects/8"},
		{"PUT", "/admin/projects/8", `{"name":"Updated Project","slug":"Custom Name"}`, "/admin/projects/8"},
		{"DELETE", "/admin/projects/8", "", "/admin/projects"},
	} {
		w := adminContentRequest(t, h, tc.method, tc.path, tc.body)
		if w.Code != http.StatusSeeOther || w.Header().Get("Location") != tc.location {
			t.Fatalf("%s %s: %d %s", tc.method, tc.path, w.Code, w.Body.String())
		}
	}
	if work.create.Slug != "a-new-work" || work.create.StartedAt.Format("2006-01-02") != "2026-01-02" || work.update.ID != 7 || work.update.Slug != "explicit-slug" || work.deleted != 7 {
		t.Fatalf("work writes: %#v %#v %d", work.create, work.update, work.deleted)
	}
	if project.create.Slug != "new-project" || len(project.create.Technologies) != 2 || project.create.Technologies[1] != " TypeScript" || project.update.ID != 8 || project.update.Slug != "custom-name" || project.deleted != 8 {
		t.Fatalf("project writes: %#v %#v %d", project.create, project.update, project.deleted)
	}
	if invalidations != 6 || authCalls != 14 {
		t.Fatalf("invalidations=%d auth=%d", invalidations, authCalls)
	}
}

func TestAdminContentValidationPreservesEditProps(t *testing.T) {
	work := &workAdminStub{entity: domains.WorkEntity{ID: 7}, err: services.ErrFeaturedWorkLimit}
	project := &projectAdminStub{entity: domains.ProjectEntity{ID: 8}, err: services.ErrFeaturedProjectLimit}
	var invalidations, authCalls int
	h := adminContentTestRouter(t, work, project, &invalidations, &authCalls)
	for _, tc := range []struct{ path, body, want string }{
		{"/admin/works/7", `{"title":"Keep work","startedAt":"bad"}`, `"startedAt":"must be a valid date"`},
		{"/admin/projects/8", `{"name":"Keep project","publishedAt":"bad"}`, `"publishedAt":"must be a valid date"`},
		{"/admin/works/7", `{"title":"Keep work","isFeatured":true}`, `"isFeatured":"No more than 3 works can be featured."`},
		{"/admin/projects/8", `{"name":"Keep project","isFeatured":true}`, `"isFeatured":"No more than 3 projects can be featured."`},
	} {
		w := adminContentRequest(t, h, "PUT", tc.path, tc.body)
		if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), tc.want) || !strings.Contains(w.Body.String(), `"Keep `) {
			t.Fatalf("%s: %d %s", tc.path, w.Code, w.Body.String())
		}
	}
	if invalidations != 0 {
		t.Fatalf("invalidated after validation failure")
	}
}
