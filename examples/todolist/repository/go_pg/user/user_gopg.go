package user

import (
	"context"
	"todolist/domain"

	"github.com/go-pg/pg/v9"
)

type userGopgRepository struct {
	DB *pg.DB
}

// NewUserGopgRepository will create new an userGopgRepository object representation of domain.UserRepository interface
func NewUserGopgRepository(DB *pg.DB) domain.UserRepository {
	return &userGopgRepository{DB}
}

func (db *userGopgRepository) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	user := &domain.User{ID: id}

	err := db.DB.Select(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *userGopgRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := new(domain.User)
	// var user domain.User
	err := db.DB.Model(user).
		Relation("Author").
		Where("user.email = ?", email).
		Select()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *userGopgRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := new(domain.User)
	// var user domain.User
	err := db.DB.Model(user).
		Relation("Author").
		Where("user.username = ?", username).
		Select()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *userGopgRepository) Register(ctx context.Context, user *domain.User) error {
	return db.DB.Insert(user)
}
