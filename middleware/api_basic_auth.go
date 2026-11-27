package middleware

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/mbvlabs/andurel/pkg/anduin"
)

// APIBasicAuth protects the legacy machine-to-machine endpoints. Empty
// credentials never authorize a request.
func APIBasicAuth(username, password string) anduin.Middleware {
	expectedUser, expectedPass := sha256.Sum256([]byte(username)), sha256.Sum256([]byte(password))
	return func(next anduin.Handler) anduin.Handler {
		return func(ctx context.Context) error {
			req := anduin.Request(ctx)
			if req != nil {
				malformed := false
				count := 0
				for _, header := range req.Header.Values("Authorization") {
					if count >= 1 {
						break
					}
					if len(header) <= 6 || !strings.EqualFold(header[:6], "Basic ") {
						continue
					}
					count++
					decoded, err := base64.StdEncoding.DecodeString(header[6:])
					if err != nil {
						malformed = true
						continue
					}
					user, pass, ok := strings.Cut(string(decoded), ":")
					userHash, passHash := sha256.Sum256([]byte(user)), sha256.Sum256([]byte(pass))
					if ok && username != "" && password != "" &&
						subtle.ConstantTimeCompare(userHash[:], expectedUser[:]) == 1 &&
						subtle.ConstantTimeCompare(passHash[:], expectedPass[:]) == 1 {
						return next(ctx)
					}
				}
				if malformed {
					return anduin.JSON(ctx, http.StatusBadRequest, map[string]string{"message": "Bad Request"})
				}
			}
			anduin.SetHeader(ctx, "WWW-Authenticate", `Basic realm="mbvlabs API"`)
			return anduin.JSON(ctx, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
		}
	}
}
