package domains_test

import (
	_ "embed"
	"errors"
	"testing"
	"time"

	"mbvlabs/domains"
	"mbvlabs/migrations"
)

//go:embed testdata/reset_content_sequence.sql
var resetContentSequence string

func TestWorksAndProjectsPersistence(t *testing.T) {
	t.Parallel()
	db := testCluster.NewTestDB(t, migrations.Migrations, ".")
	ctx := t.Context()

	works := domains.NewWorks(db)
	data := domains.CreateWorkData{
		Title: "A sample work title", Slug: "sample-work", Summary: "a valid summary",
		Status: domains.Draft, Specialisms: []string{}, Platforms: []string{}, Technologies: []string{},
	}
	work, err := works.Create(ctx, data)
	if err != nil {
		t.Fatal(err)
	}
	found, err := works.FindBySlug(ctx, work.Slug)
	if err != nil || found.ID != work.ID {
		t.Fatalf("find work = %+v, %v", found, err)
	}
	page, err := works.Paginate(ctx, 0, 0)
	if err != nil || page.TotalCount != 1 || len(page.Works) != 1 || page.PageSize != 10 {
		t.Fatalf("paginate works = %+v, %v", page, err)
	}
	_, err = works.Update(ctx, domains.UpdateWorkData{
		ID: work.ID, Title: work.Title, Slug: work.Slug, Summary: work.Summary,
		Status: domains.Draft, Specialisms: []string{}, Platforms: []string{}, Technologies: []string{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := works.Create(ctx, domains.CreateWorkData{}); !errors.Is(err, domains.ErrDomainValidation) {
		t.Fatalf("invalid work error = %v", err)
	}
	// Legacy Upsert conflicts on a generated ID, not on slug.
	if _, err := db.Exec(ctx, resetContentSequence, "works_id_seq", work.ID); err != nil {
		t.Fatal(err)
	}
	data.Slug = "updated-work"
	upsertedWork, err := works.Upsert(ctx, data)
	if err != nil || upsertedWork.ID != work.ID || upsertedWork.Slug != data.Slug {
		t.Fatalf("upsert work = %+v, %v", upsertedWork, err)
	}
	projects := domains.NewProjects(db)
	project, err := projects.Create(ctx, domains.CreateProjectData{
		Name: "Sample project", Slug: "sample-project", Tagline: "A valid tagline", ProjectType: "website",
		Technologies: []string{"Go"}, PublishedAt: time.Now(), IsFeatured: true, Description: "A longer project description",
	})
	if err != nil {
		t.Fatal(err)
	}
	foundProject, err := projects.Find(ctx, project.ID)
	if err != nil || foundProject.ID != project.ID || string(foundProject.Technologies) != `["Go"]` {
		t.Fatalf("find project = %+v, %v", foundProject, err)
	}
	featured, err := projects.FeaturedPublished(ctx, 0)
	if err != nil || len(featured) != 1 {
		t.Fatalf("featured projects = %+v, %v", featured, err)
	}
	_, err = projects.Update(ctx, domains.UpdateProjectData{
		ID: project.ID, Name: project.Name, Slug: project.Slug, Tagline: project.Tagline,
		ProjectType: project.ProjectType, Technologies: []string{"Postgres"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, resetContentSequence, "projects_id_seq", project.ID); err != nil {
		t.Fatal(err)
	}
	upsertedProject, err := projects.Upsert(ctx, domains.CreateProjectData{
		Name: project.Name, Slug: "updated-project", Tagline: project.Tagline,
		ProjectType: project.ProjectType, Technologies: []string{"Postgres"},
	})
	if err != nil || upsertedProject.ID != project.ID || upsertedProject.Slug != "updated-project" {
		t.Fatalf("upsert project = %+v, %v", upsertedProject, err)
	}
	if err := projects.Destroy(ctx, project.ID); err != nil {
		t.Fatal(err)
	}
	if err := works.Destroy(ctx, work.ID); err != nil {
		t.Fatal(err)
	}
}
