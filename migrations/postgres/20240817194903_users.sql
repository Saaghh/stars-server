-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id serial primary key,
    username varchar not null,
    password_hash varchar not null,
    created_at timestamp with time zone not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
