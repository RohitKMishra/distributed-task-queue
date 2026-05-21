package main

import (

"context"

"go.uber.org/zap"

"github.com/RohitKMishra/distributed-task-queue/internal/config"
"github.com/RohitKMishra/distributed-task-queue/internal/logger"
"github.com/RohitKMishra/distributed-task-queue/internal/storage"

"github.com/RohitKMishra/distributed-task-queue/internal/queue"
"github.com/RohitKMishra/distributed-task-queue/internal/task"
"github.com/RohitKMishra/distributed-task-queue/internal/worker"
)

func main () {
	ctx := context.Background()

	cfg := config.Load()

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	postgresPool, err := storage.NewPostgres(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatal("Failed to connect to postgres", zap.Error(err))
	}

	defer postgresPool.Close()

	redisClient := storage.NewRedis(cfg.RedisAddr)

	if err := storage.PingRedis(ctx, redisClient); err != nil {
		log.Fatal("failed to connect to redis", zap.Error(err))
	}

	taskRepository := task.NewRepository(postgresPool)
	redisQueue := queue.NewRedisQueue(redisClient)

	taskWorker := worker.NewWorker(redisQueue, taskRepository, log, 3)

	if err := taskWorker.Start(ctx); err != nil {
		log.Fatal("worker failed", zap.Error(err))
	}
}