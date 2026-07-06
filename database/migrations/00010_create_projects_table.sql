-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    tagline TEXT NOT NULL,
    description TEXT,
    project_type VARCHAR(100) NOT NULL,
    repository_url TEXT,
    live_url TEXT,
    image_url TEXT,
    technologies JSONB NOT NULL,
    started_at DATE,
    launched_at DATE,
    published_at TIMESTAMP WITH TIME ZONE,
    is_featured BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd
