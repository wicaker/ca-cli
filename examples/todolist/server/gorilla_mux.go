package server

import (
	"time"

	"todolist/middleware"
	_userGopgRepo "todolist/repository/go_pg/user"
	_userHandler "todolist/transport/http/rest_api/gorilla_mux/user"
	_userUseCase "todolist/usecase/user"

	"github.com/go-pg/pg/v9"
	"github.com/gorilla/mux"
)

// GorillaMux server
func GorillaMux(db interface{}) *mux.Router {
	database := db.(*pg.DB)

	r := mux.NewRouter()
	middL := middleware.InitGorillaMuxMiddleware()
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middL.MiddlewareLogging)
	timeoutContext := time.Duration(20) * time.Millisecond

	userRepo := _userGopgRepo.NewUserGopgRepository(database)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(r, userUcase)

	return r
}
