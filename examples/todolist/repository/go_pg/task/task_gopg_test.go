package task_test

import (
	"context"
	"log"
	"os"
	"testing"
	"todolist/domain"
	task_repository "todolist/repository/go_pg/task"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	db         *pg.DB
	repository domain.TaskRepository
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

	task1 := domain.Task{
		Title:  "Task1",
		UserID: 1,
	}
	task2 := domain.Task{
		Title:  "Task2",
		UserID: 2,
	}

	err = db.Insert(&task1)
	err = db.Insert(&task2)
	if err != nil {
		panic(err)
	}

	repository = task_repository.NewTaskGopgRepository(db)
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*domain.Task)(nil)} {
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

func TestFetch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tasks, err := repository.Fetch(context.TODO(), 1)
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
	})

	t.Run("list task empty", func(t *testing.T) {
		tasks, err := repository.Fetch(context.TODO(), 10)
		assert.NoError(t, err)
		assert.Len(t, tasks, 0)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		task, err := repository.GetByID(context.TODO(), 1)
		assert.NoError(t, err)
		assert.Equal(t, task.ID, uint64(1))

	})

	t.Run("task by id not found", func(t *testing.T) {
		task, err := repository.GetByID(context.TODO(), 10)
		assert.Error(t, err)
		assert.Nil(t, task)
	})
}

func TestStore(t *testing.T) {
	task3 := domain.Task{
		Title:  "Task3",
		UserID: 3,
	}
	t.Run("success", func(t *testing.T) {
		err := repository.Store(context.TODO(), &task3)
		assert.NoError(t, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		task3 := domain.Task{
			ID:     3,
			Title:  "UpdatedTask3",
			UserID: 3,
		}

		err := repository.Update(context.TODO(), &task3)
		assert.NoError(t, err)

		task, err := repository.GetByID(context.TODO(), 3)
		assert.Equal(t, task.Title, task3.Title)
		assert.NoError(t, err)
	})

	t.Run("error because no task.ID provided", func(t *testing.T) {
		task3 := domain.Task{
			Title:  "UpdatedTask3",
			UserID: 3,
		}

		err := repository.Update(context.TODO(), &task3)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := repository.Delete(context.TODO(), 3)
		assert.NoError(t, err)
	})

	t.Run("error no task id found", func(t *testing.T) {
		err := repository.Delete(context.TODO(), 3)
		assert.Error(t, err)
	})
}

func TestAfterTest(t *testing.T) {
	defer db.Close()
}
