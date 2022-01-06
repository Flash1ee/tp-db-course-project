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
    slug           citext NOT NULL PRIMARY KEY,
    posts          bigint DEFAULT 0,
    threads        int    DEFAULT 0
);
CREATE TABLE IF NOT EXISTS thread
(
    id      bigserial PRIMARY KEY NOT NULL,
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
    created   timestamp with time zone DEFAULT now(),
    path      bigint[]                 DEFAULT ARRAY []::INTEGER[]
);

CREATE TABLE IF NOT EXISTS vote
(
    nickname  citext NOT NULL REFERENCES users (nickname),
    thread_id int    NOT NULL REFERENCES thread (id),
    voice     int    NOT NULL
);

CREATE TABLE IF NOT EXISTS user_forum
(
    nickname citext NOT NULL UNIQUE REFERENCES users (nickname),
    forum    citext NOT NULL UNIQUE REFERENCES forum (slug)
);

CREATE OR REPLACE FUNCTION insert_votes_into_threads()
    RETURNS TRIGGER AS
$$
BEGIN
    UPDATE thread
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread_id;
    RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER insert_votes
    AFTER INSERT
    ON vote
    FOR EACH ROW
EXECUTE PROCEDURE insert_votes_into_threads();

CREATE OR REPLACE FUNCTION update_votes_in_threads()
    RETURNS TRIGGER AS
$$
BEGIN
    UPDATE thread
    SET votes = votes + NEW.voice - OLD.voice
    WHERE id = NEW.thread_id;
    RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER update_votes
    AFTER UPDATE
    ON vote
    FOR EACH ROW
EXECUTE PROCEDURE update_votes_in_threads();

CREATE OR REPLACE FUNCTION path_update() RETURNS TRIGGER AS
$$
BEGIN
    new.path = (SELECT path FROM post WHERE id = new.parent) || new.id;
    RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER path_upd
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE path_update();

SELECT id, parent, author, message, is_edited, forum, thread, created FROM post
WHERE path[1] IN (SELECT id FROM post WHERE thread = 243 ORDER BY id ASC LIMIT 65)
ORDER BY path ASC, id ASC;

SELECT id, parent, author, message, is_edited, forum, thread, created FROM post WHERE thread = 1022  and path >  (SELECT path FROM post WHERE id = 12834)  ORDER BY path asc, id LIMIT NULLIF(3, 0)