CREATE TABLE IF NOT EXISTS user_t (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(571) NOT NULL UNIQUE,
    username VARCHAR(32) NOT NULL UNIQUE,
    password CHAR(60) NOT NULL,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_account_username
ON accountz
USING gin (username gin_trgm_ops);