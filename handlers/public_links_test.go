package handlers_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"mbvlabs/config"
	"mbvlabs/domains"
	"mbvlabs/handlers"
	"mbvlabs/migrations"
	"mbvlabs/resources"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/anduin/anduintest"
	"github.com/mbvlabs/andurel/pkg/inertia"
)

func TestPublicEditorialLinksOnlyIncludePublishedTargets(t *testing.T) {
	db := testCluster.NewTestDB(t, migrations.Migrations, ".")
	works, posts := domains.NewWorks(db), domains.NewBlogPosts(db)
	now := time.Now()
	for _, item := range []domains.CreateWorkData{
		{
			Title: "Kuatro backend work", Slug: "kuatro-group-backend-infrastructure", Status: domains.Published,
			Summary: "Kuatro backend work", ClientName: "Kuatro", ClientIndustry: "Technology",
			ClientURL: "https://example.com", ClientLogoURL: "https://example.com/logo.png", CoverImageURL: "https://example.com/cover.png",
			Specialisms: []string{"Engineering"}, Platforms: []string{"Web"}, Technologies: []string{"Go"},
			Challenge: "Challenge", Approach: "Approach", Deliverables: "Deliverables", Outcome: "Outcome", Content: "Content",
			StartedAt: now, CompletedAt: now, PublishedAt: now,
		},
		{
			Title: "Everphone internal search", Slug: "everphone-internal-tool-search", Status: domains.Draft,
			Summary: "Everphone internal search", Specialisms: []string{}, Platforms: []string{}, Technologies: []string{},
		},
	} {
		if _, err := works.Create(t.Context(), item); err != nil {
			t.Fatal(err)
		}
	}
	for _, item := range []domains.CreateBlogPostData{
		{Title: "Stored rescue guide title", Slug: "software-project-rescue", Excerpt: "Stored rescue guide excerpt", Status: "published"},
		{Title: "Audit process draft", Slug: "codebase-audit", Excerpt: "Draft audit process excerpt", Status: "draft"},
		{Title: "Handover checklist", Slug: "software-project-handover-checklist", Excerpt: "Handover checklist excerpt", Status: "published"},
	} {
		item.Body = "Test article body for editorial links."
		item.CoverImageUrl = "https://example.com/cover.png"
		item.Tags = []string{"engineering"}
		if _, err := posts.Create(t.Context(), item); err != nil {
			t.Fatal(err)
		}
	}

	const assetVersion = "seo-links-test"
	renderer, err := inertia.NewRenderer(inertia.RendererConfig{
		ContainerID:         "app",
		SSRURL:              config.DefaultInertiaSSRURL,
		SSRTimeout:          config.DefaultInertiaSSRRequestTimeout,
		SSRMaxResponseBytes: config.DefaultInertiaSSRMaxResponseBytes,
	}, inertia.WithRoot(resources.Root), inertia.WithVersion(assetVersion))
	if err != nil {
		t.Fatal(err)
	}
	cfg := config.App{SiteName: "MBV Labs", Domain: "localhost", Protocol: "http"}
	h := handlers.NewPublicHandler(cfg, works, nil, posts, nil, nil, renderer)
	srv := anduintest.New(t, func(rt *anduin.Router) { h.Register(rt) }, anduintest.WithMiddleware(renderer.Middleware()))
	headers := []anduintest.Header{{Name: inertia.HeaderInertia, Value: "true"}, {Name: inertia.HeaderVersion, Value: assetVersion}}

	services := srv.Get("/services", headers...)
	services.AssertStatus(http.StatusOK)
	services.AssertInertia("Public/Services")
	if !strings.Contains(services.Text(), `"workItems":[{"Slug":"kuatro-group-backend-infrastructure"}]`) {
		t.Fatal("services should expose only the published case-study link")
	}

	post := srv.Get("/blog/software-project-rescue", headers...)
	post.AssertStatus(http.StatusOK)
	post.AssertInertia("Public/BlogShow")
	for _, expected := range []string{
		`"relatedArticles":[{"slug":"software-project-handover-checklist","title":"Handover checklist"}]`,
		`"title":"Stored rescue guide title - MBV Labs"`,
		`"description":"Stored rescue guide excerpt"`,
	} {
		if !strings.Contains(post.Text(), expected) {
			t.Errorf("blog response missing %q", expected)
		}
	}
}
