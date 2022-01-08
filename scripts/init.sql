CREATE EXTENSION IF NOT EXISTS citext;
--
-- alter table users alter column
--     nickname set data type citext COLLATE "ucs_basic";
CREATE UNLOGGED TABLE IF NOT EXISTS users
(
    id       bigserial,
    nickname citext COLLATE "ucs_basic" NOT NULL UNIQUE PRIMARY KEY,
    fullname text   NOT NULL,
    about    text,
    email    citext NOT NULL UNIQUE
);
CREATE UNLOGGED TABLE IF NOT EXISTS forum
(
    id             bigserial,
    title          text   NOT NULL,
    users_nickname citext NOT NULL REFERENCES users (nickname),
    slug           citext NOT NULL PRIMARY KEY,
    posts          bigint DEFAULT 0,
    threads        int    DEFAULT 0
);
CREATE UNLOGGED TABLE IF NOT EXISTS thread
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
CREATE UNLOGGED TABLE IF NOT EXISTS post
(
    id        bigserial PRIMARY KEY NOT NULL UNIQUE,
    parent    int                      DEFAULT 0,
    author    citext                NOT NULL REFERENCES users (nickname),
    message   text                  NOT NULL,
    is_edited bool                     DEFAULT FALSE,
    forum     citext REFERENCES forum (slug),
    thread    integer REFERENCES thread (id),
    created   timestamp with time zone  DEFAULT now(),
    path      bigint[]                 DEFAULT ARRAY []::INTEGER[]
);

CREATE UNLOGGED TABLE IF NOT EXISTS vote
(
    nickname  citext NOT NULL REFERENCES users (nickname),
    thread_id int    NOT NULL REFERENCES thread (id),
    voice     int    NOT NULL
);
CREATE UNLOGGED TABLE IF NOT EXISTS user_forum
(
    nickname citext NOT NULL REFERENCES users (nickname),
    forum    citext NOT NULL REFERENCES forum (slug)
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

CREATE OR REPLACE FUNCTION cnt_posts()
    RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forum
    SET posts = forum.posts + 1
    WHERE slug = NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_count_posts
    AFTER INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE cnt_posts();

CREATE OR REPLACE FUNCTION cnt_threads()
    RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forum
    SET threads = forum.threads + 1
    WHERE slug = NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_count_threads
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE cnt_threads();



CREATE OR REPLACE FUNCTION upd_user_forum()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO user_forum (nickname, forum)
    VALUES (NEW.author, NEW.forum)
    ON CONFLICT do nothing;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_forum
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE upd_user_forum();

CREATE TRIGGER update_users_forum
    AFTER INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE upd_user_forum();

-- -
create index if not exists post_path_parent on post ((path[1]));
create index if not exists post_thread_thread_id on post (thread, id);
create index if not exists post_path_id on post (id, (path[1]));
create index if not exists post_parent on post (thread, id, (path[1]), parent);
create index if not exists post_sorting on post ((path[1]) desc, path, id);
create index if not exists post_thread on post (thread);
create index if not exists post_thread_path_id on post (thread, path, id);
-----------
create index if not exists user_nickname_hash on users using hash (nickname);
create index if not exists user_all on users  (nickname, fullname, about, email);
-----------
create index if not exists forum_slug on forum using hash (slug);
-----------
create unique index if not exists votes on vote (nickname, thread_id, voice);
-----------
create index if not exists th_slug_hash on thread using hash (slug);
create index if not exists th_user on thread using hash (author);
create index if not exists th_created on thread (created);
create index if not exists th_forum_created on thread (forum, created);
-----------
create index if not exists forum_user on forum using hash (users_nickname);
create index if not exists f_u_nickname on user_forum using hash (nickname);
create index if not exists users_to_forums_forum on user_forum (forum);
VACUUM;
VACUUM ANALYSE;