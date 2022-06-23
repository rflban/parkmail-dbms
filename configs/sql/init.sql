CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE UNLOGGED TABLE IF NOT EXISTS users (
    id          BIGSERIAL           NOT NULL    UNIQUE,
    nickname    CITEXT COLLATE "C"  NOT NULL    PRIMARY KEY,
    fullname    TEXT                NOT NULL,
    about       TEXT,
    email       CITEXT              NOT NULL    UNIQUE
);

CREATE UNLOGGED TABLE IF NOT EXISTS forums (
    id          BIGSERIAL           NOT NULL    UNIQUE,
    title       TEXT                NOT NULL,
    "user"      CITEXT COLLATE "C"  NOT NULL    REFERENCES users(nickname),
    slug        CITEXT              NOT NULL    PRIMARY KEY,
    posts       BIGINT              DEFAULT 0,
    threads     BIGINT              DEFAULT 0
);

CREATE UNLOGGED TABLE IF NOT EXISTS forums_users (
    nickname    CITEXT COLLATE "C"  NOT NULL    REFERENCES users(nickname),
    fullname    TEXT                NOT NULL,
    about       TEXT,
    email       CITEXT              NOT NULL    UNIQUE,
    forum       CITEXT              NOT NULL    REFERENCES forums(slug),

    CONSTRAINT unique_forum_user UNIQUE(nickname, forum)
);

CREATE UNLOGGED TABLE IF NOT EXISTS threads (
    id          BIGSERIAL                   NOT NULL        PRIMARY KEY,
    title       TEXT                        NOT NULL,
    author      CITEXT COLLATE "C"          NOT NULL        REFERENCES users(nickname),
    forum       CITEXT                      NOT NULL        REFERENCES forums(slug),
    message     TEXT                        NOT NULL,
    votes       BIGINT                      DEFAULT 0,
    slug        CITEXT,
    created     TIMESTAMP WITH TIME ZONE    DEFAULT now()
);

CREATE UNLOGGED TABLE IF NOT EXISTS posts (
    id          BIGSERIAL                   NOT NULL                    PRIMARY KEY,
    parent      BIGINT                      DEFAULT 0,
    author      CITEXT COLLATE "C"          NOT NULL                    REFERENCES users(nickname),
    message     TEXT                        NOT NULL,
    is_edited   BOOLEAN                     DEFAULT FALSE,
    forum       CITEXT                      NOT NULL                    REFERENCES forums(slug),
    thread      BIGINT                      NOT NULL                    REFERENCES threads(id),
    created     TIMESTAMP WITH TIME ZONE    DEFAULT NOW(),
    path        BIGINT[]                    DEFAULT ARRAY[]::BIGINT[],
    batch_id    VARCHAR(36),
    batch_idx   INTEGER
);

CREATE UNLOGGED TABLE IF NOT EXISTS votes (
    id          BIGSERIAL           NOT NULL    PRIMARY KEY,
    nickname    CITEXT COLLATE "C"  NOT NULL    REFERENCES users(nickname),
    thread      BIGINT              NOT NULL    REFERENCES threads(id),
    voice       INT                 NOT NULL,

    CONSTRAINT unique_vote UNIQUE(nickname, thread)
);

CREATE OR REPLACE FUNCTION threads__set_votes() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE threads
           SET votes = votes + NEW.voice
         WHERE id = NEW.thread;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER votes__on_insert__threads__set_votes
    AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE threads__set_votes();

CREATE OR REPLACE FUNCTION threads__update_votes() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE threads
           SET votes = votes + NEW.voice - OLD.voice
         WHERE id = NEW.thread;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER votes__on_update__threads__update_votes
    AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE threads__update_votes();

CREATE OR REPLACE FUNCTION posts__set_path() RETURNS TRIGGER AS $$
    DECLARE
        p_path      BIGINT[];
        p_thread    BIGINT;
        p_id        BIGINT;
    BEGIN
        SELECT path, thread, id
          FROM posts
         WHERE id = NEW.parent
          INTO p_path, p_thread, p_id;

        if (p_id IS NULL AND NEW.parent != 0) THEN
            RAISE EXCEPTION SQLSTATE '23514' USING MESSAGE = 'PARENT DOES NOT EXIST';
        END IF;

        IF (p_thread != NEW.thread) THEN
            RAISE EXCEPTION SQLSTATE '23514' USING MESSAGE = 'PARENT AND NEW POSTS THREAD MISMATCH';
        END IF;

        NEW.path = p_path || NEW.id;

        RETURN NEW;
    END;
$$ LANGUAGE  plpgsql;
CREATE TRIGGER posts__on_insert__threads__update_votes
    BEFORE INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE posts__set_path();

CREATE OR REPLACE FUNCTION forums__count_threads() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE forums
           SET threads = forums.threads + 1
         WHERE slug = NEW.forum;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER threads__on_insert__forums__count_threads
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE forums__count_threads();

CREATE OR REPLACE FUNCTION forums__count_posts() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE forums
           SET posts = forums.posts + 1
         WHERE slug = NEW.forum;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER posts__on_insert__forums__count_posts
    AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE forums__count_posts();

CREATE OR REPLACE FUNCTION forums_users__update() RETURNS TRIGGER AS $$
    DECLARE
        v_nickname    CITEXT;
        v_fullname    TEXT;
        v_about       TEXT;
        v_email       CITEXT;
    BEGIN
        SELECT u.nickname, u.fullname, u.about, u.email
          FROM users u
         WHERE u.nickname = NEW.author
          INTO v_nickname, v_fullname, v_about, v_email;

        INSERT INTO forums_users (nickname, fullname, about, email, forum)
             VALUES (v_nickname, v_fullname, v_about, v_email, NEW.forum)
        ON CONFLICT DO NOTHING;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER posts__on_insert__forums_users__update
    AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE forums_users__update();
CREATE TRIGGER threads__on_insert__forums_users__update
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE forums_users__update();
