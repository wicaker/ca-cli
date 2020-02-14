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
	"google.golang.org/grpc"
)

// GRPCServer server
func GRPCServer(db *pg.DB) *grpc.Server {
	// slice of gRPC options
	// Here we can configure things like TLS
	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	timeoutContext := time.Duration(2) * time.Second

	userRepo := _userGopgRepo.NewUserGopgRepository(db)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewGrpcUserHandler(s, userUcase)

	taskRepo := _taskGopgRepo.NewGopgTaskRepository(db)
	taskUcase := _taskUseCase.NewTaskUsecase(taskRepo, userRepo, timeoutContext)
	_taskHandler.NewGrpcTaskHandler(s, taskUcase)

	return s
}
