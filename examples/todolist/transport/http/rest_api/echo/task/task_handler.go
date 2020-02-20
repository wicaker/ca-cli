package taskhandler

import (
	"context"
	"net/http"
	"strconv"
	"todolist/domain"
	"todolist/middleware"

	"github.com/labstack/echo"
)

// TaskHandler represent the httphandler for task
type TaskHandler struct {
	TaskUsecase domain.TaskUsecase
}

// NewTaskHandler will initialize the task endpoint
func NewTaskHandler(e *echo.Echo, u domain.TaskUsecase) {
	handler := &TaskHandler{
		TaskUsecase: u,
	}
	e.GET("/task", handler.FetchTask)
	e.GET("/task/:id", handler.GetByID)
	e.POST("/task", handler.Store)
	e.PUT("/task/:id", handler.Update)
	e.DELETE("/task/:id", handler.Delete)
}

// FetchTask will handle FetchTask request
func (th *TaskHandler) FetchTask(c echo.Context) error {
	tasks := []*domain.Task{}
	tokenHeader := c.Request().Header.Get("x-access-token")
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := th.TaskUsecase.Fetch(ctx, token.ID)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	if len(res) != 0 {
		tasks = res
	}

	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["tasks"] = tasks
	return c.JSON(http.StatusOK, domain.ResponseData)
}

// GetByID will handle GetByID request
func (th *TaskHandler) GetByID(c echo.Context) error {
	// get token
	tokenHeader := c.Request().Header.Get("x-access-token")
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	// get id of task
	idTask, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	task, err := th.TaskUsecase.GetByID(ctx, uint64(idTask))

	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["task"] = task
	return c.JSON(http.StatusOK, domain.ResponseData)
}

// Store will handle Store request
func (th *TaskHandler) Store(c echo.Context) error {
	var task domain.Task

	// get token
	tokenHeader := c.Request().Header.Get("x-access-token")
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	// bind req.body
	err = c.Bind(&task)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ResponseError{Message: err.Error()})
	}

	if ok, err := middleware.IsRequestValid(&task); !ok {
		return c.JSON(http.StatusBadRequest, domain.ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	resp, err := th.TaskUsecase.Store(ctx, token.ID, &task)

	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSuccess{Message: "Successfully create new task", Data: resp})
}

// Update will handle Update request
func (th *TaskHandler) Update(c echo.Context) error {
	var task domain.Task

	// bind req.body
	err := c.Bind(&task)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, domain.ResponseError{Message: err.Error()})
	}

	// get id of task
	idTask, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ResponseError{Message: err.Error()})
	}
	task.ID = uint64(idTask)

	// get token
	tokenHeader := c.Request().Header.Get("x-access-token")
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}
	task.UserID = token.ID

	if ok, err := middleware.IsRequestValid(&task); !ok {
		return c.JSON(http.StatusBadRequest, domain.ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = th.TaskUsecase.Update(ctx, &task)

	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	// return c.JSON(http.StatusOK, task)
	return c.JSON(http.StatusOK, domain.ResponseSuccess{Message: "Successfully update data", Data: task})
}

// Delete will handle Delete request
func (th *TaskHandler) Delete(c echo.Context) error {
	// get token
	tokenHeader := c.Request().Header.Get("x-access-token")
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	// get id of task
	idTask, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = th.TaskUsecase.Delete(ctx, uint64(idTask))

	if err != nil {
		return c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, &domain.ResponseSuccess{Message: "Delete item successfully"})
}
