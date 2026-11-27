package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"mbvlabs/domains"
	"mbvlabs/middleware"
	"mbvlabs/payloads"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/telemetry"
)

const contentAPITokenLifetime = time.Hour

// RegisterContentAPITokenRoutes lets a signed-in admin rotate or revoke
// their short-lived content API token. The plaintext is returned only on issue.
func RegisterContentAPITokenRoutes(rt *anduin.Router, renderer *inertia.Renderer, tokens domains.Tokens, signingKey string, sessionAuth anduin.Middleware) {
	withAdmin := func(next func(context.Context, middleware.Viewer) error) anduin.Handler {
		return func(ctx context.Context) error {
			anduin.SetHeader(ctx, "Cache-Control", "no-store")
			viewer, ok := middleware.ViewerFrom(ctx)
			if !ok || !viewer.IsAdmin {
				return apiFailure(ctx, http.StatusForbidden, "admin session required")
			}
			return next(ctx, viewer)
		}
	}

	rt.Route(routes.AdminContentTokenShow).Get(withAdmin(func(ctx context.Context, viewer middleware.Viewer) error {
		token, err := tokens.FindByUserAndScope(ctx, viewer.ID, domains.ContentAPIToken)
		if err != nil && !errors.Is(err, domains.ErrNotFound) {
			telemetry.Error(ctx, "failed to look up content API token status")
			return apiFailure(ctx, http.StatusInternalServerError, "failed to look up content API token status")
		}
		var expiresAt *time.Time
		if err == nil && time.Now().Before(token.ExpiresAt) {
			expiresAt = &token.ExpiresAt
		}
		return renderer.Page(ctx, "Admin/ContentToken", inertia.FromStruct(payloads.ContentTokenPageProps{
			ExpiresAt: expiresAt,
			CSRFToken: anduin.CSRFToken(ctx),
		})).Render()
	}), sessionAuth)

	rt.Route(routes.AdminContentTokenCreate).Post(withAdmin(func(ctx context.Context, viewer middleware.Viewer) error {
		expiresAt := time.Now().UTC().Add(contentAPITokenLifetime)
		token, err := tokens.Replace(ctx, signingKey, domains.ContentAPIToken, expiresAt, []byte("{}"), viewer.ID)
		if err != nil {
			telemetry.Error(ctx, "failed to issue content API token")
			return apiFailure(ctx, http.StatusInternalServerError, "failed to issue content API token")
		}
		return anduin.JSON(ctx, http.StatusCreated, payloads.IssuedContentToken{Token: token, ExpiresAt: expiresAt})
	}), sessionAuth)

	rt.Route(routes.AdminContentTokenDelete).Delete(withAdmin(func(ctx context.Context, viewer middleware.Viewer) error {
		if err := tokens.DestroyByUserAndScope(ctx, viewer.ID, domains.ContentAPIToken); err != nil {
			telemetry.Error(ctx, "failed to revoke content API token")
			return apiFailure(ctx, http.StatusInternalServerError, "failed to revoke content API token")
		}
		return anduin.NoContent(ctx, http.StatusNoContent)
	}), sessionAuth)
}
