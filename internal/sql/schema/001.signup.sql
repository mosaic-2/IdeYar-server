CREATE TABLE IF NOT EXISTS sign_up (
    id SERIAL PRIMARY KEY,
    email VARCHAR(571) NOT NULL,
    username VARCHAR(32) NOT NULL,
    password CHAR(60) NOT NULL,
    verification_code VARCHAR(6) NOT NULL,
    expire TIMESTAMP NOT NULL
);