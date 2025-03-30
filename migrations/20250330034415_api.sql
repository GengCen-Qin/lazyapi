-- +goose Up
-- SQL in this section is executed when the migration is applied
CREATE TABLE IF NOT EXISTS apis (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    params TEXT NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP TABLE IF EXISTS apis;
