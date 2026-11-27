package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"mbvlabs/domains"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/telemetry"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type ContentEditorHandler struct {
	posts      *domains.BlogPosts
	projects   *domains.Projects
	works      *domains.Works
	transactor storage.Transactor
}

func NewContentEditorHandler(posts *domains.BlogPosts, projects *domains.Projects, works *domains.Works, transactor storage.Transactor) ContentEditorHandler {
	return ContentEditorHandler{posts, projects, works, transactor}
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
		return anduin.JSON(ctx, http.StatusBadRequest, apiError{"validation failed", fields.ToMap()})
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
	return editorResult(ctx, v, err)
}
func (h ContentEditorHandler) projectList(ctx context.Context) error {
	page, size := adminContentPagination(ctx)
	v, err := h.projects.Paginate(ctx, page, size)
	return editorResult(ctx, v, err)
}
func (h ContentEditorHandler) workList(ctx context.Context) error {
	page, size := adminContentPagination(ctx)
	v, err := h.works.Paginate(ctx, page, size)
	return editorResult(ctx, v, err)
}
func (h ContentEditorHandler) articleGet(ctx context.Context) error {
	v, err := h.posts.FindBySlug(ctx, anduin.PathValue(ctx, "slug"))
	return editorResult(ctx, v, err)
}
func (h ContentEditorHandler) projectGet(ctx context.Context) error {
	v, err := h.projects.FindBySlug(ctx, anduin.PathValue(ctx, "slug"))
	return editorResult(ctx, v, err)
}
func (h ContentEditorHandler) workGet(ctx context.Context) error {
	v, err := h.works.FindBySlug(ctx, anduin.PathValue(ctx, "slug"))
	return editorResult(ctx, v, err)
}

type articleEdit struct {
	Title         *string   `json:"title"`
	Excerpt       *string   `json:"excerpt"`
	Body          *string   `json:"body"`
	CoverImageURL *string   `json:"coverImageUrl"`
	Tags          *[]string `json:"tags"`
}
type projectEdit struct {
	Name          *string   `json:"name"`
	Tagline       *string   `json:"tagline"`
	Description   *string   `json:"description"`
	RepositoryURL *string   `json:"repositoryUrl"`
	LiveURL       *string   `json:"liveUrl"`
	ImageURL      *string   `json:"imageUrl"`
	Technologies  *[]string `json:"technologies"`
}
type workEdit struct {
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
func setString(dst *string, src *string) {
	if src != nil {
		*dst = *src
	}
}
func setOptional(dst *sql.NullString, src *string) {
	if src != nil {
		*dst = sql.NullString{String: *src, Valid: *src != ""}
	}
}

func (h ContentEditorHandler) articlePatch(ctx context.Context) error {
	var p articleEdit
	if decodeEditor(ctx, &p) != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	var result domains.BlogPostEntity
	err := storage.AtomicOps(ctx, h.transactor, func(ctx context.Context, tx storage.Tx) error {
		var err error
		result, err = h.posts.InTx(tx).EditEditorial(ctx, anduin.PathValue(ctx, "slug"), func(e *domains.BlogPostEntity) {
			setString(&e.Title, p.Title)
			setOptional(&e.Excerpt, p.Excerpt)
			setString(&e.Body, p.Body)
			setOptional(&e.CoverImageUrl, p.CoverImageURL)
			if p.Tags != nil {
				e.Tags, _ = json.Marshal(*p.Tags)
			}
		})
		return err
	})
	return editorResult(ctx, result, err)
}
func (h ContentEditorHandler) projectPatch(ctx context.Context) error {
	var p projectEdit
	if decodeEditor(ctx, &p) != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	var result domains.ProjectEntity
	err := storage.AtomicOps(ctx, h.transactor, func(ctx context.Context, tx storage.Tx) error {
		var err error
		result, err = h.projects.InTx(tx).EditEditorial(ctx, anduin.PathValue(ctx, "slug"), func(e *domains.ProjectEntity) {
			setString(&e.Name, p.Name)
			setString(&e.Tagline, p.Tagline)
			setOptional(&e.Description, p.Description)
			setOptional(&e.RepositoryURL, p.RepositoryURL)
			setOptional(&e.LiveURL, p.LiveURL)
			setOptional(&e.ImageURL, p.ImageURL)
			if p.Technologies != nil {
				e.Technologies, _ = json.Marshal(*p.Technologies)
			}
		})
		return err
	})
	return editorResult(ctx, result, err)
}
func (h ContentEditorHandler) workPatch(ctx context.Context) error {
	var p workEdit
	if decodeEditor(ctx, &p) != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	var result domains.WorkEntity
	err := storage.AtomicOps(ctx, h.transactor, func(ctx context.Context, tx storage.Tx) error {
		var err error
		result, err = h.works.InTx(tx).EditEditorial(ctx, anduin.PathValue(ctx, "slug"), func(e *domains.WorkEntity) {
			setString(&e.Title, p.Title)
			setString(&e.Summary, p.Summary)
			setString(&e.Content, p.Content)
			setOptional(&e.ClientURL, p.ClientURL)
			setOptional(&e.ClientLogoURL, p.ClientLogoURL)
			setOptional(&e.CoverImageURL, p.CoverImageURL)
			setString(&e.Challenge, p.Challenge)
			setString(&e.Approach, p.Approach)
			setString(&e.Deliverables, p.Deliverables)
			setString(&e.Outcome, p.Outcome)
		})
		return err
	})
	return editorResult(ctx, result, err)
}
