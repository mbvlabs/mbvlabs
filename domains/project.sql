-- name: GetProject :one
SELECT * FROM projects WHERE id = $1;

-- name: GetProjectBySlug :one
SELECT * FROM projects WHERE slug = $1;

-- name: GetProjectBySlugForUpdate :one
SELECT * FROM projects WHERE slug = $1 FOR UPDATE;

-- name: ListProjects :many
SELECT * FROM projects;

-- name: LatestProjects :many
SELECT * FROM projects ORDER BY created_at DESC, id DESC LIMIT $1::bigint;

-- name: ListPublishedProjects :many
SELECT * FROM projects WHERE published_at IS NOT NULL ORDER BY published_at DESC;

-- name: ListFeaturedProjects :many
SELECT * FROM projects WHERE published_at IS NOT NULL AND is_featured ORDER BY published_at DESC LIMIT $1;

-- name: PageProjects :many
SELECT * FROM projects LIMIT $1::bigint OFFSET $2::bigint;

-- name: PagePublishedProjects :many
SELECT * FROM projects WHERE published_at IS NOT NULL ORDER BY published_at DESC LIMIT $1::bigint OFFSET $2::bigint;

-- name: CountProjects :one
SELECT count(*) FROM projects;

-- name: CountPublishedProjects :one
SELECT count(*) FROM projects WHERE published_at IS NOT NULL;

-- name: CountFeaturedProjects :one
SELECT count(*) FROM projects WHERE is_featured AND ($1::integer <= 0 OR id <> $1);

-- name: LockFeaturedProjects :exec
SELECT pg_advisory_xact_lock(hashtext('mbvlabs'), hashtext('featured-projects'));

-- name: CreateProject :one
INSERT INTO projects (created_at,updated_at,name,slug,tagline,description,project_type,repository_url,live_url,image_url,technologies,started_at,launched_at,published_at,is_featured)
VALUES ($1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING *;

-- name: UpsertProject :one
INSERT INTO projects (created_at,updated_at,name,slug,tagline,description,project_type,repository_url,live_url,image_url,technologies,started_at,launched_at,published_at,is_featured)
VALUES ($1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
ON CONFLICT (id) DO UPDATE SET name=excluded.name,slug=excluded.slug,tagline=excluded.tagline,description=excluded.description,project_type=excluded.project_type,repository_url=excluded.repository_url,live_url=excluded.live_url,image_url=excluded.image_url,technologies=excluded.technologies,started_at=excluded.started_at,launched_at=excluded.launched_at,published_at=excluded.published_at,is_featured=excluded.is_featured RETURNING *;

-- name: UpdateProject :one
UPDATE projects SET updated_at=$1,name=$2,slug=$3,tagline=$4,description=$5,project_type=$6,repository_url=$7,live_url=$8,image_url=$9,technologies=$10,started_at=$11,launched_at=$12,published_at=$13,is_featured=$14 WHERE id=$15 RETURNING *;

-- name: EditProject :one
UPDATE projects SET updated_at=$2,name=$3,tagline=$4,description=$5,repository_url=$6,live_url=$7,image_url=$8,technologies=$9 WHERE id=$1 RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id=$1;
