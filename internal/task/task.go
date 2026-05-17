package task

import "time"

type Status string

const (
	StatusPending Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted Status = "completed"
	StatusFailed Status = "failed"
)

type Task struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Payload []byte `json:"payload"`
	Status Status `json:"status"`
	RetryCount int `json:"retry_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

