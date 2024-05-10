-- +goose Up
CREATE TABLE foo (
    id   INTEGER NOT NULL PRIMARY KEY,
    bar VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS foo;
