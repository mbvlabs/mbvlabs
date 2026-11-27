-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS project_inquiries (
    id SERIAL PRIMARY KEY,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    company VARCHAR(255),
    role VARCHAR(255),
    project_type VARCHAR(100),
    timeline VARCHAR(100),
    message TEXT NOT NULL,
    source VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    metadata JSONB NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_inquiries;
-- +goose StatementEnd
