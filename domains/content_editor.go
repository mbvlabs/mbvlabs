package domains

import (
	"context"
	"errors"
	"time"
)

// EditEditorial locks and validates the current row before writing only editorial columns.
// Call it on a transaction-bound store.
func (bp *BlogPosts) EditEditorial(ctx context.Context, slug string, edit func(*BlogPostEntity)) (BlogPostEntity, error) {
	e, err := scanBlogPost(bp.db.QueryRow(ctx, "SELECT "+blogPostColumns+" FROM blog_posts WHERE slug=$1 FOR UPDATE", slug))
	if err != nil {
		return e, err
	}
	edit(&e)
	if err := e.Validate(); err != nil {
		return e, errors.Join(ErrDomainValidation, err)
	}
	return scanBlogPost(bp.db.QueryRow(ctx, `UPDATE blog_posts SET updated_at=$2,title=$3,excerpt=$4,body=$5,cover_image_url=$6,tags=$7 WHERE id=$1 RETURNING `+blogPostColumns,
		e.ID, time.Now(), e.Title, e.Excerpt, e.Body, e.CoverImageUrl, e.Tags))
}

func (p *Projects) EditEditorial(ctx context.Context, slug string, edit func(*ProjectEntity)) (ProjectEntity, error) {
	e, err := scanProject(p.db.QueryRow(ctx, "SELECT "+projectColumns+" FROM projects WHERE slug=$1 FOR UPDATE", slug))
	if err != nil {
		return e, err
	}
	edit(&e)
	if err := e.Validate(); err != nil {
		return e, errors.Join(ErrDomainValidation, err)
	}
	return scanProject(p.db.QueryRow(ctx, `UPDATE projects SET updated_at=$2,name=$3,tagline=$4,description=$5,repository_url=$6,live_url=$7,image_url=$8,technologies=$9 WHERE id=$1 RETURNING `+projectColumns,
		e.ID, time.Now(), e.Name, e.Tagline, e.Description, e.RepositoryURL, e.LiveURL, e.ImageURL, e.Technologies))
}

func (w *Works) EditEditorial(ctx context.Context, slug string, edit func(*WorkEntity)) (WorkEntity, error) {
	e, err := scanWork(w.db.QueryRow(ctx, "SELECT "+workColumns+" FROM works WHERE slug=$1 FOR UPDATE", slug))
	if err != nil {
		return e, err
	}
	edit(&e)
	if err := e.Validate(); err != nil {
		return e, errors.Join(ErrDomainValidation, err)
	}
	return scanWork(w.db.QueryRow(ctx, `UPDATE works SET updated_at=$2,title=$3,summary=$4,content=$5,client_url=$6,client_logo_url=$7,cover_image_url=$8,challenge=$9,approach=$10,deliverables=$11,outcome=$12 WHERE id=$1 RETURNING `+workColumns,
		e.ID, time.Now(), e.Title, e.Summary, e.Content, e.ClientURL, e.ClientLogoURL, e.CoverImageURL, e.Challenge, e.Approach, e.Deliverables, e.Outcome))
}
