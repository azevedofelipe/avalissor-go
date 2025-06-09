-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, username, hashed_password)
VALUES (
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

-- name: DeleteUsers :exec
DELETE FROM users;
