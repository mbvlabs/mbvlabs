package handlers

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"mbvlabs/domains"
	"mbvlabs/payloads"
)

func TestContentAPIWireShapes(t *testing.T) {
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	article := domains.BlogPostEntity{Title: "article", Excerpt: domains.OptionalString{String: "excerpt", Valid: true}, Tags: json.RawMessage(`["go"]`), PublishedAt: domains.OptionalTime{Time: now, Valid: true}}
	project := domains.ProjectEntity{Name: "project", Technologies: json.RawMessage(`["go"]`)}
	work := domains.WorkEntity{Title: "work", ClientURL: domains.OptionalString{String: "https://example.com", Valid: true}, Technologies: []string{"go"}}
	tests := []struct{ old, new any }{
		{article, articleItem(article)},
		{project, projectItem(project)},
		{work, workItem(work)},
		{domains.PaginatedBlogPosts{BlogPosts: []domains.BlogPostEntity{article}, TotalCount: 1}, payloads.APIArticlePage{BlogPosts: []payloads.ArticleItem{articleItem(article)}, TotalCount: 1}},
		{domains.PaginatedProjects{}, payloads.APIProjectPage{}},
		{domains.PaginatedWorks{}, payloads.APIWorkPage{}},
	}
	for _, tt := range tests {
		oldJSON, err := json.Marshal(tt.old)
		if err != nil {
			t.Fatal(err)
		}
		newJSON, err := json.Marshal(tt.new)
		if err != nil {
			t.Fatal(err)
		}
		var oldValue, newValue any
		if err := json.Unmarshal(oldJSON, &oldValue); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(newJSON, &newValue); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(oldValue, newValue) {
			t.Errorf("wire shape changed:\n%s\n%s", oldJSON, newJSON)
		}
	}
	if mapPayloads([]domains.BlogPostEntity(nil), articleItem) != nil {
		t.Error("nil slice must remain nil")
	}
}
