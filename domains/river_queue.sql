-- name: GetRiverJob :one
SELECT id, state::text, attempt, max_attempts, attempted_at, created_at, finalized_at, scheduled_at, priority, args::text, attempted_by, ARRAY(SELECT e::text FROM unnest(errors) e) AS errors, kind, metadata::text, queue, tags FROM river_job WHERE id=$1;

-- name: CountRiverJobs :one
SELECT count(*) FROM river_job WHERE ($1::text = '' OR state::text = $1) AND ($2::text = '' OR queue = $2) AND ($3::text = '' OR kind ILIKE '%' || $3 || '%');

-- name: ListRiverJobs :many
SELECT id, state::text, attempt, max_attempts, attempted_at, created_at, finalized_at, scheduled_at, priority, args::text, attempted_by, ARRAY(SELECT e::text FROM unnest(errors) e) AS errors, kind, metadata::text, queue, tags FROM river_job WHERE ($1::text = '' OR state::text = $1) AND ($2::text = '' OR queue = $2) AND ($3::text = '' OR kind ILIKE '%' || $3 || '%') ORDER BY created_at DESC, id DESC LIMIT $4::bigint OFFSET $5::bigint;

-- name: CountRiverJobStates :many
SELECT state::text, count(*) FROM river_job GROUP BY state ORDER BY state;

-- name: ListRiverQueueStats :many
WITH names AS (SELECT name FROM river_queue UNION SELECT queue FROM river_job),
counts AS (SELECT queue,
count(*) FILTER (WHERE state='available') available, count(*) FILTER (WHERE state='cancelled') cancelled,
count(*) FILTER (WHERE state='completed') completed, count(*) FILTER (WHERE state='discarded') discarded,
count(*) FILTER (WHERE state='pending') pending, count(*) FILTER (WHERE state='retryable') retryable,
count(*) FILTER (WHERE state='running') running, count(*) FILTER (WHERE state='scheduled') scheduled,
count(*) total FROM river_job GROUP BY queue)
SELECT names.name, q.created_at, q.updated_at, q.paused_at, coalesce(c.available,0) AS available_count, coalesce(c.cancelled,0) AS cancelled_count,
coalesce(c.completed,0) AS completed_count, coalesce(c.discarded,0) AS discarded_count, coalesce(c.pending,0) AS pending_count, coalesce(c.retryable,0) AS retryable_count,
coalesce(c.running,0) AS running_count, coalesce(c.scheduled,0) AS scheduled_count, coalesce(c.total,0) AS total_count
FROM names LEFT JOIN river_queue q ON q.name=names.name LEFT JOIN counts c ON c.queue=names.name ORDER BY names.name;

-- name: EnsureRiverQueue :exec
INSERT INTO river_queue (name) VALUES ($1) ON CONFLICT (name) DO NOTHING;

-- name: DiscardRiverJob :execrows
UPDATE river_job SET state='discarded', finalized_at=now() WHERE id=$1 AND state NOT IN ('cancelled','completed','discarded');
