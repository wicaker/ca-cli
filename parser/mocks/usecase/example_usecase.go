package usecase

import (
	"context"
	"time"

	"github.com/wicaker/cacli/parser/mocks/domain"
)

type exampleUsecase struct {
	exampleRepo    domain.ExampleRepository
	contextTimeout time.Duration
}

// NewExampleUsecase will create new an exampleUsecase object representation of domain.ExampleUsecase interface
func NewExampleUsecase(er domain.ExampleRepository, timeout time.Duration) domain.ExampleUsecase {
	return &exampleUsecase{
		contextTimeout: timeout,
		exampleRepo:    er,
	}
}
func (eu *exampleUsecase) Fetch(ctx context.Context) ([]*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}
func (eu *exampleUsecase) GetByID(ctx context.Context, id uint64) (*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}
func (eu *exampleUsecase) Store(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}
func (eu *exampleUsecase) Update(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}
func (eu *exampleUsecase) Delete(ctx context.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil
}
