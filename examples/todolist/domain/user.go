package domain

import (
	"context"
	"time"
)

// User struct, models of User table
type User struct {
	ID        uint64     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Name      string     `json:"name"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// UserUsecase represent the User's usecases contract
type UserUsecase interface {
	Register(ctx context.Context, u *User) error
	Login(ctx context.Context, u *User) (string, error)
}

// UserRepository represent the User's repository contract
type UserRepository interface {
	GetByID(ctx context.Context, id uint64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Register(ctx context.Context, t *User) error
}
