-- +goose Up
CREATE TABLE IF NOT EXISTS error
(
    id      SERIAL PRIMARY KEY,
    message TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS error;
