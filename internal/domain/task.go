package domain

import (
	"context"
	"time"
)

type Task struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"` // "pending", "processing", "completed"
	CreatedAt time.Time `json:"created_at"`
}

type TaskRepository interface {
	//положить задачу в очередь.
	Enqueue(ctx context.Context, task *Task) error
	//вытащить задачу из очереди
	Dequeue(ctx context.Context) (*Task, error)
}

type TaskUsecase interface {
	CreateTask(ctx context.Context, taskType, payload string) (*Task, error)
	ProcessNextTask(ctx context.Context) (*Task, error)
}
