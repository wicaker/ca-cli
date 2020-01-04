package taskhandler

import (
	"context"
	"net/http"
	"path"
	"strconv"
	"strings"
	"todolist/domain"
	"todolist/middleware"

	json "github.com/json-iterator/go"
)

// TaskHandler represent the httphandler for task
type TaskHandler struct {
	TaskUsecase domain.TaskUsecase
}

// NewTaskHandler will initialize the task endpoint
func NewTaskHandler(r *http.ServeMux, u domain.TaskUsecase) {
	handler := &TaskHandler{
		TaskUsecase: u,
	}
	r.Handle("/task", handler)
	// r.HandleFunc("/task/:id", handler.GetByID)
	// r.HandleFunc("/task/:id", handler.Update)
	// r.HandleFunc("/task/:id", handler.Delete)
}

func (th *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uri := ctx.Value("uri")
	s := strings.Split(r.RequestURI, "/")

	if uri == "/task" {
		if r.Method == http.MethodGet {
			th.FetchTask(w, r)
			return
		}
		if r.Method == http.MethodPost {
			th.Store(w, r)
			return
		}
	}

	if len(s) == 3 {
		// get id of task
		idTask, err := strconv.Atoi(path.Base(r.RequestURI))
		if err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(&domain.ResponseError{Message: "Method Not Allowed"})
			return
		}
		if r.Method == http.MethodGet {
			th.GetByID(w, r, uint64(idTask))
			return
		}
		if r.Method == http.MethodPut {
			th.Update(w, r, uint64(idTask))
			return
		}
		if r.Method == http.MethodDelete {
			th.Delete(w, r, uint64(idTask))
			return
		}
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(&domain.ResponseError{Message: "Method Not Allowed"})
	return
}

// FetchTask will handle FetchTask request
func (th *TaskHandler) FetchTask(w http.ResponseWriter, r *http.Request) {
	tasks := []*domain.Task{}
	tokenHeader := r.Header.Get("x-access-token")
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := th.TaskUsecase.Fetch(ctx, token.ID)
	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	if len(res) != 0 {
		tasks = res
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&domain.ResponseSuccess{Message: "Successfully load data", Data: &tasks})
	return
}

// GetByID will handle GetByID request
func (th *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request, idTask uint64) {
	// get token
	tokenHeader := r.Header.Get("x-access-token")
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	task, err := th.TaskUsecase.GetByID(ctx, uint64(idTask))

	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&domain.ResponseSuccess{Message: "Successfully load data", Data: &task})
	return
}

// Store will handle Store request
func (th *TaskHandler) Store(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: "Method Not Allowed"})
		return
	}

	var task domain.Task

	// get token
	tokenHeader := r.Header.Get("x-access-token")
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	// decode req.body
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	if ok, err := middleware.IsRequestValid(&task); !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	resp, err := th.TaskUsecase.Store(ctx, token.ID, &task)

	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&domain.ResponseSuccess{Message: "Successfully create new task", Data: resp})
	return

}

// Update will handle Update request
func (th *TaskHandler) Update(w http.ResponseWriter, r *http.Request, idTask uint64) {
	var task domain.Task

	// decode req.body
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	// get id of task
	task.ID = idTask

	// get token
	tokenHeader := r.Header.Get("x-access-token")
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}
	task.UserID = token.ID

	if ok, err := middleware.IsRequestValid(&task); !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = th.TaskUsecase.Update(ctx, &task)

	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&domain.ResponseSuccess{Message: "Successfully update data", Data: task})
	return
}

// Delete will handle Delete request
func (th *TaskHandler) Delete(w http.ResponseWriter, r *http.Request, idTask uint64) {
	// get token
	tokenHeader := r.Header.Get("x-access-token")
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = th.TaskUsecase.Delete(ctx, uint64(idTask))

	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&domain.ResponseSuccess{Message: "Delete item successfully"})
	return
}
