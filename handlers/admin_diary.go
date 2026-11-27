package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"mbvlabs/domains"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
)

type adminDiaryStore interface {
	Find(context.Context, int32) (domains.DiaryEntryEntity, error)
	FindOrCreateForDate(context.Context, time.Time) (domains.DiaryEntryEntity, error)
	Paginate(context.Context, int64, int64) (domains.PaginatedDiaryEntries, error)
	Create(context.Context, domains.CreateDiaryEntryData) (domains.DiaryEntryEntity, error)
	Update(context.Context, domains.UpdateDiaryEntryData) (domains.DiaryEntryEntity, error)
	Destroy(context.Context, int32) error
}

type AdminDiaryHandler struct {
	entries  adminDiaryStore
	renderer *inertia.Renderer
}

func NewAdminDiaryHandler(entries *domains.DiaryEntries, renderer *inertia.Renderer) AdminDiaryHandler {
	return AdminDiaryHandler{entries: entries, renderer: renderer}
}

func (h AdminDiaryHandler) Register(rt *anduin.Router, middleware ...anduin.Middleware) {
	rt.Route(routes.AdminDiaryEntryIndex).Get(h.Index, middleware...)
	rt.Route(routes.AdminDiaryEntryToday).Get(h.Today, middleware...)
	rt.Route(routes.AdminDiaryEntryNew).Get(h.New, middleware...)
	rt.Route(routes.AdminDiaryEntryShow).Get(h.Show, middleware...)
	rt.Route(routes.AdminDiaryEntryCreate).Post(h.Create, middleware...)
	rt.Route(routes.AdminDiaryEntryEdit).Get(h.Edit, middleware...)
	rt.Route(routes.AdminDiaryEntryUpdate).Put(h.Update, middleware...)
	rt.Route(routes.AdminDiaryEntryDestroy).Delete(h.Destroy, middleware...)
}

type adminDiaryData struct {
	ID              int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
	EntryDate       string
	MorningThoughts string
	EveningThoughts string
}

func adminDiaryItem(e domains.DiaryEntryEntity) adminDiaryData {
	return adminDiaryData{e.ID, e.CreatedAt, e.UpdatedAt, e.EntryDate.Format("2006-01-02"), e.MorningThoughts.String, e.EveningThoughts.String}
}

func (h AdminDiaryHandler) Index(ctx context.Context) error {
	q := anduin.Request(ctx).URL.Query()
	page, err := strconv.ParseInt(q.Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.ParseInt(q.Get("per_page"), 10, 64)
	if err != nil || size < 1 || size > 100 {
		size = 25
	}
	list, err := h.entries.Paginate(ctx, page, size)
	if err != nil {
		return err
	}
	items := make([]adminDiaryData, 0, len(list.DiaryEntries))
	for _, e := range list.DiaryEntries {
		items = append(items, adminDiaryItem(e))
	}
	return h.renderer.Page(ctx, "Admin/DiaryEntry/Index", inertia.Props{"items": items}).Render()
}

func adminToday() time.Time {
	loc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		return time.Now()
	}
	return time.Now().In(loc)
}

func (h AdminDiaryHandler) New(ctx context.Context) error {
	return h.renderer.Page(ctx, "Admin/DiaryEntry/Create", inertia.Props{"today": adminToday().Format("2006-01-02")}).Render()
}

func (h AdminDiaryHandler) Today(ctx context.Context) error {
	focus := anduin.Request(ctx).URL.Query().Get("focus")
	if focus != "morning" && focus != "evening" {
		focus = ""
	}
	e, err := h.entries.FindOrCreateForDate(ctx, adminToday())
	if err != nil {
		return err
	}
	return h.renderer.Page(ctx, "Admin/DiaryEntry/Edit", inertia.Props{"item": adminDiaryItem(e), "focus": focus}).Render()
}

func adminDiaryID(ctx context.Context) (int32, error) {
	id, err := strconv.ParseInt(anduin.PathValue(ctx, "id"), 10, 32)
	if err != nil {
		return 0, anduin.BadRequest(err)
	}
	return int32(id), nil
}

func (h AdminDiaryHandler) item(ctx context.Context, component string) error {
	id, err := adminDiaryID(ctx)
	if err != nil {
		return err
	}
	e, err := h.entries.Find(ctx, id)
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	if err != nil {
		return err
	}
	props := inertia.Props{"item": adminDiaryItem(e)}
	if component == "Admin/DiaryEntry/Edit" {
		props["focus"] = anduin.Request(ctx).URL.Query().Get("focus")
	}
	return h.renderer.Page(ctx, component, props).Render()
}
func (h AdminDiaryHandler) Show(ctx context.Context) error {
	return h.item(ctx, "Admin/DiaryEntry/Show")
}
func (h AdminDiaryHandler) Edit(ctx context.Context) error {
	return h.item(ctx, "Admin/DiaryEntry/Edit")
}

type adminDiaryPayload struct {
	EntryDate       string `json:"entryDate"`
	MorningThoughts string `json:"morningThoughts"`
	EveningThoughts string `json:"eveningThoughts"`
}

func (h AdminDiaryHandler) save(ctx context.Context, update bool) error {
	var id int32
	if update {
		var err error
		id, err = adminDiaryID(ctx)
		if err != nil {
			return err
		}
	}
	var payload adminDiaryPayload
	if err := anduin.DecodeJSON(ctx, &payload); err != nil {
		return anduin.BadRequest(err)
	}
	date, _ := time.Parse("2006-01-02", payload.EntryDate)
	var e domains.DiaryEntryEntity
	var err error
	if update {
		e, err = h.entries.Update(ctx, domains.UpdateDiaryEntryData{ID: id, EntryDate: date, MorningThoughts: payload.MorningThoughts, EveningThoughts: payload.EveningThoughts})
	} else {
		e, err = h.entries.Create(ctx, domains.CreateDiaryEntryData{EntryDate: date, MorningThoughts: payload.MorningThoughts, EveningThoughts: payload.EveningThoughts})
	}
	if err != nil {
		if errors.Is(err, domains.ErrNotFound) {
			return anduin.NotFound(err)
		}
		component := "Admin/DiaryEntry/Create"
		props := inertia.Props{"item": adminDiaryData{ID: id, EntryDate: payload.EntryDate, MorningThoughts: payload.MorningThoughts, EveningThoughts: payload.EveningThoughts}, "today": adminToday().Format("2006-01-02")}
		if update {
			component = "Admin/DiaryEntry/Edit"
			delete(props, "today")
			props["focus"] = anduin.Request(ctx).URL.Query().Get("focus")
		}
		if page, ok := h.renderer.Page(ctx, component, props).Errors(err); ok {
			return page.Render()
		}
		return err
	}
	message := "Diary entry created successfully"
	if update {
		message = "Diary entry updated successfully"
	}
	if err := kiks.AddFlash(ctx, kiks.FlashSuccess, message); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, routes.AdminDiaryEntryShow.URL(e.ID), http.StatusSeeOther)
}
func (h AdminDiaryHandler) Create(ctx context.Context) error { return h.save(ctx, false) }
func (h AdminDiaryHandler) Update(ctx context.Context) error { return h.save(ctx, true) }
func (h AdminDiaryHandler) Destroy(ctx context.Context) error {
	id, err := adminDiaryID(ctx)
	if err != nil {
		return err
	}
	if err := h.entries.Destroy(ctx, id); err != nil {
		return err
	}
	if err := kiks.AddFlash(ctx, kiks.FlashSuccess, "Diary entry deleted successfully"); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, routes.AdminDiaryEntryIndex.URL(), http.StatusSeeOther)
}
