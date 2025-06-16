package main

import (
	"CommentService/pkg/api"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	connStr := "comments://comments:priliv1337228@127.0.0.1:5432/commentdb"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer conn.Close(context.Background())

	// Проверка соединения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("Ошибка проверки соединения: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL")

	// Настройка HTTP-сервера
	server := &http.Server{
		Addr:    ":8081",
		Handler: api.Endpoints(conn),
	}

	// Graceful shutdown
	go func() {
		log.Printf("Сервис комментариев запущен на %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Получен сигнал завершения")

	// Завершение работы
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка завершения: %v", err)
	} else {
		log.Println("Сервер остановлен корректно")
	}
}
