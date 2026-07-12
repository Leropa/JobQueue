package usecase

import (
	"context"
	"jobQueue/internal/domain"
	"log"
	"time"
)

type TaskUsecase struct {
	repo domain.TaskRepository
}

func NewTaskUsecase(repo domain.TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: repo}
}

func (u *TaskUsecase) AddTask(ctx context.Context, task *domain.Task) error {
	return u.repo.Enqueue(ctx, task)
}

func (u *TaskUsecase) StartWorker(ctx context.Context) {
	log.Println("Worker started successfully...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped successfully...")
			return
		default:
		}

		task, err := u.repo.Dequeue(ctx)
		if err != nil {
			log.Printf("Error dequeuing task: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if task == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Processing task: ID=%s, Type=%s, Payload=%s", task.ID, task.Type, task.Payload)

		time.Sleep(500 * time.Millisecond)
	}
}
