-- +goose Up
ALTER TABLE blog_posts
ADD COLUMN publication_schedule JSONB;

-- +goose Down
ALTER TABLE blog_posts
DROP COLUMN publication_schedule;
