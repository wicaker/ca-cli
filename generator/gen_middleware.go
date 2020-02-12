package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
)

func (gen *caGen) GenEchoMiddleware(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	f := jen.NewFile("middleware")
	f.ImportAlias("github.com/sirupsen/logrus", "log")

	f.Comment("EchoMiddleware represent the data-struct for middleware")
	f.Type().Id("EchoMiddleware").Struct()

	f.Comment("InitEchoMiddleware intialize the middleware")
	f.Func().Id("InitEchoMiddleware").Params().Op("*").Id("EchoMiddleware").Block(
		jen.Return(jen.Op("&").Id("EchoMiddleware").Values(jen.Dict{})),
	)

	f.Comment("CORS will handle the CORS middleware")
	f.Func().Params(jen.Id("m").Op("*").Id("EchoMiddleware")).Id("CORS").Params(jen.Id("next").Qual("github.com/labstack/echo", "HandlerFunc")).Qual("github.com/labstack/echo", "HandlerFunc").Block(
		jen.Return(jen.Func().Params(jen.Id("c").Qual("github.com/labstack/echo", "Context")).Error().Block(
			jen.Id("c").Dot("Response").Call().Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Origin"), jen.Lit("*")),
			jen.Return(jen.Id("next").Call(jen.Id("c"))),
		)),
	)

	f.Comment("MiddlewareLogging for logging")
	f.Func().Params(jen.Id("m").Op("*").Id("EchoMiddleware")).Id("MiddlewareLogging").Params(jen.Id("next").Qual("github.com/labstack/echo", "HandlerFunc")).Qual("github.com/labstack/echo", "HandlerFunc").Block(
		jen.Return(jen.Func().Params(jen.Id("c").Qual("github.com/labstack/echo", "Context")).Error().Block(
			jen.Id("makeLogEntry").Call(jen.Id("c")).Dot("Info").Call(jen.Lit("incoming request")),
			jen.Return(jen.Id("next").Call(jen.Id("c"))),
		)),
	)

	f.Func().Id("makeLogEntry").Params(jen.Id("c").Qual("github.com/labstack/echo", "Context")).Op("*").Qual("github.com/sirupsen/logrus", "Entry").Block(
		jen.If(jen.Id("c").Op("==").Nil().Block(
			jen.Return(jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
				jen.Lit("at"): jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
			}))),
		)),
		jen.Return(jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
			jen.Lit("at"):     jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
			jen.Lit("method"): jen.Id("c").Dot("Request").Call().Dot("Method"),
			jen.Lit("uri"):    jen.Id("c").Dot("Request").Call().Dot("URL").Dot("String").Call(),
			jen.Lit("ip"):     jen.Id("c").Dot("Request").Call().Dot("RemoteAddr"),
		}))),
	)

	fileDir := fmt.Sprintf("%s/middleware/echo_middleware.go", dirName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGinMiddleware(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenGorillaMuxMiddleware(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenNetHTTPMiddleware(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}
