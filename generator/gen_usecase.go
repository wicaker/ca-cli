package generator

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/wicaker/cacli/domain"

	"github.com/dave/jennifer/jen"
)

func (gen *caGen) GenUsecase(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		useCase    = parser.Usecase.Name
		repository = parser.Repository.Name
		newS       = fmt.Sprintf("New%s", useCase)
		comment    = fmt.Sprintf("%s will create new an %sUsecase object representation of domain.%s interface ", newS, domainName, useCase)
		structRepo jen.Code
		paramRepo  jen.Code
	)

	f := jen.NewFile("usecase")
	f.ImportName(gomodName+"/domain", "domain")

	funcRepo := jen.Dict{
		jen.Id("contextTimeout"): jen.Id("timeout"),
	}

	if repository != "" {
		structRepo = jen.Id(domainName+"Repo").Qual(gomodName+"/domain", repository)
		paramRepo = jen.Id(string(domainName[0])+"r").Qual(gomodName+"/domain", repository)
		funcRepo = jen.Dict{
			jen.Id(domainName + "Repo"): jen.Id(string(domainName[0]) + "r"),
			jen.Id("contextTimeout"):    jen.Id("timeout"),
		}
	}

	f.Type().Id(domainName+"Usecase").Struct(
		structRepo,
		jen.Id("contextTimeout").Qual("time", "Duration"),
	)

	f.Comment(comment)
	f.Func().Id(newS).Params(
		paramRepo,
		jen.Id("timeout").Qual("time", "Duration"),
	).Qual(gomodName+"/domain", useCase).Block(
		jen.Return(jen.Op("&").Id(domainName + "Usecase").Values(funcRepo)),
	)

	for _, i := range parser.Usecase.Method {
		var (
			param            = genParamList(i)
			returnT, returnV = genReturnList(i)
		)
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"u").Op("*").Id(domainName+"Usecase")).
			Id(i.Name).Params(param[:]...).Call(returnT[:]...).Block(
			jen.List(jen.Id("ctx"), jen.Id("cancel")).Op(":=").Qual("context", "WithTimeout").Call(jen.Id("ctx"), jen.Id(string(domainName[0])+"u").Dot("contextTimeout")),
			jen.Id("defer").Id("cancel").Call(),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/%s_usecase.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}
	return nil
}
