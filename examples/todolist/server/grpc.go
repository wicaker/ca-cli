package server

import (
	"time"
	_userGopgRepo "todolist/repository/go_pg/user"
	_userHandler "todolist/transport/http2/grpc/user"
	_userUseCase "todolist/usecase/user"

	"github.com/go-pg/pg/v9"
	"google.golang.org/grpc"
)

// GRPCServer server
func GRPCServer(db interface{}) *grpc.Server {
	database := db.(*pg.DB)

	// slice of gRPC options
	// Here we can configure things like TLS
	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	timeoutContext := time.Duration(100) * time.Millisecond

	userRepo := _userGopgRepo.NewUserGopgRepository(database)
	userUcase := _userUseCase.NewUserUsecase(userRepo, timeoutContext)
	_userHandler.NewGrpcUserHandler(s, userUcase)

	return s
}
