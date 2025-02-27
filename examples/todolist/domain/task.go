package domain

import (
	"context"
	"time"
)

// Task struct, models of Task table
type Task struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     time.Time  `json:"due_date" db:"due_date"`
	Completed   bool       `json:"completed"`
	UserID      uint64     `json:"user_id" db:"user_id"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at" pg:",soft_delete"`
}

// TaskUsecase represent the Task's usecases contract
type TaskUsecase interface {
	Fetch(ctx context.Context, userID uint64) ([]*Task, error)
	GetByID(ctx context.Context, id uint64) (*Task, error)
	Store(ctx context.Context, userID uint64, t *Task) (*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id uint64) error
}

// TaskRepository represent the Task's repository contract
type TaskRepository interface {
	Fetch(ctx context.Context, userID uint64) ([]*Task, error)
	GetByID(ctx context.Context, id uint64) (*Task, error)
	Store(ctx context.Context, t *Task) (*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id uint64) error
}
