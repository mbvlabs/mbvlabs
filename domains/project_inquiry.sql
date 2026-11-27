-- name: GetProjectInquiry :one
SELECT * FROM project_inquiries WHERE id = $1;

-- name: CreateProjectInquiry :one
INSERT INTO project_inquiries (created_at, updated_at, name, email, company, role, project_type, timeline, message, source, status, metadata)
VALUES ($1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: UpdateProjectInquiry :one
UPDATE project_inquiries SET updated_at=$2, name=$3, email=$4, company=$5, role=$6, project_type=$7, timeline=$8, message=$9, source=$10, status=$11, metadata=$12 WHERE id=$1 RETURNING *;

-- name: DeleteProjectInquiry :exec
DELETE FROM project_inquiries WHERE id=$1;

-- name: ListProjectInquiries :many
SELECT * FROM project_inquiries;

-- name: PaginateProjectInquiries :many
SELECT * FROM project_inquiries ORDER BY id LIMIT $1::bigint OFFSET $2::bigint;

-- name: LatestProjectInquiries :many
SELECT * FROM project_inquiries ORDER BY created_at DESC, id DESC LIMIT $1::bigint;

-- name: CountProjectInquiries :one
SELECT count(*) FROM project_inquiries;

-- name: IsProjectInquiryDomainBanned :one
SELECT EXISTS (SELECT 1 FROM project_inquiry_domain_bans WHERE domain=$1);

-- name: BanProjectInquiryDomain :exec
INSERT INTO project_inquiry_domain_bans (domain, created_at) VALUES ($1,$2) ON CONFLICT (domain) DO NOTHING;

-- name: UnbanProjectInquiryDomain :exec
DELETE FROM project_inquiry_domain_bans WHERE domain=$1;
