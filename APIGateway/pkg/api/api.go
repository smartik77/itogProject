package api

import (
	"APIGateway/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Handler объединяет клиентов и обработчики
type Handler struct {
	NewsClient     *NewsClient
	CommentsClient *CommentsClient
}

func NewHandler(newsURL, commentsURL string) *Handler {
	return &Handler{
		NewsClient:     &NewsClient{BaseURL: newsURL},
		CommentsClient: &CommentsClient{BaseURL: commentsURL},
	}
}

// Endpoints регистрирует обработчики API
func (h *Handler) Endpoints(router *mux.Router) {
	router.HandleFunc("/news", h.GetNewsList).Methods("GET")
	router.HandleFunc("/news/{id:[0-9]+}", h.GetNewsDetail).Methods("GET")
	router.HandleFunc("/comments", h.AddComment).Methods("POST")
}

// NewsClient клиент для агрегатора новостей
type NewsClient struct {
	BaseURL string
}

func (c *NewsClient) GetNews(ctx context.Context) ([]models.NewsShort, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/posts", nil)
	if err != nil {
		return nil, err
	}
	return sendRequest[[]models.NewsShort](req)
}

func (c *NewsClient) GetNewsDetail(ctx context.Context, id int) (*models.NewsDetail, error) {
	url := fmt.Sprintf("%s/posts/%d", c.BaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return sendRequest[*models.NewsDetail](req)
}

// CommentsClient клиент для сервиса комментариев
type CommentsClient struct {
	BaseURL string
}

func (c *CommentsClient) GetComments(ctx context.Context, newsID int) ([]models.Comment, error) {
	url := fmt.Sprintf("%s/comments?news_id=%d", c.BaseURL, newsID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return sendRequest[[]models.Comment](req)
}

func (c *CommentsClient) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	body, err := json.Marshal(comment)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/comments/add", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return sendRequest[*models.Comment](req)
}

// GetNewsList Обработчики HTTP
func (h *Handler) GetNewsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	news, err := h.NewsClient.GetNews(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func (h *Handler) GetNewsDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid news ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	newsDetail, err := h.NewsClient.GetNewsDetail(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comments, err := h.CommentsClient.GetComments(ctx, id)
	if err != nil {
		log.Printf("Failed to get comments: %v", err)
		comments = []models.Comment{}
	}
	newsDetail.Comments = comments

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsDetail)
}

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if comment.NewsID <= 0 || comment.Content == "" {
		http.Error(w, "Invalid comment data", http.StatusBadRequest)
		return
	}

	comment.CreatedAt = time.Now()
	ctx := r.Context()
	newComment, err := h.CommentsClient.AddComment(ctx, &comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newComment)
}

// Общая функция для отправки запросов
func sendRequest[T any](req *http.Request) (T, error) {
	var zero T
	client := &http.Client{Timeout: 5 * time.Second}

	// Добавляем request_id из контекста
	if id, ok := req.Context().Value("request_id").(string); ok {
		q := req.URL.Query()
		q.Add("request_id", id)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return zero, fmt.Errorf("HTTP %d: failed to parse error response", resp.StatusCode)
		}
		return zero, fmt.Errorf("HTTP %d: %s", resp.StatusCode, errResp.Error)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return zero, err
	}
	return result, nil
}
