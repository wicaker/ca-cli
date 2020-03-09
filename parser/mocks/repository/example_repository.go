package repository

import (
	"context"
	pg "github.com/go-pg/pg/v9"
	domain "github.com/wicaker/tests/domain"
)

type gopgExampleRepository struct {
	Conn *pg.DB
}

// NewGopgExampleRepository will create new an gopgExampleRepository object representation of domain.ExampleRepository interface
func NewGopgExampleRepository(Conn *pg.DB) domain.ExampleRepository {
	return &gopgExampleRepository{Conn: Conn}
}
func (er *gopgExampleRepository) Fetch(ctx context.Context) ([]*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}
func (er *gopgExampleRepository) GetByID(ctx context.Context, id uint64) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}
func (er *gopgExampleRepository) Store(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}
func (er *gopgExampleRepository) Update(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}
func (er *gopgExampleRepository) Delete(ctx context.Context, id uint64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil
}
