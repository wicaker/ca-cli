package user_test

import (
	"context"
	"log"
	"os"
	"testing"
	"todolist/domain"
	user_repository "todolist/repository/go_pg/user"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	db         *pg.DB
	repository domain.UserRepository
)

func init() {
	// load .env file
	e := godotenv.Load()
	if e != nil {
		log.Print(e)
	}

	db = pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER_TEST"),
		Database: os.Getenv("DB_NAME_TEST"),
		Password: os.Getenv("DB_PASSWORD_TEST"),
	})

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	user1 := domain.User{
		Username: "test1",
		Email:    "test1@email.com",
		Password: "test123",
		Name:     "Test1",
	}
	err = db.Insert(&user1)
	if err != nil {
		panic(err)
	}

	repository = user_repository.NewUserGopgRepository(db)
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*domain.User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        true,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func TestRegister(t *testing.T) {
	user2 := domain.User{
		Username: "test2",
		Email:    "test2@email.com",
		Password: "test123",
		Name:     "Test2",
	}

	t.Run("success", func(t *testing.T) {
		err := repository.Register(context.TODO(), &user2)
		assert.NoError(t, err)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user, err := repository.GetByID(context.TODO(), 1)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, uint64(1))
		assert.Equal(t, user.Email, "test1@email.com")

	})

	t.Run("user by id not found", func(t *testing.T) {
		user, err := repository.GetByID(context.TODO(), 10)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestGetByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user, err := repository.GetByEmail(context.TODO(), "test1@email.com")
		assert.NoError(t, err)
		assert.Equal(t, user.ID, uint64(1))
		assert.Equal(t, user.Email, "test1@email.com")
	})

	t.Run("user by email not found", func(t *testing.T) {
		user, err := repository.GetByEmail(context.TODO(), "test10@email.com")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestGetByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user, err := repository.GetByUsername(context.TODO(), "test1")
		assert.NoError(t, err)
		assert.Equal(t, user.ID, uint64(1))
		assert.Equal(t, user.Email, "test1@email.com")
		assert.Equal(t, user.Username, "test1")
	})

	t.Run("user by username not found", func(t *testing.T) {
		user, err := repository.GetByUsername(context.TODO(), "test10")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestAfterTest(t *testing.T) {
	defer db.Close()
}
