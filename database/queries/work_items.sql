-- name: QueryWorkItemByID :one
select * from work_items where id=$1;

-- name: QueryWorkItems :many
select * from work_items;

-- name: QueryAllWorkItems :many
select * from work_items;

-- name: InsertWorkItem :one
insert into
    work_items (id, created_at, updated_at, title, slug, short_description, content, client, industry, project_date, project_duration, hero_image_url, hero_image_alt, external_url, is_published, is_featured, display_order, meta_title, meta_description, meta_keywords)
values
    ($1, now(), now(), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
returning *;

-- name: UpdateWorkItem :one
update work_items
    set updated_at=now(), title=$2, slug=$3, short_description=$4, content=$5, client=$6, industry=$7, project_date=$8, project_duration=$9, hero_image_url=$10, hero_image_alt=$11, external_url=$12, is_published=$13, is_featured=$14, display_order=$15, meta_title=$16, meta_description=$17, meta_keywords=$18
where id = $1
returning *;

-- name: DeleteWorkItem :exec
delete from work_items where id=$1;

-- name: QueryPaginatedWorkItems :many
select * from work_items 
order by created_at desc 
limit sqlc.arg('limit')::bigint offset sqlc.arg('offset')::bigint;

-- name: CountWorkItems :one
select count(*) from work_items;

-- name: UpsertWorkItem :one
insert into
    work_items (id, created_at, updated_at, title, slug, short_description, content, client, industry, project_date, project_duration, hero_image_url, hero_image_alt, external_url, is_published, is_featured, display_order, meta_title, meta_description, meta_keywords)
values
    ($1, now(), now(), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
on conflict (id) do update set updated_at=now(), title=excluded.title, slug=excluded.slug, short_description=excluded.short_description, content=excluded.content, client=excluded.client, industry=excluded.industry, project_date=excluded.project_date, project_duration=excluded.project_duration, hero_image_url=excluded.hero_image_url, hero_image_alt=excluded.hero_image_alt, external_url=excluded.external_url, is_published=excluded.is_published, is_featured=excluded.is_featured, display_order=excluded.display_order, meta_title=excluded.meta_title, meta_description=excluded.meta_description, meta_keywords=excluded.meta_keywords
returning *;
