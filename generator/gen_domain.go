package generator

import (
	"github.com/dave/jennifer/jen"
)

func (gen *caGen) GenDomainErrors(dirName string) error {
	f := jen.NewFile("domain")
	f.Comment("ResponseError represent the response error struct")
	f.Type().Id("ResponseError").Struct(
		jen.Id("Message").String().Tag(map[string]string{"json": "message"}),
	)

	f.Var().Defs(
		jen.Comment("ErrInternalServerError will throw if any the Internal Server Error happen"),
		jen.Id("ErrInternalServerError").Op("=").Qual("errors", "New").Call(jen.Lit("Internal Server Error")),
		jen.Comment("ErrNotFound will throw if the requested item is not exists"),
		jen.Id("ErrNotFound").Op("=").Qual("errors", "New").Call(jen.Lit("Your requested Item is not found")),
		jen.Comment("ErrConflict will throw if the current action already exists"),
		jen.Id("ErrConflict").Op("=").Qual("errors", "New").Call(jen.Lit("Your Item already exist")),
		jen.Comment("ErrBadParamInput will throw if the given request-body or params is not valid"),
		jen.Id("ErrBadParamInput").Op("=").Qual("errors", "New").Call(jen.Lit("Given Param is not valid")),
		jen.Comment("ErrUnauthorized will throw if the given request-header token is not valid"),
		jen.Id("ErrUnauthorized").Op("=").Qual("errors", "New").Call(jen.Lit("Unauthorized")),
	)

	err := f.Save(dirName + "/errors.go")
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenDomainStatusCode(dirName string) error {
	f := jen.NewFile("domain")
	f.ImportAlias("github.com/sirupsen/logrus", "log")

	f.Comment("GetStatusCode will return status code based on type of error")
	f.Func().Id("GetStatusCode").Params(jen.Id("err").Error()).Int().Block(
		jen.If(jen.Id("err").Op("==").Nil().Block(
			jen.Return(jen.Qual("net/http", "StatusOK")),
		)),
		jen.Qual("github.com/sirupsen/logrus", "Error").Call(jen.Id("err")),
		jen.Switch(jen.Id("err")).Block(
			jen.Case(jen.Id("ErrInternalServerError")).Block(jen.Return(jen.Qual("net/http", "StatusInternalServerError"))),
			jen.Case(jen.Id("ErrNotFound")).Block(jen.Return(jen.Qual("net/http", "StatusNotFound"))),
			jen.Case(jen.Id("ErrConflict")).Block(jen.Return(jen.Qual("net/http", "StatusConflict"))),
			jen.Case(jen.Id("ErrUnauthorized")).Block(jen.Return(jen.Qual("net/http", "StatusUnauthorized"))),
			jen.Default().Block(jen.Return(jen.Qual("net/http", "StatusInternalServerError"))),
		),
	)

	err := f.Save(dirName + "/status_code.go")
	if err != nil {
		return err
	}
	return nil
}

func (gen *caGen) GenDomainSuccess(dirName string) error {
	f := jen.NewFile("domain")
	f.Comment("ResponseSuccess represent the reseponse success struct")
	f.Type().Id("ResponseSuccess").Struct(
		jen.Id("Message").String().Tag(map[string]string{"json": "message"}),
		jen.Id("Data").Interface().Tag(map[string]string{"json": "data"}),
	)

	f.Var().Defs(
		jen.Comment("ResponseData with type map used to response json if no error"),
		jen.Id("ResponseData").Map(jen.String()).Interface(),
	)

	err := f.Save(dirName + "/success.go")
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenDomainExample(dirName string) error {
	f := jen.NewFile("domain")
	f.Comment("Example struct, models of example table")
	f.Type().Id("Example").Struct(
		jen.Id("ID").Uint64().Tag(map[string]string{"json": "id"}),
		jen.Id("Name").String().Tag(map[string]string{"json": "name"}),
		jen.Id("CreatedAt").Qual("time", "Time").Tag(map[string]string{"json": "created_at", "db": "created_at"}),
		jen.Id("UpdatedAt").Qual("time", "Time").Tag(map[string]string{"json": "updated_at", "db": "updated_at"}),
		jen.Id("DeletedAt").Op("*").Qual("time", "Time").Tag(map[string]string{"json": "deleted_at", "db": "deleted_at", "pg": ",soft_delete"}),
	)

	f.Comment("ExampleUsecase represent the Example's usecases contract")
	f.Type().Id("ExampleUsecase").Interface(
		jen.Id("Fetch").Params(jen.Id("ctx").Qual("context", "Context")).Call(jen.Index().Op("*").Id("Example"), jen.Error()),
		jen.Id("GetByID").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Uint64()).Call(jen.Op("*").Id("Example"), jen.Error()),
		jen.Id("Store").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("exp").Op("*").Id("Example")).Call(jen.Op("*").Id("Example"), jen.Error()),
		jen.Id("Update").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("exp").Op("*").Id("Example")).Call(jen.Op("*").Id("Example"), jen.Error()),
		jen.Id("Delete").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Uint64()).Call(jen.Error()),
	)

	f.Comment("ExampleRepository represent the Example's repository contract")
	f.Type().Id("ExampleRepository").Interface(
		jen.Id("Fetch").Params(jen.Id("ctx").Qual("context", "Context")).Call(jen.Index().Op("*").Id("Example"), jen.Error()),
		jen.Id("GetByID").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Uint64()).Call(jen.Op("*").Id("Example"), jen.Error()),
		jen.Id("Store").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("exp").Op("*").Id("Example")).Call(jen.Op("*").Id("Example"), jen.Error()),
		jen.Id("Update").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("exp").Op("*").Id("Example")).Call(jen.Op("*").Id("Example"), jen.Error()),
		jen.Id("Delete").Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("id").Uint64()).Call(jen.Error()),
	)

	err := f.Save(dirName + "/example.go")
	if err != nil {
		return err
	}

	return nil
}
