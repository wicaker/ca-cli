package task

import (
	"context"
	"time"
	"todolist/domain"
)

type taskUsecase struct {
	taskRepo       domain.TaskRepository
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewTaskUsecase will create new an taskUsecase object representation of domain.TaskUsecase interface
func NewTaskUsecase(tr domain.TaskRepository, ur domain.UserRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepo:       tr,
		userRepo:       ur,
		contextTimeout: timeout,
	}
}

func (tu *taskUsecase) Fetch(ctx context.Context, userID uint64) ([]*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	user, err := tu.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrNotFound
	}

	tasks, err := tu.taskRepo.Fetch(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (tu *taskUsecase) GetByID(ctx context.Context, id uint64) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	task, err := tu.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, domain.ErrNotFound
	}

	return task, nil
}

func (tu *taskUsecase) Store(ctx context.Context, userID uint64, t *domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	user, err := tu.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrNotFound
	}

	t.UserID = userID
	task, err := tu.taskRepo.Store(ctx, t)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (tu *taskUsecase) Update(ctx context.Context, t *domain.Task) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	task, err := tu.taskRepo.GetByID(ctx, t.ID)
	if err != nil {
		return err
	}
	if task == nil {
		return domain.ErrNotFound
	}

	t.CreatedAt = task.CreatedAt
	t.UpdatedAt = time.Now()
	err = tu.taskRepo.Update(ctx, t)
	if err != nil {
		return err
	}
	return nil
}

func (tu *taskUsecase) Delete(ctx context.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	task, err := tu.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return domain.ErrNotFound
	}

	err = tu.taskRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
