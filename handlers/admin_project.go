package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"mbvlabs/domains"
	"mbvlabs/payloads"
	"mbvlabs/routes"
	"mbvlabs/services"

	"github.com/gosimple/slug"
	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
)

type adminProjectStore interface {
	Paginate(context.Context, int64, int64) (domains.PaginatedProjects, error)
	Find(context.Context, int32) (domains.ProjectEntity, error)
	Destroy(context.Context, int32) error
}
type adminProjectWriter interface {
	CreateProject(context.Context, domains.CreateProjectData) (domains.ProjectEntity, error)
	UpdateProject(context.Context, domains.UpdateProjectData) (domains.ProjectEntity, error)
}
type AdminProjectHandler struct {
	projects   adminProjectStore
	writer     adminProjectWriter
	renderer   *inertia.Renderer
	invalidate func()
}

func NewAdminProjectHandler(projects *domains.Projects, writer services.Projects, renderer *inertia.Renderer, invalidate func()) AdminProjectHandler {
	if invalidate == nil {
		invalidate = func() {}
	}
	return AdminProjectHandler{projects, writer, renderer, invalidate}
}
func (h AdminProjectHandler) Register(rt *anduin.Router, middleware ...anduin.Middleware) {
	rt.Route(routes.AdminProjectIndex).Get(h.Index, middleware...)
	rt.Route(routes.AdminProjectNew).Get(h.New, middleware...)
	rt.Route(routes.AdminProjectShow).Get(h.Show, middleware...)
	rt.Route(routes.AdminProjectCreate).Post(h.Create, middleware...)
	rt.Route(routes.AdminProjectEdit).Get(h.Edit, middleware...)
	rt.Route(routes.AdminProjectUpdate).Put(h.Update, middleware...)
	rt.Route(routes.AdminProjectDestroy).Delete(h.Destroy, middleware...)
}

func NewAdminProjectData(e domains.ProjectEntity) payloads.AdminProjectData {
	var technologies []string
	tech := string(e.Technologies)
	if json.Unmarshal(e.Technologies, &technologies) == nil {
		tech = strings.Join(technologies, ", ")
	}
	return payloads.AdminProjectData{
		ID: e.ID, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt, Name: e.Name, Slug: e.Slug,
		Tagline: e.Tagline, Description: e.Description.String, ProjectType: e.ProjectType,
		RepositoryUrl: e.RepositoryURL.String, LiveUrl: e.LiveURL.String, ImageUrl: e.ImageURL.String,
		Technologies: tech, StartedAt: adminContentDate(e.StartedAt), LaunchedAt: adminContentDate(e.LaunchedAt),
		PublishedAt: adminContentDate(e.PublishedAt), IsFeatured: e.IsFeatured,
	}
}
func (h AdminProjectHandler) Index(ctx context.Context) error {
	page, size := adminContentPagination(ctx)
	list, err := h.projects.Paginate(ctx, page, size)
	if err != nil {
		return err
	}
	items := make([]payloads.AdminProjectData, 0, len(list.Projects))
	for _, item := range list.Projects {
		items = append(items, NewAdminProjectData(item))
	}
	return h.renderer.Page(ctx, "Admin/Project/Index", inertia.FromStruct(payloads.AdminProjectIndexProps{Items: items})).Render()
}
func (h AdminProjectHandler) New(ctx context.Context) error {
	return h.renderer.Page(ctx, "Admin/Project/Create", inertia.Props{}).Render()
}
func (h AdminProjectHandler) item(ctx context.Context, component string) error {
	id, err := adminContentID(ctx)
	if err != nil {
		return err
	}
	item, err := h.projects.Find(ctx, id)
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	if err != nil {
		return err
	}
	return h.renderer.Page(ctx, component, inertia.FromStruct(payloads.AdminProjectItemProps{Item: NewAdminProjectData(item)})).Render()
}
func (h AdminProjectHandler) Show(ctx context.Context) error {
	return h.item(ctx, "Admin/Project/Show")
}
func (h AdminProjectHandler) Edit(ctx context.Context) error {
	return h.item(ctx, "Admin/Project/Edit")
}

func adminProjectPayloadItem(p payloads.AdminProjectPayload, id int32) payloads.AdminProjectData {
	return payloads.AdminProjectData{ID: id, Name: p.Name, Slug: p.Slug, Tagline: p.Tagline, Description: p.Description,
		ProjectType: p.ProjectType, RepositoryUrl: p.RepositoryUrl, LiveUrl: p.LiveUrl, ImageUrl: p.ImageUrl,
		Technologies: p.Technologies, StartedAt: p.StartedAt, LaunchedAt: p.LaunchedAt,
		PublishedAt: p.PublishedAt, IsFeatured: p.IsFeatured}
}
func (h AdminProjectHandler) save(ctx context.Context, update bool) error {
	var id int32
	var err error
	if update {
		id, err = adminContentID(ctx)
		if err != nil {
			return err
		}
	}
	var p payloads.AdminProjectPayload
	if err := anduin.DecodeJSON(ctx, &p); err != nil {
		return anduin.BadRequest(err)
	}
	component, props := "Admin/Project/Create", inertia.Props{}
	if update {
		component = "Admin/Project/Edit"
		props["item"] = adminProjectPayloadItem(p, id)
	}
	errs := map[string]string{}
	started := adminContentParseDate(p.StartedAt, "startedAt", errs)
	launched := adminContentParseDate(p.LaunchedAt, "launchedAt", errs)
	published := adminContentParseDate(p.PublishedAt, "publishedAt", errs)
	if len(errs) > 0 {
		return h.renderer.Page(ctx, component, props).ValidationErrors(errs).Render()
	}
	slugSource := p.Slug
	if slugSource == "" {
		slugSource = p.Name
	}
	data := domains.CreateProjectData{
		Name: p.Name, Slug: slug.Make(slugSource), Tagline: p.Tagline, Description: p.Description,
		ProjectType: p.ProjectType, RepositoryUrl: p.RepositoryUrl, LiveUrl: p.LiveUrl, ImageUrl: p.ImageUrl,
		Technologies: strings.FieldsFunc(p.Technologies, func(r rune) bool { return r == ',' }),
		StartedAt:    started, LaunchedAt: launched, PublishedAt: published, IsFeatured: p.IsFeatured,
	}
	var item domains.ProjectEntity
	if update {
		item, err = h.writer.UpdateProject(ctx, domains.UpdateProjectData{
			ID: id, Name: data.Name, Slug: data.Slug, Tagline: data.Tagline, Description: data.Description,
			ProjectType: data.ProjectType, RepositoryUrl: data.RepositoryUrl, LiveUrl: data.LiveUrl, ImageUrl: data.ImageUrl,
			Technologies: data.Technologies, StartedAt: data.StartedAt, LaunchedAt: data.LaunchedAt,
			PublishedAt: data.PublishedAt, IsFeatured: data.IsFeatured,
		})
	} else {
		item, err = h.writer.CreateProject(ctx, data)
	}
	if err != nil {
		if page, ok := h.renderer.Page(ctx, component, props).Errors(err); ok {
			return page.Render()
		}
		if errors.Is(err, domains.ErrNotFound) {
			return anduin.NotFound(err)
		}
		if errors.Is(err, services.ErrFeaturedProjectLimit) {
			return h.renderer.Page(ctx, component, props).ValidationErrors(map[string]string{"isFeatured": "No more than 3 projects can be featured."}).Render()
		}
		target, verb := routes.AdminProjectNew.URL(), "create"
		if update {
			target, verb = routes.AdminProjectEdit.URL(id), "update"
		}
		return adminContentFlashRedirect(ctx, h.renderer, kiks.FlashError, fmt.Sprintf("Failed to %s project: %v", verb, err), target)
	}
	h.invalidate()
	verb := "created"
	if update {
		verb = "updated"
	}
	return adminContentFlashRedirect(ctx, h.renderer, kiks.FlashSuccess, "Project "+verb+" successfully", routes.AdminProjectShow.URL(item.ID))
}
func (h AdminProjectHandler) Create(ctx context.Context) error { return h.save(ctx, false) }
func (h AdminProjectHandler) Update(ctx context.Context) error { return h.save(ctx, true) }
func (h AdminProjectHandler) Destroy(ctx context.Context) error {
	id, err := adminContentID(ctx)
	if err != nil {
		return err
	}
	if err := h.projects.Destroy(ctx, id); err != nil {
		return adminContentFlashRedirect(ctx, h.renderer, kiks.FlashError, fmt.Sprintf("Failed to delete project: %v", err), routes.AdminProjectIndex.URL())
	}
	h.invalidate()
	return adminContentFlashRedirect(ctx, h.renderer, kiks.FlashSuccess, "Project destroyed successfully", routes.AdminProjectIndex.URL())
}
