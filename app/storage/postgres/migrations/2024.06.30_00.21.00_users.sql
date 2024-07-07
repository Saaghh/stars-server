-- +migrate Up
CREATE TABLE users(
    id serial primary key,
    username varchar not null,
    password_hash varchar not null,
    created_at timestamp with time zone not null default now()
);

-- +migrate Down
DROP TABLE users;