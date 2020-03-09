package rest

import (
	"context"
	echo "github.com/labstack/echo"
	domain "github.com/wicaker/tests/domain"
	"net/http"
)

type exampleHandler struct {
	ExampleUsecase domain.ExampleUsecase
}

// NewExampleHandler will initialize the example endpoint
func NewExampleHandler(e *echo.Echo, u domain.ExampleUsecase) {
	handler := &exampleHandler{ExampleUsecase: u}
	e.GET("/example/fetch", handler.FetchHandler)
	e.GET("/example/getbyid", handler.GetByIDHandler)
	e.GET("/example/store", handler.StoreHandler)
	e.GET("/example/update", handler.UpdateHandler)
	e.GET("/example/delete", handler.DeleteHandler)
}
func (eh *exampleHandler) FetchHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	return c.JSON(http.StatusOK, domain.ResponseData)
}
func (eh *exampleHandler) GetByIDHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	return c.JSON(http.StatusOK, domain.ResponseData)
}
func (eh *exampleHandler) StoreHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	return c.JSON(http.StatusOK, domain.ResponseData)
}
func (eh *exampleHandler) UpdateHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	return c.JSON(http.StatusOK, domain.ResponseData)
}
func (eh *exampleHandler) DeleteHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	return c.JSON(http.StatusOK, domain.ResponseData)
}
