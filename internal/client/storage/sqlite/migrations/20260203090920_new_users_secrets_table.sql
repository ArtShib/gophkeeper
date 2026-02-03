-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    user_id        INTEGER PRIMARY KEY,
    login          TEXT NOT NULL UNIQUE,
    password_hash  BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS secrets (
    secret_id  TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    type       TEXT NOT NULL,
    data       BLOB NOT NULL,
    metadata   TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER,
    status     TEXT DEFAULT 'new',
    FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_secrets_sync ON secrets (user_id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
drop table if exists secrets;
-- +goose StatementEnd
