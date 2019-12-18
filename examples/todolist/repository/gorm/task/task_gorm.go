package task

import (
	"context"
	"todolist/domain"

	"github.com/jinzhu/gorm"
)

type taskGormRepository struct {
	DB *gorm.DB
}

// NewTaskGormRepository will create new an taskGormRepository object representation of domain.TaskRepository interface
func NewTaskGormRepository(DB *gorm.DB) domain.TaskRepository {
	return &taskGormRepository{DB}
}

func (db *taskGormRepository) Fetch(ctx context.Context, userID uint64) ([]*domain.Task, error) {
	var tasks []*domain.Task
	err := db.DB.Where("UserID = ?", userID).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (db *taskGormRepository) GetByID(ctx context.Context, id uint64) (*domain.Task, error) {
	var task domain.Task

	err := db.DB.First(&task, id).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (db *taskGormRepository) Store(ctx context.Context, t *domain.Task) error {
	return db.DB.Create(t).Error
}

func (db *taskGormRepository) Update(ctx context.Context, t *domain.Task) error {
	return db.DB.Save(t).Error
}

func (db *taskGormRepository) Delete(ctx context.Context, id uint64) error {
	return db.DB.Where("id = ?", id).Delete(&domain.Task{}).Error
}
