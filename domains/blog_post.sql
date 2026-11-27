-- name: GetBlogPost :one
SELECT * FROM blog_posts WHERE id = $1;

-- name: GetBlogPostBySlug :one
SELECT * FROM blog_posts WHERE slug = $1;

-- name: GetBlogPostBySlugForUpdate :one
SELECT * FROM blog_posts WHERE slug = $1 FOR UPDATE;

-- name: ListPublishedBlogPostLinks :many
SELECT slug, title FROM blog_posts
WHERE status = 'published' AND slug = ANY($1::text[])
ORDER BY array_position($1::text[], slug::text);

-- name: CreateBlogPost :one
INSERT INTO blog_posts (created_at, updated_at, title, slug, excerpt, body, status, cover_image_url, tags, published_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING *;

-- name: UpsertBlogPost :one
INSERT INTO blog_posts (created_at, updated_at, title, slug, excerpt, body, status, cover_image_url, tags, published_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
ON CONFLICT (id) DO UPDATE SET title=excluded.title, slug=excluded.slug, excerpt=excluded.excerpt, body=excluded.body, status=excluded.status, cover_image_url=excluded.cover_image_url, tags=excluded.tags, published_at=excluded.published_at
RETURNING *;

-- name: UpdateBlogPost :one
UPDATE blog_posts SET updated_at=$2, title=$3, slug=$4, excerpt=$5, body=$6, status=$7, cover_image_url=$8, tags=$9, published_at=$10, publication_schedule=$11 WHERE id=$1 RETURNING *;

-- name: UpdateBlogPostEditorialFields :one
UPDATE blog_posts SET updated_at=$2, title=$3, excerpt=$4, body=$5, cover_image_url=$6, tags=$7 WHERE id=$1 RETURNING *;

-- name: UpdateBlogPostPublicationSchedule :one
UPDATE blog_posts SET publication_schedule=$2 WHERE id=$1 RETURNING *;

-- name: PublishBlogPostIfDraft :execrows
UPDATE blog_posts SET status='published', published_at=$2, publication_schedule=NULL, updated_at=$3 WHERE id=$1 AND status='draft' AND published_at IS NULL AND publication_schedule=$4;

-- name: DeleteBlogPost :exec
DELETE FROM blog_posts WHERE id=$1;

-- name: ListBlogPosts :many
SELECT * FROM blog_posts;

-- name: ListLatestBlogPosts :many
SELECT * FROM blog_posts ORDER BY created_at DESC, id DESC LIMIT $1::bigint;

-- name: ListDraftOrPublishedBlogPosts :many
SELECT * FROM blog_posts WHERE status IN ('draft', 'published') ORDER BY created_at DESC, id DESC;

-- name: ListPublishedBlogPosts :many
SELECT * FROM blog_posts WHERE status='published' ORDER BY published_at DESC;

-- name: PageBlogPosts :many
SELECT * FROM blog_posts LIMIT $1::bigint OFFSET $2::bigint;

-- name: PagePublishedBlogPosts :many
SELECT * FROM blog_posts WHERE status='published' ORDER BY published_at DESC LIMIT $1::bigint OFFSET $2::bigint;

-- name: CountBlogPosts :one
SELECT count(*) FROM blog_posts;

-- name: CountPublishedBlogPosts :one
SELECT count(*) FROM blog_posts WHERE status='published';
