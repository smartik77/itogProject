package comments

import (
	"CommentService/pkg/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type CommentRepository struct {
	Conn *pgx.Conn
}

func (r *CommentRepository) AddComment(ctx context.Context, comment models.Comment) (int, error) {
	// Валидация
	if comment.NewsID == 0 || comment.Content == "" || comment.Author == "" {
		return 0, errors.New("не заполнены обязательные поля")
	}

	// Проверка родительского комментария
	if comment.ParentID != nil {
		var exists bool
		err := r.Conn.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1)",
			*comment.ParentID,
		).Scan(&exists)

		if err != nil {
			return 0, fmt.Errorf("ошибка проверки родителя: %w", err)
		}
		if !exists {
			return 0, errors.New("родительский комментарий не существует")
		}
	}

	// Вставка с возвратом ID
	var id int
	err := r.Conn.QueryRow(ctx,
		`INSERT INTO comments (
			news_id, 
			parent_id, 
			content, 
			author, 
			created_at,
			moderation
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		comment.NewsID,
		comment.ParentID,
		comment.Content,
		comment.Author,
		time.Now().UTC(),
	).Scan(&id)

	return id, err
}

func (r *CommentRepository) GetCommentsByNewsID(ctx context.Context, newsID int) ([]models.Comment, error) {
	query := `
		SELECT id, news_id, parent_id, content, author, created_at, moderation
		FROM comments 
		WHERE news_id = $1 AND moderation = $2
		ORDER BY created_at ASC
	`

	// Фильтрация только approved комментариев на уровне БД
	rows, err := r.Conn.Query(ctx, query, newsID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID,
			&c.NewsID,
			&c.ParentID,
			&c.Content,
			&c.Author,
			&c.CreatedAt,
		)
		if err != nil {
			log.Printf("Ошибка сканирования: %v", err)
			continue
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentRepository) UpdateModerationStatus(ctx context.Context, id int, status string) error {
	_, err := r.Conn.Exec(ctx,
		"UPDATE comments SET moderation = $1 WHERE id = $2",
		status, id,
	)
	return err
}

func (r *CommentRepository) GetCommentContent(ctx context.Context, id int) (string, error) {
	var content string
	err := r.Conn.QueryRow(ctx,
		"SELECT content FROM comments WHERE id = $1",
		id,
	).Scan(&content)

	return content, err
}
