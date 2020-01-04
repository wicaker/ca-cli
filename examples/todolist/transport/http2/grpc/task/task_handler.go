package taskhandler

import (
	"context"
	"todolist/domain"
	"todolist/middleware"
	pb "todolist/proto"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

// GrpcTaskHandler represent the grpc handler for task
type GrpcTaskHandler struct {
	TaskUsecase domain.TaskUsecase
}

// NewGrpcTaskHandler will initialize the grpc endpoint for task entity
func NewGrpcTaskHandler(gs *grpc.Server, u domain.TaskUsecase) {
	srv := &GrpcTaskHandler{
		TaskUsecase: u,
	}

	pb.RegisterTaskServiceServer(gs, srv)
}

// FetchTask will handle FetchTask request
func (gth *GrpcTaskHandler) FetchTask(ctx context.Context, req *pb.FetchTaskReq) (*pb.FetchTaskResp, error) {
	tasks := []*pb.Task{}
	tokenHeader := req.GetXAccessToken()

	if ctx == nil {
		ctx = context.Background()
	}

	// check token
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return nil, err
	}

	// fetch data
	res, err := gth.TaskUsecase.Fetch(ctx, token.ID)

	// adjust data
	for _, task := range res {
		var t pb.Task

		dueDate, err := ptypes.TimestampProto(task.DueDate)
		createdAt, err := ptypes.TimestampProto(task.CreatedAt)
		updatedAt, err := ptypes.TimestampProto(task.UpdatedAt)
		if err != nil {
			return nil, err
		}

		t.Id = task.ID
		t.Title = task.Title
		t.Description = task.Description
		t.DueDate = dueDate
		t.Completed = task.Completed
		t.CreatedAt = createdAt
		t.UpdatedAt = updatedAt
		tasks = append(tasks, &t)
	}

	// send response
	return &pb.FetchTaskResp{Task: tasks}, nil
}

// GetByIDTask will handle GetByIDTask request
func (gth *GrpcTaskHandler) GetByIDTask(ctx context.Context, req *pb.GetByIDTaskReq) (*pb.GetByIDTaskResp, error) {
	var task pb.Task

	if ctx == nil {
		ctx = context.Background()
	}

	// get token
	tokenHeader := req.GetXAccessToken()
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return nil, err
	}

	// get data
	res, err := gth.TaskUsecase.GetByID(ctx, req.GetIdTask())
	if err != nil {
		return nil, err
	}

	// adjust data
	dueDate, err := ptypes.TimestampProto(res.DueDate)
	createdAt, err := ptypes.TimestampProto(res.CreatedAt)
	updatedAt, err := ptypes.TimestampProto(res.UpdatedAt)
	if err != nil {
		return nil, err
	}

	task.Id = res.ID
	task.Title = res.Title
	task.Description = res.Description
	task.DueDate = dueDate
	task.Completed = res.Completed
	task.CreatedAt = createdAt
	task.UpdatedAt = updatedAt

	// send response
	return &pb.GetByIDTaskResp{Task: &task}, nil
}

// CreateTask will handle CreateTask request
func (gth *GrpcTaskHandler) CreateTask(ctx context.Context, req *pb.CreateTaskReq) (*pb.CreateTaskResp, error) {
	// decode req
	taskReq := req.GetTask()
	dueDateReq, err := ptypes.Timestamp(taskReq.GetDueDate())
	if err != nil {
		return nil, err
	}

	data := domain.Task{
		Title:       taskReq.GetTitle(),
		Description: taskReq.GetDescription(),
		DueDate:     dueDateReq,
		Completed:   taskReq.GetCompleted(),
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// check jwt
	tokenHeader := req.GetXAccessToken()
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return nil, err
	}
	data.UserID = token.ID

	// store data
	res, err := gth.TaskUsecase.Store(ctx, token.ID, &data)
	if err != nil {
		return nil, err
	}

	// adjust/ encode data
	var task pb.Task

	dueDate, err := ptypes.TimestampProto(res.DueDate)
	createdAt, err := ptypes.TimestampProto(res.CreatedAt)
	updatedAt, err := ptypes.TimestampProto(res.UpdatedAt)
	if err != nil {
		return nil, err
	}

	task.Id = res.ID
	task.Title = res.Title
	task.Description = res.Description
	task.DueDate = dueDate
	task.Completed = res.Completed
	task.CreatedAt = createdAt
	task.UpdatedAt = updatedAt

	// send response
	return &pb.CreateTaskResp{Task: &task}, nil
}

// UpdateTask will handle UpdateTask request
func (gth *GrpcTaskHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskReq) (*pb.UpdateTaskResp, error) {
	// decode req
	taskReq := req.GetTask()
	dueDateReq, err := ptypes.Timestamp(taskReq.GetDueDate())
	if err != nil {
		return nil, err
	}

	data := domain.Task{
		ID:          req.GetIdTask(),
		Title:       taskReq.GetTitle(),
		Description: taskReq.GetDescription(),
		DueDate:     dueDateReq,
		Completed:   taskReq.GetCompleted(),
	}

	if ctx == nil {
		ctx = context.Background()
	}

	// check jwt
	tokenHeader := req.GetXAccessToken()
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return nil, err
	}
	data.UserID = token.ID

	// store data
	err = gth.TaskUsecase.Update(ctx, &data)
	if err != nil {
		return nil, err
	}

	// adjust/ encode data
	var task pb.Task

	dueDate, err := ptypes.TimestampProto(data.DueDate)
	createdAt, err := ptypes.TimestampProto(data.CreatedAt)
	updatedAt, err := ptypes.TimestampProto(data.UpdatedAt)
	if err != nil {
		return nil, err
	}

	task.Id = data.ID
	task.Title = data.Title
	task.Description = data.Description
	task.DueDate = dueDate
	task.Completed = data.Completed
	task.CreatedAt = createdAt
	task.UpdatedAt = updatedAt

	// send response
	return &pb.UpdateTaskResp{Task: &task}, nil
}

// DeleteTask will handle DeleteTask request
func (gth *GrpcTaskHandler) DeleteTask(ctx context.Context, req *pb.DeleteTaskReq) (*pb.DeleteTaskResp, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// check jwt
	tokenHeader := req.GetXAccessToken()
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return nil, err
	}

	// delete task
	err = gth.TaskUsecase.Delete(ctx, req.GetIdTask())
	if err != nil {
		return nil, err
	}

	return &pb.DeleteTaskResp{Task: &pb.Task{Id: req.GetIdTask()}}, nil
}
