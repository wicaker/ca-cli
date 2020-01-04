package user_test

import (
	"context"
	"testing"
	"time"

	"todolist/domain"
	"todolist/domain/mocks"
	ucase "todolist/usecase/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := domain.User{
		ID:       1,
		Username: "test",
		Email:    "test@email.com",
		Password: "test123",
		Name:     "Test",
	}

	t.Run("success", func(t *testing.T) {
		// noExistedUser := &domain.User{}
		tempMockUser := mockUser
		tempMockUser.ID = 0

		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Once()
		mockUserRepo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Once()
		mockUserRepo.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		usecase := ucase.NewUserUsecase(mockUserRepo, time.Second*2)
		err := usecase.Register(context.TODO(), &tempMockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUser.Email, tempMockUser.Email)
		assert.NotEqual(t, mockUser.Password, tempMockUser.Password)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	// Password Encryption
	pass, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	if err != nil {
		assert.Error(t, err)
	}

	mockUser := domain.User{
		ID:       1,
		Username: "test",
		Email:    "test@email.com",
		Password: string(pass),
		Name:     "Test",
	}

	t.Run("success", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.ID = 0
		tempMockUser.Password = "test123"

		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(&mockUser, nil).Once()
		// mockUserRepo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(&mockUser, nil).Once()
		// mockUserRepo.On("Login", mock.Anything, mock.AnythingOfType("*domain.User")).Return(mock.AnythingOfType("string"), nil).Once()

		usecase := ucase.NewUserUsecase(mockUserRepo, time.Second*2)
		_, err = usecase.Login(context.TODO(), &tempMockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUser.Email, tempMockUser.Email)
		assert.NotEqual(t, mockUser.Password, tempMockUser.Password)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("failed CompareHashAndPassword", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.ID = 0

		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(&mockUser, nil).Once()
		// mockUserRepo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(&mockUser, nil).Once()
		// mockUserRepo.On("Login", mock.Anything, mock.AnythingOfType("*domain.User")).Return("", mock.AnythingOfType("error")).Once()

		usecase := ucase.NewUserUsecase(mockUserRepo, time.Millisecond*100)
		_, err = usecase.Login(context.TODO(), &tempMockUser)

		assert.Error(t, err)
		assert.Equal(t, mockUser.Email, tempMockUser.Email)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("failed context deadline exceeded", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.ID = 0

		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(&mockUser, nil).Once()
		mockUserRepo.On("GetByUsername", mock.Anything, mock.AnythingOfType("string")).Return(&mockUser, nil).Once()

		usecase := ucase.NewUserUsecase(mockUserRepo, time.Millisecond*1)
		token, err := usecase.Login(context.TODO(), &tempMockUser)

		assert.Error(t, err)
		assert.Equal(t, token, "")
	})

}
