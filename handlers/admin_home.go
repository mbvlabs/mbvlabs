package handlers

import (
	"context"

	"mbvlabs/domains"
	"mbvlabs/middleware"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
)

type AdminHomeHandler struct {
	works     *domains.Works
	projects  *domains.Projects
	posts     *domains.BlogPosts
	inquiries *domains.ProjectInquiries
	renderer  *inertia.Renderer
}

func NewAdminHomeHandler(works *domains.Works, projects *domains.Projects, posts *domains.BlogPosts, inquiries *domains.ProjectInquiries, renderer *inertia.Renderer) AdminHomeHandler {
	return AdminHomeHandler{works, projects, posts, inquiries, renderer}
}

func (h AdminHomeHandler) Register(rt *anduin.Router, middleware ...anduin.Middleware) {
	rt.Route(routes.AdminHomePage).Get(h.home, middleware...)
}

func (h AdminHomeHandler) home(ctx context.Context) error {
	const latest = 5
	inquiries, err := h.inquiries.Latest(ctx, latest)
	if err != nil {
		return err
	}
	works, err := h.works.Latest(ctx, latest)
	if err != nil {
		return err
	}
	posts, err := h.posts.Latest(ctx, latest)
	if err != nil {
		return err
	}
	projects, err := h.projects.Latest(ctx, latest)
	if err != nil {
		return err
	}
	blogData := make([]adminBlogData, 0, len(posts))
	for _, post := range posts {
		item, err := adminBlogItem(post)
		if err != nil {
			return err
		}
		blogData = append(blogData, item)
	}
	viewer, _ := middleware.ViewerFrom(ctx)
	return h.renderer.Page(ctx, "Admin/Home", inertia.Props{
		"appName":               "mbvlabs",
		"projectInquiries":      adminHomeItems(inquiries, adminInquiryItem),
		"works":                 adminHomeItems(works, NewAdminWorkData),
		"blogPosts":             blogData,
		"projects":              adminHomeItems(projects, NewAdminProjectData),
		"canManageContentToken": viewer.IsAdmin,
	}).Render()
}

func adminHomeItems[T, U any](rows []T, convert func(T) U) []U {
	list := make([]U, 0, len(rows))
	for _, row := range rows {
		list = append(list, convert(row))
	}
	return list
}
