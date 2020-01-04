package middleware

import (
	"os"
	"strings"
	"todolist/domain"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwtVerify will validate a incoming jwt token in request header
func JwtVerify(token string) (*domain.Token, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		//Token is missing, returns with error code 403 Unauthorized
		return nil, domain.ErrUnauthorized
	}

	tk := &domain.Token{}
	_, err := jwt.ParseWithClaims(token, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	return tk, nil
}
