package domains

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"
)

func TestBlogPostValidationAndSchedule(t *testing.T) {
	post := BlogPostEntity{
		Title: "A useful blog post", Slug: "useful-blog-post", Status: "draft",
		Body: "A useful blog post body", Tags: json.RawMessage(`["go"]`),
	}
	if err := post.Validate(); err != nil {
		t.Fatalf("valid draft: %v", err)
	}
	post.Slug = "Not A Slug"
	if err := post.Validate(); err == nil {
		t.Fatal("invalid slug accepted")
	}
	post.Slug = "useful-blog-post"
	post.Status = "published"
	if err := post.Validate(); err == nil {
		t.Fatal("incomplete published post accepted")
	}
	post.Excerpt = sql.NullString{String: "A useful excerpt", Valid: true}
	post.CoverImageUrl = sql.NullString{String: "https://example.com/image.png", Valid: true}
	post.PublishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	if err := post.Validate(); err != nil {
		t.Fatalf("valid published post: %v", err)
	}
	if err := post.SetPublicationSchedule(BlogPostPublicationSchedule{JobID: 7, ScheduledAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	schedule, err := post.PublicationScheduleData()
	if err != nil || schedule == nil || schedule.JobID != 7 {
		t.Fatalf("schedule = %+v, %v", schedule, err)
	}
	post.PublicationSchedule = json.RawMessage(`{"job_id":`)
	if err := post.Validate(); err == nil {
		t.Fatal("invalid schedule accepted")
	}
}
