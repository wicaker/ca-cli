package server

import (
	"net/http"
	"time"
	"todolist/middleware"

	_userGopgRepo "todolist/repository/go_pg/user"
	_graphqlHandler "todolist/transport/http/graphql"
	_userUseCase "todolist/usecase/user"

	"github.com/go-pg/pg/v9"
)

// GraphQLServer server
func GraphQLServer(db interface{}) http.Handler {
	database := db.(*pg.DB)

	r := http.NewServeMux()
	middL := middleware.InitStandarMuxMiddleware()
	var handler http.Handler = r
	handler = middL.MiddlewareLogging(handler)
	handler = middL.CORS(handler)

	timeoutContext := time.Duration(90) * time.Millisecond

	userRepo := _userGopgRepo.NewUserGopgRepository(database)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_graphqlHandler.NewGraphQLHandler(r, userUcase)

	return handler
}
