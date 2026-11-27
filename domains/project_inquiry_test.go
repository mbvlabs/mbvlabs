package domains_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"mbvlabs/domains"
	"mbvlabs/migrations"
)

func TestProjectInquiryValidation(t *testing.T) {
	valid := domains.ProjectInquiryEntity{
		Name: "Jane Doe", Email: "jane@example.com", Message: "A project description",
		Status: "new", Metadata: json.RawMessage(`{}`),
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid inquiry: %v", err)
	}
	for _, mutate := range []func(*domains.ProjectInquiryEntity){
		func(e *domains.ProjectInquiryEntity) { e.Email = "invalid" },
		func(e *domains.ProjectInquiryEntity) { e.Status = "unknown" },
		func(e *domains.ProjectInquiryEntity) { e.ProjectType = sql.NullString{String: "invalid", Valid: true} },
		func(e *domains.ProjectInquiryEntity) {
			e.Company = sql.NullString{String: strings.Repeat("x", 256), Valid: true}
		},
		func(e *domains.ProjectInquiryEntity) { e.Metadata = nil },
	} {
		e := valid
		mutate(&e)
		if err := e.Validate(); err == nil {
			t.Fatalf("expected validation error for %+v", e)
		}
	}
}

func TestProjectInquiriesPersistence(t *testing.T) {
	t.Parallel()
	db := testCluster.NewTestDB(t, migrations.Migrations, ".")
	p := domains.NewProjectInquiries(db)
	ctx := t.Context()
	data := domains.CreateProjectInquiryData{
		Name: "Jane Doe", Email: "Jane@Example.COM", Message: "A project description",
		Company: sql.NullString{String: "Example", Valid: true},
		Status:  "new", Metadata: json.RawMessage(`{"key":"value"}`),
	}
	created, err := p.Create(ctx, data)
	if err != nil {
		t.Fatal(err)
	}
	found, err := p.Find(ctx, created.ID)
	if err != nil || found.Company != data.Company || found.Email != data.Email || !json.Valid(found.Metadata) {
		t.Fatalf("find = %+v, %v", found, err)
	}
	domain, err := p.BanDomain(ctx, created.ID)
	if err != nil || domain != "example.com" {
		t.Fatalf("ban = %q, %v", domain, err)
	}
	if _, err = p.BanDomain(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	banned, err := p.IsDomainBanned(ctx, "someone@EXAMPLE.com")
	if err != nil || !banned {
		t.Fatalf("banned = %v, %v", banned, err)
	}
	if _, err = p.UnbanDomain(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	banned, err = p.IsDomainBanned(ctx, "someone@example.com")
	if err != nil || banned {
		t.Fatalf("unbanned = %v, %v", banned, err)
	}
	updated, err := p.Update(ctx, domains.UpdateProjectInquiryData{
		ID: created.ID, Name: data.Name, Email: data.Email, Message: data.Message,
		Status: "contacted", Metadata: data.Metadata,
	})
	if err != nil || updated.Status != "contacted" || updated.Company.Valid || !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Fatalf("update = %+v, %v", updated, err)
	}
	if _, err = p.Upsert(ctx, data); err != nil {
		t.Fatal(err)
	}
	all, err := p.All(ctx)
	if err != nil || len(all) != 2 {
		t.Fatalf("all = %d, %v", len(all), err)
	}
	page, err := p.Paginate(ctx, 1, 1)
	if err != nil || page.TotalCount != 2 || page.TotalPages != 2 || len(page.ProjectInquiries) != 1 {
		t.Fatalf("page = %+v, %v", page, err)
	}
	if err = p.Destroy(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err = p.Find(ctx, created.ID); !errors.Is(err, domains.ErrNotFound) {
		t.Fatalf("find deleted: %v", err)
	}
}
