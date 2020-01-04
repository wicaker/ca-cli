package taskhandler_test

import (
	"context"
	"testing"
	"time"

	pb "todolist/proto"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestFetchTask(t *testing.T) {
	const address = "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTaskServiceClient(conn)

	t.Run("Success Fetch Task By User id", func(t *testing.T) {
		var tasks []*pb.Task
		r, err := c.FetchTask(context.Background(), &pb.FetchTaskReq{XAccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmFtZSI6IiIsIkVtYWlsIjoieHh4eCIsIlVzZXJuYW1lIjoiIiwiZXhwIjoxNTgzODg1Mjc2fQ.f55GU4NMy9iP2no5O2c5Vo9CG6q5DjXynhUePJX7rCo"})
		assert.NoError(t, err)
		assert.IsType(t, r.GetTask(), tasks)
	})
}

func TestGetByIDTask(t *testing.T) {
	const address = "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTaskServiceClient(conn)

	t.Run("Success get a task by id", func(t *testing.T) {
		r, err := c.GetByIDTask(context.Background(), &pb.GetByIDTaskReq{IdTask: 8, XAccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmFtZSI6IiIsIkVtYWlsIjoieHh4eCIsIlVzZXJuYW1lIjoiIiwiZXhwIjoxNTgzODg1Mjc2fQ.f55GU4NMy9iP2no5O2c5Vo9CG6q5DjXynhUePJX7rCo"})
		assert.NoError(t, err)
		assert.Equal(t, r.GetTask().GetId(), uint64(8))
	})
}

func TestCreateTask(t *testing.T) {
	const address = "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTaskServiceClient(conn)

	t.Run("Success create new task", func(t *testing.T) {
		dueDate, _ := ptypes.TimestampProto(time.Now())
		r, err := c.CreateTask(context.Background(), &pb.CreateTaskReq{
			Task: &pb.Task{
				Title:       "test grpc",
				Description: "just test grpc server",
				DueDate:     dueDate,
			},
			XAccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmFtZSI6IiIsIkVtYWlsIjoieHh4eCIsIlVzZXJuYW1lIjoiIiwiZXhwIjoxNTgzODg1Mjc2fQ.f55GU4NMy9iP2no5O2c5Vo9CG6q5DjXynhUePJX7rCo"})
		assert.NotNil(t, r.GetTask())
		if err != nil {
			t.Fatalf("could not greet: %v", err)
		}
		t.Logf("Greeting: %s", r)
	})
}

func TestUpdateTask(t *testing.T) {
	const address = "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTaskServiceClient(conn)

	t.Run("Success update a task", func(t *testing.T) {
		dueDate, _ := ptypes.TimestampProto(time.Now())
		_, err := c.UpdateTask(context.Background(), &pb.UpdateTaskReq{
			IdTask: 6,
			Task: &pb.Task{
				Title:       "test grpc",
				Description: "just test grpc server",
				DueDate:     dueDate,
				Completed:   true,
			},
			XAccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmFtZSI6IiIsIkVtYWlsIjoieHh4eCIsIlVzZXJuYW1lIjoiIiwiZXhwIjoxNTgzODg1Mjc2fQ.f55GU4NMy9iP2no5O2c5Vo9CG6q5DjXynhUePJX7rCo"})
		assert.Error(t, err)
	})
}

func TestDeleteTask(t *testing.T) {
	const address = "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTaskServiceClient(conn)

	t.Run("Success delete a task", func(t *testing.T) {
		_, err := c.DeleteTask(context.Background(), &pb.DeleteTaskReq{
			IdTask:       6,
			XAccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmFtZSI6IiIsIkVtYWlsIjoieHh4eCIsIlVzZXJuYW1lIjoiIiwiZXhwIjoxNTgzODg1Mjc2fQ.f55GU4NMy9iP2no5O2c5Vo9CG6q5DjXynhUePJX7rCo"})
		assert.Error(t, err)
		// if err != nil {
		// 	t.Fatalf("could not greet: %v", err)
		// }
		// t.Logf("Greeting: %s", r)
	})
}
