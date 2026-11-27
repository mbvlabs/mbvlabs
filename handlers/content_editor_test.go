package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mbvlabs/domains"
	"mbvlabs/handlers"
	"mbvlabs/migrations"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/storage"
)

func TestContentEditorRoutes(t *testing.T) {
	db := testCluster.NewTestDB(t, migrations.Migrations, ".")
	ctx := t.Context()
	posts, projects, works := domains.NewBlogPosts(db), domains.NewProjects(db), domains.NewWorks(db)
	if _, err := posts.Create(ctx, domains.CreateBlogPostData{Title: "Article original title", Slug: "article", Body: "Original article body", Status: "draft", Tags: []string{"go"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := projects.Create(ctx, domains.CreateProjectData{Name: "Project original", Slug: "project", Tagline: "Original project tagline", ProjectType: "web-app", Technologies: []string{"Go"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := works.Create(ctx, domains.CreateWorkData{Title: "Original work title", Slug: "work", Summary: "Original work summary", Status: domains.Draft, Specialisms: []string{}, Platforms: []string{}, Technologies: []string{}}); err != nil {
		t.Fatal(err)
	}

	bearer := func(next anduin.Handler) anduin.Handler {
		return func(ctx context.Context) error {
			if anduin.Request(ctx).Header.Get("Authorization") != "Bearer test" {
				return anduin.JSON(ctx, 401, nil)
			}
			return next(ctx)
		}
	}
	rt := anduin.NewRouter()
	handlers.NewContentEditorHandler(posts, projects, works, storage.NewTransactor(db)).Register(rt, bearer)
	request := func(method, path, body string, authorized bool) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		if authorized {
			req.Header.Set("Authorization", "Bearer test")
		}
		rec := httptest.NewRecorder()
		rt.Handler().ServeHTTP(rec, req)
		return rec
	}
	for _, tc := range []struct{ kind, slug, payload, updated string }{
		{"articles", "article", `{"title":"Updated article title","tags":["new"]}`, "Updated article title"},
		{"projects", "project", `{"name":"Updated project","technologies":["Rust"]}`, "Updated project"},
		{"works", "work", `{"title":"Updated work title","content":"New work content"}`, "Updated work title"},
	} {
		base := "/api/content/" + tc.kind
		path := base + "/" + tc.slug
		for _, target := range []string{base, path} {
			if r := request("GET", target, "", false); r.Code != 401 {
				t.Fatalf("unauthorized GET %s: %d", target, r.Code)
			}
		}
		if r := request("PATCH", path, tc.payload, false); r.Code != 401 {
			t.Fatalf("unauthorized PATCH %s: %d", path, r.Code)
		}
		if r := request("GET", base+"?page=1&per_page=1", "", true); r.Code != 200 || !strings.Contains(r.Body.String(), tc.slug) || !strings.Contains(r.Body.String(), `"PageSize":1`) {
			t.Fatalf("list: %d %s", r.Code, r.Body)
		}
		if r := request("GET", path, "", true); r.Code != 200 {
			t.Fatalf("get: %d %s", r.Code, r.Body)
		}
		for _, body := range []string{`{"status":"published"}`, `{"slug":"changed"}`, `{"isFeatured":true}`, `{"publishedAt":"2026-01-01"}`, `{"unknown":1}`, `{"title":null}`, `{} {}`, strings.Repeat(" ", 1<<20+1)} {
			if r := request("PATCH", path, body, true); r.Code != 400 {
				t.Fatalf("invalid patch %s: %d %s", body[:min(len(body), 40)], r.Code, r.Body)
			}
		}
		if r := request("PATCH", path, tc.payload, true); r.Code != 200 || !strings.Contains(r.Body.String(), tc.updated) {
			t.Fatalf("patch: %d %s", r.Code, r.Body)
		}
		if r := request("GET", path, "", true); r.Code != 200 || !strings.Contains(r.Body.String(), tc.updated) {
			t.Fatalf("persisted: %d %s", r.Code, r.Body)
		}
		if r := request("GET", base+"/missing", "", true); r.Code != http.StatusNotFound {
			t.Fatalf("missing: %d %s", r.Code, r.Body)
		}
	}
	post, err := posts.FindBySlug(ctx, "article")
	if err != nil || post.Status != "draft" || post.Body != "Original article body" || post.PublishedAt.Valid || string(post.Tags) != `["new"]` {
		t.Fatalf("article preserved: %+v %v", post, err)
	}
	project, err := projects.FindBySlug(ctx, "project")
	if err != nil || project.Tagline != "Original project tagline" || project.PublishedAt.Valid || project.IsFeatured || string(project.Technologies) != `["Rust"]` {
		t.Fatalf("project preserved: %+v %v", project, err)
	}
	work, err := works.FindBySlug(ctx, "work")
	if err != nil || work.Status != domains.Draft || work.Summary != "Original work summary" || work.PublishedAt.Valid || work.IsFeatured || work.Content != "New work content" {
		t.Fatalf("work preserved: %+v %v", work, err)
	}
	if r := request("PATCH", "/api/content/articles/article", `{"title":"x"}`, true); r.Code != 400 {
		t.Fatalf("validation: %d %s", r.Code, r.Body)
	}
}
