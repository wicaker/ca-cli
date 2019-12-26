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
	c := make(chan error)
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()
	go func() {
		if t.Email == "" && t.Username == "" {
			c <- errors.New("Email or Username was required")
			return
		}
		if t.Password == "" {
			c <- errors.New("Password was required")
			return
		}

		if t.Email != "" {
			checkEmail, err := tu.userRepo.GetByEmail(ctx, t.Email)
			if err == nil {
				c <- errors.New("Email already in database, use another name")
				return
			}
			if checkEmail != nil {
				c <- errors.New("Email already in database, use another name")
				return
			}
		}

		if t.Username != "" {
			checkUsername, err := tu.userRepo.GetByUsername(ctx, t.Username)
			if err == nil {
				c <- errors.New("Username already in database, use another name")
				return
			}
			if checkUsername != nil {
				c <- errors.New("Username already in database, use another name")
				return
			}
		}

		// Password Encryption
		pass, err := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
		if err != nil {
			c <- errors.New("Password Encryption failed")
			return
		}
		t.Password = string(pass)

		// Save new user data
		err = tu.userRepo.Register(ctx, t)
		if err != nil {
			c <- err
			return
		}

		c <- nil
		return
	}()

	// stop run service when timeout
	go func() {
		select {
		case <-ctx.Done():
			c <- ctx.Err()
			return
		}
	}()

	// catch return value from go channel
	y := <-c
	return y
}

func (tu *userUsecase) Login(ctx context.Context, t *domain.User) (string, error) {
	c := make(chan func() (string, error))
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	go func() {
		var user *domain.User

		// check Email if user login using email
		if t.Email != "" {
			checkEmail, err := tu.userRepo.GetByEmail(ctx, t.Email)
			if err != nil {
				c <- (func() (string, error) { return "", errors.New("user not found") })
				return
			}
			if checkEmail == nil {
				c <- (func() (string, error) { return "", errors.New("user not found") })
				return
			}
			user = checkEmail
		} else if t.Username != "" {
			checkUsername, err := tu.userRepo.GetByUsername(ctx, t.Username)

			if err != nil {
				c <- (func() (string, error) { return "", errors.New("user not found") })
				return
			}
			if checkUsername == nil {
				c <- (func() (string, error) { return "", errors.New("user not found") })
				return
			}

			user = checkUsername
		} else {
			c <- (func() (string, error) { return "", errors.New("username or email must provided") })
			return
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
		errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(t.Password))
		if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
			c <- (func() (string, error) { return "", errors.New("Invalid login credentials. Please try again") })
			return
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
			c <- (func() (string, error) { return "", err })
		}

		c <- (func() (string, error) { return tokenString, nil })
		return
	}()

	// stop run service when timeout
	go func() {
		select {
		case <-ctx.Done():
			c <- (func() (string, error) { return "", ctx.Err() })
			return
		}
	}()

	// catch return value from go channel
	y, z := (<-c)()
	return y, z
}
