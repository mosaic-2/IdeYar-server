CREATE TABLE IF NOT EXISTS fund (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    amount DECIMAL(26, 6),
    FOREIGN KEY (user_id)
        REFERENCES user_t(id),
    FOREIGN KEY (post_id)
        REFERENCES post(id)
);

CREATE INDEX user_id_idx ON fund(user_id);
CREATE INDEX post_id_idx ON fund(post_id);