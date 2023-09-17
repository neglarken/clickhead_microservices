CREATE TABLE if not EXISTS users (
    id serial PRIMARY KEY,
    login text NOT NULL UNIQUE,
    hashed_password text NOT NULL
);

CREATE TABLE if not EXISTS items (
    id serial PRIMARY KEY,
    info text NOT NULL,
    price integer not null
);

CREATE TABLE if not EXISTS session (
    user_id integer UNIQUE REFERENCES users (id),
    refresh_token text not null UNIQUE,
    exires_at timestamp not null
);