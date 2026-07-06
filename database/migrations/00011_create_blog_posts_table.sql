-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS blog_posts (
    id SERIAL PRIMARY KEY,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    excerpt TEXT,
    body TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    cover_image_url TEXT,
    tags JSONB NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blog_posts;
-- +goose StatementEnd
