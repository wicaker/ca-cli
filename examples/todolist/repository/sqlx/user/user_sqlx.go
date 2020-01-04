package user

import (
	"context"
	"database/sql"
	"time"
	"todolist/domain"

	"github.com/jmoiron/sqlx"
)

type userSqlxRepository struct {
	Conn *sqlx.DB
}

// NewUserSqlxRepository will create new an userSqlxRepository object representation of domain.UserRepository interface
func NewUserSqlxRepository(Conn *sqlx.DB) domain.UserRepository {
	return &userSqlxRepository{Conn}
}

func (db *userSqlxRepository) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	var (
		sqlUsername sql.NullString
		sqlEmail    sql.NullString
		sqlName     sql.NullString
	)

	user := new(domain.User)

	sqlStatement := `SELECT id,username,email,password,name,created_at,updated_at  FROM users WHERE id = ? `
	rows, err := db.Conn.QueryContext(ctx, sqlStatement, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&sqlUsername,
			&sqlEmail,
			&user.Password,
			&sqlName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.Username = sqlUsername.String
		user.Email = sqlEmail.String
		user.Name = sqlName.String
	}
	defer rows.Close()

	if user.ID == 0 {
		return nil, nil
	}

	return user, nil
}

func (db *userSqlxRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var (
		sqlUsername sql.NullString
		sqlEmail    sql.NullString
		sqlName     sql.NullString
	)
	user := new(domain.User)

	sqlStatement := `SELECT id,username,email,password,name,created_at,updated_at  FROM users WHERE email = ? `
	rows, err := db.Conn.QueryContext(ctx, sqlStatement, email)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&sqlUsername,
			&sqlEmail,
			&user.Password,
			&sqlName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.Username = sqlUsername.String
		user.Email = sqlEmail.String
		user.Name = sqlName.String

	}
	defer rows.Close()

	if user.ID == 0 {
		return nil, nil
	}

	return user, nil
}

func (db *userSqlxRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var (
		sqlUsername sql.NullString
		sqlEmail    sql.NullString
		sqlName     sql.NullString
	)

	user := new(domain.User)

	sqlStatement := `SELECT id,username,email,password,name,created_at,updated_at  FROM users WHERE username = ? `
	rows, err := db.Conn.QueryContext(ctx, sqlStatement, username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&sqlUsername,
			&sqlEmail,
			&user.Password,
			&sqlName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.Username = sqlUsername.String
		user.Email = sqlEmail.String
		user.Name = sqlName.String
	}
	defer rows.Close()

	if user.ID == 0 {
		return nil, nil
	}

	return user, nil
}

func (db *userSqlxRepository) Register(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := db.Conn.NamedExecContext(ctx, `INSERT INTO users (username, email, password, name, created_at, updated_at) VALUES (:username,:email,:password,:name,:created_at,:updated_at)`,
		map[string]interface{}{
			"username":   newNullString(user.Username),
			"email":      newNullString(user.Email),
			"password":   newNullString(user.Password),
			"name":       newNullString(user.Name),
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	if err != nil {
		return err
	}

	return nil
}

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
