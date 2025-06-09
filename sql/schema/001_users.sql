-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;
