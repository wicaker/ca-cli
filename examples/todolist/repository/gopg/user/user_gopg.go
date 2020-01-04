package user

import (
	"context"
	"time"
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
	user := new(domain.User)

	err := db.DB.ModelContext(ctx, user).Where("id = ?", id).Select()
	if err != nil && err != pg.ErrNoRows {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}

	return user, nil
}

func (db *userGopgRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := new(domain.User)

	err := db.DB.ModelContext(ctx, user).Where("email = ?", email).Select()
	if err != nil && err != pg.ErrNoRows {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}

	return user, nil
}

func (db *userGopgRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := new(domain.User)

	err := db.DB.ModelContext(ctx, user).Where("username = ?", username).Select()
	if err != nil && err != pg.ErrNoRows {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}

	return user, nil
}

func (db *userGopgRepository) Register(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return db.DB.WithContext(ctx).Insert(user)
}
