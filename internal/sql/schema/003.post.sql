CREATE TABLE IF NOT EXISTS post (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT,
	image TEXT,
	user_id BIGINT NOT NULL,
	minimum_fund DECIMAL(28, 6) NOT NULL,
	fund_raised DECIMAL(28, 6) NOT NULL DEFAULT 0,
	deadline_date DATE NOT NULL,
	created_at TIMESTAMPTZ,
    category TEXT NOT NULL,
	FOREIGN KEY (user_id)
		REFERENCES user_t(id)
);