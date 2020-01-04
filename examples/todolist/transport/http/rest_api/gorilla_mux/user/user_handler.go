package userhandler

import (
	"context"
	"net/http"
	"todolist/domain"
	"todolist/middleware"

	"github.com/gorilla/mux"
	json "github.com/json-iterator/go"
)

// UserHandler represent the httphandler for user
type UserHandler struct {
	UserUsecase domain.UserUsecase
}

// NewUserHandler will initialize the user endpoint
func NewUserHandler(r *mux.Router, u domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: u,
	}

	r.HandleFunc("/user/register", handler.Register).Methods("POST")
	r.HandleFunc("/user/login", handler.Login).Methods("POST")
}

// Register will handle register request
func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}
	if ok, err := middleware.IsRequestValid(&user); !ok {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if err := uh.UserUsecase.Register(ctx, &user); err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&domain.ResponseSuccess{Message: "Successfully register new user"})
	return
}

// Login will handle login request
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}
	if ok, err := middleware.IsRequestValid(&user); !ok {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	type data struct {
		Token string `json:"token"`
	}
	d := &data{}
	d.Token, err = uh.UserUsecase.Login(ctx, &user)

	if err != nil {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(&domain.ResponseError{Message: err.Error()})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&domain.ResponseSuccess{Message: "Login successfully", Data: d})
	return
}
