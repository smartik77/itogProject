DROP TABLE IF EXISTS comments;

CREATE TABLE comments (
      id SERIAL PRIMARY KEY,
      news_id INTEGER NOT NULL,
      parent_id INTEGER,
      content TEXT NOT NULL,
      author VARCHAR(100) NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);