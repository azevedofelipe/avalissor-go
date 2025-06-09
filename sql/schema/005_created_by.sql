-- +goose Up
ALTER TABLE college
ADD COLUMN created_by UUID NOT NULL REFERENCES users(id);

-- +goose Down
ALTER TABLE college
DROP COLUMN created_by;
