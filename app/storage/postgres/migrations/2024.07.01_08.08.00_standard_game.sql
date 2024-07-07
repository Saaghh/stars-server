-- +migrate Up
INSERT INTO users (id, username, password_hash) VALUES (1, 'system', 'system');
INSERT INTO games (id, owner, world_time, name) VALUES (1, 1, now(), 'standard_game');

INSERT INTO systems (game_id, name) VALUES (1, 'Solar System');

INSERT INTO stellar_bodies_types (id, name) VALUES (1, 'star');
INSERT INTO stellar_bodies_types (id, name) VALUES (2, 'planet');
INSERT INTO stellar_bodies_types (id, name) VALUES (3, 'asteroid');
INSERT INTO stellar_bodies_types (id, name) VALUES (4, 'craft');

-- +migrate Down

TRUNCATE TABLE stellar_bodies_types, systems, games, users;