package task

import (
	"context"
	"errors"
	"todolist/domain"

	"github.com/go-pg/pg/v9"
)

type taskGopgRepository struct {
	DB *pg.DB
}

// NewTaskGopgRepository will create new an taskGopgRepository object representation of domain.TaskRepository interface
func NewTaskGopgRepository(DB *pg.DB) domain.TaskRepository {
	return &taskGopgRepository{DB}
}

func (db *taskGopgRepository) Fetch(ctx context.Context, userID uint64) ([]*domain.Task, error) {
	var tasks []*domain.Task

	err := db.DB.Model(&tasks).Where("user_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (db *taskGopgRepository) GetByID(ctx context.Context, id uint64) (*domain.Task, error) {
	task := &domain.Task{ID: id}

	err := db.DB.Select(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (db *taskGopgRepository) Store(ctx context.Context, t *domain.Task) error {
	return db.DB.Insert(t)
}

func (db *taskGopgRepository) Update(ctx context.Context, t *domain.Task) error {
	if t.ID == uint64(0) {
		return errors.New("Task ID required")
	}
	return db.DB.Update(t)
}

func (db *taskGopgRepository) Delete(ctx context.Context, id uint64) error {
	task := domain.Task{ID: id}
	return db.DB.Delete(&task)
}
