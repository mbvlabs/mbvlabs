package domains

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
)

func TestInquiryValidationAndDomain(t *testing.T) {
	e := ProjectInquiryEntity{Name: "Jane Doe", Email: "Jane@Example.COM", Message: "Project description", Status: "new", Metadata: json.RawMessage(`{}`)}
	if err := e.Validate(); err != nil {
		t.Fatalf("valid inquiry: %v", err)
	}
	e.ProjectType = sql.NullString{String: "invalid", Valid: true}
	if err := e.Validate(); err == nil {
		t.Fatal("invalid project type accepted")
	}
	e.ProjectType.Valid = false
	e.Company = sql.NullString{String: strings.Repeat("x", 256), Valid: true}
	if err := e.Validate(); err == nil {
		t.Fatal("oversized company accepted")
	}
	e.Company.Valid = false
	e.Metadata = nil
	if err := e.Validate(); err == nil {
		t.Fatal("missing metadata accepted")
	}
	domain, err := projectInquiryEmailDomain("Jane <Jane@EXAMPLE.COM>")
	if err != nil || domain != "example.com" {
		t.Fatalf("domain = %q, %v", domain, err)
	}
}
