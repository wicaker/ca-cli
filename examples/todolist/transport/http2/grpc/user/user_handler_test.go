package userhandler_test

import (
	"context"
	"fmt"
	"testing"

	pb "todolist/proto"

	"google.golang.org/grpc"
)

func TestRegister(t *testing.T) {
	const address = "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	fmt.Println(conn)
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
		// r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			t.Fatalf("could not greet: %v", err)
		}
		t.Logf("Greeting: %s", r)
	})
}
