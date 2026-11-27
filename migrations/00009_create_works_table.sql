-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS works (
    id SERIAL PRIMARY KEY,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    client_name VARCHAR(255),
    client_industry VARCHAR(255),
    client_url TEXT,
    client_logo_url TEXT,
    summary TEXT NOT NULL,
    cover_image_url TEXT,
    specialisms TEXT[] NOT NULL,
    platforms TEXT[] NOT NULL,
    technologies TEXT[] NOT NULL,
    challenge TEXT NOT NULL,
    approach TEXT NOT NULL,
    deliverables TEXT NOT NULL,
    outcome TEXT NOT NULL,
    content TEXT NOT NULL,
    started_at DATE,
    completed_at DATE,
    status VARCHAR(50) NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE,
    is_featured BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS works;
-- +goose StatementEnd
