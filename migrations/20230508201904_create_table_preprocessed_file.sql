-- +goose Up
CREATE TABLE IF NOT EXISTS preprocessed_file
(
    name TEXT NOT NULL PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS preprocessed_file;
