-- +migrate Up
CREATE TABLE IF NOT EXISTS user_usage
(
    user_id         VARCHAR(128) NOT NULL PRIMARY KEY,
    quota           BIGINT       NOT NULL,
    quota_usage     BIGINT       NOT NULL DEFAULT 0,
    start_timestamp BIGINT       NOT NULL,
    end_timestamp   BIGINT       NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +migrate Down
DROP TABLE IF EXISTS user_usage;