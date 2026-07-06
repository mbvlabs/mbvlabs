package validation

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type testPostStatus string

type testValidatable struct {
	err error
}

func (v testValidatable) Validate() error {
	return v.err
}

func TestBuilderUsesExplicitFieldNames(t *testing.T) {
	builder := NewBuilder()
	status := testPostStatus("archived")

	builder.Required("name", "")
	builder.OneOf("status", status, "draft", "published")

	errs := builder.Errors()
	if len(errs) != 2 {
		t.Fatalf("len(Errors()) = %d, want 2", len(errs))
	}

	if errs[0].Field != "name" {
		t.Fatalf("first error field = %q, want name", errs[0].Field)
	}
	if errs[0].Code != "required" {
		t.Fatalf("first error code = %q, want required", errs[0].Code)
	}

	if errs[1].Field != "status" {
		t.Fatalf("second error field = %q, want status", errs[1].Field)
	}
	if errs[1].Code != "one_of" {
		t.Fatalf("second error code = %q, want one_of", errs[1].Code)
	}
}

func TestEmptyFieldNameUsesUnknown(t *testing.T) {
	builder := NewBuilder()

	builder.Required("", "")

	errs := builder.Errors()
	if len(errs) != 1 {
		t.Fatalf("len(Errors()) = %d, want 1", len(errs))
	}

	if errs[0].Field != UnknownField {
		t.Fatalf("error field = %q, want %q", errs[0].Field, UnknownField)
	}
}

func TestRequiredSupportsCommonValues(t *testing.T) {
	var nilString *string
	nonZeroTime := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{name: "blank string", value: "  ", wantErr: true},
		{name: "nonblank string", value: "server-01", wantErr: false},
		{name: "nil string pointer", value: nilString, wantErr: true},
		{name: "invalid null string", value: sql.NullString{}, wantErr: true},
		{
			name:    "valid null string",
			value:   sql.NullString{String: "ok", Valid: true},
			wantErr: false,
		},
		{name: "zero int32", value: int32(0), wantErr: true},
		{name: "nonzero int32", value: int32(22), wantErr: false},
		{name: "zero uuid", value: uuid.Nil, wantErr: true},
		{name: "nonzero uuid", value: uuid.New(), wantErr: false},
		{name: "zero time", value: time.Time{}, wantErr: true},
		{name: "nonzero time", value: nonZeroTime, wantErr: false},
		{name: "empty slice", value: []string{}, wantErr: true},
		{name: "nonempty slice", value: []string{"go"}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder()

			builder.Required("field", tt.value)

			gotErr := len(builder.Errors()) > 0
			if gotErr != tt.wantErr {
				t.Fatalf("Required() error = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestOptionalStringValidatorsSkipMissingValues(t *testing.T) {
	builder := NewBuilder()

	builder.MinLen("name", "", 3)
	builder.MaxLen("summary", sql.NullString{}, 10)
	builder.OneOf("status", (*string)(nil), "draft")
	builder.URL("website", " ")

	if len(builder.Errors()) != 0 {
		t.Fatalf("len(Errors()) = %d, want 0", len(builder.Errors()))
	}
}

func TestFieldHelpersAddStructuredErrors(t *testing.T) {
	start := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	end := time.Date(2026, 7, 2, 11, 0, 0, 0, time.UTC)
	builder := NewBuilder()

	builder.MinLen("name", "ab", 3)
	builder.MaxLen("summary", "abcd", 3)
	builder.URL("website", "not-a-url")
	builder.MinItems("applications", []string{}, 1)
	builder.MaxItems("ports", []int{1, 2}, 1)
	builder.NoBlankItems("labels", []string{"api", " "})
	builder.TimeBeforeOrEqual("started_at", start, "completed_at", end)
	builder.True("accepted_terms", false)

	errs := builder.Errors()
	if len(errs) != 8 {
		t.Fatalf("len(Errors()) = %d, want 8", len(errs))
	}

	wantFields := []string{
		"name",
		"summary",
		"website",
		"applications",
		"ports",
		"labels[1]",
		"started_at",
		"accepted_terms",
	}
	for i, want := range wantFields {
		if errs[i].Field != want {
			t.Fatalf("error %d field = %q, want %q", i, errs[i].Field, want)
		}
	}

	if errs[0].Params["min"] != 3 {
		t.Fatalf("min params = %#v, want min=3", errs[0].Params)
	}
	if errs[6].Params["other"] != "completed_at" {
		t.Fatalf("time params = %#v, want other=completed_at", errs[6].Params)
	}
}

func TestFieldHelpersAcceptCustomMessages(t *testing.T) {
	start := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	end := time.Date(2026, 7, 2, 11, 0, 0, 0, time.UTC)
	builder := NewBuilder()

	builder.Required("name", "", "enter a name")
	builder.RequiredWhen(true, "headline", "", "enter a headline")
	builder.MinLen("title", "ab", 3, "title is too short")
	builder.MaxLen("summary", "abcd", 3, "summary is too long")
	builder.LenBetween("slug", "ab", 3, 10, "slug length is invalid")
	builder.OneOfWithMessage("status", "archived", "choose a valid status", "draft", "published")
	builder.URL("website", "not-a-url", "enter a valid website")
	builder.RequiredURL("callbackUrl", "", "enter a callback URL")
	builder.MinItems("applications", []string{}, 1, "add an application")
	builder.MaxItems("ports", []int{1, 2}, 1, "remove a port")
	builder.NoBlankItems("labels", []string{"api", " "}, "label cannot be blank")
	builder.TimeBeforeOrEqual(
		"started_at",
		start,
		"completed_at",
		end,
		"start must be before completion",
	)
	builder.True("accepted_terms", false, "accept the terms")

	errs := builder.Errors()
	wantMessages := []string{
		"enter a name",
		"enter a headline",
		"title is too short",
		"summary is too long",
		"slug length is invalid",
		"choose a valid status",
		"enter a valid website",
		"enter a callback URL",
		"add an application",
		"remove a port",
		"label cannot be blank",
		"start must be before completion",
		"accept the terms",
	}

	if len(errs) != len(wantMessages) {
		t.Fatalf("len(Errors()) = %d, want %d", len(errs), len(wantMessages))
	}

	for i, want := range wantMessages {
		if errs[i].Message != want {
			t.Fatalf("error %d message = %q, want %q", i, errs[i].Message, want)
		}
	}
}

func TestBlankCustomMessageFallsBackToDefault(t *testing.T) {
	builder := NewBuilder()

	builder.Required("name", "", "")
	builder.OneOfWithMessage("status", "archived", "", "draft", "published")

	errs := builder.Errors()
	if len(errs) != 2 {
		t.Fatalf("len(Errors()) = %d, want 2", len(errs))
	}

	if errs[0].Message != "is required" {
		t.Fatalf("required message = %q, want default", errs[0].Message)
	}
	if errs[1].Message != "has an invalid value" {
		t.Fatalf("one_of message = %q, want default", errs[1].Message)
	}
}

func TestErrAndAs(t *testing.T) {
	builder := NewBuilder()
	if builder.Err() != nil {
		t.Fatal("Err() for empty builder returned non-nil error")
	}

	builder.Required("name", "")

	err := builder.Err()
	errs, ok := As(err)
	if !ok {
		t.Fatal("As() did not identify ValidationErrors")
	}
	if errs.Len() != 1 || errs.Empty() {
		t.Fatalf(
			"ValidationErrors state = len %d empty %v, want len 1 empty false",
			errs.Len(),
			errs.Empty(),
		)
	}
}

func TestValidationErrorsToMap(t *testing.T) {
	errs := ValidationErrors{
		{Field: "title", Message: "is required"},
		{Field: "coverImageUrl", Message: "is required"},
		{Field: "specialisms[1]", Message: "is required"},
		{Field: "title", Message: "must be at least 10 characters"},
	}

	got := errs.ToMap()

	if got["title"] != "is required" {
		t.Fatalf("title = %v, want first title error", got["title"])
	}
	if got["coverImageUrl"] != "is required" {
		t.Fatalf("coverImageUrl = %v, want field error", got["coverImageUrl"])
	}
	if got["specialisms[1]"] != "is required" {
		t.Fatalf("specialisms[1] = %v, want indexed field error", got["specialisms[1]"])
	}
}

func TestValidateCallsValidatableValues(t *testing.T) {
	wantErr := errors.New("invalid")

	if err := Validate(testValidatable{err: wantErr}); !errors.Is(err, wantErr) {
		t.Fatalf("Validate() error = %v, want %v", err, wantErr)
	}

	if err := Validate(struct{}{}); err != nil {
		t.Fatalf("Validate() non-validatable error = %v, want nil", err)
	}
}
