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

create unlogged table if not exists forums_users (
    "user"      CITEXT COLLATE "C"  NOT NULL    REFERENCES users(nickname),
    forum       CITEXT              NOT NULL    REFERENCES forums(slug),

    CONSTRAINT
    unique_forum_user UNIQUE("user", forum)
);

CREATE UNLOGGED TABLE IF NOT EXISTS threads (
    id          BIGSERIAL                   NOT NULL        PRIMARY KEY,
    title       TEXT                        NOT NULL,
    author      CITEXT COLLATE "C"          NOT NULL        REFERENCES users(nickname),
    forum       CITEXT                      NOT NULL        REFERENCES forums(slug),
    message     TEXT                        NOT NULL,
    votes       BIGINT                      DEFAULT 0,
    slug        CITEXT                      NOT NULL,
    created     TIMESTAMP WITH TIME ZONE    DEFAULT now()
);

CREATE UNLOGGED TABLE IF NOT EXISTS posts (
    id          BIGSERIAL                   NOT NULL                    PRIMARY KEY,
    parent      BIGSERIAL                                               REFERENCES posts(id),
    author      CITEXT COLLATE "C"          NOT NULL                    REFERENCES users(nickname),
    message     TEXT                        NOT NULL,
    is_edited   BOOLEAN                     DEFAULT FALSE,
    forum       CITEXT                      NOT NULL                    REFERENCES forums(slug),
    thread      BIGINT                      NOT NULL                    REFERENCES threads(id),
    created     TIMESTAMP WITH TIME ZONE    DEFAULT now(),
    path        BIGINT[]                    DEFAULT ARRAY[]::BIGINT[]
);

CREATE UNLOGGED TABLE IF NOT EXISTS votes (
    id          BIGSERIAL           NOT NULL    PRIMARY KEY,
    nickname    CITEXT COLLATE "C"  NOT NULL    REFERENCES users(nickname),
    thread      BIGINT              NOT NULL    REFERENCES threads(id),
    voice       INT                 NOT NULL,

    CONSTRAINT
    unique_vote UNIQUE(nickname, thread)
);
