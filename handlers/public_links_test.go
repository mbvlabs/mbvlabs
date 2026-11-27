package handlers_test

import (
	"net/http"
	"strings"
	"testing"

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
	if _, err := db.Exec(t.Context(), `INSERT INTO works
		(created_at, updated_at, title, slug, summary, specialisms, platforms, technologies, challenge, approach, deliverables, outcome, content, status, is_featured)
		SELECT now(), now(), title, slug, title, '{}', '{}', '{}', '', '', '', '', '', status, false FROM (VALUES
			('Kuatro backend work', 'kuatro-group-backend-infrastructure', 'published'),
			('Everphone internal search', 'everphone-internal-tool-search', 'draft')
		) AS fixtures(title, slug, status)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(t.Context(), `INSERT INTO blog_posts (created_at, updated_at, title, slug, excerpt, body, status, tags) VALUES
		(now(), now(), 'Stored rescue guide title', 'software-project-rescue', 'Stored rescue guide excerpt', '', 'published', '[]'),
		(now(), now(), 'Audit process draft', 'codebase-audit', '', '', 'draft', '[]'),
		(now(), now(), 'Handover checklist', 'software-project-handover-checklist', '', '', 'published', '[]')`); err != nil {
		t.Fatal(err)
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
	h := handlers.NewPublicHandler(cfg, domains.NewWorks(db), nil, domains.NewBlogPosts(db), nil, nil, renderer)
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
