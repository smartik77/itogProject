package posts

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	PubTime int64  `json:"pubTime"`
	Link    string `json:"link"`
}

type DB struct {
	pool *pgxpool.Pool
}

func Connect() (*DB, error) {
	connStr := "comments://testdbuser:pass@127.0.0.1:5432/aggregator"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return &DB{pool: pool}, nil
}

func (db *DB) InsertPost(post *Post) error {
	query := `
		INSERT INTO posts(title, content, pub_time, link) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (link) DO NOTHING`
	_, err := db.pool.Exec(
		context.Background(),
		query,
		post.Title,
		post.Content,
		time.Unix(post.PubTime, 0),
		post.Link,
	)
	return err
}

func (db *DB) SearchPosts(search string, offset, limit int) ([]Post, int, error) {
	countQuery := "SELECT COUNT(*) FROM posts WHERE title ILIKE $1"
	var total int
	err := db.pool.QueryRow(
		context.Background(),
		countQuery,
		"%"+search+"%",
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, title, content, pub_time, link 
		FROM posts 
		WHERE title ILIKE $1
		ORDER BY pub_time DESC
		LIMIT $2 OFFSET $3`

	rows, err := db.pool.Query(
		context.Background(),
		query,
		"%"+search+"%",
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var pubTime time.Time
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &pubTime, &p.Link)
		if err != nil {
			return nil, 0, err
		}
		p.PubTime = pubTime.Unix()
		posts = append(posts, p)
	}

	return posts, total, nil
}
