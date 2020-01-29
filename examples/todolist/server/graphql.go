package server

import (
	"net/http"
	"time"
	"todolist/middleware"

	_userGopgRepo "todolist/repository/gopg/user"
	_graphqlHandler "todolist/transport/http/graphql"
	_userUseCase "todolist/usecase/user"

	_taskGopgRepo "todolist/repository/gopg/task"
	// _taskHandler "todolist/transport/http/rgraphql"
	_taskUseCase "todolist/usecase/task"

	"github.com/go-pg/pg/v9"
	"github.com/spf13/viper"
)

// GraphQLServer server
func GraphQLServer(db interface{}) http.Handler {
	database := db.(*pg.DB)

	r := http.NewServeMux()
	middL := middleware.InitStandarMuxMiddleware()
	var handler http.Handler = r
	handler = middL.MiddlewareLogging(handler)
	handler = middL.CORS(handler)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	userRepo := _userGopgRepo.NewUserGopgRepository(database)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)

	taskRepo := _taskGopgRepo.NewGopgTaskRepository(database)
	taskUcase := _taskUseCase.NewTaskUsecase(taskRepo, userRepo, timeoutContext)

	_graphqlHandler.NewGraphQLHandler(r, userUcase, taskUcase)

	return handler
}
