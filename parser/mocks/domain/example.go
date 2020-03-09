package domain

import (
	"context"
	"time"
)

// Example struct, models of example table
type Example struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at" pg:",soft_delete"`
}

// ExampleUsecase represent the Example's usecases contract
type ExampleUsecase interface {
	Fetch(ctx context.Context) ([]*Example, error)
	GetByID(ctx context.Context, id uint64) (*Example, error)
	Store(ctx context.Context, exp *Example) (*Example, error)
	Update(ctx context.Context, exp *Example) (*Example, error)
	Delete(ctx context.Context, id uint64) error
}

// ExampleRepository represent the Example's repository contract
type ExampleRepository interface {
	Fetch(ctx context.Context) ([]*Example, error)
	GetByID(ctx context.Context, id uint64) (*Example, error)
	Store(ctx context.Context, exp *Example) (*Example, error)
	Update(ctx context.Context, exp *Example) (*Example, error)
	Delete(ctx context.Context, id uint64) error
}
