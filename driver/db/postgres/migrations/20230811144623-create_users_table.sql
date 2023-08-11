-- +migrate Up
CREATE TABLE IF NOT EXISTS users
(
    id         VARCHAR(128) NOT NULL PRIMARY KEY,
    rate_limit BIGINT       NOT NULL,
    quota      BIGINT       NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL,
    deleted_at TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS users;