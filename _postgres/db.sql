CREATE TABLE IF NOT EXISTS "users"
(
    id       SERIAL PRIMARY KEY NOT NULL,
    username VARCHAR(255)       NOT NULL,
    password VARCHAR(255)       NOT NULL,
    role     VARCHAR(15)        NOT NULL
);


CREATE TABLE IF NOT EXISTS "films"
(
    id              SERIAL PRIMARY KEY NOT NULL,
    name            VARCHAR(150)       NOT NULL,
    description     TEXT               NOT NULL,
    date_of_release TIMESTAMP          NOT NULL,
    rating          NUMERIC(3, 1)      NOT NULL
);



CREATE TABLE IF NOT EXISTS "actors"
(
    id       SERIAL PRIMARY KEY NOT NULL,
    name     VARCHAR(100)       NOT NULL,
    surname  VARCHAR(100)       NOT NULL,
    gender   VARCHAR(6)         NOT NULL,
    birthday TIMESTAMP          NOT NULL
);

CREATE TABLE IF NOT EXISTS film_actors
(
    film_id  INT REFERENCES films (id) ON DELETE CASCADE,
    actor_id INT REFERENCES actors (id) ON DELETE CASCADE,
    PRIMARY KEY (film_id, actor_id)
);

CREATE INDEX IF NOT EXISTS idx_actor_id ON actors (id);

CREATE INDEX IF NOT EXISTS idx_film_id ON films (id);
