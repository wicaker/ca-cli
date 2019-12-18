package user_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"todolist/domain"
	user_repository "todolist/repository/gorm/user"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository domain.UserRepository
	user       domain.User
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

	s.repository = user_repository.NewUserGormRepository(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetByID() {
	id := uint64(1)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 1)`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(id))

	res, err := s.repository.GetByID(context.TODO(), id)

	assert.Equal(s.T(), &domain.User{ID: id}, res)
	require.NoError(s.T(), err)
}

func (s *Suite) TestGetByEmail() {
	email := "test@email.com"

	s.mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(0, email))

	res, err := s.repository.GetByEmail(context.TODO(), email)
	assert.Equal(s.T(), &domain.User{Email: email}, res)
	require.NoError(s.T(), err)
}

func (s *Suite) TestGetByUsername() {
	username := "test123"

	s.mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow(username))

	res, err := s.repository.GetByUsername(context.TODO(), username)
	assert.Equal(s.T(), &domain.User{Username: username}, res)
	require.NoError(s.T(), err)
}

func (s *Suite) TestRegister() {
	// user := &domain.User{
	// 	ID:        1,
	// 	Username:  "test123",
	// 	Email:     "test123@email.com",
	// 	Password:  "xxx",
	// 	Name:      "Test",
	// 	UpdatedAt: time.Now(),
	// 	CreatedAt: time.Now(),
	// 	DeletedAt: nil,
	// }
	// s.mock.ExpectQuery(regexp.QuoteMeta(
	// 	`INSERT INTO "users" ("id", "username","email","password","name","updated_at","created_at","deleted_at")
	// 		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "users"."id"`)).
	// 	WithArgs(user.ID, user.Username, user.Email, user.Password, user.Name, user.UpdatedAt, user.CreatedAt, user.DeletedAt).
	// 	WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

	// err := s.repository.Register(context.TODO(), user)
	// require.Error(s.T(), err)
}
