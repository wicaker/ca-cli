package task

import (
	"context"
	"time"
	"todolist/domain"

	"database/sql"
)

type taskSQLRepository struct {
	Conn *sql.DB
}

// NewTaskSQLRepository will create new an taskSQLRepository object representation of domain.TaskRepository interface
func NewTaskSQLRepository(Conn *sql.DB) domain.TaskRepository {
	return &taskSQLRepository{Conn}
}

func (db *taskSQLRepository) Fetch(ctx context.Context, userID uint64) ([]*domain.Task, error) {
	tasks := make([]*domain.Task, 0)

	query := `SELECT id,title,description,due_date,completed,user_id,created_at,updated_at FROM tasks WHERE (user_id=? AND deleted_at is NULL) `
	rows, err := db.Conn.QueryContext(ctx, query, userID)

	if err != nil {
		return tasks, err
	}

	for rows.Next() {
		t := new(domain.Task)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.DueDate,
			&t.Completed,
			&t.UserID,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	defer rows.Close()

	return tasks, nil
}

func (db *taskSQLRepository) GetByID(ctx context.Context, id uint64) (*domain.Task, error) {
	task := new(domain.Task)

	query := `SELECT id,title,description,due_date,completed,user_id,created_at,updated_at FROM tasks WHERE (id=? AND deleted_at is NULL) `
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Completed,
			&task.UserID,
			&task.UpdatedAt,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	if task.ID == 0 {
		return nil, nil
	}

	return task, nil
}

func (db *taskSQLRepository) Store(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	query := `INSERT  tasks SET title=? , description=? , due_date=?, completed=? , user_id=?, created_at=?, updated_at=?`
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(
		ctx,
		task.Title,
		task.Description,
		task.DueDate,
		task.Completed,
		task.UserID,
		task.CreatedAt,
		task.UpdatedAt,
	)
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

func (db *taskSQLRepository) Update(ctx context.Context, task *domain.Task) error {
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
func (db *taskSQLRepository) Delete(ctx context.Context, id uint64) error {
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
