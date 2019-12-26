package server

import (
	"time"
	"todolist/middleware"

	_userGopgRepo "todolist/repository/go_pg/user"
	_userHandler "todolist/transport/http/rest_api/echo/user"
	_userUseCase "todolist/usecase/user"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

// Echo server
func Echo(db interface{}) *echo.Echo {
	database := db.(*pg.DB)

	e := echo.New()
	middL := middleware.InitEchoMiddleware()
	e.Use(middL.MiddlewareLogging)
	e.Use(middL.CORS)

	timeoutContext := time.Duration(2) * time.Millisecond

	userRepo := _userGopgRepo.NewUserGopgRepository(database)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(e, userUcase)

	return e
}
