-- +goose Up
ALTER TABLE users ADD COLUMN session_version INTEGER NOT NULL DEFAULT 1;

-- Keep unmatched legacy tokens intact. Only password-reset tokens with a
-- matching user can be used by the V2 reset flow.
ALTER TABLE tokens ADD COLUMN user_id uuid REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE tokens ADD COLUMN attempts INTEGER NOT NULL DEFAULT 0;
UPDATE tokens AS t
SET user_id = u.id
FROM users AS u
WHERE t.scope = 'user_password_reset'
  AND lower(t.meta_data->>'email') = lower(u.email)
  AND 1 = (SELECT count(*) FROM users AS candidate
           WHERE lower(candidate.email) = lower(u.email));

CREATE INDEX tokens_user_id_scope_created_at_idx
    ON tokens (user_id, scope, created_at DESC);
CREATE INDEX tokens_scope_hash_idx ON tokens (scope, hash);

-- +goose Down
DROP INDEX tokens_scope_hash_idx;
DROP INDEX tokens_user_id_scope_created_at_idx;
ALTER TABLE tokens DROP COLUMN attempts;
ALTER TABLE tokens DROP COLUMN user_id;
ALTER TABLE users DROP COLUMN session_version;
