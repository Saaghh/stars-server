-- +migrate Up
CREATE TABLE games(
    id serial primary key,
    owner_id int references users(id) not null,
    world_time timestamp not null,
    name varchar not null,
    is_archive boolean not null default true,
    created_at timestamp with time zone not null default now()
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
    mass numeric not null, -- in earths
    diameter numeric not null, -- in earths

    -- orbital params
    parent_body_id uuid references stellar_bodies(id),
    orbital_radius numeric, -- in au
    angle numeric, -- in degree
    angle_speed numeric, -- in degree/day

    -- linear params
    linear_speed numeric,
    coordinate_x numeric, -- in au
    coordinate_y numeric -- in au
);

-- +migrate Down
DROP TABLE stellar_bodies, stellar_bodies_types, systems, games;