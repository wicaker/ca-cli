package domain

import (
	"context"
	"time"
)

// Example4 struct, models of example table
type Example4 struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at" pg:",soft_delete"`
}

// Example5Usecase represent the Example's usecases contract
type Example5Usecase interface {
	Fetch(ctx context.Context) ([]*Example, error)
	GetByID(ctx context.Context, id uint64) (*Example, error)
	Store(ctx context.Context, exp *Example) (*Example, error)
	Update(ctx context.Context, exp *Example) (*Example, error)
	Delete(ctx context.Context, id uint64) error
}

// Example5Repository represent the Example's repository contract
type Example5Repository interface {
	Fetch(ctx context.Context) ([]*Example, error)
	GetByID(ctx context.Context, id uint64) (*Example, error)
	Store(ctx context.Context, exp *Example) (*Example, error)
	Update(ctx context.Context, exp *Example) (*Example, error)
	Delete(ctx context.Context, id uint64) error
}
