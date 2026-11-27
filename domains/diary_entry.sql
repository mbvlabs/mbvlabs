-- name: GetDiaryEntry :one
SELECT * FROM diary_entries WHERE id=$1;

-- name: GetDiaryEntryByDate :one
SELECT * FROM diary_entries WHERE entry_date=$1;

-- name: CreateDiaryEntry :one
INSERT INTO diary_entries (created_at, updated_at, entry_date, morning_thoughts, evening_thoughts)
VALUES ($1,$1,$2,$3,$4) RETURNING *;

-- name: UpsertDiaryEntry :one
INSERT INTO diary_entries (created_at, updated_at, entry_date, morning_thoughts, evening_thoughts)
VALUES ($1,$1,$2,$3,$4)
ON CONFLICT (entry_date) DO UPDATE SET updated_at=excluded.updated_at, morning_thoughts=excluded.morning_thoughts, evening_thoughts=excluded.evening_thoughts RETURNING *;

-- name: FindOrCreateDiaryEntry :one
INSERT INTO diary_entries (created_at, updated_at, entry_date, morning_thoughts, evening_thoughts)
VALUES ($1,$1,$2,$3,$4)
ON CONFLICT (entry_date) DO UPDATE SET entry_date=excluded.entry_date RETURNING *;

-- name: UpdateDiaryEntry :one
UPDATE diary_entries SET updated_at=$2, entry_date=$3, morning_thoughts=$4, evening_thoughts=$5 WHERE id=$1 RETURNING *;

-- name: DeleteDiaryEntry :exec
DELETE FROM diary_entries WHERE id=$1;

-- name: ListDiaryEntries :many
SELECT * FROM diary_entries;

-- name: ListDiaryEntriesBetweenDates :many
SELECT * FROM diary_entries WHERE entry_date >= $1 AND entry_date <= $2 ORDER BY entry_date ASC;

-- name: ListDiaryEntriesNewestFirst :many
SELECT * FROM diary_entries ORDER BY entry_date DESC LIMIT $1::bigint OFFSET $2::bigint;

-- name: CountDiaryEntries :one
SELECT count(*) FROM diary_entries;
