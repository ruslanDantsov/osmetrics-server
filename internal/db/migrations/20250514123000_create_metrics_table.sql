-- +goose Up
-- SQL-запрос для создания таблицы
CREATE TABLE metrics (
    Id     VARCHAR(50) PRIMARY KEY,
    Type   VARCHAR(10) NOT NULL,
    Delta  BIGINT,
    Value  DOUBLE PRECISION
);

-- +goose Down
-- SQL-запрос для отката (удаления таблицы)
DROP TABLE IF EXISTS metrics;