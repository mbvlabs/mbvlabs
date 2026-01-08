-- name: QueryTagByID :one
select * from tags where id=$1;

-- name: QueryTags :many
select * from tags;

-- name: QueryAllTags :many
select * from tags;

-- name: InsertTag :one
insert into
    tags (id, created_at, updated_at, name, slug, description, color)
values
    ($1, now(), now(), $2, $3, $4, $5)
returning *;

-- name: UpdateTag :one
update tags
    set updated_at=now(), name=$2, slug=$3, description=$4, color=$5
where id = $1
returning *;

-- name: DeleteTag :exec
delete from tags where id=$1;

-- name: QueryPaginatedTags :many
select * from tags 
order by created_at desc 
limit sqlc.arg('limit')::bigint offset sqlc.arg('offset')::bigint;

-- name: CountTags :one
select count(*) from tags;

-- name: UpsertTag :one
insert into
    tags (id, created_at, updated_at, name, slug, description, color)
values
    ($1, now(), now(), $2, $3, $4, $5)
on conflict (id) do update set updated_at=now(), name=excluded.name, slug=excluded.slug, description=excluded.description, color=excluded.color
returning *;
