DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
   id SERIAL PRIMARY KEY,
   title TEXT NOT NULL,
   content TEXT NOT NULL,
   pub_time BIGINT DEFAULT 0,
   link TEXT NOT NULL UNIQUE
);