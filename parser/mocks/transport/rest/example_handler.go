package rest

import (
	"context"
	"net/http"

	echo "github.com/labstack/echo"
	"github.com/wicaker/cacli/parser/mocks/domain"
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

	return c.JSON(http.StatusOK, "response")
}
func (eh *exampleHandler) GetByIDHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	return c.JSON(http.StatusOK, "response")
}
func (eh *exampleHandler) StoreHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	return c.JSON(http.StatusOK, "response")
}
func (eh *exampleHandler) UpdateHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	return c.JSON(http.StatusOK, "response")
}
func (eh *exampleHandler) DeleteHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	return c.JSON(http.StatusOK, "response")
}
