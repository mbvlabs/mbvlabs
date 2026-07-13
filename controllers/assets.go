package controllers

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"mbvlabs/assets"
	"mbvlabs/config"
	"mbvlabs/internal/routing"
	"mbvlabs/internal/server"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/routes"

	"github.com/labstack/echo/v5"
)

const (
	threeMonthsCache = "7776000"
	sitemapCacheKey  = "assets:sitemap"
)

type Assets struct {
	cache *Cache[string]
	db    storage.Pool
}

func NewAssets(cache *Cache[string], db storage.Pool) Assets {
	return Assets{cache: cache, db: db}
}

func (a Assets) RegisterRoutes(r *router.Router) error {
	errs := []error{}

	_, err := r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Robots.Path(),
		Name:    routes.Robots.Name(),
		Handler: a.Robots,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Sitemap.Path(),
		Name:    routes.Sitemap.Name(),
		Handler: a.Sitemap,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Stylesheet.Path(),
		Name:    routes.Stylesheet.Name(),
		Handler: a.Stylesheet,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Scripts.Path(),
		Name:    routes.Scripts.Name(),
		Handler: a.Scripts,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Script.Path(),
		Name:    routes.Script.Name(),
		Handler: a.Script,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Fonts.Path(),
		Name:    routes.Fonts.Name(),
		Handler: a.Font,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.ViteBuild.Path(),
		Name:    routes.ViteBuild.Name(),
		Handler: a.ViteBuild,
	})
	if err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (a Assets) enableCaching(etx *echo.Context, content []byte) *echo.Context {
	if config.Env == server.ProdEnvironment {
		//nolint:gosec //only needed for browser caching
		hash := md5.Sum(content)
		etag := fmt.Sprintf(`"%x-%x"`, hash, len(content))

		if match := etx.Request().Header.Get("If-None-Match"); match == etag {
			etx.Response().
				Header().
				Set("Cache-Control", fmt.Sprintf("public, max-age=%s, immutable", threeMonthsCache))
			etx.Response().
				Header().
				Set("ETag", etag)
			etx.NoContent(http.StatusNotModified)
			return etx
		}

		etx.Response().
			Header().
			Set("Cache-Control", fmt.Sprintf("public, max-age=%s, immutable", threeMonthsCache))
		etx.Response().
			Header().
			Set("Vary", "Accept-Encoding")
		etx.Response().
			Header().
			Set("ETag", etag)
	}

	return etx
}

func createRobotsTxt() (string, error) {
	return fmt.Sprintf(
		"User-agent: *\nAllow: /\nSitemap: %s\n",
		absoluteURL(routes.Sitemap.URL()),
	), nil
}

func (a Assets) Robots(etx *echo.Context) error {
	cacheKey := "assets:robots"

	robotsTxt, err := a.cache.Get(cacheKey, func() (string, error) {
		return createRobotsTxt()
	})
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to get robots.txt from cache",
			"error", err,
		)
		result, _ := createRobotsTxt()
		return etx.String(http.StatusOK, result)
	}

	return etx.String(http.StatusOK, robotsTxt)
}

func (a Assets) Sitemap(etx *echo.Context) error {
	sitemap, err := a.cache.Get(sitemapCacheKey, func() (string, error) {
		return createSitemap(etx.Request().Context(), a.db.Executor(), []routing.Route{})
	})
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to get sitemap from cache",
			"error", err,
		)

		result, err := createSitemap(etx.Request().Context(), a.db.Executor(), []routing.Route{})
		if err != nil {
			return err
		}

		return etx.Blob(http.StatusOK, "application/xml", []byte(result))
	}

	return etx.Blob(http.StatusOK, "application/xml", []byte(sitemap))
}

type URL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	ChangeFreq string   `xml:"changefreq"`
	LastMod    string   `xml:"lastmod,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URL     []URL    `xml:"url"`
}

func createSitemap(ctx context.Context, db storage.Executor, extraRoutes []routing.Route) (string, error) {
	var urls []URL

	urls = append(urls, URL{
		Loc:        absoluteURL(routes.HomePage.URL()),
		ChangeFreq: "monthly",
		LastMod:    "2024-10-22T09:43:09+00:00",
		Priority:   "1",
	})
	for _, route := range []routing.Route{
		routes.ServiceOfferingIndex,
		routes.WorkIndex,
		routes.ProjectIndex,
		routes.BlogPostIndex,
		routes.AboutMe,
		routes.ProjectInquiryIndex,
	} {
		urls = append(urls, URL{
			Loc:        absoluteURL(route.URL()),
			ChangeFreq: "monthly",
			Priority:   "0.8",
		})
	}
	for _, route := range extraRoutes {
		urls = append(urls, URL{
			Loc:        absoluteURL(route.URL()),
			ChangeFreq: "monthly",
			Priority:   "0.6",
		})
	}

	dynamicURLs, err := publishedSitemapURLs(ctx, db)
	if err != nil {
		return "", err
	}
	urls = append(urls, dynamicURLs...)

	sitemap := Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL:   urls,
	}

	xmlBytes, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(xmlBytes), nil
}

func publishedSitemapURLs(ctx context.Context, db storage.Executor) ([]URL, error) {
	var urls []URL

	works, err := models.Work.AllPublished(ctx, db)
	if err != nil {
		return nil, err
	}
	for _, work := range works {
		urls = append(urls, URL{
			Loc:        absoluteURL(routes.WorkShow.URL(work.Slug)),
			ChangeFreq: "monthly",
			LastMod:    lastMod(work.UpdatedAt, work.PublishedAt),
			Priority:   "0.6",
		})
	}

	projects, err := models.Project.AllPublished(ctx, db)
	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		urls = append(urls, URL{
			Loc:        absoluteURL(routes.ProjectShow.URL(project.Slug)),
			ChangeFreq: "monthly",
			LastMod:    lastMod(project.UpdatedAt, project.PublishedAt),
			Priority:   "0.6",
		})
	}

	posts, err := models.BlogPost.AllPublished(ctx, db)
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		urls = append(urls, URL{
			Loc:        absoluteURL(routes.BlogPostShow.URL(post.Slug)),
			ChangeFreq: "monthly",
			LastMod:    lastMod(post.UpdatedAt, post.PublishedAt),
			Priority:   "0.6",
		})
	}

	return urls, nil
}

func absoluteURL(path string) string {
	return strings.TrimRight(config.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func lastMod(updatedAt time.Time, publishedAt sql.NullTime) string {
	if !updatedAt.IsZero() {
		return updatedAt.Format(time.RFC3339)
	}
	if publishedAt.Valid {
		return publishedAt.Time.Format(time.RFC3339)
	}
	return ""
}

func (a Assets) Stylesheet(etx *echo.Context) error {
	stylesheet, err := assets.Files.ReadFile(
		"css/style.css",
	)
	if err != nil {
		return err
	}

	etx = a.enableCaching(etx, stylesheet)
	return etx.Blob(http.StatusOK, "text/css", stylesheet)
}

func (a Assets) Scripts(etx *echo.Context) error {
	stylesheet, err := assets.Files.ReadFile(
		"js/scripts.js",
	)
	if err != nil {
		return err
	}

	etx = a.enableCaching(etx, stylesheet)
	return etx.Blob(http.StatusOK, "application/javascript", stylesheet)
}

func (a Assets) Script(etx *echo.Context) error {
	param := etx.Param("file")
	stylesheet, err := assets.Files.ReadFile(
		fmt.Sprintf("js/%s", param),
	)
	if err != nil {
		return err
	}

	etx = a.enableCaching(etx, stylesheet)
	return etx.Blob(http.StatusOK, "application/javascript", stylesheet)
}

func (a Assets) Font(etx *echo.Context) error {
	param := path.Clean(strings.TrimPrefix(etx.Param("*"), "/"))
	if param == "." || param == "" || strings.HasPrefix(param, "../") || !strings.HasSuffix(param, ".woff2") {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	font, err := assets.Files.ReadFile("fonts/" + param)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	etx = a.enableCaching(etx, font)
	return etx.Blob(http.StatusOK, "font/woff2", font)
}

func contentTypeByExt(name string) string {
	switch {
	case strings.HasSuffix(name, ".js"):
		return "application/javascript"
	case strings.HasSuffix(name, ".css"):
		return "text/css"
	case strings.HasSuffix(name, ".json"):
		return "application/json"
	case strings.HasSuffix(name, ".svg"):
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

func (a Assets) ViteBuild(etx *echo.Context) error {
	param := strings.TrimPrefix(etx.Param("*"), "/")
	param = path.Clean(param)
	if param == "." || param == "" || strings.HasPrefix(param, "../") {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	data, err := assets.Files.ReadFile(fmt.Sprintf("dist/%s", param))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	etx = a.enableCaching(etx, data)
	return etx.Blob(http.StatusOK, contentTypeByExt(param), data)
}
