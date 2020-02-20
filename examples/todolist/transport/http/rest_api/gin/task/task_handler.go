package taskhandler

import (
	"context"
	"net/http"
	"strconv"
	"todolist/domain"
	"todolist/middleware"

	"github.com/gin-gonic/gin"
)

// TaskHandler represent the httphandler for task
type TaskHandler struct {
	TaskUsecase domain.TaskUsecase
}

// NewTaskHandler will initialize the task endpoint
func NewTaskHandler(r *gin.Engine, u domain.TaskUsecase) {
	handler := &TaskHandler{
		TaskUsecase: u,
	}
	r.GET("/task", handler.FetchTask)
	r.GET("/task/:id", handler.GetByID)
	r.POST("/task", handler.Store)
	r.PUT("/task/:id", handler.Update)
	r.DELETE("/task/:id", handler.Delete)
}

// FetchTask will handle FetchTask request
func (th *TaskHandler) FetchTask(c *gin.Context) {
	tasks := []*domain.Task{}
	domain.ResponseData = make(map[string]interface{})
	tokenHeader := c.Request.Header.Get("x-access-token")
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := th.TaskUsecase.Fetch(ctx, token.ID)
	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	if len(res) != 0 {
		tasks = res
	}
	domain.ResponseData["tasks"] = tasks
	c.JSON(http.StatusOK, domain.ResponseData)
	return
}

// GetByID will handle GetByID request
func (th *TaskHandler) GetByID(c *gin.Context) {
	// get token
	tokenHeader := c.Request.Header.Get("x-access-token")
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	// get id of task
	idTask, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	task, err := th.TaskUsecase.GetByID(ctx, uint64(idTask))

	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.ResponseSuccess{Message: "Successfully load data", Data: task})
	return
}

// Store will handle Store request
func (th *TaskHandler) Store(c *gin.Context) {
	var task domain.Task

	// get token
	tokenHeader := c.Request.Header.Get("x-access-token")
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	// bind req.body
	err = c.Bind(&task)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, domain.ResponseError{Message: err.Error()})
		return
	}

	if ok, err := middleware.IsRequestValid(&task); !ok {
		c.JSON(http.StatusBadRequest, domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	resp, err := th.TaskUsecase.Store(ctx, token.ID, &task)

	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.ResponseSuccess{Message: "Successfully create new task", Data: resp})
	return

}

// Update will handle Update request
func (th *TaskHandler) Update(c *gin.Context) {
	var task domain.Task

	// bind req.body
	err := c.Bind(&task)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, domain.ResponseError{Message: err.Error()})
		return
	}

	// get id of task
	idTask, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ResponseError{Message: err.Error()})
		return
	}
	task.ID = uint64(idTask)

	// get token
	tokenHeader := c.Request.Header.Get("x-access-token")
	token, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}
	task.UserID = token.ID

	if ok, err := middleware.IsRequestValid(&task); !ok {
		c.JSON(http.StatusBadRequest, domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = th.TaskUsecase.Update(ctx, &task)

	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.ResponseSuccess{Message: "Successfully update data", Data: task})
	return
}

// Delete will handle Delete request
func (th *TaskHandler) Delete(c *gin.Context) {
	// get token
	tokenHeader := c.Request.Header.Get("x-access-token")
	_, err := middleware.JwtVerify(tokenHeader)
	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	// get id of task
	idTask, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ResponseError{Message: err.Error()})
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = th.TaskUsecase.Delete(ctx, uint64(idTask))

	if err != nil {
		c.JSON(domain.GetStatusCode(err), domain.ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &domain.ResponseSuccess{Message: "Delete item successfully"})
	return
}
