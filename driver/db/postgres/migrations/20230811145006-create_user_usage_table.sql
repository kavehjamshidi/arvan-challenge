-- +migrate Up
CREATE TABLE IF NOT EXISTS user_usage
(
    id          BIGSERIAL    NOT NULL PRIMARY KEY,
    user_id     VARCHAR(128) NOT NULL UNIQUE,
    quota       BIGINT       NOT NULL,
    quota_usage BIGINT       NOT NULL DEFAULT 0,
    start_date  TIMESTAMP    NOT NULL,
    end_date    TIMESTAMP    NOT NULL,
    created_at  TIMESTAMP    NOT NULL,
    updated_at  TIMESTAMP    NOT NULL,
    deleted_at  TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +migrate Down
DROP TABLE IF EXISTS user_usage;