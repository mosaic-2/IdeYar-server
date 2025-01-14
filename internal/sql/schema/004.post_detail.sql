CREATE TABLE IF NOT EXISTS post_detail (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	order_c INT NOT NULL,
	title TEXT,
	description TEXT,
	image TEXT,
	post_id BIGINT NOT NULL,
	FOREIGN KEY (post_id)
		REFERENCES post(id)
);

CREATE INDEX idx_post_id ON post_detail(post_id);