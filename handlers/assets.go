package handlers

import (
	"fmt"
	"io/fs"
	"net/http"

	"mbvlabs/assets"
	"mbvlabs/config"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
)

// AssetsHandler serves the embedded Vite build files.
type AssetsHandler struct {
	dist anduin.Handler
}

func NewAssetsHandler(appCfg config.App) (AssetsHandler, error) {
	dist, err := fs.Sub(assets.Files, "dist")
	if err != nil {
		return AssetsHandler{}, fmt.Errorf("assets: dist: %w", err)
	}

	distHandler, err := anduin.StaticFS(dist, anduin.Immutable(appCfg.IsProduction()))
	if err != nil {
		return AssetsHandler{}, fmt.Errorf("assets: static files: %w", err)
	}

	return AssetsHandler{
		dist: distHandler,
	}, nil
}

func (a AssetsHandler) Register(rt *anduin.Router) {
	// Vite build files are public: a secondary-host page loads them without the
	// visitor's session, so this CORS response is wildcard and uncredentialed.
	rt.Route(routes.ViteBuild).Get(a.dist, anduin.CORS(anduin.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead},
		AllowCredentials: false,
	}))
}
