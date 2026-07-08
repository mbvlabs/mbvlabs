package routes

import "mbvlabs/internal/routing"

const AdminQueuePrefix = "/admin/queue"

var AdminQueueIndex = routing.NewSimpleRoute(
	"",
	"admin.queue.index",
	AdminQueuePrefix,
)
var AdminQueueJobs = routing.NewSimpleRoute(
	"/jobs",
	"admin.queue.jobs",
	AdminQueuePrefix,
)
var AdminQueueJobShow = routing.NewRouteWithBigSerialID(
	"/jobs/:id",
	"admin.queue.jobs.show",
	AdminQueuePrefix,
)
var AdminQueuePause = routing.NewSimpleRoute(
	"/pause",
	"admin.queue.pause",
	AdminQueuePrefix,
)
var AdminQueueResume = routing.NewSimpleRoute(
	"/resume",
	"admin.queue.resume",
	AdminQueuePrefix,
)
var AdminQueueJobCancel = routing.NewRouteWithBigSerialID(
	"/jobs/:id/cancel",
	"admin.queue.jobs.cancel",
	AdminQueuePrefix,
)
var AdminQueueJobRetry = routing.NewRouteWithBigSerialID(
	"/jobs/:id/retry",
	"admin.queue.jobs.retry",
	AdminQueuePrefix,
)
var AdminQueueJobDiscard = routing.NewRouteWithBigSerialID(
	"/jobs/:id/discard",
	"admin.queue.jobs.discard",
	AdminQueuePrefix,
)
