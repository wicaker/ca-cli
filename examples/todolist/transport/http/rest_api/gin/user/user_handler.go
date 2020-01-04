package userhandler

import (
	"context"
	"net/http"
	"todolist/domain"
	"todolist/middleware"

	"github.com/gin-gonic/gin"
)

// UserHandler represent the httphandler for user
type UserHandler struct {
	UserUsecase domain.UserUsecase
}

// NewUserHandler will initialize the user endpoint
func NewUserHandler(r *gin.Engine, u domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: u,
	}
	r.POST("/user/register", handler.Register)
	r.POST("/user/login", handler.Login)
}

// Register will handle register request
func (uh *UserHandler) Register(c *gin.Context) {
	var user domain.User

	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if ok, err := middleware.IsRequestValid(&user); !ok {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = uh.UserUsecase.Register(ctx, &user)

	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.ResponseSuccess{Message: "Successfully register new user"})
}

// Login will handle login request
func (uh *UserHandler) Login(c *gin.Context) {
	var user domain.User
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if ok, err := middleware.IsRequestValid(&user); !ok {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	type data struct {
		Token string `json:"token"`
	}

	d := &data{}
	d.Token, err = uh.UserUsecase.Login(ctx, &user)

	if err != nil {
		c.JSON(401, domain.ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.ResponseSuccess{Message: "Login successfully", Data: d})
}
