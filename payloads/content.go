package payloads

import (
	"encoding/json"
	"time"
)

// These two wire types preserve the existing content API's nullable objects.
type ContentString struct {
	String string
	Valid  bool
}

type ContentTime struct {
	Time  time.Time
	Valid bool
}

type WorkItem struct {
	ID                                                   int32
	CreatedAt, UpdatedAt                                 time.Time
	Title, Slug                                          string
	ClientName, ClientIndustry, ClientURL, ClientLogoURL ContentString
	Summary                                              string
	CoverImageURL                                        ContentString
	Specialisms, Platforms, Technologies                 []string
	Challenge, Approach, Deliverables, Outcome, Content  string
	StartedAt, CompletedAt                               ContentTime
	Status                                               string
	PublishedAt                                          ContentTime
	IsFeatured                                           bool
}

type ProjectItem struct {
	ID                                 int32
	CreatedAt, UpdatedAt               time.Time
	Name, Slug, Tagline                string
	Description                        ContentString
	ProjectType                        string
	RepositoryURL, LiveURL, ImageURL   ContentString
	Technologies                       json.RawMessage
	StartedAt, LaunchedAt, PublishedAt ContentTime
	IsFeatured                         bool
}

type ArticleItem struct {
	ID                  int32
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Title               string
	Slug                string
	Excerpt             ContentString
	Body                string
	Status              string
	CoverImageUrl       ContentString
	Tags                json.RawMessage
	PublishedAt         ContentTime
	PublicationSchedule json.RawMessage
}
