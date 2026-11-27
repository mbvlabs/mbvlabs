package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"mbvlabs/domains"
	"mbvlabs/payloads"
	"mbvlabs/routes"

	"github.com/gosimple/slug"
	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type adminBlogStore interface {
	Find(context.Context, int32) (domains.BlogPostEntity, error)
	Paginate(context.Context, int64, int64) (domains.PaginatedBlogPosts, error)
	Update(context.Context, domains.UpdateBlogPostData) (domains.BlogPostEntity, error)
	Destroy(context.Context, int32) error
}

type AdminBlogPostHandler struct {
	posts      adminBlogStore
	creator    BlogPostCreator
	renderer   *inertia.Renderer
	invalidate func()
}

type BlogPostCreator interface {
	Create(context.Context, domains.CreateBlogPostData, *time.Time) (domains.BlogPostEntity, error)
}

func NewAdminBlogPostHandler(posts *domains.BlogPosts, creator BlogPostCreator, renderer *inertia.Renderer, invalidate func()) AdminBlogPostHandler {
	if invalidate == nil {
		invalidate = func() {}
	}
	return AdminBlogPostHandler{posts: posts, creator: creator, renderer: renderer, invalidate: invalidate}
}

func (h AdminBlogPostHandler) Register(rt *anduin.Router, middleware ...anduin.Middleware) {
	rt.Route(routes.AdminBlogPostIndex).Get(h.Index, middleware...)
	rt.Route(routes.AdminBlogPostNew).Get(h.New, middleware...)
	rt.Route(routes.AdminBlogPostShow).Get(h.Show, middleware...)
	rt.Route(routes.AdminBlogPostCreate).Post(h.Create, middleware...)
	rt.Route(routes.AdminBlogPostEdit).Get(h.Edit, middleware...)
	rt.Route(routes.AdminBlogPostUpdate).Put(h.Update, middleware...)
	rt.Route(routes.AdminBlogPostDestroy).Delete(h.Destroy, middleware...)
}

func adminBlogItem(e domains.BlogPostEntity) (payloads.AdminBlogData, error) {
	schedule, err := e.PublicationScheduleData()
	if err != nil {
		return payloads.AdminBlogData{}, err
	}
	var tags []string
	tagString := string(e.Tags)
	if json.Unmarshal(e.Tags, &tags) == nil {
		tagString = strings.Join(tags, ", ")
	}
	published := ""
	if e.PublishedAt.Valid {
		published = e.PublishedAt.Time.Format("2006-01-02")
	}
	var resultSchedule *payloads.AdminBlogSchedule
	if schedule != nil {
		resultSchedule = &payloads.AdminBlogSchedule{JobID: schedule.JobID, ScheduledAt: schedule.ScheduledAt}
	}
	return payloads.AdminBlogData{ID: e.ID, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt, Title: e.Title, Slug: e.Slug, Excerpt: e.Excerpt.String, Body: e.Body, Status: e.Status, CoverImageUrl: e.CoverImageUrl.String, Tags: tagString, PublishedAt: published, PublicationSchedule: resultSchedule}, nil
}
func (h AdminBlogPostHandler) Index(ctx context.Context) error {
	page, size := adminPagination(ctx)
	list, err := h.posts.Paginate(ctx, page, size)
	if err != nil {
		return err
	}
	items := make([]payloads.AdminBlogData, 0, len(list.BlogPosts))
	for _, e := range list.BlogPosts {
		item, err := adminBlogItem(e)
		if err != nil {
			return err
		}
		items = append(items, item)
	}
	return h.renderer.Page(ctx, "Admin/BlogPost/Index", inertia.FromStruct(payloads.AdminBlogIndexProps{Items: items})).Render()
}
func (h AdminBlogPostHandler) New(ctx context.Context) error {
	return h.renderer.Page(ctx, "Admin/BlogPost/Create", inertia.Props{}).Render()
}
func (h AdminBlogPostHandler) item(ctx context.Context, component string) error {
	id, err := adminDiaryID(ctx)
	if err != nil {
		return err
	}
	e, err := h.posts.Find(ctx, id)
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	if err != nil {
		return err
	}
	item, err := adminBlogItem(e)
	if err != nil {
		return err
	}
	return h.renderer.Page(ctx, component, inertia.FromStruct(payloads.AdminBlogItemProps{Item: item})).Render()
}
func (h AdminBlogPostHandler) Show(ctx context.Context) error {
	return h.item(ctx, "Admin/BlogPost/Show")
}
func (h AdminBlogPostHandler) Edit(ctx context.Context) error {
	return h.item(ctx, "Admin/BlogPost/Edit")
}

func (h AdminBlogPostHandler) save(ctx context.Context, update bool) error {
	var id int32
	var err error
	if update {
		id, err = adminDiaryID(ctx)
		if err != nil {
			return err
		}
	}
	var p payloads.AdminBlogPayload
	if err := anduin.DecodeJSON(ctx, &p); err != nil {
		return anduin.BadRequest(err)
	}
	source := p.Slug
	if source == "" {
		source = p.Title
	}
	tags := strings.FieldsFunc(p.Tags, func(r rune) bool { return r == ',' })
	var e domains.BlogPostEntity
	if update {
		e, err = h.posts.Update(ctx, domains.UpdateBlogPostData{ID: id, Title: p.Title, Slug: slug.Make(source), Excerpt: p.Excerpt, Body: p.Body, Status: p.Status, CoverImageUrl: p.CoverImageUrl, Tags: tags})
	} else {
		data := domains.CreateBlogPostData{Title: p.Title, Slug: slug.Make(source), Excerpt: p.Excerpt, Body: p.Body, Status: p.Status, CoverImageUrl: p.CoverImageUrl, Tags: tags}
		if p.Status == "scheduled" {
			var at time.Time
			if p.ScheduledAt == "" {
				err = validation.Field("scheduledAt", "required", "is required")
			} else {
				at, err = time.Parse(time.RFC3339, p.ScheduledAt)
				if err != nil {
					err = validation.Field("scheduledAt", "datetime", "must be a valid date and time")
				}
			}
			if err == nil {
				e, err = h.creator.Create(ctx, data, &at)
			}
		} else {
			e, err = h.creator.Create(ctx, data, nil)
		}
	}
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	if err != nil {
		component := "Admin/BlogPost/Create"
		props := inertia.Props{}
		if update {
			component = "Admin/BlogPost/Edit"
			props["item"] = payloads.AdminBlogData{ID: id, Title: p.Title, Slug: p.Slug, Excerpt: p.Excerpt, Body: p.Body, Status: p.Status, CoverImageUrl: p.CoverImageUrl, Tags: p.Tags}
		}
		if page, ok := h.renderer.Page(ctx, component, props).Errors(err); ok {
			return page.Render()
		}
		return err
	}
	h.invalidate()
	message := "Blog post created successfully"
	if update {
		message = "Blog post updated successfully"
	}
	if err := kiks.AddFlash(ctx, kiks.FlashSuccess, message); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, routes.AdminBlogPostShow.URL(e.ID), http.StatusSeeOther)
}
func (h AdminBlogPostHandler) Create(ctx context.Context) error { return h.save(ctx, false) }
func (h AdminBlogPostHandler) Update(ctx context.Context) error { return h.save(ctx, true) }
func (h AdminBlogPostHandler) Destroy(ctx context.Context) error {
	id, err := adminDiaryID(ctx)
	if err != nil {
		return err
	}
	if err := h.posts.Destroy(ctx, id); err != nil {
		return err
	}
	h.invalidate()
	if err := kiks.AddFlash(ctx, kiks.FlashSuccess, "BlogPost destroyed successfully"); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, routes.AdminBlogPostIndex.URL(), http.StatusSeeOther)
}
