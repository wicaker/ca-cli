package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
)

func (gen *caGen) GenGopgRepository(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	repository := parser.Repository.Name

	newR := fmt.Sprintf("NewGopg%s", repository)

	comment := fmt.Sprintf("NewGopg%s will create new an gopg%s object representation of domain.%s interface", repository, repository, repository)
	f := jen.NewFile("repository")
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")

	f.Type().Id("gopg" + repository).Struct(
		jen.Id("Conn").Op("*").Qual("github.com/go-pg/pg/v9", "DB"),
	)

	f.Comment(comment)
	f.Func().Id(newR).Params(
		jen.Id("Conn").Op("*").Qual("github.com/go-pg/pg/v9", "DB"),
	).Qual(gomodName+"/domain", repository).Block(
		jen.Return(jen.Op("&").Id("gopg" + repository).Values(jen.Dict{
			jen.Id("Conn"): jen.Id("Conn"),
		})),
	)

	for _, i := range parser.Repository.Method {
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
			Params(jen.Id(string(domainName[0]) + "r").Op("*").Id("gopg" + repository)).
			Id(i.Name).Params(param[:]...).Call(result[:]...).Block(

			// jen.List(jen.Id("ctx"), jen.Id("cancel")).Op(":=").Qual("context", "WithTimeout").Call(jen.Id("ctx"), jen.Qual(string(domainName[0])+"u", "contextTimeout")),
			// jen.Id("defer").Id("cancel").Call(),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/repository/%s_repository.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGormRepository(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	repository := parser.Repository.Name

	newR := fmt.Sprintf("NewGorm%s", repository)

	comment := fmt.Sprintf("NewGorm%s will create new an gorm%s object representation of domain.%s interface", repository, repository, repository)
	f := jen.NewFile("repository")

	f.Type().Id("gorm" + repository).Struct(
		jen.Id("Conn").Op("*").Qual("github.com/jinzhu/gorm", "DB"),
	)

	f.Comment(comment)
	f.Func().Id(newR).Params(
		jen.Id("Conn").Op("*").Qual("github.com/jinzhu/gorm", "DB"),
	).Qual(gomodName+"/domain", repository).Block(
		jen.Return(jen.Op("&").Id("gorm" + repository).Values(jen.Dict{
			jen.Id("Conn"): jen.Id("Conn"),
		})),
	)

	for _, i := range parser.Repository.Method {
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
			Params(jen.Id(string(domainName[0]) + "r").Op("*").Id("gorm" + repository)).
			Id(i.Name).Params(param[:]...).Call(result[:]...).Block(
			// jen.List(jen.Id("ctx"), jen.Id("cancel")).Op(":=").Qual("context", "WithTimeout").Call(jen.Id("ctx"), jen.Qual(string(domainName[0])+"u", "contextTimeout")),
			// jen.Id("defer").Id("cancel").Call(),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/repository/%s_repository.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenSQLRepository(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	repository := parser.Repository.Name

	newR := fmt.Sprintf("NewSQL%s", repository)

	comment := fmt.Sprintf("NewSQL%s will create new an sql%s object representation of domain.%s interface", repository, repository, repository)
	f := jen.NewFile("repository")

	f.Type().Id("sql" + repository).Struct(
		jen.Id("Conn").Op("*").Qual("database/sql", "DB"),
	)

	f.Comment(comment)
	f.Func().Id(newR).Params(
		jen.Id("Conn").Op("*").Qual("database/sql", "DB"),
	).Qual(gomodName+"/domain", repository).Block(
		jen.Return(jen.Op("&").Id("sql" + repository).Values(jen.Dict{
			jen.Id("Conn"): jen.Id("Conn"),
		})),
	)

	for _, i := range parser.Repository.Method {
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
			Params(jen.Id(string(domainName[0]) + "r").Op("*").Id("sql" + repository)).
			Id(i.Name).Params(param[:]...).Call(result[:]...).Block(
			// jen.List(jen.Id("ctx"), jen.Id("cancel")).Op(":=").Qual("context", "WithTimeout").Call(jen.Id("ctx"), jen.Qual(string(domainName[0])+"u", "contextTimeout")),
			// jen.Id("defer").Id("cancel").Call(),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/repository/%s_repository.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenSqlxRepository(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	repository := parser.Repository.Name

	newR := fmt.Sprintf("NewSqlx%s", repository)

	comment := fmt.Sprintf("NewSqlx%s will create new an sqlx%s object representation of domain.%s interface", repository, repository, repository)
	f := jen.NewFile("repository")

	f.Type().Id("sqlx" + repository).Struct(
		jen.Id("Conn").Op("*").Qual("github.com/jmoiron/sqlx", "DB"),
	)

	f.Comment(comment)
	f.Func().Id(newR).Params(
		jen.Id("Conn").Op("*").Qual("github.com/jmoiron/sqlx", "DB"),
	).Qual(gomodName+"/domain", repository).Block(
		jen.Return(jen.Op("&").Id("sqlx" + repository).Values(jen.Dict{
			jen.Id("Conn"): jen.Id("Conn"),
		})),
	)

	for _, i := range parser.Repository.Method {
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
			Params(jen.Id(string(domainName[0]) + "r").Op("*").Id("sqlx" + repository)).
			Id(i.Name).Params(param[:]...).Call(result[:]...).Block(
			// jen.List(jen.Id("ctx"), jen.Id("cancel")).Op(":=").Qual("context", "WithTimeout").Call(jen.Id("ctx"), jen.Qual(string(domainName[0])+"u", "contextTimeout")),
			// jen.Id("defer").Id("cancel").Call(),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/repository/%s_repository.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}
