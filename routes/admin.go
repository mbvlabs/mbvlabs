package routes

import "github.com/mbvlabs/andurel/pkg/routing"

const AdminPrefix = "/admin"
const AdminWorkPrefix = "/admin/works"
const AdminProjectPrefix = "/admin/projects"
const AdminBlogPostPrefix = "/admin/blog-posts"
const AdminProjectInquiryPrefix = "/admin/project-inquiries"
const AdminDiaryEntryPrefix = "/admin/diary"
const AdminQueuePrefix = "/admin/queue"

var AdminHomePage = routing.NewSimpleRoute("", "admin.home", AdminPrefix, routing.TypeScriptRoute())

var AdminWorkIndex = routing.NewSimpleRoute("", "admin.works.index", AdminWorkPrefix, routing.TypeScriptRoute())
var AdminWorkShow = routing.NewRouteWithSerialID("/:id", "admin.works.show", AdminWorkPrefix, routing.TypeScriptRoute())
var AdminWorkNew = routing.NewSimpleRoute("/new", "admin.works.new", AdminWorkPrefix, routing.TypeScriptRoute())
var AdminWorkCreate = routing.NewSimpleRoute("", "admin.works.create", AdminWorkPrefix, routing.TypeScriptRoute())
var AdminWorkEdit = routing.NewRouteWithSerialID("/:id/edit", "admin.works.edit", AdminWorkPrefix, routing.TypeScriptRoute())
var AdminWorkUpdate = routing.NewRouteWithSerialID("/:id", "admin.works.update", AdminWorkPrefix, routing.TypeScriptRoute())
var AdminWorkDestroy = routing.NewRouteWithSerialID("/:id", "admin.works.destroy", AdminWorkPrefix, routing.TypeScriptRoute())

var AdminProjectIndex = routing.NewSimpleRoute("", "admin.projects.index", AdminProjectPrefix, routing.TypeScriptRoute())
var AdminProjectShow = routing.NewRouteWithSerialID("/:id", "admin.projects.show", AdminProjectPrefix, routing.TypeScriptRoute())
var AdminProjectNew = routing.NewSimpleRoute("/new", "admin.projects.new", AdminProjectPrefix, routing.TypeScriptRoute())
var AdminProjectCreate = routing.NewSimpleRoute("", "admin.projects.create", AdminProjectPrefix, routing.TypeScriptRoute())
var AdminProjectEdit = routing.NewRouteWithSerialID("/:id/edit", "admin.projects.edit", AdminProjectPrefix, routing.TypeScriptRoute())
var AdminProjectUpdate = routing.NewRouteWithSerialID("/:id", "admin.projects.update", AdminProjectPrefix, routing.TypeScriptRoute())
var AdminProjectDestroy = routing.NewRouteWithSerialID("/:id", "admin.projects.destroy", AdminProjectPrefix, routing.TypeScriptRoute())

var AdminBlogPostIndex = routing.NewSimpleRoute("", "admin.blog_posts.index", AdminBlogPostPrefix, routing.TypeScriptRoute())
var AdminBlogPostShow = routing.NewRouteWithSerialID("/:id", "admin.blog_posts.show", AdminBlogPostPrefix, routing.TypeScriptRoute())
var AdminBlogPostNew = routing.NewSimpleRoute("/new", "admin.blog_posts.new", AdminBlogPostPrefix, routing.TypeScriptRoute())
var AdminBlogPostCreate = routing.NewSimpleRoute("", "admin.blog_posts.create", AdminBlogPostPrefix, routing.TypeScriptRoute())
var AdminBlogPostEdit = routing.NewRouteWithSerialID("/:id/edit", "admin.blog_posts.edit", AdminBlogPostPrefix, routing.TypeScriptRoute())
var AdminBlogPostUpdate = routing.NewRouteWithSerialID("/:id", "admin.blog_posts.update", AdminBlogPostPrefix, routing.TypeScriptRoute())
var AdminBlogPostDestroy = routing.NewRouteWithSerialID("/:id", "admin.blog_posts.destroy", AdminBlogPostPrefix, routing.TypeScriptRoute())

var AdminProjectInquiryIndex = routing.NewSimpleRoute("", "admin.project_inquiries.index", AdminProjectInquiryPrefix, routing.TypeScriptRoute())
var AdminProjectInquiryShow = routing.NewRouteWithSerialID("/:id", "admin.project_inquiries.show", AdminProjectInquiryPrefix, routing.TypeScriptRoute())
var AdminProjectInquiryNew = routing.NewSimpleRoute("/new", "admin.project_inquiries.new", AdminProjectInquiryPrefix, routing.TypeScriptRoute())
var AdminProjectInquiryCreate = routing.NewSimpleRoute("", "admin.project_inquiries.create", AdminProjectInquiryPrefix, routing.TypeScriptRoute())
var AdminProjectInquiryEdit = routing.NewRouteWithSerialID("/:id/edit", "admin.project_inquiries.edit", AdminProjectInquiryPrefix, routing.TypeScriptRoute())
var AdminProjectInquiryUpdate = routing.NewRouteWithSerialID("/:id", "admin.project_inquiries.update", AdminProjectInquiryPrefix, routing.TypeScriptRoute())
var AdminProjectInquiryBanDomain = routing.NewRouteWithSerialID("/:id/domain-ban", "admin.project_inquiries.ban_domain", AdminProjectInquiryPrefix, routing.TypeScriptRoute())
var AdminProjectInquiryUnbanDomain = routing.NewRouteWithSerialID("/:id/domain-ban", "admin.project_inquiries.unban_domain", AdminProjectInquiryPrefix, routing.TypeScriptRoute())
var AdminProjectInquiryDestroy = routing.NewRouteWithSerialID("/:id", "admin.project_inquiries.destroy", AdminProjectInquiryPrefix, routing.TypeScriptRoute())

var AdminDiaryEntryIndex = routing.NewSimpleRoute("", "admin.diary_entries.index", AdminDiaryEntryPrefix, routing.TypeScriptRoute())
var AdminDiaryEntryShow = routing.NewRouteWithSerialID("/:id", "admin.diary_entries.show", AdminDiaryEntryPrefix, routing.TypeScriptRoute())
var AdminDiaryEntryNew = routing.NewSimpleRoute("/new", "admin.diary_entries.new", AdminDiaryEntryPrefix, routing.TypeScriptRoute())
var AdminDiaryEntryToday = routing.NewSimpleRoute("/today", "admin.diary_entries.today", AdminDiaryEntryPrefix, routing.TypeScriptRoute())
var AdminDiaryEntryCreate = routing.NewSimpleRoute("", "admin.diary_entries.create", AdminDiaryEntryPrefix, routing.TypeScriptRoute())
var AdminDiaryEntryEdit = routing.NewRouteWithSerialID("/:id/edit", "admin.diary_entries.edit", AdminDiaryEntryPrefix, routing.TypeScriptRoute())
var AdminDiaryEntryUpdate = routing.NewRouteWithSerialID("/:id", "admin.diary_entries.update", AdminDiaryEntryPrefix, routing.TypeScriptRoute())
var AdminDiaryEntryDestroy = routing.NewRouteWithSerialID("/:id", "admin.diary_entries.destroy", AdminDiaryEntryPrefix, routing.TypeScriptRoute())

var AdminQueueIndex = routing.NewSimpleRoute("", "admin.queue.index", AdminQueuePrefix, routing.TypeScriptRoute())
var AdminQueueJobs = routing.NewSimpleRoute("/jobs", "admin.queue.jobs", AdminQueuePrefix, routing.TypeScriptRoute())
var AdminQueueJobShow = routing.NewRouteWithBigSerialID("/jobs/:id", "admin.queue.jobs.show", AdminQueuePrefix, routing.TypeScriptRoute())
var AdminQueuePause = routing.NewSimpleRoute("/pause", "admin.queue.pause", AdminQueuePrefix, routing.TypeScriptRoute())
var AdminQueueResume = routing.NewSimpleRoute("/resume", "admin.queue.resume", AdminQueuePrefix, routing.TypeScriptRoute())
var AdminQueueJobCancel = routing.NewRouteWithBigSerialID("/jobs/:id/cancel", "admin.queue.jobs.cancel", AdminQueuePrefix, routing.TypeScriptRoute())
var AdminQueueJobRetry = routing.NewRouteWithBigSerialID("/jobs/:id/retry", "admin.queue.jobs.retry", AdminQueuePrefix, routing.TypeScriptRoute())
var AdminQueueJobDiscard = routing.NewRouteWithBigSerialID("/jobs/:id/discard", "admin.queue.jobs.discard", AdminQueuePrefix, routing.TypeScriptRoute())
