package domain

import (
	jwt "github.com/dgrijalva/jwt-go"
)

//Token struct declaration
type Token struct {
	ID       uint64
	Name     string
	Email    string
	Username string
	*jwt.StandardClaims
}
