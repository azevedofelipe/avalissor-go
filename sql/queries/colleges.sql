-- name: GetColleges :many
SELECT * FROM college;

-- name: CreateCollege :one
INSERT INTO college (created_at, updated_at, name_college, created_by) 
VALUES (
	NOW(),
	NOW(),
	$1,
	$2
)
RETURNING *;

-- name: GetCollegeByID :one
SELECT * FROM college
WHERE id = $1;

-- name: DeleteCollegeID :exec
DELETE FROM college
WHERE id = $1;
