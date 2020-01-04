package task

import (
	"context"
	"errors"
	"time"
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

	err := db.DB.ModelContext(ctx, &tasks).Where("user_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (db *taskGopgRepository) GetByID(ctx context.Context, id uint64) (*domain.Task, error) {
	task := &domain.Task{ID: id}

	err := db.DB.WithContext(ctx).Select(task)
	if err != nil {

		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return task, nil
}

func (db *taskGopgRepository) Store(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	return task, db.DB.WithContext(ctx).Insert(task)
}

func (db *taskGopgRepository) Update(ctx context.Context, task *domain.Task) error {
	task.UpdatedAt = time.Now()
	if task.ID == uint64(0) {
		return errors.New("Task ID required")
	}

	// _, err := db.DB.ModelContext(ctx, task).WherePK().UpdateNotZero()
	// if err != nil {
	// 	return err
	// }

	return db.DB.WithContext(ctx).Update(task)
}

func (db *taskGopgRepository) Delete(ctx context.Context, id uint64) error {
	task := domain.Task{ID: id}
	return db.DB.WithContext(ctx).Delete(&task)
}
