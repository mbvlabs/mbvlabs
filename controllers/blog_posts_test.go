package controllers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"mbvlabs/config"
	"mbvlabs/controllers"
	"mbvlabs/database"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/models/factories"

	"github.com/labstack/echo/v5"
)

var testCluster *storage.TestCluster

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	testCluster, err = storage.NewTestCluster(ctx)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	if err := testCluster.Close(ctx); err != nil && code == 0 {
		panic(err)
	}
	os.Exit(code)
}

type rssDocument struct {
	Version string `xml:"version,attr"`
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	GUID  struct {
		IsPermaLink string `xml:"isPermaLink,attr"`
		Value       string `xml:",chardata"`
	} `xml:"guid"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func TestBlogPosts_Feed(t *testing.T) {
	originalBaseURL := config.BaseURL
	config.BaseURL = "https://example.test"
	t.Cleanup(func() { config.BaseURL = originalBaseURL })

	t.Run("returns the latest 25 published posts", func(t *testing.T) {
		ctx := context.Background()
		db := testCluster.NewTestDB(t, database.Migrations, "migrations")
		publishedAt := time.Date(2026, time.July, 1, 12, 0, 0, 0, time.UTC)

		for i := range 27 {
			if _, err := factories.CreateBlogPost(ctx, db.Executor(),
				factories.WithBlogPostsTitle(fmt.Sprintf("Published post %02d", i)),
				factories.WithBlogPostsSlug(fmt.Sprintf("published-post-%02d", i)),
				factories.WithBlogPostsExcerpt(sql.NullString{String: fmt.Sprintf("Excerpt %02d & details", i), Valid: true}),
				factories.WithBlogPostsBody("A complete published article body."),
				factories.WithBlogPostsStatus(models.Published.String()),
				factories.WithBlogPostsTags(json.RawMessage(`["systems"]`)),
				factories.WithBlogPostsPublishedAt(sql.NullTime{Time: publishedAt.Add(time.Duration(i) * time.Hour), Valid: true}),
				factories.WithBlogPostsPublicationSchedule(nil),
			); err != nil {
				t.Fatalf("create published fixture %d: %v", i, err)
			}
		}

		for _, status := range []string{models.Draft.String(), "archived"} {
			if _, err := factories.CreateBlogPost(ctx, db.Executor(),
				factories.WithBlogPostsTitle(status+" post"),
				factories.WithBlogPostsSlug(status+"-post"),
				factories.WithBlogPostsExcerpt(sql.NullString{String: status + " excerpt", Valid: true}),
				factories.WithBlogPostsBody("A non-published article body."),
				factories.WithBlogPostsStatus(status),
				factories.WithBlogPostsTags(json.RawMessage(`["systems"]`)),
				factories.WithBlogPostsPublishedAt(sql.NullTime{Time: publishedAt.Add(100 * time.Hour), Valid: true}),
				factories.WithBlogPostsPublicationSchedule(nil),
			); err != nil {
				t.Fatalf("create %s fixture: %v", status, err)
			}
		}

		controller := controllers.NewBlogPosts(db)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/rss.xml", nil)
		rec := httptest.NewRecorder()
		if err := controller.Feed(e.NewContext(req, rec)); err != nil {
			t.Fatalf("render feed: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("feed status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
		}
		if got, want := rec.Header().Get(echo.HeaderContentType), "application/rss+xml; charset=utf-8"; got != want {
			t.Fatalf("content type = %q, want %q", got, want)
		}

		var feed rssDocument
		if err := xml.Unmarshal(rec.Body.Bytes(), &feed); err != nil {
			t.Fatalf("decode feed: %v\n%s", err, rec.Body.String())
		}
		if feed.Version != "2.0" || feed.Channel.Title != "MBV Labs Blog" || feed.Channel.Link != "https://example.test/blog" || feed.Channel.Language != "en-US" {
			t.Fatalf("unexpected channel metadata: %#v", feed)
		}
		if len(feed.Channel.Items) != 25 {
			t.Fatalf("item count = %d, want 25", len(feed.Channel.Items))
		}

		first := feed.Channel.Items[0]
		if first.Title != "Published post 26" || first.Link != "https://example.test/blog/published-post-26" || first.GUID.Value != first.Link || first.GUID.IsPermaLink != "true" || first.Description != "Excerpt 26 & details" || first.PubDate != publishedAt.Add(26*time.Hour).Format(time.RFC1123Z) {
			t.Errorf("first item = %#v", first)
		}
		last := feed.Channel.Items[24]
		if last.Link != "https://example.test/blog/published-post-02" {
			t.Errorf("last item link = %q, want newest 25 posts only", last.Link)
		}
		for _, item := range feed.Channel.Items {
			if item.Link == "https://example.test/blog/draft-post" || item.Link == "https://example.test/blog/archived-post" {
				t.Errorf("non-published item included: %q", item.Link)
			}
		}
	})

	t.Run("returns a valid empty feed", func(t *testing.T) {
		db := testCluster.NewTestDB(t, database.Migrations, "migrations")
		controller := controllers.NewBlogPosts(db)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/rss.xml", nil)
		rec := httptest.NewRecorder()
		if err := controller.Feed(e.NewContext(req, rec)); err != nil {
			t.Fatalf("render empty feed: %v", err)
		}

		var feed rssDocument
		if err := xml.Unmarshal(rec.Body.Bytes(), &feed); err != nil {
			t.Fatalf("decode empty feed: %v", err)
		}
		if len(feed.Channel.Items) != 0 {
			t.Fatalf("empty feed item count = %d, want 0", len(feed.Channel.Items))
		}
	})
}
