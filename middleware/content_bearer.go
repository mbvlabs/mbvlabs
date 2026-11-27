package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"mbvlabs/domains"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/telemetry"
)

// ContentBearerAuth accepts only scoped content API tokens, never login
// cookies or tokens issued for password reset or email verification.
func ContentBearerAuth(tokens domains.Tokens, users domains.Users, signingKey string) anduin.Middleware {
	return func(next anduin.Handler) anduin.Handler {
		return func(ctx context.Context) error {
			anduin.SetHeader(ctx, "Cache-Control", "no-store")
			plain := contentBearer(anduin.Request(ctx))
			if plain == "" || signingKey == "" {
				return rejectContentBearer(ctx)
			}

			token, err := tokens.FindByScopeAndHash(ctx, signingKey, domains.ContentAPIToken, plain)
			if errors.Is(err, domains.ErrNotFound) {
				return rejectContentBearer(ctx)
			}
			if err != nil {
				return contentAuthUnavailable(ctx)
			}
			if !time.Now().Before(token.ExpiresAt) {
				return rejectContentBearer(ctx)
			}

			user, err := users.Find(ctx, token.UserID)
			if errors.Is(err, domains.ErrNotFound) {
				return rejectContentBearer(ctx)
			}
			if err != nil {
				return contentAuthUnavailable(ctx)
			}
			if !user.IsAdmin {
				return rejectContentBearer(ctx)
			}
			return next(WithViewer(ctx, user))
		}
	}
}

func contentBearer(req *http.Request) string {
	if req == nil {
		return ""
	}
	values := req.Header.Values("Authorization")
	if len(values) != 1 {
		return ""
	}
	scheme, token, ok := strings.Cut(values[0], " ")
	token = strings.TrimSpace(token)
	if !ok || !strings.EqualFold(scheme, "Bearer") || len(token) > 256 || strings.ContainsAny(token, " \t\r\n,") {
		return ""
	}
	return token
}

func rejectContentBearer(ctx context.Context) error {
	anduin.SetHeader(ctx, "WWW-Authenticate", `Bearer realm="mbvlabs content"`)
	return anduin.JSON(ctx, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
}

func contentAuthUnavailable(ctx context.Context) error {
	telemetry.Error(ctx, "content API authentication unavailable")
	return anduin.JSON(ctx, http.StatusInternalServerError, map[string]string{"message": "Authentication unavailable"})
}
