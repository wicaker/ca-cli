package userhandler

import (
	"context"
	"todolist/domain"
	"todolist/middleware"
	pb "todolist/proto"

	"google.golang.org/grpc"
)

// GrpcUserHandler represent the grpc handler for user
type GrpcUserHandler struct {
	UserUsecase domain.UserUsecase
}

// NewGrpcUserHandler will initialize the grpc endpoint for user entity
func NewGrpcUserHandler(gs *grpc.Server, u domain.UserUsecase) {
	srv := &GrpcUserHandler{
		UserUsecase: u,
	}

	pb.RegisterUserServiceServer(gs, srv)
}

// Register will handle register request
func (guh *GrpcUserHandler) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	user := req.GetUser()
	data := domain.User{
		Email:    user.GetEmail(),
		Username: user.GetUsername(),
		Name:     user.GetName(),
		Password: user.GetPassword(),
	}

	if ok, err := middleware.IsRequestValid(&data); !ok {
		return nil, err
	}

	if err := guh.UserUsecase.Register(ctx, &data); err != nil {
		return nil, err
	}

	return &pb.RegisterResp{Message: "Successfully register new user"}, nil
}

// Login will handle login request
func (guh *GrpcUserHandler) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	return nil, nil
}
