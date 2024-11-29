CREATE TABLE IF NOT EXISTS user_t (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(571) NOT NULL UNIQUE,
    username VARCHAR(32) NOT NULL UNIQUE,
    password CHAR(60) NOT NULL,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    bio TEXT NOT NULL,
    birthday DATE,
    phone VARCHAR(15),
    profile_image_url VARCHAR(2048)
);
