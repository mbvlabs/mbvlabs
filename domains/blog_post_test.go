package domains_test

import (
	"errors"
	"testing"
	"time"

	"mbvlabs/domains"
	"mbvlabs/migrations"
)

func TestBlogPostsLifecycle(t *testing.T) {
	db := testCluster.NewTestDB(t, migrations.Migrations, ".")
	posts := domains.NewBlogPosts(db)
	ctx := t.Context()
	data := domains.CreateBlogPostData{
		Title: "A useful blog post", Slug: "useful-blog-post",
		Excerpt: "A useful excerpt", Body: "A useful blog post body",
		Status: "draft", CoverImageUrl: "https://example.com/image.png", Tags: []string{"go"},
	}
	if _, err := posts.Create(ctx, domains.CreateBlogPostData{}); !errors.Is(err, domains.ErrDomainValidation) {
		t.Fatalf("invalid create: %v", err)
	}
	post, err := posts.Create(ctx, data)
	if err != nil {
		t.Fatal(err)
	}
	found, err := posts.FindBySlug(ctx, data.Slug)
	if err != nil || found.ID != post.ID {
		t.Fatalf("find by slug: %+v, %v", found, err)
	}
	if err := post.SetPublicationSchedule(domains.BlogPostPublicationSchedule{JobID: 42, ScheduledAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	post, err = posts.UpdatePublicationSchedule(ctx, post)
	if err != nil {
		t.Fatal(err)
	}
	schedule, err := post.PublicationScheduleData()
	if err != nil || schedule == nil || schedule.JobID != 42 {
		t.Fatalf("schedule: %+v, %v", schedule, err)
	}
	ok, err := posts.PublishIfDraft(ctx, post.ID, time.Now())
	if err != nil || !ok {
		t.Fatalf("publish: %v, %v", ok, err)
	}
	ok, err = posts.PublishIfDraft(ctx, post.ID, time.Now())
	if err != nil || ok {
		t.Fatalf("republish: %v, %v", ok, err)
	}
	post, err = posts.Find(ctx, post.ID)
	if err != nil || post.Status != "published" || !post.PublishedAt.Valid || len(post.PublicationSchedule) != 0 {
		t.Fatalf("published post: %+v, %v", post, err)
	}
	page, err := posts.PaginatePublished(ctx, 0, 0)
	if err != nil || page.TotalCount != 1 || page.Page != 1 || page.PageSize != 10 {
		t.Fatalf("page: %+v, %v", page, err)
	}
	if err := posts.Destroy(ctx, post.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := posts.Find(ctx, post.ID); !errors.Is(err, domains.ErrNotFound) {
		t.Fatalf("find deleted: %v", err)
	}
}
