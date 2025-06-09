-- +goose Up
CREATE TABLE campus(
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	location TEXT,
	college_id INT NOT NULL REFERENCES college(id) ON DELETE CASCADE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE campus;
