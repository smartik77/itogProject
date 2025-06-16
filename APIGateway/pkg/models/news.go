package models

import "time"

// NewsShort Модели для новостей
type NewsShort struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Summary string    `json:"summary"`
	Date    time.Time `json:"date"`
}

type NewsDetail struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Date     time.Time `json:"date"`
	Link     string    `json:"link"`
	Comments []Comment `json:"comments"`
}

// Comment Модель для комментариев
type Comment struct {
	ID        int       `json:"id"`
	NewsID    int       `json:"news_id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse Модель для ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
