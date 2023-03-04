CREATE TABLE test (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	title VARCHAR(100),
	description TEXT,
	price FLOAT,
	stock BOOL
);
