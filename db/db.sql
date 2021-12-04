CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id       bigserial,
    nickname citext NOT NULL unique primary key,
    fullname text   NOT NULL,
    about    text,
    email    citext NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS forum
(
    id             bigserial,
    title          text   NOT NULL,
    users_nickname citext NOT NULL REFERENCES users (nickname),
    slug           citext NOT NULL PRIMARY KEY UNIQUE,
    posts          bigint default 0,
    threads        int    default 0
);
