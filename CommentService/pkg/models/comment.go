package models

import (
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	NewsID    int       `json:"news_id"`
	ParentID  *int      `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}
