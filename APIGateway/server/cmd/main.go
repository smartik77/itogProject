package main

import (
	"APIGateway/pkg/api"
	"APIGateway/pkg/middleware"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	newsURL := os.Getenv("NEWS_SERVICE_URL")
	if newsURL == "" {
		newsURL = "http://localhost:8080"
	}

	commentsURL := os.Getenv("COMMENTS_SERVICE_URL")
	if commentsURL == "" {
		commentsURL = "http://localhost:8081"
	}

	router := mux.NewRouter()

	// Применение middleware
	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.LoggingMiddleware) // Исправлено на правильный пакет

	// Инициализация и регистрация обработчиков
	handler := api.NewHandler(newsURL, commentsURL)
	handler.Endpoints(router)

	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Println("API Gateway запущен на порту 8000")
	log.Fatal(server.ListenAndServe())
}
