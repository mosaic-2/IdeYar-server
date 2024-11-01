CREATE TABLE IF NOT EXISTS account (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(571) NOT NULL UNIQUE,
    username VARCHAR(32) NOT NULL UNIQUE,
    creation_date DATE NOT NULL DEFAULT CURRENT_DATE,
    password CHAR(60) NOT NULL
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_account_username
ON account
USING gin (username gin_trgm_ops);