package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"mbvlabs/middleware"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
)

func TestAPIBasicAuth(t *testing.T) {
	for _, tc := range []struct {
		name, configuredUser, configuredPass, user, pass string
		want                                             int
	}{
		{"valid", "u", "p", "u", "p", 200},
		{"wrong", "u", "p", "u", "wrong", 401},
		{"empty-config", "", "", "", "", 401},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rt := anduin.NewRouter()
			rt.Route(routes.ApiSearch).Post(func(ctx context.Context) error {
				return anduin.NoContent(ctx, 200)
			}, middleware.APIBasicAuth(tc.configuredUser, tc.configuredPass))
			req := httptest.NewRequest("POST", routes.ApiSearch.Path(), nil)
			req.SetBasicAuth(tc.user, tc.pass)
			rec := httptest.NewRecorder()
			rt.Handler().ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status %d, want %d: %s", rec.Code, tc.want, rec.Body)
			}
			if tc.want == http.StatusUnauthorized && rec.Header().Get("WWW-Authenticate") != `Basic realm="mbvlabs API"` {
				t.Fatalf("challenge = %q", rec.Header().Get("WWW-Authenticate"))
			}
		})
	}
}
