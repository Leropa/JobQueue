package main

import (
	"context"
	myHttp "jobQueue/internal/delivery/http"
	myRedis "jobQueue/internal/repository/redis"
	usecase "jobQueue/internal/usecase"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	r := chi.NewRouter()

	repo := myRedis.NewTaskRepository(rdb)
	uc := usecase.NewTaskUsecase(repo)
	handler := myHttp.NewTaskHandler(uc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go uc.StartWorker(ctx)

	r.Post("/task", handler.CreateTask)

	log.Println("🌐 HTTP Server started on http://localhost:8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
