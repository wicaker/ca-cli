package server

import (
	"net/http"
	"time"
	"todolist/middleware"

	_userGopgRepo "todolist/repository/go_pg/user"
	_userHandler "todolist/transport/http/rest_api/net_http/user"
	_userUseCase "todolist/usecase/user"

	"github.com/go-pg/pg/v9"
)

// ServeMux server
func ServeMux(db interface{}) http.Handler {
	database := db.(*pg.DB)

	r := http.NewServeMux()
	middL := middleware.InitStandarMuxMiddleware()
	var handler http.Handler = r
	handler = middL.MiddlewareLogging(handler)
	handler = middL.CORS(handler)

	timeoutContext := time.Duration(90) * time.Millisecond

	userRepo := _userGopgRepo.NewUserGopgRepository(database)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(r, userUcase)

	return handler
}
