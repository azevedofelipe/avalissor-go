-- name: CreateCampus :one
INSERT INTO campus (created_at, updated_at, name, location, college_id) 
VALUES (
	NOW(),
	NOW(),
	$1,
	$2,
	$3
)
RETURNING *;

-- name: DeleteCampus :exec
DELETE FROM campus
WHERE id = $1;
