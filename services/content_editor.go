package services

import (
	"context"
	"encoding/json"

	"mbvlabs/domains"

	"github.com/mbvlabs/andurel/pkg/storage"
)

type ContentEditor struct {
	posts      *domains.BlogPosts
	projects   *domains.Projects
	works      *domains.Works
	transactor storage.Transactor
}

func NewContentEditor(posts *domains.BlogPosts, projects *domains.Projects, works *domains.Works, transactor storage.Transactor) ContentEditor {
	return ContentEditor{posts, projects, works, transactor}
}

type ArticleEdit struct {
	Title         *string
	Excerpt       *string
	Body          *string
	CoverImageURL *string
	Tags          *[]string
}
type ProjectEdit struct {
	Name          *string
	Tagline       *string
	Description   *string
	RepositoryURL *string
	LiveURL       *string
	ImageURL      *string
	Technologies  *[]string
}
type WorkEdit struct {
	Title         *string
	Summary       *string
	Content       *string
	ClientURL     *string
	ClientLogoURL *string
	CoverImageURL *string
	Challenge     *string
	Approach      *string
	Deliverables  *string
	Outcome       *string
}

func setString(dst *string, src *string) {
	if src != nil {
		*dst = *src
	}
}
func setOptional(dst *domains.OptionalString, src *string) {
	if src != nil {
		*dst = domains.OptionalString{String: *src, Valid: *src != ""}
	}
}

func (h ContentEditor) EditArticle(ctx context.Context, slug string, p ArticleEdit) (domains.BlogPostEntity, error) {
	var result domains.BlogPostEntity
	err := storage.AtomicOps(ctx, h.transactor, func(ctx context.Context, tx storage.Tx) error {
		var err error
		result, err = h.posts.InTx(tx).EditEditorial(ctx, slug, func(e *domains.BlogPostEntity) {
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
	return result, err
}
func (h ContentEditor) EditProject(ctx context.Context, slug string, p ProjectEdit) (domains.ProjectEntity, error) {
	var result domains.ProjectEntity
	err := storage.AtomicOps(ctx, h.transactor, func(ctx context.Context, tx storage.Tx) error {
		var err error
		result, err = h.projects.InTx(tx).EditEditorial(ctx, slug, func(e *domains.ProjectEntity) {
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
	return result, err
}
func (h ContentEditor) EditWork(ctx context.Context, slug string, p WorkEdit) (domains.WorkEntity, error) {
	var result domains.WorkEntity
	err := storage.AtomicOps(ctx, h.transactor, func(ctx context.Context, tx storage.Tx) error {
		var err error
		result, err = h.works.InTx(tx).EditEditorial(ctx, slug, func(e *domains.WorkEntity) {
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
	return result, err
}
