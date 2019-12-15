package user

import (
	"context"
	"time"
	"todolist/domain"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an userUsecase object representation of domain.UserUsecase interface
func NewUserUsecase(ur domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       ur,
		contextTimeout: timeout,
	}
}

func (tu *userUsecase) Register(ctx context.Context, t *domain.User) error {
	return nil
}

func (tu *userUsecase) Login(ctx context.Context, t *domain.User) (string, error) {
	return "", nil
}
