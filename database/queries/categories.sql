-- name: QueryCategoryByID :one
select * from categories where id=$1;

-- name: QueryCategorys :many
select * from categories;

-- name: QueryAllCategorys :many
select * from categories;

-- name: InsertCategory :one
insert into
    categories (id, created_at, updated_at, name, slug, description, color)
values
    ($1, now(), now(), $2, $3, $4, $5)
returning *;

-- name: UpdateCategory :one
update categories
    set updated_at=now(), name=$2, slug=$3, description=$4, color=$5
where id = $1
returning *;

-- name: DeleteCategory :exec
delete from categories where id=$1;

-- name: QueryPaginatedCategorys :many
select * from categories 
order by created_at desc 
limit sqlc.arg('limit')::bigint offset sqlc.arg('offset')::bigint;

-- name: CountCategorys :one
select count(*) from categories;

-- name: UpsertCategory :one
insert into
    categories (id, created_at, updated_at, name, slug, description, color)
values
    ($1, now(), now(), $2, $3, $4, $5)
on conflict (id) do update set updated_at=now(), name=excluded.name, slug=excluded.slug, description=excluded.description, color=excluded.color
returning *;
