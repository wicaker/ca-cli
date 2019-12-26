package server

import (
	"os"
	"time"
	"todolist/middleware"
	_userGopgRepo "todolist/repository/go_pg/user"
	_userHandler "todolist/transport/http/rest_api/gin/user"
	_userUseCase "todolist/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
)

// Gin server
func Gin(db interface{}) *gin.Engine {
	database := db.(*pg.DB)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	middL := middleware.InitGinMiddleware()
	r.Use(middL.CORS())
	r.Use(middL.MiddlewareLogging())
	r.Use(gin.Recovery())

	timeoutContext := time.Duration(1) * time.Millisecond

	userRepo := _userGopgRepo.NewUserGopgRepository(database)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(r, userUcase)
	return r

}
