// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: campuses.sql

package database

import (
	"context"
	"database/sql"
)

const createCampus = `-- name: CreateCampus :one
INSERT INTO campus (created_at, updated_at, name, location, college_id) 
VALUES (
	NOW(),
	NOW(),
	$1,
	$2,
	$3
)
RETURNING id, name, location, college_id, created_at, updated_at
`

type CreateCampusParams struct {
	Name      string
	Location  sql.NullString
	CollegeID int32
}

func (q *Queries) CreateCampus(ctx context.Context, arg CreateCampusParams) (Campus, error) {
	row := q.db.QueryRowContext(ctx, createCampus, arg.Name, arg.Location, arg.CollegeID)
	var i Campus
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Location,
		&i.CollegeID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCampus = `-- name: DeleteCampus :exec
DELETE FROM campus
WHERE id = $1
`

func (q *Queries) DeleteCampus(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteCampus, id)
	return err
}

const getCampusID = `-- name: GetCampusID :one
SELECT id, name, location, college_id, created_at, updated_at FROM campus
WHERE id = $1
`

func (q *Queries) GetCampusID(ctx context.Context, id int32) (Campus, error) {
	row := q.db.QueryRowContext(ctx, getCampusID, id)
	var i Campus
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Location,
		&i.CollegeID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCampuses = `-- name: GetCampuses :many
SELECT id, name, location, college_id, created_at, updated_at FROM campus
`

func (q *Queries) GetCampuses(ctx context.Context) ([]Campus, error) {
	rows, err := q.db.QueryContext(ctx, getCampuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Campus
	for rows.Next() {
		var i Campus
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Location,
			&i.CollegeID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCollegeCampuses = `-- name: GetCollegeCampuses :many
SELECT id, name, location, college_id, created_at, updated_at FROM campus
WHERE college_id = $1
`

func (q *Queries) GetCollegeCampuses(ctx context.Context, collegeID int32) ([]Campus, error) {
	rows, err := q.db.QueryContext(ctx, getCollegeCampuses, collegeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Campus
	for rows.Next() {
		var i Campus
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Location,
			&i.CollegeID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
