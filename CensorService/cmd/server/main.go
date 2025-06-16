package main

import (
	"CensorService/pkg/api"
	"CensorService/pkg/middleware"
	"CensorService/pkg/moderation"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Инициализация сервиса цензурирования
	censor := moderation.NewCensorService([]string{
		"qwerty",
		"йцукен",
		"zxvbnm",
	})

	// Создание API
	api := api.NewAPI(censor)

	// Настройка роутера
	router := mux.NewRouter()

	// Применение middleware
	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.LoggingMiddleware)

	// Регистрация обработчиков
	api.Endpoints(router)

	// Получение порта из переменных окружения

	server := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	log.Printf("Censor service started on port 8082")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
