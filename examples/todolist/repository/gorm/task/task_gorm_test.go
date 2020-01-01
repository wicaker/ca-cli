package task_test

import (
	"context"
	"database/sql"
	"testing"
	"todolist/domain"
	task_repository "todolist/repository/gorm/task"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gopkg.in/go-playground/assert.v1"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository domain.TaskRepository
	task       domain.Task
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.repository = task_repository.NewTaskGormRepository(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestFetch() {
	s.mock.ExpectQuery("").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(12))

	res, err := s.repository.Fetch(context.TODO(), 12)

	assert.Equal(s.T(), res[0].ID, uint64(12))
	require.NoError(s.T(), err)
}
