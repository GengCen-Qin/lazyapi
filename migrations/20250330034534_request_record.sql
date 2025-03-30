-- +goose Up
-- SQL in this section is executed when the migration is applied
CREATE TABLE IF NOT EXISTS request_records (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    params TEXT NOT NULL,
    respond TEXT NOT NULL,
    request_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP TABLE IF EXISTS request_records;
