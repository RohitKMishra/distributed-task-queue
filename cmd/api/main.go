package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/RohitKMishra/distributed-task-queue/internal/config"
	"github.com/RohitKMishra/distributed-task-queue/internal/logger"
	"github.com/RohitKMishra/distributed-task-queue/internal/storage"

	"github.com/RohitKMishra/distributed-task-queue/internal/api"
	"github.com/RohitKMishra/distributed-task-queue/internal/queue"
	"github.com/RohitKMishra/distributed-task-queue/internal/task"
)

func main(){
	ctx := context.Background()

	cfg := config.Load()

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	postgresPool, err := storage.NewPostgres(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal("Failed to connect to Postgres", zap.Error(err))
	}

	defer postgresPool.Close()

	redisClient := storage.NewRedis(cfg.RedisAddr)

	if err := storage.PingRedis(ctx, redisClient); err != nil {
		log.Fatal("Failed to connect to redis", zap.Error(err))
	}

	taskRepository := task.NewRepository(postgresPool)
	redisQueue := queue.NewRedisQueue(redisClient)

	taskHandler := api.NewTaskHandler(taskRepository, redisQueue, log)

	router := chi.NewRouter()

	router.Post("/tasks", taskHandler.CreateTask)
	router.Get("/health", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Info("Api server started", zap.String("port", cfg.HTTPPort))

	if err:= http.ListenAndServe(cfg.HTTPPort, router); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}