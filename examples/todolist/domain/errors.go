package domain

import "errors"

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Given Param is not valid")
	// ErrUnauthorized will throw if the given request-header token is not valid
	ErrUnauthorized = errors.New("Unauthorized")
)
