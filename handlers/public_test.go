package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/cache"
	"github.com/mbvlabs/andurel/pkg/kiks"
	"mbvlabs/config"
	"mbvlabs/middleware"
	"mbvlabs/resources"

	"github.com/mbvlabs/andurel/pkg/inertia"
)

type publicTestSSR struct{}

func (publicTestSSR) Render(_ context.Context, page inertia.Page) (*inertia.SSRResponse, error) {
	names := map[string]string{
		"Public/About":   "A fractional tech lead who still ships",
		"Public/Contact": `<input name="_csrf">`,
	}
	body := names[page.Component]
	auth, _ := json.Marshal(page.Props["auth"])
	if strings.Contains(string(auth), `"isAuthenticated":true`) {
		body += `<a href="/admin">Admin</a>`
	}
	if flash, _ := page.Flash.([]kiks.FlashMessage); len(flash) > 0 {
		body += flash[0].Message
	}
	return &inertia.SSRResponse{Body: `<div data-server-rendered="true" data-page="{}"><main>` + body + `</main></div>`}, nil
}

func setupPublicTest(t *testing.T, mw ...anduin.Middleware) (PublicHandler, *anduin.Router) {
	t.Helper()
	renderer, err := inertia.NewRenderer(inertia.RendererConfig{ContainerID: "app"},
		inertia.WithRoot(resources.Root), inertia.WithSSRRenderer(publicTestSSR{}),
		inertia.WithFlashProvider(kiks.Flashes),
		inertia.WithSharedProvider(func(ctx context.Context) (inertia.Props, error) {
			session, _ := kiks.Get[*middleware.Session](ctx)
			return inertia.Props{"auth": map[string]bool{"isAuthenticated": session != nil && session.IsAuthenticated}}, nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	router := anduin.NewRouter()
	if len(mw) > 0 {
		router.Use(mw...)
	}
	router.Use(renderer.Middleware())
	return PublicHandler{cfg: config.App{SiteName: "MBV Labs", Domain: "example.com", Protocol: "https"}, renderer: renderer}, router
}

func TestPublicLayoutSessionAndFlash(t *testing.T) {
	keys := kiks.Keys{Authentication: []byte("01234567890123456789012345678901"), Encryption: []byte("01234567890123456789012345678901")}
	store, err := kiks.NewCookieStore(keys)
	if err != nil {
		t.Fatal(err)
	}
	jar, err := kiks.NewJar(keys, store, kiks.NewSession[*middleware.Session]("session"))
	if err != nil {
		t.Fatal(err)
	}
	h, router := setupPublicTest(t, jar.Middleware())
	router.Get("/sign-in-test", func(ctx context.Context) error {
		if err := kiks.Set(ctx, &middleware.Session{IsAuthenticated: true}); err != nil {
			return err
		}
		if err := kiks.AddFlash(ctx, kiks.FlashSuccess, "Successfully logged in!"); err != nil {
			return err
		}
		return anduin.Redirect(ctx, "/about", 303)
	})
	router.Get("/about", h.about)
	server := router.Handler()
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/sign-in-test", nil))
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	for _, cookie := range rec.Result().Cookies() {
		req.AddCookie(cookie)
	}
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Successfully logged in!") || !strings.Contains(rec.Body.String(), `href="/admin"`) {
		t.Fatalf("status=%d: session or flash missing in %s", rec.Code, rec.Body.String())
	}
}

func TestPublicRegistration(t *testing.T) {
	store, err := cache.New(cache.Config{})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	h, router := setupPublicTest(t)
	h.cache = store
	h.Register(router)
	server := router.Handler()
	for _, tc := range []struct {
		path               string
		status             int
		location, contains string
	}{
		{"/about", 200, "", "A fractional tech lead who still ships"},
		{"/contact", 200, "", `name="_csrf"`},
		{"/about-me", 301, "/about", ""},
		{"/blog-posts", 301, "/blog", ""},
		{"/blog-posts/example", 301, "/blog/example", ""},
		{"/project-inquiries", 301, "/contact", ""},
		{"/assets/css/style.css", 404, "", ""},
		{"/assets/fonts/dm-mono-v16-latin-400.woff2", 404, "", ""},
	} {
		t.Run(tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if rec.Code != tc.status || rec.Header().Get("Location") != tc.location || !strings.Contains(rec.Body.String(), tc.contains) {
				t.Fatalf("status=%d location=%q", rec.Code, rec.Header().Get("Location"))
			}
		})
	}
}

func TestPublicStaticAndContentNegotiation(t *testing.T) {
	h, router := setupPublicTest(t)
	router.Get("/about", h.about)
	router.Get("/robots.txt", h.robots)
	server := router.Handler()

	for _, tc := range []struct{ path, accept, contentType, contains string }{
		{"/about", "", "text/html", "A fractional tech lead who still ships"},
		{"/about", "text/markdown", "text/markdown", "A fractional tech lead who still ships"},
		{"/robots.txt", "", "text/plain", "Sitemap: https://example.com/sitemap.xml"},
	} {
		t.Run(tc.path+tc.accept, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			req.Header.Set("Accept", tc.accept)
			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)
			if rec.Code != 200 || !strings.HasPrefix(rec.Header().Get("Content-Type"), tc.contentType) || !strings.Contains(rec.Body.String(), tc.contains) {
				t.Fatalf("status=%d content-type=%q body prefix=%q", rec.Code, rec.Header().Get("Content-Type"), rec.Body.String()[:min(200, rec.Body.Len())])
			}
		})
	}
}

func TestPublicRedirectAndMarkdownQuality(t *testing.T) {
	h, router := setupPublicTest(t)
	router.Get("/about", h.about)
	router.Get("/about-me", h.redirect("/about"))
	server := router.Handler()
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	req.Header.Set("Accept", "text/markdown;q=0, text/html")
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	if !strings.HasPrefix(rec.Header().Get("Content-Type"), "text/html") || !strings.Contains(rec.Header().Get("Vary"), "Accept") {
		t.Fatalf("unexpected headers: %v", rec.Header())
	}
	rec = httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/about-me", nil))
	if rec.Code != http.StatusMovedPermanently || rec.Header().Get("Location") != "/about" {
		t.Fatalf("status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
}
