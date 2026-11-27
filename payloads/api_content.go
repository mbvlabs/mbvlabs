package payloads

type APIArticlePage struct {
	BlogPosts                              []ArticleItem
	TotalCount, Page, PageSize, TotalPages int64
}

type APIProjectPage struct {
	Projects                               []ProjectItem
	TotalCount, Page, PageSize, TotalPages int64
}

type APIWorkPage struct {
	Works                                  []WorkItem
	TotalCount, Page, PageSize, TotalPages int64
}

type APICreateArticle struct {
	Title         string   `json:"title"`
	Slug          string   `json:"slug"`
	Excerpt       string   `json:"excerpt"`
	Body          string   `json:"body"`
	CoverImageURL string   `json:"coverImageUrl"`
	Tags          []string `json:"tags"`
	ScheduledAt   *string  `json:"scheduledAt"`
}
type APICreateWork struct {
	Title          string   `json:"title"`
	Slug           string   `json:"slug"`
	ClientName     string   `json:"clientName"`
	ClientIndustry string   `json:"clientIndustry"`
	ClientURL      string   `json:"clientUrl"`
	ClientLogoURL  string   `json:"clientLogoUrl"`
	Summary        string   `json:"summary"`
	CoverImageURL  string   `json:"coverImageUrl"`
	Specialisms    []string `json:"specialisms"`
	Platforms      []string `json:"platforms"`
	Technologies   []string `json:"technologies"`
	Challenge      string   `json:"challenge"`
	Approach       string   `json:"approach"`
	Deliverables   string   `json:"deliverables"`
	Outcome        string   `json:"outcome"`
	Content        string   `json:"content"`
	StartedAt      string   `json:"startedAt"`
	CompletedAt    string   `json:"completedAt"`
	Status         string   `json:"status"`
	PublishedAt    string   `json:"publishedAt"`
	IsFeatured     bool     `json:"isFeatured"`
}
type APIArticleEdit struct {
	Title         *string   `json:"title"`
	Excerpt       *string   `json:"excerpt"`
	Body          *string   `json:"body"`
	CoverImageURL *string   `json:"coverImageUrl"`
	Tags          *[]string `json:"tags"`
}
type APIProjectEdit struct {
	Name          *string   `json:"name"`
	Tagline       *string   `json:"tagline"`
	Description   *string   `json:"description"`
	RepositoryURL *string   `json:"repositoryUrl"`
	LiveURL       *string   `json:"liveUrl"`
	ImageURL      *string   `json:"imageUrl"`
	Technologies  *[]string `json:"technologies"`
}
type APIWorkEdit struct {
	Title         *string `json:"title"`
	Summary       *string `json:"summary"`
	Content       *string `json:"content"`
	ClientURL     *string `json:"clientUrl"`
	ClientLogoURL *string `json:"clientLogoUrl"`
	CoverImageURL *string `json:"coverImageUrl"`
	Challenge     *string `json:"challenge"`
	Approach      *string `json:"approach"`
	Deliverables  *string `json:"deliverables"`
	Outcome       *string `json:"outcome"`
}
type APIThought struct {
	EntryDate       string `json:"entryDate"`
	MorningThoughts string `json:"morningThoughts"`
	EveningThoughts string `json:"eveningThoughts"`
}

type APIThoughtWeek struct {
	WeekStart string       `json:"weekStart"`
	WeekEnd   string       `json:"weekEnd"`
	Items     []APIThought `json:"items"`
}
