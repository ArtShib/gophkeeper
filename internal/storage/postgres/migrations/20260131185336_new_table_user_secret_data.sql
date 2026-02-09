-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users
(
    id            BIGSERIAL PRIMARY KEY,
    login         TEXT      NOT NULL UNIQUE,
    pass_hash     TEXT      NOT NULL,
    created_at    BIGINT    NOT NULL
);
CREATE INDEX IF NOT EXISTS gophkeeper_users_login_idx ON users (login);

DO $$ BEGIN
CREATE TYPE data_type AS ENUM
    ('credit_card', 'text_data', 'credentials', 'binary_data');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS data
(
    id            BIGSERIAL,
    secret_id     TEXT      NOT NULL,
    owner_id      BIGINT    NOT NULL REFERENCES users(id),
    type          data_type NOT NULL,
    data          BYTEA     NOT NULL,
    metadata      JSONB,
    created_at    BIGINT    NOT NULL,
    updated_at    BIGINT,
    is_deleted    BOOL      DEFAULT FALSE,
    PRIMARY KEY (id, type),
    UNIQUE (owner_id, secret_id, type)
    ) PARTITION BY LIST (type);
CREATE INDEX IF NOT EXISTS gophkeeper_data_owner_idx ON data (owner_id);
CREATE INDEX IF NOT EXISTS gophkeeper_data_secret_id_idx ON data (secret_id);


CREATE TABLE IF NOT EXISTS data_credit_card PARTITION OF data
    FOR VALUES IN ('credit_card');
CREATE TABLE IF NOT EXISTS data_text_data PARTITION OF data
    FOR VALUES IN ('text_data');
CREATE TABLE IF NOT EXISTS data_credentials PARTITION OF data
    FOR VALUES IN ('credentials');
CREATE TABLE IF NOT EXISTS data_binary_data PARTITION OF data
    FOR VALUES IN ('binary_data');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists data;
drop table if exists users;
drop type if exists data_type;
-- +goose StatementEnd
