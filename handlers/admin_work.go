package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mbvlabs/domains"
	"mbvlabs/routes"
	"mbvlabs/services"

	"github.com/gosimple/slug"
	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
)

type adminWorkStore interface {
	Paginate(context.Context, int64, int64) (domains.PaginatedWorks, error)
	Find(context.Context, int32) (domains.WorkEntity, error)
	Destroy(context.Context, int32) error
}

type adminWorkWriter interface {
	CreateWork(context.Context, domains.CreateWorkData) (domains.WorkEntity, error)
	UpdateWork(context.Context, domains.UpdateWorkData) (domains.WorkEntity, error)
}

type AdminWorkHandler struct {
	works      adminWorkStore
	writer     adminWorkWriter
	renderer   *inertia.Renderer
	invalidate func()
}

func NewAdminWorkHandler(works *domains.Works, writer services.Works, renderer *inertia.Renderer, invalidate func()) AdminWorkHandler {
	if invalidate == nil {
		invalidate = func() {}
	}
	return AdminWorkHandler{works, writer, renderer, invalidate}
}

func (h AdminWorkHandler) Register(rt *anduin.Router, middleware ...anduin.Middleware) {
	rt.Route(routes.AdminWorkIndex).Get(h.Index, middleware...)
	rt.Route(routes.AdminWorkNew).Get(h.New, middleware...)
	rt.Route(routes.AdminWorkShow).Get(h.Show, middleware...)
	rt.Route(routes.AdminWorkCreate).Post(h.Create, middleware...)
	rt.Route(routes.AdminWorkEdit).Get(h.Edit, middleware...)
	rt.Route(routes.AdminWorkUpdate).Put(h.Update, middleware...)
	rt.Route(routes.AdminWorkDestroy).Delete(h.Destroy, middleware...)
}

type AdminWorkData struct {
	ID                                                                                        int32
	CreatedAt, UpdatedAt                                                                      time.Time
	Title, Slug, ClientName, ClientIndustry, ClientUrl, ClientLogoUrl, Summary, CoverImageUrl string
	Specialisms, Platforms, Technologies                                                      []string
	Challenge, Approach, Deliverables, Outcome, Content                                       string
	TestimonialQuote, TestimonialName, TestimonialRole                                        string
	StartedAt, CompletedAt, Status, PublishedAt                                               string
	IsFeatured                                                                                bool
}

func NewAdminWorkData(e domains.WorkEntity) AdminWorkData {
	return AdminWorkData{
		ID: e.ID, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
		Title: e.Title, Slug: e.Slug, ClientName: e.ClientName.String,
		ClientIndustry: e.ClientIndustry.String, ClientUrl: e.ClientURL.String,
		ClientLogoUrl: e.ClientLogoURL.String, Summary: e.Summary, CoverImageUrl: e.CoverImageURL.String,
		Specialisms: e.Specialisms, Platforms: e.Platforms, Technologies: e.Technologies,
		Challenge: e.Challenge, Approach: e.Approach, Deliverables: e.Deliverables,
		Outcome: e.Outcome, Content: e.Content,
		StartedAt: adminContentDate(e.StartedAt), CompletedAt: adminContentDate(e.CompletedAt),
		Status: e.Status.String(), PublishedAt: adminContentDate(e.PublishedAt), IsFeatured: e.IsFeatured,
	}
}

func (h AdminWorkHandler) Index(ctx context.Context) error {
	page, size := adminContentPagination(ctx)
	list, err := h.works.Paginate(ctx, page, size)
	if err != nil {
		return err
	}
	items := make([]AdminWorkData, 0, len(list.Works))
	for _, item := range list.Works {
		items = append(items, NewAdminWorkData(item))
	}
	return h.renderer.Page(ctx, "Admin/Work/Index", inertia.Props{"items": items}).Render()
}

func (h AdminWorkHandler) New(ctx context.Context) error {
	return h.renderer.Page(ctx, "Admin/Work/Create", inertia.Props{}).Render()
}

func (h AdminWorkHandler) item(ctx context.Context, component string) error {
	id, err := adminContentID(ctx)
	if err != nil {
		return err
	}
	item, err := h.works.Find(ctx, id)
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	if err != nil {
		return err
	}
	return h.renderer.Page(ctx, component, inertia.Props{"item": NewAdminWorkData(item)}).Render()
}
func (h AdminWorkHandler) Show(ctx context.Context) error { return h.item(ctx, "Admin/Work/Show") }
func (h AdminWorkHandler) Edit(ctx context.Context) error { return h.item(ctx, "Admin/Work/Edit") }

type adminWorkPayload struct {
	Title            string   `json:"title"`
	Slug             string   `json:"slug"`
	ClientName       string   `json:"clientName"`
	ClientIndustry   string   `json:"clientIndustry"`
	ClientURL        string   `json:"clientUrl"`
	ClientLogoURL    string   `json:"clientLogoUrl"`
	Summary          string   `json:"summary"`
	CoverImageURL    string   `json:"coverImageUrl"`
	Specialisms      []string `json:"specialisms"`
	Platforms        []string `json:"platforms"`
	Technologies     []string `json:"technologies"`
	Challenge        string   `json:"challenge"`
	Approach         string   `json:"approach"`
	Deliverables     string   `json:"deliverables"`
	Outcome          string   `json:"outcome"`
	Content          string   `json:"content"`
	TestimonialQuote string   `json:"testimonialQuote"`
	TestimonialName  string   `json:"testimonialName"`
	TestimonialRole  string   `json:"testimonialRole"`
	StartedAt        string   `json:"startedAt"`
	CompletedAt      string   `json:"completedAt"`
	Status           string   `json:"status"`
	IsFeatured       bool     `json:"isFeatured"`
}

func (p adminWorkPayload) item(id int32) AdminWorkData {
	return AdminWorkData{ID: id, Title: p.Title, Slug: p.Slug, ClientName: p.ClientName, ClientIndustry: p.ClientIndustry,
		ClientUrl: p.ClientURL, ClientLogoUrl: p.ClientLogoURL, Summary: p.Summary, CoverImageUrl: p.CoverImageURL,
		Specialisms: p.Specialisms, Platforms: p.Platforms, Technologies: p.Technologies,
		Challenge: p.Challenge, Approach: p.Approach, Deliverables: p.Deliverables, Outcome: p.Outcome, Content: p.Content,
		TestimonialQuote: p.TestimonialQuote, TestimonialName: p.TestimonialName, TestimonialRole: p.TestimonialRole,
		StartedAt: p.StartedAt, CompletedAt: p.CompletedAt, Status: p.Status, IsFeatured: p.IsFeatured}
}

func (h AdminWorkHandler) save(ctx context.Context, update bool) error {
	var id int32
	var err error
	if update {
		id, err = adminContentID(ctx)
		if err != nil {
			return err
		}
	}
	var p adminWorkPayload
	if err := anduin.DecodeJSON(ctx, &p); err != nil {
		return anduin.BadRequest(err)
	}
	component, props := "Admin/Work/Create", inertia.Props{}
	if update {
		component = "Admin/Work/Edit"
		props["item"] = p.item(id)
	}
	errs := map[string]string{}
	started := adminContentParseDate(p.StartedAt, "startedAt", errs)
	completed := adminContentParseDate(p.CompletedAt, "completedAt", errs)
	if len(errs) > 0 {
		return h.renderer.Page(ctx, component, props).ValidationErrors(errs).Render()
	}
	slugSource := p.Slug
	if slugSource == "" {
		slugSource = p.Title
	}
	data := domains.CreateWorkData{
		Title: p.Title, Slug: slug.Make(slugSource), ClientName: p.ClientName, ClientIndustry: p.ClientIndustry,
		ClientURL: p.ClientURL, ClientLogoURL: p.ClientLogoURL, Summary: p.Summary, CoverImageURL: p.CoverImageURL,
		Specialisms: p.Specialisms, Platforms: p.Platforms, Technologies: p.Technologies,
		Challenge: p.Challenge, Approach: p.Approach, Deliverables: p.Deliverables, Outcome: p.Outcome,
		Content: p.Content, TestimonialQuote: p.TestimonialQuote, TestimonialName: p.TestimonialName, TestimonialRole: p.TestimonialRole,
		StartedAt: started, CompletedAt: completed, Status: domains.StatusEnum(p.Status), IsFeatured: p.IsFeatured,
	}
	var item domains.WorkEntity
	if update {
		item, err = h.writer.UpdateWork(ctx, domains.UpdateWorkData{
			ID: id, Title: data.Title, Slug: data.Slug, ClientName: data.ClientName, ClientIndustry: data.ClientIndustry,
			ClientURL: data.ClientURL, ClientLogoURL: data.ClientLogoURL, Summary: data.Summary, CoverImageURL: data.CoverImageURL,
			Specialisms: data.Specialisms, Platforms: data.Platforms, Technologies: data.Technologies, Challenge: data.Challenge,
			Approach: data.Approach, Deliverables: data.Deliverables, Outcome: data.Outcome, Content: data.Content,
			TestimonialQuote: data.TestimonialQuote, TestimonialName: data.TestimonialName, TestimonialRole: data.TestimonialRole,
			StartedAt: data.StartedAt, CompletedAt: data.CompletedAt, Status: data.Status, IsFeatured: data.IsFeatured,
		})
	} else {
		item, err = h.writer.CreateWork(ctx, data)
	}
	if err != nil {
		if page, ok := h.renderer.Page(ctx, component, props).Errors(err); ok {
			return page.Render()
		}
		if errors.Is(err, domains.ErrNotFound) {
			return anduin.NotFound(err)
		}
		if errors.Is(err, services.ErrFeaturedWorkLimit) {
			return h.renderer.Page(ctx, component, props).ValidationErrors(map[string]string{"isFeatured": "No more than 3 works can be featured."}).Render()
		}
		target, verb := routes.AdminWorkNew.URL(), "create"
		if update {
			target, verb = routes.AdminWorkEdit.URL(id), "update"
		}
		return adminContentFlashRedirect(ctx, h.renderer, kiks.FlashError, fmt.Sprintf("Failed to %s work: %v", verb, err), target)
	}
	h.invalidate()
	verb := "created"
	if update {
		verb = "updated"
	}
	return adminContentFlashRedirect(ctx, h.renderer, kiks.FlashSuccess, "Work "+verb+" successfully", routes.AdminWorkShow.URL(item.ID))
}
func (h AdminWorkHandler) Create(ctx context.Context) error { return h.save(ctx, false) }
func (h AdminWorkHandler) Update(ctx context.Context) error { return h.save(ctx, true) }
func (h AdminWorkHandler) Destroy(ctx context.Context) error {
	id, err := adminContentID(ctx)
	if err != nil {
		return err
	}
	if err := h.works.Destroy(ctx, id); err != nil {
		return adminContentFlashRedirect(ctx, h.renderer, kiks.FlashError, fmt.Sprintf("Failed to delete work: %v", err), routes.AdminWorkIndex.URL())
	}
	h.invalidate()
	return adminContentFlashRedirect(ctx, h.renderer, kiks.FlashSuccess, "Work destroyed successfully", routes.AdminWorkIndex.URL())
}
