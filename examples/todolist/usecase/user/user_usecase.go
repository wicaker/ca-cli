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
			if err != nil {
				c <- err
				return
			}
			if checkEmail != nil {
				c <- domain.ErrConflict
				return
			}
		}

		if t.Username != "" {
			checkUsername, err := tu.userRepo.GetByUsername(ctx, t.Username)
			if err != nil {
				c <- err
				return
			}
			if checkUsername != nil {
				c <- domain.ErrConflict
				return
			}
		}

		// Password Encryption
		if ctx.Err() == nil {
			pass, err := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
			if err != nil {
				c <- errors.New("Password Encryption failed")
				return
			}
			t.Password = string(pass)

			// Save new user data
			if ctx.Err() == nil {
				err := tu.userRepo.Register(ctx, t)
				if err != nil {
					c <- err
					return
				}
				c <- nil
				return
			}
		}
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
	return <-c
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
				c <- (func() (string, error) { return "", err })
				return
			}
			if checkEmail == nil {
				c <- (func() (string, error) { return "", domain.ErrNotFound })
				return
			}

			user = checkEmail
		} else if t.Username != "" {
			checkUsername, err := tu.userRepo.GetByUsername(ctx, t.Username)
			if err != nil {
				c <- (func() (string, error) { return "", err })
				return
			}
			if checkUsername == nil {
				c <- (func() (string, error) { return "", domain.ErrNotFound })
				return
			}

			user = checkUsername
		} else {
			c <- (func() (string, error) { return "", errors.New("username or email must provided") })
			return
		}

		// process of creating token
		if ctx.Err() == nil {
			expiresAt := time.Now().Add(time.Minute * 100000).Unix()
			errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(t.Password))
			if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
				c <- (func() (string, error) { return "", errors.New("Invalid login credentials. Please try again") })
				return
			}

			tk := &domain.Token{
				ID:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Username: user.Username,
				StandardClaims: &jwt.StandardClaims{
					ExpiresAt: expiresAt,
				},
			}
			if ctx.Err() == nil {
				token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
				tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
				if err != nil {
					c <- (func() (string, error) { return "", err })
				}

				c <- (func() (string, error) { return tokenString, nil })
				return
			}
		}
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
