package generator

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
)

func (gen *caGen) GenGopgRepository(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		repository = parser.Repository.Name
		newR       = fmt.Sprintf("NewGopg%s", repository)
		comment    = fmt.Sprintf("NewGopg%s will create new an gopg%s object representation of domain.%s interface", repository, repository, repository)
		f          = jen.NewFile("repository")
	)

	f.ImportAlias("github.com/go-pg/pg/v9", "pg")
	f.ImportName(gomodName+"/domain", "domain")

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
		var (
			param            = genParamList(i)
			returnT, returnV = genReturnList(i)
		)
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"r").Op("*").Id("gopg"+repository)).
			Id(i.Name).Params(param[:]...).Call(returnT[:]...).Block(
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/%s_repository.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGormRepository(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		repository = parser.Repository.Name
		newR       = fmt.Sprintf("NewGorm%s", repository)
		comment    = fmt.Sprintf("NewGorm%s will create new an gorm%s object representation of domain.%s interface", repository, repository, repository)
		f          = jen.NewFile("repository")
		importName = map[string]string{
			"github.com/jinzhu/gorm": "gorm",
			gomodName + "/domain":    "domain",
		}
	)

	f.ImportNames(importName)

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
		var (
			param            = genParamList(i)
			returnT, returnV = genReturnList(i)
		)
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"r").Op("*").Id("gorm"+repository)).
			Id(i.Name).Params(param[:]...).Call(returnT[:]...).Block(
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/%s_repository.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenSQLRepository(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		repository = parser.Repository.Name
		newR       = fmt.Sprintf("NewSQL%s", repository)
		comment    = fmt.Sprintf("NewSQL%s will create new an sql%s object representation of domain.%s interface", repository, repository, repository)
		f          = jen.NewFile("repository")
	)

	f.ImportName(gomodName+"/domain", "domain")

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
		var (
			param            = genParamList(i)
			returnT, returnV = genReturnList(i)
		)
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"r").Op("*").Id("sql"+repository)).
			Id(i.Name).Params(param[:]...).Call(returnT[:]...).Block(
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/%s_repository.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenSqlxRepository(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		repository = parser.Repository.Name
		newR       = fmt.Sprintf("NewSqlx%s", repository)
		comment    = fmt.Sprintf("NewSqlx%s will create new an sqlx%s object representation of domain.%s interface", repository, repository, repository)
		f          = jen.NewFile("repository")
		importName = map[string]string{
			"github.com/jmoiron/sqlx": "sqlx",
			gomodName + "/domain":     "domain",
		}
	)

	f.ImportNames(importName)

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
		var (
			param            = genParamList(i)
			returnT, returnV = genReturnList(i)
		)
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"r").Op("*").Id("sqlx"+repository)).
			Id(i.Name).Params(param[:]...).Call(returnT[:]...).Block(
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/%s_repository.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenMongodRepository(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		repository = parser.Repository.Name
		newR       = fmt.Sprintf("NewMongod%s", repository)
		comment    = fmt.Sprintf("NewMongod%s will create new an mongod%s object representation of domain.%s interface", repository, repository, repository)
		f          = jen.NewFile("repository")
		importName = map[string]string{
			"go.mongodb.org/mongo-driver/mongo": "mongo",
			gomodName + "/domain":               "domain",
		}
	)

	f.ImportNames(importName)

	f.Type().Id("mongod" + repository).Struct(
		jen.Id("Conn").Op("*").Qual("go.mongodb.org/mongo-driver/mongo", "Database"),
	)

	f.Comment(comment)
	f.Func().Id(newR).Params(
		jen.Id("Conn").Op("*").Qual("go.mongodb.org/mongo-driver/mongo", "Database"),
	).Qual(gomodName+"/domain", repository).Block(
		jen.Return(jen.Op("&").Id("mongod" + repository).Values(jen.Dict{
			jen.Id("Conn"): jen.Id("Conn"),
		})),
	)

	for _, i := range parser.Repository.Method {
		var (
			param            = genParamList(i)
			returnT, returnV = genReturnList(i)
		)
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"r").Op("*").Id("mongod"+repository)).
			Id(i.Name).Params(param[:]...).Call(returnT[:]...).Block(
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
			jen.Return(returnV[:]...),
		)
	}

	fileDir := fmt.Sprintf("%s/%s_repository.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}
