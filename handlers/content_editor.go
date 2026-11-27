package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"mbvlabs/domains"
	"mbvlabs/payloads"
	"mbvlabs/payloads/siteapi"
	"mbvlabs/routes"
	"mbvlabs/services"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/telemetry"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type ContentEditorHandler struct {
	posts    *domains.BlogPosts
	projects *domains.Projects
	works    *domains.Works
	editor   services.ContentEditor
}

func NewContentEditorHandler(posts *domains.BlogPosts, projects *domains.Projects, works *domains.Works, editor services.ContentEditor) ContentEditorHandler {
	return ContentEditorHandler{posts, projects, works, editor}
}

func (h ContentEditorHandler) Register(rt *anduin.Router, bearer anduin.Middleware) {
	rt.Route(routes.ContentArticles).Get(h.articleList, bearer)
	rt.Route(routes.ContentArticle).Get(h.articleGet, bearer)
	rt.Route(routes.ContentArticle).Patch(h.articlePatch, bearer)
	rt.Route(routes.ContentProjects).Get(h.projectList, bearer)
	rt.Route(routes.ContentProject).Get(h.projectGet, bearer)
	rt.Route(routes.ContentProject).Patch(h.projectPatch, bearer)
	rt.Route(routes.ContentWorks).Get(h.workList, bearer)
	rt.Route(routes.ContentWork).Get(h.workGet, bearer)
	rt.Route(routes.ContentWork).Patch(h.workPatch, bearer)
}

func editorResult(ctx context.Context, value any, err error) error {
	if err == nil {
		return anduin.JSON(ctx, http.StatusOK, value)
	}
	if errors.Is(err, domains.ErrNotFound) {
		return apiFailure(ctx, http.StatusNotFound, "content not found")
	}
	if fields, ok := validation.As(err); ok {
		return anduin.JSON(ctx, http.StatusBadRequest, siteapi.APIError{Error: "validation failed", Fields: fields.ToMap()})
	}
	if errors.Is(err, domains.ErrDomainValidation) {
		return apiFailure(ctx, http.StatusBadRequest, "validation failed")
	}
	telemetry.Error(ctx, "content editor request failed", "route", anduin.RoutePath(anduin.RoutePattern(ctx)))
	return apiFailure(ctx, http.StatusInternalServerError, "content request failed")
}

func (h ContentEditorHandler) articleList(ctx context.Context) error {
	page, size := adminContentPagination(ctx)
	v, err := h.posts.Paginate(ctx, page, size)
	return editorResult(ctx, payloads.APIArticlePage{BlogPosts: mapPayloads(v.BlogPosts, articleItem), TotalCount: v.TotalCount, Page: v.Page, PageSize: v.PageSize, TotalPages: v.TotalPages}, err)
}
func (h ContentEditorHandler) projectList(ctx context.Context) error {
	page, size := adminContentPagination(ctx)
	v, err := h.projects.Paginate(ctx, page, size)
	return editorResult(ctx, payloads.APIProjectPage{Projects: mapPayloads(v.Projects, projectItem), TotalCount: v.TotalCount, Page: v.Page, PageSize: v.PageSize, TotalPages: v.TotalPages}, err)
}
func (h ContentEditorHandler) workList(ctx context.Context) error {
	page, size := adminContentPagination(ctx)
	v, err := h.works.Paginate(ctx, page, size)
	return editorResult(ctx, payloads.APIWorkPage{Works: mapPayloads(v.Works, workItem), TotalCount: v.TotalCount, Page: v.Page, PageSize: v.PageSize, TotalPages: v.TotalPages}, err)
}
func (h ContentEditorHandler) articleGet(ctx context.Context) error {
	v, err := h.posts.FindBySlug(ctx, anduin.PathValue(ctx, "slug"))
	return editorResult(ctx, articleItem(v), err)
}
func (h ContentEditorHandler) projectGet(ctx context.Context) error {
	v, err := h.projects.FindBySlug(ctx, anduin.PathValue(ctx, "slug"))
	return editorResult(ctx, projectItem(v), err)
}
func (h ContentEditorHandler) workGet(ctx context.Context) error {
	v, err := h.works.FindBySlug(ctx, anduin.PathValue(ctx, "slug"))
	return editorResult(ctx, workItem(v), err)
}

// Reject null as well as unknown fields: omitted fields are the only no-op.
func decodeEditor(ctx context.Context, dst any) error {
	body, err := io.ReadAll(io.LimitReader(anduin.Request(ctx).Body, (1<<20)+1))
	if err != nil || len(body) > 1<<20 {
		return errors.New("invalid body")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil || len(fields) == 0 {
		return errors.New("invalid object")
	}
	for _, value := range fields {
		if bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			return errors.New("null field")
		}
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
func (h ContentEditorHandler) articlePatch(ctx context.Context) error {
	var p payloads.APIArticleEdit
	if decodeEditor(ctx, &p) != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	result, err := h.editor.EditArticle(ctx, anduin.PathValue(ctx, "slug"), services.ArticleEdit(p))
	return editorResult(ctx, articleItem(result), err)
}
func (h ContentEditorHandler) projectPatch(ctx context.Context) error {
	var p payloads.APIProjectEdit
	if decodeEditor(ctx, &p) != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	result, err := h.editor.EditProject(ctx, anduin.PathValue(ctx, "slug"), services.ProjectEdit(p))
	return editorResult(ctx, projectItem(result), err)
}
func (h ContentEditorHandler) workPatch(ctx context.Context) error {
	var p payloads.APIWorkEdit
	if decodeEditor(ctx, &p) != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	result, err := h.editor.EditWork(ctx, anduin.PathValue(ctx, "slug"), services.WorkEdit(p))
	return editorResult(ctx, workItem(result), err)
}
