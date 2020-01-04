package taskhandler

import (
	"context"
	"net/http"
	"strconv"
	"todolist/domain"
	"todolist/middleware"

	"github.com/gorilla/mux"
	json "github.com/json-iterator/go"
)

// TaskHandler represent the httphandler for task
type TaskHandler struct {
	TaskUsecase domain.TaskUsecase
}

// NewTaskHandler will initialize the task endpoint
func NewTaskHandler(r *mux.Router, u domain.TaskUsecase) {
	handler := &TaskHandler{
		TaskUsecase: u,
	}
	r.HandleFunc("/task", handler.FetchTask).Methods("GET")
	r.HandleFunc("/task/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/task", handler.Store).Methods("POST")
	r.HandleFunc("/task/{id}", handler.Update).Methods("PUT")
	r.HandleFunc("/task/{id}", handler.Delete).Methods("DELETE")
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
func (th *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// get token
	tokenHeader := r.Header.Get("x-access-token")
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	// get id of task
	idTask, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
func (th *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	var task domain.Task
	vars := mux.Vars(r)

	// decode req.body
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	// get id of task
	idTask, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return

	}
	task.ID = uint64(idTask)

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
func (th *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// get token
	tokenHeader := r.Header.Get("x-access-token")
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		w.WriteHeader(domain.GetStatusCode(err))
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	// get id of task
	vars := mux.Vars(r)
	idTask, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
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
