package api_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"mbvlabs/config"
	"mbvlabs/controllers/api"
	"mbvlabs/database"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/models/factories"
	"mbvlabs/queue"
	"mbvlabs/queue/jobs"
	"mbvlabs/services"

	"github.com/labstack/echo/v5"
	"github.com/riverqueue/river"
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

func TestArticles_CreateScheduledAndPublish(t *testing.T) {
	ctx := context.Background()
	db := testCluster.NewTestDB(t, database.Migrations, "migrations")
	workers := river.NewWorkers()
	invalidations := 0
	worker := queue.NewPublishBlogPostWorker(db, func() { invalidations++ })
	if err := worker.Register(workers); err != nil {
		t.Fatalf("register publish worker: %v", err)
	}
	insertOnly, err := queue.NewInsertOnly(db, workers)
	if err != nil {
		t.Fatalf("create insert queue: %v", err)
	}
	controller := api.NewArticles(services.NewBlogPosts(db, insertOnly), config.Config{})

	publishAt := time.Now().UTC().Add(time.Hour).Truncate(time.Microsecond)
	body, err := json.Marshal(api.CreateArticlePayload{
		Title:         "Scheduled article title",
		Slug:          "scheduled-article",
		Excerpt:       "A publication-ready excerpt.",
		Body:          "A publication-ready article body.",
		CoverImageURL: "https://example.com/cover.jpg",
		Tags:          []string{"systems"},
		ScheduledAt:   pointer(publishAt.Format(time.RFC3339Nano)),
	})
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/articles", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	if err := controller.Create(e.NewContext(req, rec)); err != nil {
		t.Fatalf("create scheduled article: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var created models.BlogPostEntity
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created article: %v", err)
	}
	created, err = models.BlogPost.Find(ctx, db.Executor(), created.ID)
	if err != nil {
		t.Fatalf("find created article: %v", err)
	}
	if created.Status != models.Draft.String() || created.PublishedAt.Valid {
		t.Fatalf("created article status = %q, published_at = %v", created.Status, created.PublishedAt)
	}
	if schedule, err := created.PublicationScheduleData(); err != nil || schedule == nil || schedule.JobID <= 0 || !schedule.ScheduledAt.Equal(publishAt) {
		t.Fatalf("created schedule = %#v, %v", schedule, err)
	}
	if err := worker.Work(ctx, &river.Job[jobs.PublishBlogPostArgs]{Args: jobs.PublishBlogPostArgs{
		BlogPostID: created.ID,
		PublishAt:  publishAt,
	}}); err != nil {
		t.Fatalf("publish scheduled article: %v", err)
	}

	published, err := models.BlogPost.Find(ctx, db.Executor(), created.ID)
	if err != nil {
		t.Fatalf("find published article: %v", err)
	}
	if published.Status != models.Published.String() || !published.PublishedAt.Time.Equal(publishAt) {
		t.Fatalf("published article status = %q, published_at = %v", published.Status, published.PublishedAt)
	}
	if published.PublicationSchedule != nil || invalidations != 1 {
		t.Fatalf("publication schedule = %s, invalidations = %d", published.PublicationSchedule, invalidations)
	}

	manualPublishAt := publishAt.Add(time.Hour)
	fixture := models.BlogPostEntity{}
	if err := fixture.SetPublicationSchedule(models.BlogPostPublicationSchedule{
		JobID:       999,
		ScheduledAt: manualPublishAt,
	}); err != nil {
		t.Fatal(err)
	}
	manual, err := factories.CreateBlogPost(ctx, db.Executor(),
		factories.WithBlogPostsTitle("Manually published article"),
		factories.WithBlogPostsSlug("manually-published-article"),
		factories.WithBlogPostsExcerpt(sql.NullString{String: "A publication-ready excerpt.", Valid: true}),
		factories.WithBlogPostsBody("A publication-ready article body."),
		factories.WithBlogPostsStatus(models.Draft.String()),
		factories.WithBlogPostsCoverImageUrl(sql.NullString{String: "https://example.com/manual.jpg", Valid: true}),
		factories.WithBlogPostsTags(json.RawMessage(`["systems"]`)),
		factories.WithBlogPostsPublishedAt(sql.NullTime{}),
		factories.WithBlogPostsPublicationSchedule(fixture.PublicationSchedule),
	)
	if err != nil {
		t.Fatalf("create manual article fixture: %v", err)
	}
	manual, err = models.BlogPost.Update(ctx, db.Executor(), models.UpdateBlogPostData{
		ID:            manual.ID,
		Title:         manual.Title,
		Slug:          manual.Slug,
		Excerpt:       manual.Excerpt.String,
		Body:          manual.Body,
		Status:        models.Published.String(),
		CoverImageUrl: manual.CoverImageUrl.String,
		Tags:          []string{"systems"},
	})
	if err != nil {
		t.Fatalf("manually publish article: %v", err)
	}

	if err := worker.Work(ctx, &river.Job[jobs.PublishBlogPostArgs]{Args: jobs.PublishBlogPostArgs{
		BlogPostID: manual.ID,
		PublishAt:  manualPublishAt,
	}}); err != nil {
		t.Fatalf("run no-op publication: %v", err)
	}
	unchanged, err := models.BlogPost.Find(ctx, db.Executor(), manual.ID)
	if err != nil {
		t.Fatalf("find manually published article: %v", err)
	}
	if !unchanged.PublishedAt.Time.Equal(manual.PublishedAt.Time) || invalidations != 1 {
		t.Fatalf("manual published_at changed from %v to %v, invalidations = %d", manual.PublishedAt, unchanged.PublishedAt, invalidations)
	}
}

func pointer[T any](value T) *T { return &value }
