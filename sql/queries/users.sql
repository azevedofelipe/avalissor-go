-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, username, hashed_password)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2,
	$3

)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;
