package handlers

import (
	"testing"
	"unicode/utf8"
)

func TestProjectSearchDescription(t *testing.T) {
	description := projectSearchDescription("deploycrate-community-edition", "original")
	if n := utf8.RuneCountInString(description); description == "original" || n < 120 || n > 155 {
		t.Errorf("invalid project description (%d characters): %q", n, description)
	}
	if got := projectSearchDescription("unknown", "original"); got != "original" {
		t.Errorf("unknown project description changed: %q", got)
	}
}
