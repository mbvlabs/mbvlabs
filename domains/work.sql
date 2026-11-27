-- name: GetWork :one
SELECT * FROM works WHERE id = $1;

-- name: GetWorkBySlug :one
SELECT * FROM works WHERE slug = $1;

-- name: GetWorkBySlugForUpdate :one
SELECT * FROM works WHERE slug = $1 FOR UPDATE;

-- name: ListWorks :many
SELECT * FROM works;

-- name: LatestWorks :many
SELECT * FROM works ORDER BY created_at DESC, id DESC LIMIT $1::bigint;

-- name: ListPublishedWorks :many
SELECT * FROM works WHERE status = 'published'
ORDER BY completed_at DESC NULLS LAST, started_at DESC NULLS LAST, published_at DESC;

-- name: ListFeaturedWorks :many
SELECT * FROM works WHERE status = 'published' AND is_featured
ORDER BY completed_at DESC NULLS LAST, started_at DESC NULLS LAST, published_at DESC LIMIT $1;

-- name: PageWorks :many
SELECT * FROM works LIMIT $1::bigint OFFSET $2::bigint;

-- name: PagePublishedWorks :many
SELECT * FROM works WHERE status = 'published'
ORDER BY completed_at DESC NULLS LAST, started_at DESC NULLS LAST, published_at DESC LIMIT $1::bigint OFFSET $2::bigint;

-- name: CountWorks :one
SELECT count(*) FROM works;

-- name: CountPublishedWorks :one
SELECT count(*) FROM works WHERE status = 'published';

-- name: CountFeaturedWorks :one
SELECT count(*) FROM works WHERE is_featured AND ($1::integer <= 0 OR id <> $1);

-- name: LockFeaturedWorks :exec
SELECT pg_advisory_xact_lock(hashtext('mbvlabs'), hashtext('featured-works'));

-- name: CreateWork :one
INSERT INTO works (created_at, updated_at, title, slug, client_name, client_industry, client_url, client_logo_url, summary, cover_image_url, specialisms, platforms, technologies, challenge, approach, deliverables, outcome, content, started_at, completed_at, status, published_at, is_featured)
VALUES ($1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
RETURNING *;

-- name: UpsertWork :one
INSERT INTO works (created_at, updated_at, title, slug, client_name, client_industry, client_url, client_logo_url, summary, cover_image_url, specialisms, platforms, technologies, challenge, approach, deliverables, outcome, content, started_at, completed_at, status, published_at, is_featured)
VALUES ($1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
ON CONFLICT (id) DO UPDATE SET title=excluded.title,slug=excluded.slug,client_name=excluded.client_name,client_industry=excluded.client_industry,client_url=excluded.client_url,client_logo_url=excluded.client_logo_url,summary=excluded.summary,cover_image_url=excluded.cover_image_url,specialisms=excluded.specialisms,platforms=excluded.platforms,technologies=excluded.technologies,challenge=excluded.challenge,approach=excluded.approach,deliverables=excluded.deliverables,outcome=excluded.outcome,content=excluded.content,started_at=excluded.started_at,completed_at=excluded.completed_at,status=excluded.status,published_at=excluded.published_at,is_featured=excluded.is_featured
RETURNING *;

-- name: UpdateWork :one
UPDATE works SET updated_at=$1,title=$2,slug=$3,client_name=$4,client_industry=$5,client_url=$6,client_logo_url=$7,summary=$8,cover_image_url=$9,specialisms=$10,platforms=$11,technologies=$12,challenge=$13,approach=$14,deliverables=$15,outcome=$16,content=$17,started_at=$18,completed_at=$19,status=$20,published_at=$21,is_featured=$22 WHERE id=$23 RETURNING *;

-- name: EditWork :one
UPDATE works SET updated_at=$2,title=$3,summary=$4,content=$5,client_url=$6,client_logo_url=$7,cover_image_url=$8,challenge=$9,approach=$10,deliverables=$11,outcome=$12 WHERE id=$1 RETURNING *;

-- name: DeleteWork :exec
DELETE FROM works WHERE id=$1;
