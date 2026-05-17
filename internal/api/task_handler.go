package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/RohitKMishra/distributed-task-queue/internal/queue"
	"github.com/RohitKMishra/distributed-task-queue/internal/task"
)

type TaskHandler struct {
	repository *task.Repository
	queue *queue.RedisQueue
	logger *zap.Logger
}

func NewTaskHandler(repository *task.Repository, queue *queue.RedisQueue, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{
		repository: repository,
		queue: queue,
		logger: logger,
	}
}

type CreateTaskRequest struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()
	var request CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	taskID := uuid.New().String()

	newTask := &task.Task{
		ID: taskID,
		Type: request.Type,
		Payload: request.Payload,
		Status: task.StatusPending,
		RetryCount: 0,
	}

	if err:= h.repository.Create(ctx, newTask); err != nil {
		h.logger.Error("Failed to persist task", zap.Error(err))
		http.Error(w, "failed to persist task", http.StatusInternalServerError)
		return
	}

	if err := h.queue.Enqueue(ctx, taskID); err != nil {
		h.logger.Error("Failed to enqueue task", zap.Error(err))
		http.Error(w, "failed to enqueue task", http.StatusInternalServerError)
		return
	}

	response := CreateTaskResponse {
		ID: taskID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}