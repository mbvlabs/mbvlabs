package domains

import "context"

func (w *Works) Latest(ctx context.Context, limit int) ([]WorkEntity, error) {
	return w.list(ctx, "ORDER BY created_at DESC, id DESC LIMIT $1", limit)
}

func (p *Projects) Latest(ctx context.Context, limit int) ([]ProjectEntity, error) {
	return p.list(ctx, "ORDER BY created_at DESC, id DESC LIMIT $1", limit)
}

func (bp *BlogPosts) Latest(ctx context.Context, limit int) ([]BlogPostEntity, error) {
	return bp.list(ctx, "ORDER BY created_at DESC, id DESC LIMIT $1", limit)
}

func (p *ProjectInquiries) Latest(ctx context.Context, limit int) ([]ProjectInquiryEntity, error) {
	return p.list(ctx, "SELECT "+inquiryColumns+" FROM project_inquiries ORDER BY created_at DESC, id DESC LIMIT $1", limit)
}
