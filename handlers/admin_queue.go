package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"mbvlabs/domains"
	"mbvlabs/payloads"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
)

type adminQueueStore interface {
	Find(context.Context, int64) (domains.RiverJob, error)
	Jobs(context.Context, domains.RiverFilters, int64, int64) ([]domains.RiverJob, domains.RiverPagination, error)
	Stats(context.Context) ([]domains.RiverQueueStats, error)
	StateCounts(context.Context) ([]domains.RiverStateCount, error)
	SetPaused(context.Context, string, bool) error
	Action(context.Context, int64, string) error
}

type AdminQueueHandler struct {
	queues   adminQueueStore
	renderer *inertia.Renderer
}

func NewAdminQueueHandler(queues *domains.RiverQueues, renderer *inertia.Renderer) AdminQueueHandler {
	return AdminQueueHandler{queues: queues, renderer: renderer}
}

func (h AdminQueueHandler) Register(rt *anduin.Router, middleware ...anduin.Middleware) {
	rt.Route(routes.AdminQueueIndex).Get(h.Index, middleware...)
	rt.Route(routes.AdminQueueJobs).Get(h.Jobs, middleware...)
	rt.Route(routes.AdminQueueJobShow).Get(h.Show, middleware...)
	rt.Route(routes.AdminQueuePause).Post(h.Pause, middleware...)
	rt.Route(routes.AdminQueueResume).Post(h.Resume, middleware...)
	rt.Route(routes.AdminQueueJobCancel).Post(h.Cancel, middleware...)
	rt.Route(routes.AdminQueueJobRetry).Post(h.Retry, middleware...)
	rt.Route(routes.AdminQueueJobDiscard).Post(h.Discard, middleware...)
}

func (h AdminQueueHandler) Index(ctx context.Context) error {
	queues, err := h.queues.Stats(ctx)
	if err != nil {
		return err
	}
	counts, err := h.queues.StateCounts(ctx)
	if err != nil {
		return err
	}
	jobs, _, err := h.queues.Jobs(ctx, domains.RiverFilters{}, 1, 10)
	if err != nil {
		return err
	}
	return h.renderer.Page(ctx, "Admin/Queue/Index", inertia.FromStruct(payloads.AdminQueueIndexProps{
		Queues:      mapPayloads(queues, func(q domains.RiverQueueStats) payloads.AdminQueueStats { return payloads.AdminQueueStats(q) }),
		StateCounts: mapPayloads(counts, func(c domains.RiverStateCount) payloads.AdminQueueStateCount { return payloads.AdminQueueStateCount(c) }),
		RecentJobs:  mapPayloads(jobs, func(j domains.RiverJob) payloads.AdminQueueJob { return payloads.AdminQueueJob(j) }),
	})).Render()
}

func queueInt(values url.Values, key string, fallback, min, max int64) int64 {
	n, err := strconv.ParseInt(values.Get(key), 10, 64)
	if err != nil {
		return fallback
	}
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func (h AdminQueueHandler) Jobs(ctx context.Context) error {
	values := anduin.Request(ctx).URL.Query()
	state := values.Get("state")
	switch state {
	case "available", "cancelled", "completed", "discarded", "pending", "retryable", "running", "scheduled":
	default:
		state = ""
	}
	filters := domains.RiverFilters{State: state, Queue: values.Get("queue"), Kind: values.Get("kind")}
	jobs, pagination, err := h.queues.Jobs(ctx, filters, queueInt(values, "page", 1, 1, 100000), queueInt(values, "per_page", 25, 1, 100))
	if err != nil {
		return err
	}
	queues, err := h.queues.Stats(ctx)
	if err != nil {
		return err
	}
	return h.renderer.Page(ctx, "Admin/Queue/Jobs", inertia.FromStruct(payloads.AdminQueueJobsProps{
		Items:      mapPayloads(jobs, func(j domains.RiverJob) payloads.AdminQueueJob { return payloads.AdminQueueJob(j) }),
		Queues:     mapPayloads(queues, func(q domains.RiverQueueStats) payloads.AdminQueueStats { return payloads.AdminQueueStats(q) }),
		Filters:    payloads.AdminQueueFilters(filters),
		Pagination: payloads.AdminQueuePagination(pagination),
	})).Render()
}

func queueID(ctx context.Context) (int64, error) {
	id, err := strconv.ParseInt(anduin.PathValue(ctx, "id"), 10, 64)
	if err != nil {
		return 0, anduin.BadRequest(err)
	}
	return id, nil
}

func (h AdminQueueHandler) Show(ctx context.Context) error {
	id, err := queueID(ctx)
	if err != nil {
		return err
	}
	job, err := h.queues.Find(ctx, id)
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	if err != nil {
		return err
	}
	return h.renderer.Page(ctx, "Admin/Queue/Show", inertia.FromStruct(payloads.AdminQueueShowProps{Item: payloads.AdminQueueJob(job)})).Render()
}

func (h AdminQueueHandler) setPaused(ctx context.Context, paused bool) error {
	var payload payloads.AdminQueueActionPayload
	if err := anduin.DecodeJSON(ctx, &payload); err != nil {
		return anduin.BadRequest(err)
	}
	verb := "resumed"
	if paused {
		verb = "paused"
	}
	err := h.queues.SetPaused(ctx, payload.Name, paused)
	kind, msg := kiks.FlashSuccess, fmt.Sprintf("Queue %q %s", payload.Name, verb)
	if err != nil {
		kind, msg = kiks.FlashError, fmt.Sprintf("Failed to %s queue: %v", map[bool]string{true: "pause", false: "resume"}[paused], err)
	}
	if err := kiks.AddFlash(ctx, kind, msg); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, routes.AdminQueueIndex.URL(), http.StatusSeeOther)
}
func (h AdminQueueHandler) Pause(ctx context.Context) error  { return h.setPaused(ctx, true) }
func (h AdminQueueHandler) Resume(ctx context.Context) error { return h.setPaused(ctx, false) }

func (h AdminQueueHandler) action(ctx context.Context, action, label string) error {
	id, err := queueID(ctx)
	if err != nil {
		return err
	}
	target := routes.AdminQueueJobs.URL()
	// Only redirect to a local referrer, never to a different origin.
	if ref, e := url.Parse(anduin.Request(ctx).Header.Get("Referer")); e == nil && ref != nil && ref.Path != "" && (ref.Host == "" || ref.Host == anduin.Request(ctx).Host) {
		target = ref.RequestURI()
	}
	err = h.queues.Action(ctx, id, action)
	if errors.Is(err, domains.ErrNotFound) {
		return anduin.NotFound(err)
	}
	kind, msg := kiks.FlashSuccess, fmt.Sprintf("Job %d %s", id, label)
	if errors.Is(err, domains.ErrRiverJobInvalidTransition) {
		kind, msg = kiks.FlashError, fmt.Sprintf("Job %d cannot be %s from its current state", id, label)
	} else if err != nil {
		kind, msg = kiks.FlashError, fmt.Sprintf("Failed to update job %d: %v", id, err)
	}
	if err := kiks.AddFlash(ctx, kind, msg); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, target, http.StatusSeeOther)
}
func (h AdminQueueHandler) Cancel(ctx context.Context) error {
	return h.action(ctx, "cancel", "cancelled")
}
func (h AdminQueueHandler) Retry(ctx context.Context) error {
	return h.action(ctx, "retry", "queued for retry")
}
func (h AdminQueueHandler) Discard(ctx context.Context) error {
	return h.action(ctx, "discard", "discarded")
}
