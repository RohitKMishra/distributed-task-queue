package task

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, task *Task) error {
	query := `INSERT INTO tasks (id, type, payload, status, retry_count) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query, task.ID, task.Type, task.Payload, task.Status, task.RetryCount)
	return err
}

func (r *Repository) GetById(ctx context.Context, id string) (*Task, error) {
	query := `SELECT id, type, payload, status, retry_count, created_at, updated_at FROM tasks WHERE id = $1`

	var task Task

	err:= r.db.QueryRow(ctx, query, id).Scan(&task.ID, &task.Type, &task.Payload, &task.Status, &task.RetryCount, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &task, nil

}

func (r *Repository) UpdateStatus(ctx context.Context, taskID string, status Status) error {
	query := `UPDATE tasks SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, taskID)
	return err
}