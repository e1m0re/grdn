-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE metrics
(
    Id    SERIAL PRIMARY KEY,
    Name  VARCHAR(50) NOT NULL,
    Type  VARCHAR(50) NOT NULL,
    Delta INT,
    Value DOUBLE PRECISION,
    UNIQUE (Name, Type)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE metrics;
-- +goose StatementEnd
