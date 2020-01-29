package server

import (
	"time"
	_userGopgRepo "todolist/repository/gopg/user"
	_userHandler "todolist/transport/http2/grpc/user"
	_userUseCase "todolist/usecase/user"

	_taskGopgRepo "todolist/repository/gopg/task"
	_taskHandler "todolist/transport/http2/grpc/task"
	_taskUseCase "todolist/usecase/task"

	"github.com/go-pg/pg/v9"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// GRPCServer server
func GRPCServer(db interface{}) *grpc.Server {
	database := db.(*pg.DB)

	// slice of gRPC options
	// Here we can configure things like TLS
	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	userRepo := _userGopgRepo.NewUserGopgRepository(database)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewGrpcUserHandler(s, userUcase)

	taskRepo := _taskGopgRepo.NewGopgTaskRepository(database)
	taskUcase := _taskUseCase.NewTaskUsecase(taskRepo, userRepo, timeoutContext)
	_taskHandler.NewGrpcTaskHandler(s, taskUcase)

	return s
}
