-- +goose Up
-- +goose StatementBegin
CREATE TABLE structure_types (
    id serial primary key,
    name varchar not null
);

CREATE TABLE structures (
    id uuid primary key default uuid_generate_v4(),
    stockpile_id int references stockpiles(id) not null,
    structure_type_id int references structure_types(id) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE structures, structure_types;
-- +goose StatementEnd
