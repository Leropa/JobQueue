package redis

import (
	"context"
	"encoding/json"
	"jobQueue/internal/domain"

	"github.com/redis/go-redis/v9"
)

const queueKey = "tasks_queue"

type TaskRepository struct {
	Client *redis.Client
}

func NewTaskRepository(client *redis.Client) *TaskRepository {
	return &TaskRepository{Client: client}
}

func (t *TaskRepository) Enqueue(ctx context.Context, task *domain.Task) error {
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return t.Client.LPush(ctx, queueKey, taskBytes).Err()
}

func (t *TaskRepository) Dequeue(ctx context.Context) (*domain.Task, error) {
	val, err := t.Client.RPop(ctx, queueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No tasks in the queue
		}
		return nil, err
	}

	var task domain.Task
	err = json.Unmarshal([]byte(val), &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
