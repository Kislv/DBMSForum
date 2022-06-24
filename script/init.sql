DROP TABLE users CASCADE;
DROP TABLE forums CASCADE;
DROP TABLE threads CASCADE;
DROP TABLE posts CASCADE;
DROP TABLE votes CASCADE;
DROP TABLE forum_users CASCADE;

DROP FUNCTION update_threads_count() CASCADE;
DROP FUNCTION update_forum_users_by_insert_th_or_post() CASCADE;
DROP FUNCTION update_path() CASCADE;
DROP FUNCTION insert_votes() CASCADE;
DROP FUNCTION update_votes() CASCADE;

CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE UNLOGGED TABLE "users" (
  "about" text NOT NULL,
  "email" citext collate "C" NOT NULL,
  "fullname" text NOT NULL,
  "nickname" citext collate "C" PRIMARY KEY
);

-- CREATE INDEX index_users_all ON users (nickname);
CREATE INDEX index_users_allH ON users USING hash (nickname);

CREATE UNLOGGED TABLE "forums" (
  "username" citext collate "C" NOT null,
  "posts" BIGINT DEFAULT 0,
  "threads" int DEFAULT 0,
  "slug" citext collate "C" PRIMARY KEY,
  "title" TEXT NOT NULL,
  FOREIGN KEY ("username") REFERENCES "users" (nickname)
);

CREATE INDEX index_forums_slug ON forums USING hash (slug);
-- CREATE INDEX index_users_fk ON forums (username);
-- CREATE INDEX index_forum_all ON forums (slug, title, author, posts, threads);

CREATE UNLOGGED TABLE "threads" (
  "id" SERIAL PRIMARY KEY,
  "author" citext collate "C" NOT NULL,
  "created" timestamptz DEFAULT now(),
  "forum" citext collate "C" NOT NULL,
  "message" TEXT NOT NULL,
  "slug" citext collate "C",
  "title" TEXT NOT NULL,
  "votes" int DEFAULT 0,
  FOREIGN KEY (author) REFERENCES "users" (nickname)
--   FOREIGN KEY (forum) REFERENCES "forums" (slug)
);

-- CREATE INDEX index_threads_slug ON threads (slug);
CREATE INDEX index_thread_slug_hash ON threads USING hash (slug);
-- CREATE INDEX index_thread_users_fk ON threads (author);
-- CREATE INDEX index_thread_forum_fk ON threads (forum);

CREATE INDEX index_thread_forum_created ON threads (forum, created);
-- CREATE INDEX index_thread_all ON threads (title, message, created, slug, author, forum, votes);

CREATE UNLOGGED TABLE "posts" (
  "author" citext collate "C" NOT NULL,
  "created" timestamp DEFAULT now(),
  "forum" citext collate "C" NOT NULL,
  "id" BIGSERIAL PRIMARY KEY,
  "is_edited" BOOL DEFAULT false,
  "message" TEXT NOT NULL,
  "parent" BIGINT DEFAULT 0,
  "thread" int,
  "path" BIGINT[] DEFAULT ARRAY []::INTEGER[],
  
  FOREIGN KEY (author) REFERENCES "users" (nickname)
--   FOREIGN KEY (forum) REFERENCES "forums" (slug)
--   FOREIGN KEY (thread) REFERENCES "threads" (id)
--   FOREIGN KEY (parent) REFERENCES "posts" (id)
);

CREATE INDEX index_posts_authorid ON posts (thread, id, path);
CREATE INDEX index_posts_authorp ON posts (thread, path);
CREATE INDEX index_post_path1_path ON posts ((path[1]), path);
CREATE INDEX index_post_thread_created_id ON posts (thread, created, id);

CREATE UNLOGGED TABLE "votes" (
  "thread" int,
  "nickname" citext collate "C" NOT NULL,
  "voice" int,
  
   FOREIGN KEY (nickname) REFERENCES "users" (nickname),
--    FOREIGN KEY (thread) REFERENCES "threads" (id),
   UNIQUE (nickname, thread)
);

-- CREATE INDEX index_votes_thread_nick ON votes (thread, nickname);

CREATE UNLOGGED TABLE forum_users
(
    author citext collate "C" REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    slug   citext collate "C" NOT NULL,
    UNIQUE (author, slug)
);
-- CREATE INDEX on forum_users (slug);

CREATE OR REPLACE FUNCTION update_threads_count() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    UPDATE forums SET Threads=(Threads+1) WHERE slug=NEW.forum;
    return NEW;
END
$update_users_forum$ LANGUAGE plpgsql;

CREATE TRIGGER add_thread
    BEFORE INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_threads_count();

CREATE OR REPLACE FUNCTION update_forum_users_by_insert_th_or_post()
RETURNS TRIGGER AS
$BODY$
BEGIN
    INSERT INTO forum_users values (NEW.author, NEW.forum)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER thread_insert_forum
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_users_by_insert_th_or_post();

CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS
$update_path$
DECLARE
    parent_path         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(NEW.path, NEW.id);
    ELSE
        SELECT thread, path FROM posts
        WHERE thread = NEW.thread AND id = NEW.parent
        INTO first_parent_thread, parent_path;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'bad parent thread';
        END IF;

        NEW.path := parent_path || NEW.id;
    END IF;
    UPDATE forums SET Posts=Posts + 1 WHERE forums.slug = new.forum;
    RETURN NEW;
END
$update_path$ LANGUAGE plpgsql;

CREATE TRIGGER path_update
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();


CREATE OR REPLACE FUNCTION insert_votes() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    IF NEW.voice > 0 THEN
        UPDATE threads SET votes = (votes + 1) WHERE id = NEW.thread;
    ELSE
        UPDATE threads SET votes = (votes - 1) WHERE id = NEW.thread;
    END IF;
    return NEW;
END
$update_users_forum$ LANGUAGE plpgsql;

CREATE TRIGGER add_vote
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insert_votes();


CREATE OR REPLACE FUNCTION update_votes() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    IF NEW.voice = OLD.voice THEN
        RETURN NEW;
    END IF;
    IF NEW.voice > 0 THEN
        UPDATE threads SET votes = (votes + 2) WHERE id = NEW.thread;
    ELSE
        UPDATE threads SET votes = (votes - 2) WHERE id = NEW.thread;
    END IF;
    return NEW;
END
$update_users_forum$ LANGUAGE plpgsql;

CREATE TRIGGER edit_vote
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();
