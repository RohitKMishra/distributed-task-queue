package worker

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/RohitKMishra/distributed-task-queue/internal/queue"
	"github.com/RohitKMishra/distributed-task-queue/internal/task"
)

type Worker struct {
	queue *queue.RedisQueue
	repository *task.Repository
	logger *zap.Logger

	workerCount int
}

func NewWorker (queue *queue.RedisQueue, repository *task.Repository, logger *zap.Logger, workerCount int) *Worker {
	return &Worker{
		queue: queue,
		repository: repository,
		logger: logger,

		workerCount: workerCount,
	}
}

func (w *Worker) Start(ctx context.Context) error {

	w.logger.Info("Worker started", zap.Int("worker_count", w.workerCount))
	taskChannel := make(chan string)

	for i := 0; i < w.workerCount; i++ {
		go w.workerLoop(ctx, i, taskChannel)
	}

	for {
		taskID, err := w.queue.Dequeue(ctx)
		if err != nil {
			w.logger.Error("Failed to dequeue task", zap.Error(err))
			continue
		}

		taskChannel <- taskID
	}
}

func (w *Worker) workerLoop(ctx context.Context, workerID int, taskChannel <-chan string) {
	w.logger.Info("Worker goroutine started ", zap.Int("Worker id", workerID))

	for taskID := range taskChannel {
		w.logger.Info("task received ", zap.Int("worker_id", workerID), zap.String("task_id", taskID))

		taskData, err := w.repository.GetById(ctx, taskID)

		if err != nil {
			w.logger.Error("Failed to fetch task data ", zap.String("task_id ", taskID), zap.Error(err))
		}
		err = w.repository.UpdateStatus(ctx, taskID, task.StatusProcessing)
		if err != nil {
			w.logger.Error("Failed to update task status ", zap.String("task_id ", taskID), zap.Error(err))
			continue
		}

		w.logger.Info("Processing task ",zap.Int("worker_id", workerID), zap.String("task_id", taskID), zap.String("task_type", taskData.Type))

		time.Sleep(5 * time.Second)

		err = w.repository.UpdateStatus(ctx, taskID, task.StatusCompleted)

		if err != nil {
			w.logger.Error("Failed to update task status ", zap.String("task_id ", taskID), zap.Error(err))
			continue
		}
		
		w.logger.Info("Task completed ", zap.Int("worker_id", workerID), zap.String("task_id", taskID))

	}
}