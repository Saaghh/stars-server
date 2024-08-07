-- +migrate Up
CREATE TABLE game_events_types (
    id serial primary key,
    name varchar not null,
    is_stopping boolean not null
);

CREATE TABLE game_events (
    id serial primary key,
    game_id int references games(id) not null,
    timestamp timestamp not null,
    type_id int references game_events_types(id) not null,
    civilization_id int references civilizations(id)
);

CREATE TABLE actor (
    id uuid primary key,
    name varchar not null,
    civilization_id int references civilizations(id)
);

-- +migrate Down
