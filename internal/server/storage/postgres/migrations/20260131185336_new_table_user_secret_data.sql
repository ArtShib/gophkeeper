-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS gophkeeper;

CREATE TABLE IF NOT EXISTS gophkeeper.users
(
    id            BIGSERIAL PRIMARY KEY,
    login         TEXT      NOT NULL UNIQUE,
    pass_hash     TEXT      NOT NULL,
    created_at    BIGINT    NOT NULL
);
CREATE INDEX IF NOT EXISTS gophkeeper_users_login_idx ON gophkeeper.users (login);

DO $$ BEGIN
CREATE TYPE gophkeeper.data_type AS ENUM
    ('credit_card', 'text_data', 'credentials', 'binary_data');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS gophkeeper.data
(
    id            BIGSERIAL,
    secret_id     TEXT      NOT NULL,
    owner_id      BIGINT    NOT NULL REFERENCES gophkeeper.users(id),
    type          gophkeeper.data_type NOT NULL,
    data          BYTEA     NOT NULL,
    metadata      JSONB,
    created_at    BIGINT    NOT NULL,
    updated_at    BIGINT,
    is_deleted    BOOL      DEFAULT FALSE,
    PRIMARY KEY (id, type),
    UNIQUE (owner_id, secret_id, type)
    ) PARTITION BY LIST (type);
CREATE INDEX IF NOT EXISTS gophkeeper_data_owner_idx ON gophkeeper.data (owner_id);
CREATE INDEX IF NOT EXISTS gophkeeper_data_secret_id_idx ON gophkeeper.data (secret_id);


CREATE TABLE IF NOT EXISTS gophkeeper.data_credit_card PARTITION OF gophkeeper.data
    FOR VALUES IN ('credit_card');
CREATE TABLE IF NOT EXISTS gophkeeper.data_text_data PARTITION OF gophkeeper.data
    FOR VALUES IN ('text_data');
CREATE TABLE IF NOT EXISTS gophkeeper.data_credentials PARTITION OF gophkeeper.data
    FOR VALUES IN ('credentials');
CREATE TABLE IF NOT EXISTS gophkeeper.data_binary_data PARTITION OF gophkeeper.data
    FOR VALUES IN ('binary_data');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists gophkeeper.data;
drop table if exists gophkeeper.users;
drop type if exists gophkeeper.data_type;
drop schema if exists gophkeeper cascade;
-- +goose StatementEnd
