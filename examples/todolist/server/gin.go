package server

import (
	"os"
	"time"
	"todolist/middleware"

	_userGopgRepo "todolist/repository/gopg/user"
	_userHandler "todolist/transport/http/rest_api/gin/user"
	_userUseCase "todolist/usecase/user"

	_taskGopgRepo "todolist/repository/gopg/task"
	_taskHandler "todolist/transport/http/rest_api/gin/task"
	_taskUseCase "todolist/usecase/task"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/spf13/viper"
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

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	userRepo := _userGopgRepo.NewUserGopgRepository(database)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewUserHandler(r, userUcase)

	taskRepo := _taskGopgRepo.NewTaskGopgRepository(database)
	taskUcase := _taskUseCase.NewTaskUsecase(taskRepo, userRepo, timeoutContext)
	_taskHandler.NewTaskHandler(r, taskUcase)

	return r

}
