package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"mbvlabs/assets"
	"mbvlabs/config"
	"mbvlabs/domains"
	"mbvlabs/routes"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/cache"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
	"github.com/mbvlabs/andurel/pkg/routing"
)

type PublicHandler struct {
	cfg       config.App
	works     *domains.Works
	projects  *domains.Projects
	posts     *domains.BlogPosts
	inquiries *domains.ProjectInquiries
	cache     cache.Store
	renderer  *inertia.Renderer
}

func NewPublicHandler(cfg config.App, works *domains.Works, projects *domains.Projects, posts *domains.BlogPosts, inquiries *domains.ProjectInquiries, store cache.Store, renderer *inertia.Renderer) PublicHandler {
	return PublicHandler{cfg, works, projects, posts, inquiries, store, renderer}
}

func (h PublicHandler) Register(rt *anduin.Router) {
	rt.Route(routes.HomePage).Get(h.home)
	rt.Route(routes.AboutMe).Get(h.about)
	rt.Route(routes.LegacyAboutMe).Get(h.redirect(routes.AboutMe.URL()))
	rt.Route(routes.ServiceOfferingIndex).Get(h.services)
	rt.Route(routes.WorkIndex).Get(h.workIndex)
	rt.Route(routes.WorkShow).Get(h.workShow)
	rt.Route(routes.ProjectIndex).Get(h.projectIndex)
	rt.Route(routes.ProjectShow).Get(h.projectShow)
	rt.Route(routes.BlogPostIndex).Get(h.blogIndex)
	rt.Route(routes.BlogPostFeed).Get(h.feed)
	rt.Route(routes.BlogPostShow).Get(h.blogShow)
	rt.Route(routes.LegacyBlogPostIndex).Get(h.redirect(routes.BlogPostIndex.URL()))
	rt.Route(routes.LegacyBlogPostShow).Get(func(ctx context.Context) error {
		return anduin.Redirect(ctx, routes.BlogPostShow.URL(anduin.Request(ctx).PathValue("slug")), http.StatusMovedPermanently)
	})
	rt.Route(routes.ProjectInquiryIndex).Get(h.contact)
	limited := anduin.RateLimit(h.cache, 3, 10*time.Minute, func(ctx context.Context) error {
		return anduin.Text(ctx, http.StatusTooManyRequests, routes.ProjectInquiryIndex.URL())
	})
	rt.Route(routes.ProjectInquiryCreate).Post(h.createInquiry, limited)
	rt.Route(routes.LegacyProjectInquiryIndex).Get(h.redirect(routes.ProjectInquiryIndex.URL()))
	rt.Route(routes.LegacyProjectInquiryCreate).Post(h.createInquiry, limited)
	rt.Route(routes.Robots).Get(h.robots)
	rt.Route(routes.Sitemap).Get(h.sitemap)
	rt.Route(routes.IndexNowKey).Get(func(ctx context.Context) error { return anduin.Text(ctx, http.StatusOK, routes.IndexNowKeyValue) })
	rt.Route(routes.Stylesheet).Get(h.static("css/style.css", "text/css; charset=utf-8"))
	rt.Route(routes.PublicFonts).Get(h.font)
}

func (h PublicHandler) render(ctx context.Context, status int, component string, seo HeadData, props inertia.Props) error {
	if props == nil {
		props = inertia.Props{}
	}
	props["seo"] = seo
	props["publicSite"] = inertia.Props{"production": h.cfg.IsProduction(), "year": time.Now().UTC().Year()}
	anduin.AddHeader(ctx, "Vary", "Accept")
	if anduin.Header(ctx, inertia.HeaderInertia) != "" || !acceptsPublicMarkdown(anduin.Request(ctx).Header.Values("Accept")) {
		return h.renderer.Page(ctx, component, props).SSR().Status(status).Render()
	}

	// Negotiate Markdown from the same server-rendered React document. The
	// document isn't committed until conversion succeeds.
	rec := httptest.NewRecorder()
	err := func() error {
		original := anduin.ResponseWriter(ctx)
		defer anduin.ReplaceResponseWriter(ctx, original)
		return h.renderer.Page(anduin.ReplaceResponseWriter(ctx, rec), component, props).SSR().Status(status).Render()
	}()
	if err != nil {
		return err
	}
	if !strings.Contains(rec.Body.String(), `data-server-rendered="true"`) {
		return anduin.Text(ctx, http.StatusServiceUnavailable, "Server-rendered content is unavailable")
	}
	for _, vary := range rec.Header().Values("Vary") {
		anduin.AddHeader(ctx, "Vary", vary)
	}
	text, err := htmltomarkdown.ConvertString(rec.Body.String())
	if err != nil {
		return err
	}
	return anduin.Blob(ctx, status, "text/markdown; charset=utf-8", []byte(text))
}

func acceptsPublicMarkdown(accept []string) bool {
	for part := range strings.SplitSeq(strings.Join(accept, ","), ",") {
		kind, params, err := mime.ParseMediaType(strings.TrimSpace(part))
		q := params["q"]
		quality, qerr := strconv.ParseFloat(q, 64)
		if err == nil && strings.EqualFold(kind, "text/markdown") && (q == "" || qerr == nil && quality > 0) {
			return true
		}
	}
	return false
}

func (h PublicHandler) redirect(path string) anduin.Handler {
	return func(ctx context.Context) error { return anduin.Redirect(ctx, path, http.StatusMovedPermanently) }
}
func (h PublicHandler) home(ctx context.Context) error {
	works, err := h.works.FeaturedPublished(ctx, 3)
	if err != nil {
		return err
	}
	projects, err := h.projects.FeaturedPublished(ctx, 3)
	if err != nil {
		return err
	}
	anduin.SetHeader(ctx, "Link", "<"+routes.Sitemap.URL()+">; rel=\"describedby\"; type=\"application/xml\"")
	return h.render(ctx, 200, "Public/Home", h.metadata("Fractional Tech Lead",
		"MBV Labs helps founders and lean teams turn ambiguous software goals into shipped production systems.",
		routes.HomePage.URL(), "", "", HomePageSchema()), inertia.Props{"works": works, "projects": projects})
}
func (h PublicHandler) about(ctx context.Context) error {
	return h.render(ctx, 200, "Public/About", h.metadata("About",
		"MBV Labs is a fractional tech lead and engineering studio for founders and lean teams that need senior judgment and hands-on execution.",
		routes.AboutMe.URL(), "", "", PageSchema("AboutPage", []SchemaBreadcrumb{
			{Name: "Home", Path: routes.HomePage.URL()}, {Name: "About", Path: routes.AboutMe.URL()},
		})), nil)
}
func (h PublicHandler) services(ctx context.Context) error {
	return h.render(ctx, 200, "Public/Services", h.metadata("Services",
		"Project-based fractional tech lead and senior software engineering for backend systems, Go consulting, AI workflows, and technical leadership.",
		routes.ServiceOfferingIndex.URL(), "", "", ServiceOfferingIndexSchema()), nil)
}
func (h PublicHandler) contact(ctx context.Context) error {
	return h.render(ctx, 200, "Public/Contact", h.metadata("Contact",
		"Tell MBV Labs about your software project, product engineering needs, or fractional tech lead support.",
		routes.ProjectInquiryIndex.URL(), "", "", ContactPageSchema()),
		inertia.Props{"csrfToken": anduin.CSRFToken(ctx)})
}
func (h PublicHandler) notFound(ctx context.Context) error {
	return h.renderError(ctx, http.StatusNotFound)
}
func (h PublicHandler) renderError(ctx context.Context, status int) error {
	return h.render(ctx, status, "Public/PublicError", h.metadata("", "", "/", "", ""), inertia.Props{"status": status})
}

func pagination(ctx context.Context) (int64, int64) {
	q := anduin.Request(ctx).URL.Query()
	page, _ := strconv.ParseInt(q.Get("page"), 10, 64)
	per, _ := strconv.ParseInt(q.Get("per_page"), 10, 64)
	if page < 1 {
		page = 1
	}
	if per < 1 || per > 100 {
		per = 25
	}
	return page, per
}
func (h PublicHandler) workIndex(ctx context.Context) error {
	page, per := pagination(ctx)
	result, err := h.works.PaginatePublished(ctx, page, per)
	if err != nil {
		return err
	}
	return h.render(ctx, 200, "Public/WorkIndex", h.metadata("Client Work and Software Engineering Case Studies",
		"Client work and full-time engineering experience from MBV Labs, a fractional tech lead for founders and lean teams.",
		routes.WorkIndex.URL(), "", "", WorkIndexSchema(result.Works)), inertia.Props{"items": result.Works})
}
func (h PublicHandler) workShow(ctx context.Context) error {
	item, err := h.works.FindBySlug(ctx, anduin.Request(ctx).PathValue("slug"))
	if errors.Is(err, domains.ErrNotFound) {
		return h.notFound(ctx)
	}
	if err != nil {
		return err
	}
	if item.Status != domains.Published {
		return h.notFound(ctx)
	}
	return h.render(ctx, 200, "Public/WorkShow", h.metadata(item.Title, item.Summary,
		routes.WorkShow.URL(item.Slug), item.CoverImageURL.String, "article", WorkShowSchema(item)),
		inertia.Props{"item": item, "markdownHTML": renderPublicMarkdown(item.Content)})
}
func (h PublicHandler) projectIndex(ctx context.Context) error {
	page, per := pagination(ctx)
	result, err := h.projects.PaginatePublished(ctx, page, per)
	if err != nil {
		return err
	}
	return h.render(ctx, 200, "Public/ProjectIndex", h.metadata("Projects",
		"Explore open-source and commercial projects built by MBV Labs, including developer tools, platforms, infrastructure, and product experiments.",
		routes.ProjectIndex.URL(), "", "", ProjectIndexSchema(result.Projects)), inertia.Props{"items": result.Projects})
}
func (h PublicHandler) projectShow(ctx context.Context) error {
	item, err := h.projects.FindBySlug(ctx, anduin.Request(ctx).PathValue("slug"))
	if errors.Is(err, domains.ErrNotFound) {
		return h.notFound(ctx)
	}
	if err != nil {
		return err
	}
	if !item.PublishedAt.Valid {
		return h.notFound(ctx)
	}
	return h.render(ctx, 200, "Public/ProjectShow", h.metadata(item.Name, item.Tagline,
		routes.ProjectShow.URL(item.Slug), item.ImageURL.String, "", ProjectShowSchema(item, publicStringList(item.Technologies))),
		inertia.Props{"item": item, "markdownHTML": renderPublicMarkdown(item.Description.String)})
}
func (h PublicHandler) blogIndex(ctx context.Context) error {
	page, per := pagination(ctx)
	result, err := h.posts.PaginatePublished(ctx, page, per)
	if err != nil {
		return err
	}
	return h.render(ctx, 200, "Public/BlogIndex", h.metadata("Blog",
		"Pragmatic engineering notes from MBV Labs on product delivery, software architecture, and building maintainable systems.",
		routes.BlogPostIndex.URL(), "", "", BlogPostIndexSchema(result.BlogPosts)), inertia.Props{"items": result.BlogPosts})
}
func (h PublicHandler) blogShow(ctx context.Context) error {
	item, err := h.posts.FindBySlug(ctx, anduin.Request(ctx).PathValue("slug"))
	if errors.Is(err, domains.ErrNotFound) {
		return h.notFound(ctx)
	}
	if err != nil {
		return err
	}
	if item.Status != "published" {
		return h.notFound(ctx)
	}
	return h.render(ctx, 200, "Public/BlogShow", h.metadata(item.Title, item.Excerpt.String,
		routes.BlogPostShow.URL(item.Slug), item.CoverImageUrl.String, "article", BlogPostShowSchema(item, publicStringList(item.Tags))),
		inertia.Props{"item": item, "markdownHTML": renderPublicMarkdown(item.Body)})
}
func optionalPublicString(s string) domains.OptionalString {
	s = strings.TrimSpace(s)
	return domains.OptionalString{String: s, Valid: s != ""}
}

type publicInquiryPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Company     string `json:"company"`
	Role        string `json:"role"`
	ProjectType string `json:"projectType"`
	Timeline    string `json:"timeline"`
	Message     string `json:"message"`
}

func (h PublicHandler) createInquiry(ctx context.Context) error {
	r := anduin.Request(ctx)
	var payload publicInquiryPayload
	contentType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			return h.renderError(ctx, http.StatusBadRequest)
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return h.renderError(ctx, http.StatusBadRequest)
		}
		payload = publicInquiryPayload{
			Name: r.FormValue("name"), Email: r.FormValue("email"), Company: r.FormValue("company"),
			Role: r.FormValue("role"), ProjectType: r.FormValue("projectType"),
			Timeline: r.FormValue("timeline"), Message: r.FormValue("message"),
		}
	}
	email := strings.TrimSpace(payload.Email)
	banned, err := h.inquiries.IsDomainBanned(ctx, email)
	if err != nil {
		return err
	}
	if !banned {
		metadata, err := json.Marshal(map[string]string{"referer": r.Referer(), "user_agent": r.UserAgent(), "remote_addr": anduin.ClientIP(ctx)})
		if err != nil {
			return err
		}
		_, err = h.inquiries.Create(ctx, domains.CreateProjectInquiryData{
			Name: strings.TrimSpace(payload.Name), Email: email,
			Company: optionalPublicString(payload.Company), Role: optionalPublicString(payload.Role),
			ProjectType: optionalPublicString(payload.ProjectType), Timeline: optionalPublicString(payload.Timeline),
			Message: strings.TrimSpace(payload.Message), Source: domains.OptionalString{String: "contact", Valid: true},
			Status: "new", Metadata: metadata,
		})
		if err != nil {
			if err := kiks.AddFlash(ctx, kiks.FlashError, "Please check the form and try again."); err != nil {
				return err
			}
			return anduin.Redirect(ctx, routes.ProjectInquiryIndex.URL(), http.StatusSeeOther)
		}
	}
	if err := kiks.AddFlash(ctx, kiks.FlashSuccess, "Thanks. I'll get back to you soon."); err != nil {
		return err
	}
	return anduin.Redirect(ctx, routes.ProjectInquiryIndex.URL(), http.StatusSeeOther)
}
func (h PublicHandler) static(name, contentType string) anduin.Handler {
	return func(ctx context.Context) error {
		data, err := assets.Files.ReadFile(name)
		if err != nil {
			return err
		}
		return anduin.Blob(ctx, 200, contentType, data)
	}
}
func (h PublicHandler) font(ctx context.Context) error {
	name := anduin.Request(ctx).PathValue(routing.WildcardParam)
	if name == "" || strings.Contains(name, "/") || !strings.HasSuffix(name, ".woff2") {
		return anduin.NoContent(ctx, 404)
	}
	data, err := assets.Files.ReadFile("fonts/" + name)
	if err != nil {
		return anduin.NoContent(ctx, 404)
	}
	return anduin.Blob(ctx, 200, "font/woff2", data)
}
