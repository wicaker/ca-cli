package user

import (
	"context"
	"database/sql"
	"time"

	"todolist/domain"
)

type userSQLRepository struct {
	Conn *sql.DB
}

// NewUserSQLRepository will create new an userSQLRepository object representation of domain.UserRepository interface
func NewUserSQLRepository(Conn *sql.DB) domain.UserRepository {
	return &userSQLRepository{Conn}
}

func (db *userSQLRepository) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
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

func (db *userSQLRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
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

func (db *userSQLRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
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

func (db *userSQLRepository) Register(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `INSERT  users SET username=? , email=? , password=?, name=? , created_at=?, updated_at=?`
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		newNullString(user.Username),
		newNullString(user.Email),
		newNullString(user.Password),
		newNullString(user.Name),
		user.CreatedAt,
		user.UpdatedAt,
	)
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
