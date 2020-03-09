package domain

import (
	"context"
	"time"
)

// Example3 struct, models of example3 table
type Example3 struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at" pg:",soft_delete"`
}

// Example3Usecase represent the Example's usecases contract
type example3Usecase interface {
	Fetch(ctx context.Context) ([]*Example, error)
	GetByID(ctx context.Context, id uint64) (*Example, error)
	Store(ctx context.Context, exp *Example) (*Example, error)
	Update(ctx context.Context, exp *Example) (*Example, error)
	Delete(ctx context.Context, id uint64) error
}

// Example3Repository represent the Example's repository contract
type example3Repository interface {
	Fetch(ctx context.Context) ([]*Example, error)
	GetByID(ctx context.Context, id uint64) (*Example, error)
	Store(ctx context.Context, exp *Example) (*Example, error)
	Update(ctx context.Context, exp *Example) (*Example, error)
	Delete(ctx context.Context, id uint64) error
}
