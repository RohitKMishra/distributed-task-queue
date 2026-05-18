package queue

import (
	"context"
	"github.com/redis/go-redis/v9"
)

const TaskQueueName = "tasks"

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{
		client:client,
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, taskID string) error {
	return q.client.LPush(
		ctx,
		TaskQueueName,
		taskID,
	).Err()
}

func (q *RedisQueue) Dequeue(ctx context.Context) (string, error) {

	result, err := q.client.BRPop(ctx, 0, TaskQueueName).Result()

	if err != nil {
		return "", err
	}

	return result[1], nil
}