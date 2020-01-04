package userhandler_test

import (
	"context"
	"fmt"
	"testing"

	pb "todolist/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestRegister(t *testing.T) {
	const address = "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	t.Run("Success Register", func(t *testing.T) {
		user := &pb.User{
			Email:    "tttt@email.com",
			Password: "123",
		}
		r, err := c.Register(context.Background(), &pb.RegisterReq{User: user})
		fmt.Println("rrrr", r.GetError())
		// r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			t.Fatalf("could not greet: %v", err)
		}
		t.Logf("Greeting: %s", r)
	})
}

func TestLogin(t *testing.T) {
	const address = "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	t.Run("Success Login", func(t *testing.T) {
		user := &pb.User{
			Email:    "tttt@email.com",
			Password: "123k",
		}
		r, err := c.Login(context.Background(), &pb.LoginReq{Email: user.Email, Password: user.Password})
		// r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		fmt.Println(r.GetError())
		assert.Equal(t, "mockUser.Email", r.GetToken())
		if err != nil {
			t.Fatalf("could not greet: %v", err)
		}
		t.Logf("Greeting: %s", r)
	})
}
