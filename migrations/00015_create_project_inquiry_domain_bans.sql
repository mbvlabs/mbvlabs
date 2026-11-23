-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_inquiry_domain_bans (
    domain VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_inquiry_domain_bans;
-- +goose StatementEnd
