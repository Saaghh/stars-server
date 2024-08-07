-- +migrate Up
CREATE TABLE games(
    id serial primary key,
    owner_id int references users(id) not null,
    world_time timestamp not null,
    name varchar not null,
    is_archive boolean not null default true,
    created_at timestamp with time zone not null default now()
);

CREATE TABLE civilizations(
    id serial primary key,
    game int references games(id) not null,
    owner int references users(id) not null,
    name varchar not null
);

CREATE TABLE systems(
    id serial primary key,
    game_id int references games(id) not null,
    name varchar not null
);

CREATE TABLE stellar_bodies_types(
    id serial primary key,
    name varchar not null
);

CREATE TABLE stellar_bodies(
    id uuid primary key,
    system_id int references systems(id) not null,

    -- body properties
    name varchar not null,
    type_id int references stellar_bodies_types not null,
    mass numeric not null,
    diameter numeric not null,

    -- orbital params
    parent_body_id uuid references stellar_bodies(id),
    orbital_radius numeric,
    angle numeric,
    angle_speed numeric,

    -- linear params
    linear_speed numeric,
    coordinate_x numeric,
    coordinate_y numeric
);

CREATE TABLE colonies(
    id uuid primary key,
    stellar_body uuid references stellar_bodies(id) not null,
    name varchar not null,
    civilization int references civilizations(id) not null
);

CREATE TABLE system_connections(
    id uuid primary key,
    systems_from int references systems(id) not null,
    systems_to int references systems(id) not null
);

-- +migrate Down
DROP TABLE system_connections, colonies, stellar_bodies, stellar_bodies_types, systems, civilizations, games;