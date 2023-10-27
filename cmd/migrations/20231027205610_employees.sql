-- +goose Up
-- +goose StatementBegin
CREATE TABLE employees (
    id serial PRIMARY KEY,
    first_name text,
    last_name text,
    email text,
    hire_date date
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table employees;
-- +goose StatementEnd
