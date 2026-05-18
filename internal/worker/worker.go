package worker

import (
	"context"

	"go.uber.org/zap"

	"github.com/RohitKMishra/distributed-task-queue/internal/queue"
	"github.com/RohitKMishra/distributed-task-queue/internal/task"
)

type Worker struct {
	queue *queue.RedisQueue
	repository *task.Repository
	logger *zap.Logger
}

func NewWorker (queue *queue.RedisQueue, repository *task.Repository, logger *zap.Logger) *Worker {
	return &Worker{
		queue: queue,
		repository: repository,
		logger: logger,
	}
}

func (w *Worker) Start(ctx context.Context) error {

	w.logger.Info("Worker started")

	for {
		taskID, err := w.queue.Dequeue(ctx)
		if err != nil {
			w.logger.Error("Failed to dequeue task", zap.Error(err))
			continue
		}

		w.logger.Info("task dequed", zap.String("task_id", taskID))

		taskData, err := w.repository.GetById(ctx, taskID)
		if err != nil {
			w.logger.Error("Failed to fetch task", zap.Error(err))
			continue
		}

		err = w.repository.UpdateStatus(ctx, taskID, task.StatusProcessing)

		if err != nil {
			w.logger.Error("Failed to update task status", zap.Error(err))
			continue
		}

		w.logger.Info("Processing task", zap.String("task_id", taskID), zap.String("task_type", taskData.Type))

		err = w.repository.UpdateStatus(ctx, taskID, task.StatusCompleted)

		if err != nil {
			w.logger.Error ("Failed to complete task", zap.Error(err))
			continue
		}

		w.logger.Info("Task Completed", zap.String("task_id", taskID))
	}
}