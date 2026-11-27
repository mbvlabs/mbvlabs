-- +goose Up
UPDATE river_job
SET state = 'cancelled', finalized_at = now()
WHERE kind = 'send_diary_reminder'
  AND state IN ('available', 'pending', 'retryable', 'scheduled');

-- +goose Down
-- Do not revive cancelled reminders on rollback.
