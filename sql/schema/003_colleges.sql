-- +goose Up
CREATE TABLE college(
	id SERIAL PRIMARY KEY,
	name_college TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE college;
