package api

import (
	"CommentService/pkg/comments"
	"CommentService/pkg/models"
	"CommentService/pkg/moderation"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

type CommentHandler struct {
	Repo  *comments.CommentRepository
	Queue chan int // Канал для очереди модерации
}

func NewCommentHandler(conn *pgx.Conn) *CommentHandler {
	return &CommentHandler{
		Repo:  &comments.CommentRepository{Conn: conn},
		Queue: make(chan int, 100),
	}
}

func (h *CommentHandler) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка на запрещенные слова перед сохранением
	if moderation.CheckContent(comment.Content) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Комментарий содержит запрещенные слова",
		})
		return
	}

	id, err := h.Repo.AddComment(r.Context(), comment)
	if err != nil {
		http.Error(w, "Ошибка добавления: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка в очередь на асинхронную модерацию
	go func() {
		h.Queue <- id
	}()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"status":  "approved", // Теперь сразу approved после проверки
		"message": "Комментарий успешно добавлен",
	})
}

func (h *CommentHandler) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	newsIDStr := r.URL.Query().Get("news_id")
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		http.Error(w, "Некорректный news_id", http.StatusBadRequest)
		return
	}

	comments, err := h.Repo.GetCommentsByNewsID(r.Context(), newsID)
	if err != nil {
		http.Error(w, "Ошибка получения: "+err.Error(), http.StatusInternalServerError)
		return
	}

	approvedComments := make([]models.Comment, 0)
	for _, c := range comments {
		if c.Moderation == models.ModerationApproved {
			approvedComments = append(approvedComments, c)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(approvedComments)
}

// ModerateCommentHandler API для ручной модерации
func (h *CommentHandler) ModerateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка допустимых статусов
	if req.Status != models.ModerationApproved && req.Status != models.ModerationRejected {
		http.Error(w, "Недопустимый статус модерации", http.StatusBadRequest)
		return
	}

	if err := h.Repo.UpdateModerationStatus(r.Context(), req.ID, req.Status); err != nil {
		http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Статус комментария обновлен",
	})
}

// Endpoints настраивает маршруты
func Endpoints(conn *pgx.Conn) http.Handler {
	handler := NewCommentHandler(conn)

	// Запускаем фоновую горутину для модерации
	go startModerationWorker(handler.Repo, handler.Queue)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /comments/add", handler.AddCommentHandler)
	mux.HandleFunc("GET /comments", handler.GetCommentsHandler)
	mux.HandleFunc("POST /comments/moderate", handler.ModerateCommentHandler)

	return mux
}

// Фоновая модерация
func startModerationWorker(repo *comments.CommentRepository, queue chan int) {
	log.Println("Запущен фоновый модератор комментариев")

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case commentID := <-queue:
			content, err := repo.GetCommentContent(context.Background(), commentID)
			if err != nil {
				log.Printf("Ошибка модерации #%d: %v", commentID, err)
				continue
			}

			if moderation.CheckContent(content) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				if err := repo.UpdateModerationStatus(ctx, commentID, models.ModerationRejected); err != nil {
					log.Printf("Ошибка обновления статуса #%d: %v", commentID, err)
				} else {
					log.Printf("Комментарий #%d отклонен: содержит запрещенные слова", commentID)
				}
			}
		case <-ticker.C:
			// Периодическая проверка очереди
		}
	}
}
