package userhandler

import (
	"context"
	"net/http"
	"todolist/domain"
	"todolist/middleware"

	"github.com/labstack/echo"
)

// UserHandler represent the httphandler for user
type UserHandler struct {
	UserUsecase domain.UserUsecase
}

// NewUserHandler will initialize the user endpoint
func NewUserHandler(e *echo.Echo, u domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: u,
	}
	e.POST("/user/register", handler.Register)
	e.POST("/user/login", handler.Login)
}

// Register will handle register request
func (uh *UserHandler) Register(c echo.Context) error {
	var user domain.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ResponseError{Message: err.Error()})
	}

	if ok, err := middleware.IsRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, domain.ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = uh.UserUsecase.Register(ctx, &user)

	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSuccess{Message: "Successfully register new user"})
}

// Login will handle login request
func (uh *UserHandler) Login(c echo.Context) error {
	var user domain.User

	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ResponseError{Message: err.Error()})
	}

	if ok, err := middleware.IsRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, domain.ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	type data struct {
		Token string `json:"token"`
	}

	d := &data{}
	d.Token, err = uh.UserUsecase.Login(ctx, &user)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSuccess{Message: "Login successfully", Data: d})
}
