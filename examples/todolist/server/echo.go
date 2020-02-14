package server

import (
	"time"
	"todolist/middleware"

	_userSqlxRepo "todolist/repository/sqlx/user"
	_userHandler "todolist/transport/http/rest_api/echo/user"
	_userUseCase "todolist/usecase/user"

	_taskSqlxRepo "todolist/repository/sqlx/task"
	_taskHandler "todolist/transport/http/rest_api/echo/task"
	_taskUseCase "todolist/usecase/task"

	// "github.com/go-pg/pg/v9"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

// Echo server
func Echo(db *sqlx.DB) *echo.Echo {
	e := echo.New()
	middL := middleware.InitEchoMiddleware()
	e.Use(middL.MiddlewareLogging)
	e.Use(middL.CORS)

	timeoutContext := time.Duration(2) * time.Second

	userRepo := _userSqlxRepo.NewUserSqlxRepository(db)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(e, userUcase)

	taskRepo := _taskSqlxRepo.NewTaskSqlxRepository(db)
	taskUcase := _taskUseCase.NewTaskUsecase(taskRepo, userRepo, timeoutContext)
	_taskHandler.NewTaskHandler(e, taskUcase)

	return e
}
