package payloads

import "time"

type AdminBlogSchedule struct {
	JobID       int64     `json:"job_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
}

type AdminBlogData struct {
	ID                                                                   int32
	CreatedAt, UpdatedAt                                                 time.Time
	Title, Slug, Excerpt, Body, Status, CoverImageUrl, Tags, PublishedAt string
	PublicationSchedule                                                  *AdminBlogSchedule
}

type AdminBlogPayload struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Excerpt       string `json:"excerpt"`
	Body          string `json:"body"`
	Status        string `json:"status"`
	CoverImageUrl string `json:"coverImageUrl"`
	Tags          string `json:"tags"`
	ScheduledAt   string `json:"scheduledAt"`
}

type AdminDiaryData struct {
	ID                                          int32
	CreatedAt, UpdatedAt                        time.Time
	EntryDate, MorningThoughts, EveningThoughts string
}

type AdminDiaryPayload struct {
	EntryDate       string `json:"entryDate"`
	MorningThoughts string `json:"morningThoughts"`
	EveningThoughts string `json:"eveningThoughts"`
}

type AdminInquiryData struct {
	ID                                                                                   int32
	CreatedAt, UpdatedAt                                                                 time.Time
	Name, Email, Company, Role, ProjectType, Timeline, Message, Source, Status, Metadata string
	DomainBanned                                                                         bool
}

type AdminInquiryPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Company     string `json:"company"`
	Role        string `json:"role"`
	ProjectType string `json:"projectType"`
	Timeline    string `json:"timeline"`
	Message     string `json:"message"`
	Source      string `json:"source"`
	Status      string `json:"status"`
	Metadata    string `json:"metadata"`
}

type AdminProjectData struct {
	ID                                                                              int32
	CreatedAt, UpdatedAt                                                            time.Time
	Name, Slug, Tagline, Description, ProjectType, RepositoryUrl, LiveUrl, ImageUrl string
	Technologies, StartedAt, LaunchedAt, PublishedAt                                string
	IsFeatured                                                                      bool
}

type AdminProjectPayload struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Tagline       string `json:"tagline"`
	Description   string `json:"description"`
	ProjectType   string `json:"projectType"`
	RepositoryUrl string `json:"repositoryUrl"`
	LiveUrl       string `json:"liveUrl"`
	ImageUrl      string `json:"imageUrl"`
	Technologies  string `json:"technologies"`
	StartedAt     string `json:"startedAt"`
	LaunchedAt    string `json:"launchedAt"`
	PublishedAt   string `json:"publishedAt"`
	IsFeatured    bool   `json:"isFeatured"`
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

type AdminWorkPayload struct {
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

type AdminBlogIndexProps struct {
	Items []AdminBlogData `json:"items"`
}
type AdminBlogItemProps struct {
	Item AdminBlogData `json:"item"`
}
type AdminDiaryIndexProps struct {
	Items []AdminDiaryData `json:"items"`
}
type AdminDiaryCreateProps struct {
	Today string `json:"today"`
}
type AdminDiaryItemProps struct {
	Item AdminDiaryData `json:"item"`
}
type AdminDiaryEditProps struct {
	Item  AdminDiaryData `json:"item"`
	Focus string         `json:"focus"`
}
type AdminInquiryIndexProps struct {
	Items []AdminInquiryData `json:"items"`
}
type AdminInquiryItemProps struct {
	Item AdminInquiryData `json:"item"`
}
type AdminProjectIndexProps struct {
	Items []AdminProjectData `json:"items"`
}
type AdminProjectItemProps struct {
	Item AdminProjectData `json:"item"`
}
type AdminWorkIndexProps struct {
	Items []AdminWorkData `json:"items"`
}
type AdminWorkItemProps struct {
	Item AdminWorkData `json:"item"`
}
type AdminHomeProps struct {
	AppName               string             `json:"appName"`
	ProjectInquiries      []AdminInquiryData `json:"projectInquiries"`
	Works                 []AdminWorkData    `json:"works"`
	BlogPosts             []AdminBlogData    `json:"blogPosts"`
	Projects              []AdminProjectData `json:"projects"`
	CanManageContentToken bool               `json:"canManageContentToken"`
}

type AdminQueueJob struct {
	ID                    int64
	State                 string
	Attempt, MaxAttempts  int16
	AttemptedAt           *time.Time
	CreatedAt             time.Time
	FinalizedAt           *time.Time
	ScheduledAt           time.Time
	Priority              int16
	Args                  string
	AttemptedBy, Errors   []string
	Kind, Metadata, Queue string
	Tags                  []string
}
type AdminQueueStats struct {
	Name                                                                                                                                   string
	CreatedAt, UpdatedAt, PausedAt                                                                                                         *time.Time
	IsPaused                                                                                                                               bool
	AvailableCount, CancelledCount, CompletedCount, DiscardedCount, PendingCount, RetryableCount, RunningCount, ScheduledCount, TotalCount int64
	ActiveClients, MaxWorkers, NumJobsRunning, NumJobsCompleted                                                                            int64
}
type AdminQueueStateCount struct {
	State string
	Count int64
}
type AdminQueueFilters struct{ State, Queue, Kind string }
type AdminQueuePagination struct{ Page, PageSize, TotalCount, TotalPages int64 }
type AdminQueueActionPayload struct {
	Name string `json:"name"`
}
type AdminQueueIndexProps struct {
	Queues      []AdminQueueStats      `json:"queues"`
	StateCounts []AdminQueueStateCount `json:"stateCounts"`
	RecentJobs  []AdminQueueJob        `json:"recentJobs"`
}
type AdminQueueJobsProps struct {
	Items      []AdminQueueJob      `json:"items"`
	Queues     []AdminQueueStats    `json:"queues"`
	Filters    AdminQueueFilters    `json:"filters"`
	Pagination AdminQueuePagination `json:"pagination"`
}
type AdminQueueShowProps struct {
	Item AdminQueueJob `json:"item"`
}
