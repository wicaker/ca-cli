package generator

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"

	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
	"github.com/wicaker/cacli/parser"

	"github.com/spf13/afero"
)

func getRepository(path string, gomodName string) []jen.Code {
	var repo []jen.Code

	fs := afero.NewOsFs()
	res, err := afero.ReadDir(fs, path+"/repository")
	if err != nil {
		fmt.Println(err)
	}
	for i := range res {
		if filepath.Ext(res[i].Name()) == ".go" {
			p := parser.NewParserRepoService()
			par, err := p.RepoParser(path + "/repository/" + res[i].Name())
			if err != nil {
				fmt.Println(err)
			}
			reg, err := regexp.Compile("[^a-zA-Z0-9]+")
			if err != nil {
				log.Fatal(err)
			}

			fileName := reg.ReplaceAllString(res[i].Name(), "")
			repoName := par.Repository.Method[0].Name
			repo = append(repo, jen.Id(fileName[:len(fileName)-2]).Op(":=").Qual(gomodName+"/repository", repoName).Call(jen.Id("db")))
		}
	}
	// userRepo := _userSqlxRepo.NewUserSqlxRepository(db)
	return repo
}

func getUsecase(path string, gomodName string) []jen.Code {
	var repo []jen.Code

	fs := afero.NewOsFs()
	res, err := afero.ReadDir(fs, path+"/usecase")
	if err != nil {
		fmt.Println(err)
	}
	for i := range res {
		if filepath.Ext(res[i].Name()) == ".go" {
			p := parser.NewParserRepoService()
			par, err := p.RepoParser(path + "/usecase/" + res[i].Name())
			if err != nil {
				fmt.Println(err)
			}

			reg, err := regexp.Compile("[^a-zA-Z0-9]+")
			if err != nil {
				log.Fatal(err)
			}

			fileName := reg.ReplaceAllString(res[i].Name(), "")
			usecaseName := par.Usecase.Method[0].Name
			repoName := fileName[:len(fileName)-9]
			repo = append(repo, jen.Id(fileName[:len(fileName)-2]).Op(":=").Qual(gomodName+"/usecase", usecaseName).Call(jen.Id(repoName+"repository"), jen.Id("timeoutContext")))
		}
	}
	// userRepo := _userSqlxRepo.NewUserSqlxRepository(db)
	return repo
}

func getHandler(path string, gomodName string) []jen.Code {
	var repo []jen.Code

	fs := afero.NewOsFs()
	res, err := afero.ReadDir(fs, path)
	if err != nil {
		fmt.Println(err)
	}
	for i := range res {
		if filepath.Ext(res[i].Name()) == ".go" {
			p := parser.NewParserRepoService()
			par, err := p.RepoParser(path + "/" + res[i].Name())
			if err != nil {
				fmt.Println(err)
			}

			reg, err := regexp.Compile("[^a-zA-Z0-9]+")
			if err != nil {
				log.Fatal(err)
			}

			fileName := reg.ReplaceAllString(res[i].Name(), "")
			handlerName := par.Handler.Method[0].Name
			usecaseName := fileName[:len(fileName)-9]
			repo = append(repo, jen.Qual(path, handlerName).Call(jen.Id("e"), jen.Id(usecaseName+"usecase")))
		}
	}
	return repo
}

func (gen *caGen) GenEchoServer(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	f := jen.NewFile("server")
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")

	f.Comment("EchoServer server")
	f.Func().Id("EchoServer").Params(jen.Id("db").Op("*").Qual("github.com/go-pg/pg/v9", "DB")).Op("*").Qual("github.com/labstack/echo", "Echo").Block(
		jen.Id("e").Op(":=").Qual("github.com/labstack/echo", "New").Call(),
		jen.Id("middl").Op(":=").Qual(gomodName+"/middleware", "InitEchoMiddleware").Call(),
		jen.Id("e").Dot("Use").Call(jen.Id("middl").Dot("MiddlewareLogging")),
		jen.Id("e").Dot("Use").Call(jen.Id("middl").Dot("CORS")),
		jen.Id("timeoutContext").Op(":=").Qual("time", "Duration").Call(jen.Lit(2).Call(jen.Lit("context.timeout"))).Op("*").Qual("time", "Second"),
		getRepository(dirName, gomodName)[0],
		getUsecase(dirName, gomodName)[0],
		getHandler(dirName+"/transport/rest", gomodName)[0],

		// read repository folder
		// read usecase folder
		// read echo folder
		// f.Qual("fmt", "Println").Call(f.Id("db")),

		jen.Return(jen.Id("e")),
	)

	fileDir := fmt.Sprintf("%s/server/echo_server.go", dirName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGinServer(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenGorillaMuxServer(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenNetHTTPServer(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenGraphqlServer(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}

func (gen *caGen) GenGrpcServer(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	return nil
}
