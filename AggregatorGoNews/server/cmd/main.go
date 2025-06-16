package main

import (
	"aggregator/pkg/api"
	"aggregator/pkg/posts"
	"aggregator/pkg/rss"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := posts.Connect()
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	b, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		log.Fatal(err)
	}

	chanPost := make(chan []posts.Post)
	chErr := make(chan error)
	for _, url := range conf.URLS {
		go rss.ParseURL(url, chanPost, chErr, conf.Period)
	}

	go func() {
		for newPosts := range chanPost {
			for _, post := range newPosts {
				if err := db.InsertPost(&post); err != nil {
					log.Printf("Ошибка сохранения новости: %v", err)
				}
			}
		}
	}()

	go func() {
		for err := range chErr {
			log.Println("Ошибка RSS:", err)
		}
	}()

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", api.Router()))
}
