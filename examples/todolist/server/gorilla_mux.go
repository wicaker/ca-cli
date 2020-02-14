package server

import (
	"database/sql"
	"time"

	"todolist/middleware"
	_userSQLRepo "todolist/repository/sql/user"
	_userHandler "todolist/transport/http/rest_api/gorilla_mux/user"
	_userUseCase "todolist/usecase/user"

	_taskSQLRepo "todolist/repository/sql/task"
	_taskHandler "todolist/transport/http/rest_api/gorilla_mux/task"
	_taskUseCase "todolist/usecase/task"

	"github.com/gorilla/mux"
)

// GorillaMux server
func GorillaMux(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	middL := middleware.InitGorillaMuxMiddleware()
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middL.MiddlewareLogging)
	timeoutContext := time.Duration(2) * time.Second

	userRepo := _userSQLRepo.NewUserSQLRepository(db)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(r, userUcase)

	taskRepo := _taskSQLRepo.NewTaskSQLRepository(db)
	taskUcase := _taskUseCase.NewTaskUsecase(taskRepo, userRepo, timeoutContext)
	_taskHandler.NewTaskHandler(r, taskUcase)

	return r
}
