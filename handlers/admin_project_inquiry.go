package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"mbvlabs/domains"
	"mbvlabs/payloads"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type adminInquiryStore interface {
	Find(context.Context, int32) (domains.ProjectInquiryEntity, error)
	Paginate(context.Context, int64, int64) (domains.PaginatedProjectInquiries, error)
	Create(context.Context, domains.CreateProjectInquiryData) (domains.ProjectInquiryEntity, error)
	Update(context.Context, domains.UpdateProjectInquiryData) (domains.ProjectInquiryEntity, error)
	Destroy(context.Context, int32) error
	IsDomainBanned(context.Context, string) (bool, error)
	BanDomain(context.Context, int32) (string, error)
	UnbanDomain(context.Context, int32) (string, error)
}

type AdminProjectInquiryHandler struct {
	inquiries adminInquiryStore
	renderer  *inertia.Renderer
}

func NewAdminProjectInquiryHandler(inquiries *domains.ProjectInquiries, renderer *inertia.Renderer) AdminProjectInquiryHandler {
	return AdminProjectInquiryHandler{inquiries, renderer}
}

func (h AdminProjectInquiryHandler) Register(rt *anduin.Router, middleware ...anduin.Middleware) {
	rt.Route(routes.AdminProjectInquiryIndex).Get(h.Index, middleware...)
	rt.Route(routes.AdminProjectInquiryNew).Get(h.New, middleware...)
	rt.Route(routes.AdminProjectInquiryShow).Get(h.Show, middleware...)
	rt.Route(routes.AdminProjectInquiryCreate).Post(h.Create, middleware...)
	rt.Route(routes.AdminProjectInquiryEdit).Get(h.Edit, middleware...)
	rt.Route(routes.AdminProjectInquiryUpdate).Put(h.Update, middleware...)
	rt.Route(routes.AdminProjectInquiryBanDomain).Post(h.BanDomain, middleware...)
	rt.Route(routes.AdminProjectInquiryUnbanDomain).Delete(h.UnbanDomain, middleware...)
	rt.Route(routes.AdminProjectInquiryDestroy).Delete(h.Destroy, middleware...)
}

func adminInquiryItem(e domains.ProjectInquiryEntity) payloads.AdminInquiryData {
	metadata := string(e.Metadata)
	if metadata == "" {
		metadata = "{}"
	} else {
		var obj map[string]any
		if json.Unmarshal(e.Metadata, &obj) == nil {
			if data, err := json.MarshalIndent(obj, "", "  "); err == nil {
				metadata = string(data)
			}
		}
	}
	return payloads.AdminInquiryData{ID: e.ID, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt, Name: e.Name, Email: e.Email, Company: e.Company.String, Role: e.Role.String, ProjectType: e.ProjectType.String, Timeline: e.Timeline.String, Message: e.Message, Source: e.Source.String, Status: e.Status, Metadata: metadata, DomainBanned: false}
}

func adminPagination(ctx context.Context) (int64, int64) {
	q := anduin.Request(ctx).URL.Query()
	page, err := strconv.ParseInt(q.Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.ParseInt(q.Get("per_page"), 10, 64)
	if err != nil || size < 1 || size > 100 {
		size = 25
	}
	return page, size
}

func (h AdminProjectInquiryHandler) Index(ctx context.Context) error {
	page, size := adminPagination(ctx)
	list, err := h.inquiries.Paginate(ctx, page, size)
	if err != nil {
		return err
	}
	items := make([]payloads.AdminInquiryData, 0, len(list.ProjectInquiries))
	for _, e := range list.ProjectInquiries {
		items = append(items, adminInquiryItem(e))
	}
	return h.renderer.Page(ctx, "Admin/ProjectInquiry/Index", inertia.FromStruct(payloads.AdminInquiryIndexProps{Items: items})).Render()
}
func (h AdminProjectInquiryHandler) New(ctx context.Context) error {
	return h.renderer.Page(ctx, "Admin/ProjectInquiry/Create", inertia.Props{}).Render()
}
func (h AdminProjectInquiryHandler) item(ctx context.Context, show bool) error {
	id, err := adminDiaryID(ctx)
	if err != nil {
		return err
	}
	e, err := h.inquiries.Find(ctx, id)
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	if err != nil {
		return err
	}
	item := adminInquiryItem(e)
	component := "Admin/ProjectInquiry/Edit"
	if show {
		item.DomainBanned, err = h.inquiries.IsDomainBanned(ctx, e.Email)
		if err != nil {
			return err
		}
		component = "Admin/ProjectInquiry/Show"
	}
	return h.renderer.Page(ctx, component, inertia.FromStruct(payloads.AdminInquiryItemProps{Item: item})).Render()
}
func (h AdminProjectInquiryHandler) Show(ctx context.Context) error { return h.item(ctx, true) }
func (h AdminProjectInquiryHandler) Edit(ctx context.Context) error { return h.item(ctx, false) }

func (h AdminProjectInquiryHandler) save(ctx context.Context, update bool) error {
	var id int32
	var err error
	if update {
		id, err = adminDiaryID(ctx)
		if err != nil {
			return err
		}
	}
	var p payloads.AdminInquiryPayload
	if err := anduin.DecodeJSON(ctx, &p); err != nil {
		return anduin.BadRequest(err)
	}
	raw := p.Metadata
	if strings.TrimSpace(raw) == "" {
		raw = "{}"
	}
	var object map[string]any
	if err := json.Unmarshal([]byte(raw), &object); err != nil || object == nil {
		err = validation.Field("metadata", "json", "must be a JSON object")
		component := "Admin/ProjectInquiry/Create"
		props := inertia.Props{}
		if update {
			component = "Admin/ProjectInquiry/Edit"
			props["item"] = payloads.AdminInquiryData{ID: id, Name: p.Name, Email: p.Email, Company: p.Company, Role: p.Role, ProjectType: p.ProjectType, Timeline: p.Timeline, Message: p.Message, Source: p.Source, Status: p.Status, Metadata: p.Metadata}
		}
		if page, ok := h.renderer.Page(ctx, component, props).Errors(err); ok {
			return page.Render()
		}
		return err
	}
	metadata, err := json.Marshal(object)
	if err != nil {
		return err
	}
	ns := func(s string) domains.OptionalString { return domains.OptionalString{String: s, Valid: true} }
	var e domains.ProjectInquiryEntity
	if update {
		e, err = h.inquiries.Update(ctx, domains.UpdateProjectInquiryData{ID: id, Name: p.Name, Email: p.Email, Company: ns(p.Company), Role: ns(p.Role), ProjectType: ns(p.ProjectType), Timeline: ns(p.Timeline), Message: p.Message, Source: ns(p.Source), Status: p.Status, Metadata: metadata})
	} else {
		e, err = h.inquiries.Create(ctx, domains.CreateProjectInquiryData{Name: p.Name, Email: p.Email, Company: ns(p.Company), Role: ns(p.Role), ProjectType: ns(p.ProjectType), Timeline: ns(p.Timeline), Message: p.Message, Source: ns(p.Source), Status: p.Status, Metadata: metadata})
	}
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	if err != nil {
		component := "Admin/ProjectInquiry/Create"
		props := inertia.Props{}
		if update {
			component = "Admin/ProjectInquiry/Edit"
			props["item"] = payloads.AdminInquiryData{ID: id, Name: p.Name, Email: p.Email, Company: p.Company, Role: p.Role, ProjectType: p.ProjectType, Timeline: p.Timeline, Message: p.Message, Source: p.Source, Status: p.Status, Metadata: p.Metadata}
		}
		if page, ok := h.renderer.Page(ctx, component, props).Errors(err); ok {
			return page.Render()
		}
		return err
	}
	message := "ProjectInquiry created successfully"
	if update {
		message = "ProjectInquiry updated successfully"
	}
	if err := kiks.AddFlash(ctx, kiks.FlashSuccess, message); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, routes.AdminProjectInquiryShow.URL(e.ID), http.StatusSeeOther)
}
func (h AdminProjectInquiryHandler) Create(ctx context.Context) error { return h.save(ctx, false) }
func (h AdminProjectInquiryHandler) Update(ctx context.Context) error { return h.save(ctx, true) }
func (h AdminProjectInquiryHandler) ban(ctx context.Context, banned bool) error {
	id, err := adminDiaryID(ctx)
	if err != nil {
		return err
	}
	var domain string
	if banned {
		domain, err = h.inquiries.BanDomain(ctx, id)
	} else {
		domain, err = h.inquiries.UnbanDomain(ctx, id)
	}
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	if err != nil {
		return err
	}
	action := "unbanned"
	if banned {
		action = "banned"
	}
	if err := kiks.AddFlash(ctx, kiks.FlashSuccess, fmt.Sprintf("Domain %s %s", domain, action)); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, routes.AdminProjectInquiryShow.URL(id), http.StatusSeeOther)
}
func (h AdminProjectInquiryHandler) BanDomain(ctx context.Context) error   { return h.ban(ctx, true) }
func (h AdminProjectInquiryHandler) UnbanDomain(ctx context.Context) error { return h.ban(ctx, false) }
func (h AdminProjectInquiryHandler) Destroy(ctx context.Context) error {
	id, err := adminDiaryID(ctx)
	if err != nil {
		return err
	}
	if err := h.inquiries.Destroy(ctx, id); err != nil {
		return err
	}
	if err := kiks.AddFlash(ctx, kiks.FlashSuccess, "ProjectInquiry destroyed successfully"); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, routes.AdminProjectInquiryIndex.URL(), http.StatusSeeOther)
}
