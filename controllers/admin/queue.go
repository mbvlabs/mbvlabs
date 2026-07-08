package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mbvlabs/internal/inertia"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/cookies"
	"mbvlabs/router/routes"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

type Queue struct {
	db storage.Pool
}

func NewQueue(db storage.Pool) Queue {
	return Queue{db}
}

type RiverJobData struct {
	ID          int64
	State       string
	Attempt     int16
	MaxAttempts int16
	AttemptedAt *time.Time
	CreatedAt   time.Time
	FinalizedAt *time.Time
	ScheduledAt time.Time
	Priority    int16
	Args        string
	AttemptedBy []string
	Errors      []string
	Kind        string
	Metadata    string
	Queue       string
	Tags        []string
}

type RiverQueueData struct {
	Name             string
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	PausedAt         *time.Time
	IsPaused         bool
	AvailableCount   int64
	CancelledCount   int64
	CompletedCount   int64
	DiscardedCount   int64
	PendingCount     int64
	RetryableCount   int64
	RunningCount     int64
	ScheduledCount   int64
	TotalCount       int64
	ActiveClients    int64
	MaxWorkers       int64
	NumJobsRunning   int64
	NumJobsCompleted int64
}

type RiverJobStateCountData struct {
	State string
	Count int64
}

type RiverJobFiltersData struct {
	State string
	Queue string
	Kind  string
}

type RiverPaginationData struct {
	Page       int64
	PageSize   int64
	TotalCount int64
	TotalPages int64
}

func (q Queue) Index(etx *echo.Context) error {
	dashboard, err := models.RiverQueue.Dashboard(etx.Request().Context(), q.db.Executor())
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Queue/Index", inertia.Props{
		"queues":      newRiverQueueDataList(dashboard.Queues),
		"stateCounts": newRiverJobStateCountDataList(dashboard.StateCounts),
		"recentJobs":  newRiverJobDataList(dashboard.RecentJobs),
	})
}

func (q Queue) Jobs(etx *echo.Context) error {
	page := parseInt64Query(etx, "page", 1, 1, 100000)
	pageSize := parseInt64Query(etx, "per_page", 25, 1, 100)
	filters := riverJobFiltersFromRequest(etx)

	jobs, err := models.RiverJob.Paginate(
		etx.Request().Context(),
		q.db.Executor(),
		models.RiverJobFilters(filters),
		page,
		pageSize,
	)
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	queues, err := models.RiverQueue.WithStats(etx.Request().Context(), q.db.Executor())
	if err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Queue/Jobs", inertia.Props{
		"items":   newRiverJobDataList(jobs.RiverJobs),
		"queues":  newRiverQueueDataList(queues),
		"filters": RiverJobFiltersData(filters),
		"pagination": RiverPaginationData{
			Page:       jobs.Page,
			PageSize:   jobs.PageSize,
			TotalCount: jobs.TotalCount,
			TotalPages: jobs.TotalPages,
		},
	})
}

func (q Queue) Show(etx *echo.Context) error {
	id, err := parseRiverJobID(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	job, err := models.RiverJob.Find(etx.Request().Context(), q.db.Executor(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}

	return inertia.Page(etx, "Admin/Queue/Show", inertia.Props{
		"item": newRiverJobData(job),
	})
}

type QueueActionPayload struct {
	Name string `json:"name"`
}

func (q Queue) Pause(etx *echo.Context) error {
	var payload QueueActionPayload
	if err := etx.Bind(&payload); err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	if err := models.RiverQueue.Pause(etx.Request().Context(), q.db.Executor(), payload.Name); err != nil {
		return queueActionError(etx, routes.AdminQueueIndex.URL(), fmt.Sprintf("Failed to pause queue: %v", err))
	}

	if err := cookies.AddFlash(etx, cookies.FlashSuccess, fmt.Sprintf("Queue %q paused", payload.Name)); err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminQueueIndex.URL())
}

func (q Queue) Resume(etx *echo.Context) error {
	var payload QueueActionPayload
	if err := etx.Bind(&payload); err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	if err := models.RiverQueue.Resume(etx.Request().Context(), q.db.Executor(), payload.Name); err != nil {
		return queueActionError(etx, routes.AdminQueueIndex.URL(), fmt.Sprintf("Failed to resume queue: %v", err))
	}

	if err := cookies.AddFlash(etx, cookies.FlashSuccess, fmt.Sprintf("Queue %q resumed", payload.Name)); err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, routes.AdminQueueIndex.URL())
}

func (q Queue) Cancel(etx *echo.Context) error {
	return q.jobAction(etx, "cancelled", models.RiverJob.Cancel)
}

func (q Queue) Retry(etx *echo.Context) error {
	return q.jobAction(etx, "queued for retry", models.RiverJob.Retry)
}

func (q Queue) Discard(etx *echo.Context) error {
	return q.jobAction(etx, "discarded", models.RiverJob.Discard)
}

func (q Queue) jobAction(
	etx *echo.Context,
	label string,
	action func(context.Context, storage.Executor, int64) error,
) error {
	id, err := parseRiverJobID(etx)
	if err != nil {
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	redirectURL := etx.Request().Header.Get("Referer")
	if redirectURL == "" {
		redirectURL = routes.AdminQueueJobs.URL()
	}

	if err := action(etx.Request().Context(), q.db.Executor(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return inertia.Page(etx, "Errors/NotFound", inertia.Props{})
		}
		if errors.Is(err, models.ErrRiverJobInvalidTransition) {
			return queueActionError(etx, redirectURL, fmt.Sprintf("Job %d cannot be %s from its current state", id, label))
		}
		return queueActionError(etx, redirectURL, fmt.Sprintf("Failed to update job %d: %v", id, err))
	}

	if err := cookies.AddFlash(etx, cookies.FlashSuccess, fmt.Sprintf("Job %d %s", id, label)); err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, redirectURL)
}

func (q Queue) RegisterRoutes(r *router.Router) error {
	var errs []error
	routeDefs := []echo.Route{
		{
			Method:      http.MethodGet,
			Path:        routes.AdminQueueIndex.Path(),
			Name:        routes.AdminQueueIndex.Name(),
			Handler:     q.Index,
			Middlewares: authOnly,
		},
		{
			Method:      http.MethodGet,
			Path:        routes.AdminQueueJobs.Path(),
			Name:        routes.AdminQueueJobs.Name(),
			Handler:     q.Jobs,
			Middlewares: authOnly,
		},
		{
			Method:      http.MethodGet,
			Path:        routes.AdminQueueJobShow.Path(),
			Name:        routes.AdminQueueJobShow.Name(),
			Handler:     q.Show,
			Middlewares: authOnly,
		},
		{
			Method:      http.MethodPost,
			Path:        routes.AdminQueuePause.Path(),
			Name:        routes.AdminQueuePause.Name(),
			Handler:     q.Pause,
			Middlewares: authOnly,
		},
		{
			Method:      http.MethodPost,
			Path:        routes.AdminQueueResume.Path(),
			Name:        routes.AdminQueueResume.Name(),
			Handler:     q.Resume,
			Middlewares: authOnly,
		},
		{
			Method:      http.MethodPost,
			Path:        routes.AdminQueueJobCancel.Path(),
			Name:        routes.AdminQueueJobCancel.Name(),
			Handler:     q.Cancel,
			Middlewares: authOnly,
		},
		{
			Method:      http.MethodPost,
			Path:        routes.AdminQueueJobRetry.Path(),
			Name:        routes.AdminQueueJobRetry.Name(),
			Handler:     q.Retry,
			Middlewares: authOnly,
		},
		{
			Method:      http.MethodPost,
			Path:        routes.AdminQueueJobDiscard.Path(),
			Name:        routes.AdminQueueJobDiscard.Name(),
			Handler:     q.Discard,
			Middlewares: authOnly,
		},
	}

	for _, route := range routeDefs {
		if _, err := r.AddRoute(route); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func newRiverJobData(entity models.RiverJobEntity) RiverJobData {
	return RiverJobData{
		ID:          entity.ID,
		State:       entity.State,
		Attempt:     entity.Attempt,
		MaxAttempts: entity.MaxAttempts,
		AttemptedAt: nullableTime(entity.AttemptedAt),
		CreatedAt:   entity.CreatedAt,
		FinalizedAt: nullableTime(entity.FinalizedAt),
		ScheduledAt: entity.ScheduledAt,
		Priority:    entity.Priority,
		Args:        string(entity.Args),
		AttemptedBy: entity.AttemptedBy,
		Errors:      entity.Errors,
		Kind:        entity.Kind,
		Metadata:    string(entity.Metadata),
		Queue:       entity.Queue,
		Tags:        entity.Tags,
	}
}

func newRiverJobDataList(entities []models.RiverJobEntity) []RiverJobData {
	items := make([]RiverJobData, 0, len(entities))
	for _, entity := range entities {
		items = append(items, newRiverJobData(entity))
	}
	return items
}

func newRiverQueueData(stats models.RiverQueueStats) RiverQueueData {
	return RiverQueueData{
		Name:             stats.Name,
		CreatedAt:        nullableTime(stats.CreatedAt),
		UpdatedAt:        nullableTime(stats.UpdatedAt),
		PausedAt:         nullableTime(stats.PausedAt),
		IsPaused:         stats.PausedAt.Valid,
		AvailableCount:   stats.AvailableCount,
		CancelledCount:   stats.CancelledCount,
		CompletedCount:   stats.CompletedCount,
		DiscardedCount:   stats.DiscardedCount,
		PendingCount:     stats.PendingCount,
		RetryableCount:   stats.RetryableCount,
		RunningCount:     stats.RunningCount,
		ScheduledCount:   stats.ScheduledCount,
		TotalCount:       stats.TotalCount,
		ActiveClients:    stats.ActiveClients,
		MaxWorkers:       stats.MaxWorkers,
		NumJobsRunning:   stats.NumJobsRunning,
		NumJobsCompleted: stats.NumJobsCompleted,
	}
}

func newRiverQueueDataList(stats []models.RiverQueueStats) []RiverQueueData {
	items := make([]RiverQueueData, 0, len(stats))
	for _, item := range stats {
		items = append(items, newRiverQueueData(item))
	}
	return items
}

func newRiverJobStateCountDataList(counts []models.RiverJobStateCount) []RiverJobStateCountData {
	items := make([]RiverJobStateCountData, 0, len(counts))
	for _, count := range counts {
		items = append(items, RiverJobStateCountData(count))
	}
	return items
}

func nullableTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}

func riverJobFiltersFromRequest(etx *echo.Context) RiverJobFiltersData {
	state := etx.QueryParam("state")
	if !models.IsRiverJobState(state) {
		state = ""
	}

	return RiverJobFiltersData{
		State: state,
		Queue: etx.QueryParam("queue"),
		Kind:  etx.QueryParam("kind"),
	}
}

func parseRiverJobID(etx *echo.Context) (int64, error) {
	return strconv.ParseInt(etx.Param("id"), 10, 64)
}

func parseInt64Query(etx *echo.Context, key string, fallback int64, min int64, max int64) int64 {
	value := fallback
	if raw := etx.QueryParam(key); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil {
			value = parsed
		}
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func queueActionError(etx *echo.Context, redirectURL string, msg string) error {
	if err := cookies.AddFlash(etx, cookies.FlashError, msg); err != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Redirect(etx, redirectURL)
}
