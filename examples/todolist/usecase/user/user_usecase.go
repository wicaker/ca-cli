package user

import (
	"context"
	"errors"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"todolist/domain"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an userUsecase object representation of domain.UserUsecase interface
func NewUserUsecase(ur domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       ur,
		contextTimeout: timeout,
	}
}

func (tu *userUsecase) Register(ctx context.Context, t *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	if t.Email == "" && t.Username == "" {
		return errors.New("Email and Username was required")
	}
	if t.Password == "" {
		return errors.New("Password was required")
	}

	checkEmail, err := tu.userRepo.GetByEmail(ctx, t.Email)
	if err != nil {
		return err
	}
	if checkEmail.Email != "" {
		return errors.New("Email already in database, use another name")
	}

	checkUsername, err := tu.userRepo.GetByUsername(ctx, t.Username)
	if err != nil {
		return err
	}
	if checkUsername.Username != "" {
		return errors.New("Username already in database, use another name")
	}

	// Password Encryption
	pass, err := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Password Encryption failed")
	}
	t.Password = string(pass)

	// Save new user data
	err = tu.userRepo.Register(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

func (tu *userUsecase) Login(ctx context.Context, t *domain.User) (string, error) {
	var user *domain.User

	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	checkEmail, err := tu.userRepo.GetByEmail(ctx, t.Email)
	if err != nil {
		return "", err
	}
	if checkEmail != nil {
		user = checkEmail
	}

	checkUsername, err := tu.userRepo.GetByUsername(ctx, t.Username)
	if err != nil {
		return "", err
	}
	if checkUsername != nil {
		user = checkUsername
	}

	//Token struct declaration
	type Token struct {
		ID       uint64
		Name     string
		Email    string
		Username string
		*jwt.StandardClaims
	}

	// process of creating token
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	errf := bcrypt.CompareHashAndPassword([]byte(t.Password), []byte(user.Password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return "", errors.New("Invalid login credentials. Please try again")
	}

	tk := &Token{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
