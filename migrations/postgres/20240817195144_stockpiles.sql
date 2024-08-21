-- +goose Up
-- +goose StatementBegin
CREATE TABLE resource_types (
    id serial primary key,
    density numeric default 1000 not null, -- kg/m^3. default - water density
    name varchar not null unique
);

CREATE TABLE stockpiles (
    id serial primary key,
    resource_type_id int references resource_types(id) not null,
    quantity numeric not null -- in m^3
);

CREATE TABLE stellar_bodies_stockpiles (
    stellar_body_id uuid references stellar_bodies(id) on delete cascade not null,
    stockpile_id int references stockpiles(id) on delete cascade unique not null,
    CONSTRAINT stellar_bodies_stockpiles_stellar_body_id_stockpile_id UNIQUE (stellar_body_id, stockpile_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stellar_bodies_stockpiles, stockpiles, resource_types;
-- +goose StatementEnd
