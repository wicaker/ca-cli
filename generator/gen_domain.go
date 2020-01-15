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

	err := f.Save(dirName + "/domain/errors.go")
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

	err := f.Save(dirName + "/domain/status_code.go")
	if err != nil {
		return err
	}
	return nil
}
