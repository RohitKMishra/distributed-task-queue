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