package generator

import (
	"fmt"

	"github.com/wicaker/cacli/domain"

	"github.com/dave/jennifer/jen"
)

func (gen *caGen) GenUsecase(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	useCase := parser.Usecase.Name
	repository := parser.Repository.Name
	// domainNameUpper := strings.ToUpper(string(domainName[0])) + string(domainName[1:])
	newS := fmt.Sprintf("New%s", useCase)
	comment := fmt.Sprintf("%s will create new an %sUsecase object representation of domain.%s interface ", newS, domainName, useCase)

	f := jen.NewFile("usecase")
	f.Type().Id(domainName+"Usecase").Struct(
		jen.Id(domainName+"Repo").Qual(gomodName+"/domain", repository),
		jen.Id("contextTimeout").Qual("time", "Duration"),
	)

	f.Comment(comment)
	f.Func().Id(newS).Params(
		jen.Id(string(domainName[0])+"r").Qual(gomodName+"/domain", repository),
		jen.Id("timeout").Qual("time", "Duration"),
	).Qual(gomodName+"/domain", useCase).Block(
		jen.Return(jen.Op("&").Id(domainName + "Usecase").Values(jen.Dict{
			jen.Id(domainName + "Repo"): jen.Id(string(domainName[0]) + "r"),
			jen.Id("contextTimeout"):    jen.Id("timeout"),
		})),
	)

	for _, i := range parser.Usecase.Method {
		var param []jen.Code
		var result []jen.Code
		var returnV []jen.Code
		for _, j := range i.ParameterList {
			param = append(param, jen.Id(j.Name).Op(j.Type))
		}
		for _, k := range i.ResultList {
			result = append(result, jen.Id(k.Name).Op(k.Type))
			switch k.Type {
			case "string":
				returnV = append(returnV, jen.Op(`""`))
			case "bool":
				returnV = append(returnV, jen.Op("false"))
			case "float32", "float64", "complex64", "complex128":
				returnV = append(returnV, jen.Op("0.0"))
			case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintpr", "byte", "rune":
				returnV = append(returnV, jen.Op("0"))
			default:
				if len(k.Type) > 7 {
					if k.Type[:7] == "domain." {
						returnV = append(returnV, jen.Op(k.Type+"{}"))
					} else {
						returnV = append(returnV, jen.Nil())
					}
				} else {
					returnV = append(returnV, jen.Nil())
				}
			}

		}
		f.Func().
			Params(jen.Id(string(domainName[0])+"u").Op("*").Id(domainName+"Usecase")).
			Id(i.Name).Params(param[:]...).Call(result[:]...).Block(
			jen.List(jen.Id("ctx"), jen.Id("cancel")).Op(":=").Qual("context", "WithTimeout").Call(jen.Id("ctx"), jen.Qual(string(domainName[0])+"u", "contextTimeout")),
			jen.Id("defer").Id("cancel").Call(),
			jen.Return(returnV[:]...),
		)
	}

	// f.Func().
	// 	Params(jen.Id(string(domainName[0])+"u").Op("*").Id(domainName+"Usecase")).
	// 	Id("Fetch").Params(jen.Id("ctx").Qual("context", "Context")).Call(jen.Index().Op("*").Id("Example"), jen.Error()).Block(
	// 	jen.List(jen.Id("ctx"), jen.Id("cancel")).Op(":=").Qual("context", "WithTimeout").Call(jen.Id("ctx"), jen.Qual(string(domainName[0])+"u", "contextTimeout")),
	// 	jen.Id("defer").Id("cancel").Call(),
	// )
	// func (tu *taskUsecase) Fetch(ctx context.Context, userID uint64) ([]*domain.Task, error) {
	// 	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	// 	defer cancel()

	// 	user, err := tu.userRepo.GetByID(ctx, userID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if user == nil {
	// 		return nil, domain.ErrNotFound
	// 	}

	// 	tasks, err := tu.taskRepo.Fetch(ctx, user.ID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return tasks, nil
	// }

	fileDir := fmt.Sprintf("%s/usecase/%s_usecase.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}
