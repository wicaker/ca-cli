package middleware

import "gopkg.in/go-playground/validator.v9"

// IsRequestValid will validate a incoming request
func IsRequestValid(m interface{}) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
