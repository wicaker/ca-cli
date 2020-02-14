package server

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"
	"todolist/middleware"

	_userGormRepo "todolist/repository/gorm/user"
	_userHandler "todolist/transport/http/rest_api/net_http/user"
	_userUseCase "todolist/usecase/user"

	_taskGormRepo "todolist/repository/gorm/task"
	_taskHandler "todolist/transport/http/rest_api/net_http/task"
	_taskUseCase "todolist/usecase/task"

	"github.com/jinzhu/gorm"
)

// ServeMux server
func ServeMux(db *gorm.DB) http.Handler {
	r := http.NewServeMux()
	middL := middleware.InitStandarMuxMiddleware()
	var handler http.Handler = r
	handler = middL.MiddlewareLogging(handler)
	handler = middL.CORS(handler)
	handler = parseURL(handler)

	timeoutContext := time.Duration(2) * time.Millisecond

	userRepo := _userGormRepo.NewUserGormRepository(db)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(r, userUcase)

	taskRepo := _taskGormRepo.NewTaskGormRepository(db)
	taskUcase := _taskUseCase.NewTaskUsecase(taskRepo, userRepo, timeoutContext)
	_taskHandler.NewTaskHandler(r, taskUcase)

	return handler
}

func parseURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		ctx = context.WithValue(ctx, "uri", r.RequestURI)

		s := strings.Split(r.RequestURI, "/")
		if s[1] == "task" {
			r.URL = &url.URL{Path: "/task"}
		}
		if s[1] == "user" {
			r.URL = &url.URL{Path: "/user"}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
