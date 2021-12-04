CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id       bigserial,
    nickname citext NOT NULL UNIQUE PRIMARY KEY,
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
    posts          bigint DEFAULT 0,
    threads        int    DEFAULT 0
);

CREATE TABLE IF NOT EXISTS thread
(
    id      bigserial PRIMARY KEY NOT NULL UNIQUE,
    title   text                  NOT NULL,
    author  citext                NOT NULL REFERENCES users (nickname),
    forum   citext                NOT NULL REFERENCES forum (slug),
    message text                  NOT NULL,
    votes   integer                  DEFAULT 0,
    slug    citext,
    created timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS post
(
    id        bigserial PRIMARY KEY NOT NULL UNIQUE,
    parent    int                      DEFAULT 0,
    author    citext                NOT NULL REFERENCES users (nickname),
    message   text                  NOT NULL,
    is_edited bool                     DEFAULT FALSE,
    forum     citext REFERENCES forum (slug),
    thread    integer REFERENCES thread (id),
    created   timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vote
(
    nickname  citext NOT NULL UNIQUE REFERENCES users (nickname),
    thread_id int    NOT NULL UNIQUE REFERENCES thread (id),
    voice     int    NOT NULL
);

CREATE TABLE IF NOT EXISTS user_forum
(
    nickname citext NOT NULL UNIQUE REFERENCES  users (nickname),
    forum    citext NOT NULL UNIQUE REFERENCES forum (slug)
);