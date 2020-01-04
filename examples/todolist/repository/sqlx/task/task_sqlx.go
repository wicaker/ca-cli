package task

import (
	"context"
	"database/sql"
	"time"
	"todolist/domain"

	"github.com/jmoiron/sqlx"
)

type taskSqlxRepository struct {
	Conn *sqlx.DB
}

// NewTaskSqlxRepository will create new an taskSqlxRepository object representation of domain.TaskRepository interface
func NewTaskSqlxRepository(Conn *sqlx.DB) domain.TaskRepository {
	return &taskSqlxRepository{Conn}
}

func (db *taskSqlxRepository) Fetch(ctx context.Context, userID uint64) ([]*domain.Task, error) {
	var tasks []*domain.Task

	err := db.Conn.SelectContext(ctx, &tasks, "SELECT * FROM tasks WHERE (user_id=? AND deleted_at is NULL)", userID)
	if err != nil {
		return tasks, err
	}
	return tasks, nil
}

func (db *taskSqlxRepository) GetByID(ctx context.Context, id uint64) (*domain.Task, error) {
	task := new(domain.Task)
	err := db.Conn.GetContext(ctx, task, "SELECT * FROM tasks WHERE (id=? AND deleted_at is NULL)", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return task, nil
}

func (db *taskSqlxRepository) Store(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	res, err := db.Conn.NamedExecContext(ctx, `INSERT INTO tasks (title, description, due_date, completed, user_id, created_at, updated_at) VALUES (:title,:description,:due_date,:completed,:user_id,:created_at,:updated_at)`,
		map[string]interface{}{
			"title":       task.Title,
			"description": task.Description,
			"due_date":    task.DueDate,
			"completed":   task.Completed,
			"user_id":     task.UserID,
			"created_at":  task.CreatedAt,
			"updated_at":  task.UpdatedAt,
		})

	if err != nil {
		return nil, err
	}

	idTask, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	task.ID = uint64(idTask)
	return task, nil
}

func (db *taskSqlxRepository) Update(ctx context.Context, task *domain.Task) error {
	task.UpdatedAt = time.Now()

	query := `UPDATE tasks SET title=? , description=? , due_date=?, completed=? , user_id=?, updated_at=? WHERE (id=? AND deleted_at is NULL)`
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		task.Title,
		task.Description,
		task.DueDate,
		task.Completed,
		task.UserID,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// soft delete
func (db *taskSqlxRepository) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	task := &domain.Task{DeletedAt: &now}

	query := `UPDATE tasks SET deleted_at=? WHERE (id=? AND deleted_at is NULL)`
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctx,
		task.DeletedAt,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
