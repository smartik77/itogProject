package api

import (
	"CensorService/pkg/moderation"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// API структура для обработчиков API
type API struct {
	censor *moderation.CensorService
}

// NewAPI создает новый экземпляр API
func NewAPI(censor *moderation.CensorService) *API {
	return &API{censor: censor}
}

// CheckRequest структура запроса на проверку
type CheckRequest struct {
	Text string `json:"text"`
}

// CheckResponse структура ответа на проверку
type CheckResponse struct {
	OK     bool     `json:"ok"`
	Errors []string `json:"errors,omitempty"`
}

// HealthResponse структура ответа для health check
type HealthResponse struct {
	Status string `json:"status"`
}

// Endpoints регистрирует обработчики API в роутере
func (api *API) Endpoints(router *mux.Router) {
	router.HandleFunc("/check", api.CheckHandler).Methods("POST")
	router.HandleFunc("/health", api.HealthHandler).Methods("GET")
	router.HandleFunc("/forbidden", api.AddForbiddenWordHandler).Methods("POST")
}

// CheckHandler обрабатывает запрос на проверку текста
func (api *API) CheckHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	if api.censor.CheckText(req.Text) {
		json.NewEncoder(w).Encode(CheckResponse{OK: true})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CheckResponse{
			OK:     false,
			Errors: []string{"Forbidden words detected"},
		})
	}
}

// HealthHandler обрабатывает health check
func (api *API) HealthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(HealthResponse{Status: "OK"})
}

// AddForbiddenWordHandler добавляет новое запрещенное слово
func (api *API) AddForbiddenWordHandler(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("word")
	if word == "" {
		http.Error(w, "Word parameter is required", http.StatusBadRequest)
		return
	}

	api.censor.AddForbiddenWord(word)
	w.WriteHeader(http.StatusCreated)
}
