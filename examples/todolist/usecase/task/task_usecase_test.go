package task_test

import (
	"context"
	"testing"
	"time"
	"todolist/domain"
	"todolist/domain/mocks"
	ucase "todolist/usecase/task"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	dueDate, err := time.Parse("2006-01-02", "2020-05-20")
	if err != nil {
		assert.Error(t, err)
	}

	mockTaskRepo := new(mocks.TaskRepository)
	mockUserRepo := new(mocks.UserRepository)

	mockTask := &domain.Task{
		ID:          1,
		Title:       "Coding Challenge",
		Description: "Coding Challenge routine at HackerRank",
		DueDate:     dueDate,
		Completed:   false,
		UserID:      1,
	}
	mockUser := &domain.User{
		ID:       1,
		Username: "test",
		Email:    "test@email.com",
		Password: "test123",
		Name:     "Test",
	}

	mockListTask := make([]*domain.Task, 0)
	mockListTask = append(mockListTask, mockTask)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint64")).Return(mockUser, nil).Once()
		mockTaskRepo.On("Fetch", mock.Anything, mock.AnythingOfType("uint64")).Return(mockListTask, nil).Once()

		usecase := ucase.NewTaskUsecase(mockTaskRepo, mockUserRepo, time.Millisecond*100)
		listTask, err := usecase.Fetch(context.TODO(), mockUser.ID)

		assert.NotEmpty(t, listTask)
		assert.NoError(t, err)
		assert.Len(t, listTask, len(mockListTask))

		mockUserRepo.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("error no user id", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint64")).Return(nil, nil).Once()

		usecase := ucase.NewTaskUsecase(mockTaskRepo, mockUserRepo, time.Millisecond*100)
		listTask, err := usecase.Fetch(context.TODO(), mockUser.ID)

		assert.Empty(t, listTask)
		assert.Error(t, err)
		assert.Len(t, listTask, 0)

		mockUserRepo.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	dueDate, err := time.Parse("2006-01-02", "2020-05-20")
	if err != nil {
		assert.Error(t, err)
	}

	mockTaskRepo := new(mocks.TaskRepository)
	mockUserRepo := new(mocks.UserRepository)

	mockTask := &domain.Task{
		ID:          1,
		Title:       "Coding Challenge",
		Description: "Coding Challenge routine at HackerRank",
		DueDate:     dueDate,
		Completed:   false,
		UserID:      1,
	}

	t.Run("success", func(t *testing.T) {
		mockTaskRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uint64")).Return(mockTask, nil).Once()

		usecase := ucase.NewTaskUsecase(mockTaskRepo, mockUserRepo, time.Millisecond*100)
		listTask, err := usecase.GetByID(context.TODO(), 1)

		assert.NotEmpty(t, listTask)
		assert.NoError(t, err)

		mockTaskRepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {}

func TestUpdate(t *testing.T) {
}

func TestDelete(t *testing.T) {
}
