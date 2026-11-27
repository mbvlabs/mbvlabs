package payloads

import "encoding/json"

type PublicSEO struct {
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Image          string          `json:"image"`
	ImageAlt       string          `json:"imageAlt"`
	Canonical      string          `json:"canonical"`
	MetaType       string          `json:"metaType"`
	StructuredData json.RawMessage `json:"structuredData,omitempty"`
}

type PublicSite struct {
	Production bool `json:"production"`
	Year       int  `json:"year"`
}

type PublicLayoutProps struct {
	SEO        PublicSEO  `json:"seo"`
	PublicSite PublicSite `json:"publicSite"`
}

type PublicHomeProps struct {
	Works    []WorkItem    `json:"works"`
	Projects []ProjectItem `json:"projects"`
}

type PublicWorkIndexProps struct {
	Items []WorkItem `json:"items"`
}

type PublicWorkShowProps struct {
	Item         WorkItem `json:"item"`
	MarkdownHTML string   `json:"markdownHTML"`
}

type PublicProjectIndexProps struct {
	Items []ProjectItem `json:"items"`
}

type PublicProjectShowProps struct {
	Item         ProjectItem `json:"item"`
	MarkdownHTML string      `json:"markdownHTML"`
}

type PublicArticleIndexProps struct {
	Items []ArticleItem `json:"items"`
}

type RelatedArticle struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type PublicArticleShowProps struct {
	Item            ArticleItem      `json:"item"`
	MarkdownHTML    string           `json:"markdownHTML"`
	RelatedArticles []RelatedArticle `json:"relatedArticles"`
}

type WorkSlug struct {
	Slug string
}

type PublicServicesProps struct {
	WorkItems []WorkSlug `json:"workItems"`
}

type PublicContactProps struct {
	CSRFToken string `json:"csrfToken"`
}

type PublicErrorProps struct {
	Status int `json:"status"`
}

type PublicInquiryPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Company     string `json:"company"`
	Role        string `json:"role"`
	ProjectType string `json:"projectType"`
	Timeline    string `json:"timeline"`
	Message     string `json:"message"`
}
