package handlers

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"mbvlabs/domains"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
)

type apiTestArticles struct{ created domains.CreateBlogPostData }

func (a *apiTestArticles) AllDraftOrPublished(context.Context) ([]domains.BlogPostEntity, error) {
	return []domains.BlogPostEntity{}, nil
}
func (a *apiTestArticles) Create(_ context.Context, data domains.CreateBlogPostData, _ *time.Time) (domains.BlogPostEntity, error) {
	a.created = data
	return domains.BlogPostEntity{}, nil
}

type apiTestWorks struct{ created domains.CreateWorkData }

func (w *apiTestWorks) CreateWork(_ context.Context, data domains.CreateWorkData) (domains.WorkEntity, error) {
	w.created = data
	return domains.WorkEntity{}, nil
}

type apiTestDiary struct{}

func (apiTestDiary) FindBetweenDates(_ context.Context, start, _ time.Time) ([]domains.DiaryEntryEntity, error) {
	return []domains.DiaryEntryEntity{{EntryDate: start, MorningThoughts: domains.OptionalString{String: "morning", Valid: true}}}, nil
}

func TestAPIContentRoutes(t *testing.T) {
	articles, works := &apiTestArticles{}, &apiTestWorks{}
	rt := anduin.NewRouter()
	NewAPIContentHandler(articles, articles, works, nil, "user", "pass").Register(rt)
	NewAPIDiaryHandler(apiTestDiary{}, "user", "pass").Register(rt)
	request := func(method, path, body string, auth bool) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if auth {
			req.SetBasicAuth("user", "pass")
		}
		rec := httptest.NewRecorder()
		rt.Handler().ServeHTTP(rec, req)
		return rec
	}
	if r := request("GET", routes.ApiArticleIndex.Path(), "", false); r.Code != 401 {
		t.Fatalf("articles auth = %d", r.Code)
	}
	if r := request("POST", routes.ApiArticleCreate.Path(), `{"title":"Hello World"}`, true); r.Code != 201 || articles.created.Slug != "hello-world" || articles.created.Status != "draft" {
		t.Fatalf("article = %d %s, data=%+v", r.Code, r.Body, articles.created)
	}
	if r := request("POST", routes.ApiArticleCreate.Path(), `{"title":"Hi","scheduledAt":"tomorrow"}`, true); r.Code != 400 || !strings.Contains(r.Body.String(), "scheduledAt") {
		t.Fatalf("schedule validation = %d %s", r.Code, r.Body)
	}
	if r := request("POST", routes.ApiWorkCreate.Path(), `{"title":"Work","startedAt":"bad"}`, false); r.Code != 400 || !strings.Contains(r.Body.String(), "startedAt") {
		t.Fatalf("work validation = %d %s", r.Code, r.Body)
	}
	if r := request("POST", routes.ApiWorkCreate.Path(), `{"title":"Work"}`, false); r.Code != 201 || works.created.Status != domains.Draft {
		t.Fatalf("work = %d %s, data=%+v", r.Code, r.Body, works.created)
	}
	if r := request("GET", routes.ApiDiaryThoughtsPreviousWeek.Path(), "", true); r.Code != 200 || !strings.Contains(r.Body.String(), `"morningThoughts":"morning"`) {
		t.Fatalf("diary = %d %s", r.Code, r.Body)
	}
}
