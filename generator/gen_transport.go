package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
)

func (gen *caGen) GenEchoTransport(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	useCase := parser.Usecase.Name

	f := jen.NewFile("rest")
	f.Type().Id(domainName + "Handler").Struct(
		jen.Id(useCase).Qual(gomodName+"/domain", useCase),
	)

	fileDir := fmt.Sprintf("%s/transport/rest/%s_handler.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGinTransport(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenGorillaMuxTransport(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenNetHTTPTransport(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenGraphqlTransport(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenGrpcTransport(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

// // TaskHandler represent the httphandler for task
// type TaskHandler struct {
// 	TaskUsecase domain.TaskUsecase
// }

// // NewTaskHandler will initialize the task endpoint
// func NewTaskHandler(e *echo.Echo, u domain.TaskUsecase) {
// 	handler := &TaskHandler{
// 		TaskUsecase: u,
// 	}
// 	e.GET("/task", handler.FetchTask)
// 	e.GET("/task/:id", handler.GetByID)
// 	e.POST("/task", handler.Store)
// 	e.PUT("/task/:id", handler.Update)
// 	e.DELETE("/task/:id", handler.Delete)
// }

// // FetchTask will handle FetchTask request
// func (th *TaskHandler) FetchTask(c echo.Context) error {
