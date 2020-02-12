package generator

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
)

func (gen *caGen) GenEchoTransport(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	useCase := parser.Usecase.Name
	newS := fmt.Sprintf("New%sHandler", strings.ToUpper(string(domainName[0]))+domainName[1:])
	comment := fmt.Sprintf("%s will initialize the %s endpoint", newS, domainName)
	var handler []jen.Code
	handler = append(handler, jen.Id("handler").Op(":=").Op("&").Id(domainName+"Handler").Values(jen.Dict{
		jen.Id(useCase): jen.Id("u"),
	}))

	for _, i := range parser.Usecase.Method {
		path := fmt.Sprintf("/%s/%s", domainName, strings.ToLower(i.Name))
		handler = append(handler, jen.Id("e").Dot("GET").Call(jen.Lit(path), jen.Id("handler").Dot(i.Name+"Handler")))
	}

	f := jen.NewFile("rest")
	f.Type().Id(domainName + "Handler").Struct(
		jen.Id(useCase).Qual(gomodName+"/domain", useCase),
	)

	f.Comment(comment)
	f.Func().Id(newS).Params(
		jen.Id("e").Op("*").Qual("github.com/labstack/echo", "Echo"),
		jen.Id("u").Qual(gomodName+"/domain", useCase),
	).Block(handler[:]...)

	for _, i := range parser.Usecase.Method {
		f.Func().
			Params(jen.Id(string(domainName[0])+"h").Op("*").Id(domainName+"Handler")).
			Id(i.Name+"Handler").Params(jen.Id("c").Qual("github.com/labstack/echo", "Context")).Call(jen.Error()).Block(
			jen.Id("ctx").Op(":=").Id("c").Dot("Request").Call().Dot("Context").Call(),
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),

			jen.Qual(gomodName+"/domain", "ResponseData").Op("=").Make(jen.Map(jen.String()).Interface()),
			jen.Return(jen.Id("c").Dot("JSON").Call(jen.Qual("net/http", "StatusOK"), jen.Qual(gomodName+"/domain", "ResponseData"))),
		)
	}

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
