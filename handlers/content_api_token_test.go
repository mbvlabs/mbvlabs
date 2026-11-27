package handlers_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"mbvlabs/domains"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin/anduintest"
)

func TestContentTokenLifecycle(t *testing.T) {
	app := newAuthTestApp(t)
	server := app.server
	path := routes.AdminContentTokenShow.Path()
	const protected = "/test/content-api"

	server.PostJSON(path, nil).AssertStatus(http.StatusSeeOther)
	user := app.verifiedUser(t, "content-editor@example.com")
	app.signIn(t, user.Email)
	server.Get(path).AssertStatus(http.StatusForbidden)
	server.PostJSON(path, nil).AssertStatus(http.StatusForbidden)
	server.Get(protected).AssertStatus(http.StatusUnauthorized)

	isAdmin := true
	if _, err := app.users.Update(t.Context(), domains.UpdateUserData{ID: user.ID, IsAdmin: &isAdmin}); err != nil {
		t.Fatal(err)
	}

	checkBearer := func(token string, status int) {
		t.Helper()
		resp := server.Get(protected, anduintest.Header{Name: "Authorization", Value: "Bearer " + token})
		resp.AssertStatus(status)
		if resp.Header.Get("Cache-Control") != "no-store" {
			t.Error("bearer response should not be cached")
		}
		if status == http.StatusUnauthorized && !strings.HasPrefix(resp.Header.Get("WWW-Authenticate"), "Bearer ") {
			t.Error("unauthorized response missing Bearer challenge")
		}
	}

	for _, tc := range []struct {
		scope   domains.TokenScope
		expires time.Time
	}{
		{domains.UserResetPassword, time.Now().Add(time.Hour)},
		{domains.ContentAPIToken, time.Now().Add(-time.Hour)},
	} {
		token, err := app.tokens.Replace(t.Context(), "test-token-signing-key", tc.scope, tc.expires, []byte("{}"), user.ID)
		if err != nil {
			t.Fatal(err)
		}
		checkBearer(token, http.StatusUnauthorized)
	}

	issue := func() string {
		t.Helper()
		resp := server.PostJSON(path, nil)
		resp.AssertStatus(http.StatusCreated)
		if resp.Header.Get("Cache-Control") != "no-store" {
			t.Error("issuance response should not be cached")
		}
		var result struct {
			Token     string
			ExpiresAt time.Time
		}
		resp.JSON(&result)
		if result.Token == "" || !result.ExpiresAt.After(time.Now()) || result.ExpiresAt.After(time.Now().Add(time.Hour)) {
			t.Fatal("invalid content token response")
		}
		return result.Token
	}
	old := issue()
	checkBearer(old, http.StatusNoContent)
	current := issue()
	checkBearer(old, http.StatusUnauthorized)
	checkBearer(current, http.StatusNoContent)

	page := server.Get(path, anduintest.Header{Name: "X-Inertia", Value: "true"})
	page.AssertStatus(http.StatusOK)
	page.AssertInertia("Admin/ContentToken")
	if strings.Contains(page.Text(), current) {
		t.Fatal("admin status page exposes token plaintext")
	}

	req, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, server.URL(path), nil)
	if err != nil {
		t.Fatal(err)
	}
	server.Do(req).AssertStatus(http.StatusNoContent)
	checkBearer(current, http.StatusUnauthorized)

	current = issue()
	isAdmin = false
	if _, err := app.users.Update(t.Context(), domains.UpdateUserData{ID: user.ID, IsAdmin: &isAdmin}); err != nil {
		t.Fatal(err)
	}
	checkBearer(current, http.StatusUnauthorized)
}
