package user

import (
	"context"
	"todolist/domain"

	"github.com/jinzhu/gorm"
)

type userGormRepository struct {
	DB *gorm.DB
}

// NewUserGormRepository will create new an userGormRepository object representation of domain.UserRepository interface
func NewUserGormRepository(DB *gorm.DB) domain.UserRepository {
	return &userGormRepository{DB}
}

func (db *userGormRepository) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User

	err := db.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (db *userGormRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	err := db.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (db *userGormRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	err := db.DB.Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (db *userGormRepository) Register(ctx context.Context, user *domain.User) error {
	return db.DB.Create(user).Error
}
