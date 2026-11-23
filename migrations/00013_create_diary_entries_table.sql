-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diary_entries (
    id SERIAL PRIMARY KEY,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    entry_date DATE NOT NULL UNIQUE,
    morning_thoughts TEXT,
    evening_thoughts TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diary_entries;
-- +goose StatementEnd
