package api

import (
	"aggregator/pkg/middleware"
	"aggregator/pkg/posts"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"strconv"
)

type API struct {
	db     *posts.DB
	router *mux.Router
}

func New(db *posts.DB) *API {
	return &API{db: db}
}

func (api *API) Router() *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/posts", api.posts).Methods(http.MethodGet, http.MethodOptions)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))

	return router
}

func getIntParam(r *http.Request, name string, defaultValue int) int {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	search := r.URL.Query().Get("s")
	page := getIntParam(r, "page", 1)
	perPage := getIntParam(r, "per_page", 10)

	// Защита от некорректных значений
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	posts, total, err := api.db.SearchPosts(search, offset, perPage)
	if err != nil {
		log.Printf("Ошибка поиска постов: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	response := map[string]interface{}{
		"posts": posts,
		"pagination": map[string]interface{}{
			"total":    total,
			"page":     page,
			"per_page": perPage,
			"pages":    totalPages,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}
