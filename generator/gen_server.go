package generator

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/parser"
)

type genServer struct{}

func (gen *caGen) GenEchoServer(dirName string, serviceName string, repoLib string, gomodName string, parser *domain.Parser) error {
	var (
		genServer  genServer
		f          = jen.NewFile("server")
		genCode    []jen.Code
		importName = map[string]string{
			gomodName + "/middleware":           "middleware",
			gomodName + "/repository":           "repository",
			gomodName + "/usecase":              "usecase",
			gomodName + "/transport/rest":       "rest",
			"github.com/labstack/echo":          "echo",
			"go.mongodb.org/mongo-driver/mongo": "mongo",
			"github.com/jinzhu/gorm":            "gorm",
		}
	)

	usecaseFile, repoFile, handlerFile, err := genServer.getAllLayer(serviceName, gomodName, "rest")
	if err != nil {
		return err
	}

	libRepo, err := genServer.checkRepoLib(repoLib)
	if err != nil {
		return err
	}

	genCode = append(genCode, jen.Id("r").Op(":=").Qual("github.com/labstack/echo", "New").Call())
	genCode = append(genCode, jen.Id("middl").Op(":=").Qual(gomodName+"/middleware", "InitEchoMiddleware").Call())
	genCode = append(genCode, jen.Id("r").Dot("Use").Call(jen.Id("middl").Dot("MiddlewareLogging")))
	genCode = append(genCode, jen.Id("r").Dot("Use").Call(jen.Id("middl").Dot("CORS")))
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Id("timeoutContext").Op(":=").Qual("time", "Duration").Call(jen.Lit(2)).Op("*").Qual("time", "Second"))
	genCode = append(genCode, jen.Line())
	for _, i := range repoFile {
		genCode = append(genCode, i)
	}
	for _, i := range usecaseFile {
		genCode = append(genCode, i)
	}
	for _, i := range handlerFile {
		genCode = append(genCode, i)
	}
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Return(jen.Id("r")))

	f.ImportNames(importName)
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")
	f.Comment("EchoServer /")
	f.Func().Id("EchoServer").Params(libRepo).Op("*").Qual("github.com/labstack/echo", "Echo").Block(
		genCode[:]...,
	)

	fileDir := fmt.Sprintf("%s/echo_server.go", dirName)
	err = f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGinServer(dirName string, serviceName string, repoLib string, gomodName string, parser *domain.Parser) error {
	var (
		genServer  genServer
		f          = jen.NewFile("server")
		genCode    []jen.Code
		importName = map[string]string{
			gomodName + "/middleware":           "middleware",
			gomodName + "/repository":           "repository",
			gomodName + "/usecase":              "usecase",
			gomodName + "/transport/rest":       "rest",
			"github.com/gin-gonic/gin":          "gin",
			"go.mongodb.org/mongo-driver/mongo": "mongo",
			"github.com/jinzhu/gorm":            "gorm",
		}
	)

	usecaseFile, repoFile, handlerFile, err := genServer.getAllLayer(serviceName, gomodName, "rest")
	if err != nil {
		return err
	}

	libRepo, err := genServer.checkRepoLib(repoLib)
	if err != nil {
		return err
	}

	genCode = append(genCode, jen.If(jen.Qual("os", "Getenv").Call(jen.Lit("GIN_MODE")).Op("==").Lit("release").Block(
		jen.Qual("github.com/gin-gonic/gin", "SetMode").Call(jen.Qual("github.com/gin-gonic/gin", "ReleaseMode")),
	)))
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Id("r").Op(":=").Qual("github.com/gin-gonic/gin", "New").Call())
	genCode = append(genCode, jen.Id("middl").Op(":=").Qual(gomodName+"/middleware", "InitGinMiddleware").Call())
	genCode = append(genCode, jen.Id("r").Dot("Use").Call(jen.Id("middl").Dot("MiddlewareLogging").Call()))
	genCode = append(genCode, jen.Id("r").Dot("Use").Call(jen.Id("middl").Dot("CORS").Call()))
	genCode = append(genCode, jen.Id("r").Dot("Use").Call(jen.Qual("github.com/gin-gonic/gin", "Recovery").Call()))
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Id("timeoutContext").Op(":=").Qual("time", "Duration").Call(jen.Lit(2)).Op("*").Qual("time", "Second"))
	genCode = append(genCode, jen.Line())
	for _, i := range repoFile {
		genCode = append(genCode, i)
	}
	for _, i := range usecaseFile {
		genCode = append(genCode, i)
	}
	for _, i := range handlerFile {
		genCode = append(genCode, i)
	}
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Return(jen.Id("r")))

	f.ImportNames(importName)
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")
	f.Comment("GinServer /")
	f.Func().Id("GinServer").Params(libRepo).Op("*").Qual("github.com/gin-gonic/gin", "Engine").Block(
		genCode[:]...,
	)

	fileDir := fmt.Sprintf("%s/gin_server.go", dirName)
	err = f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGorillaMuxServer(dirName string, serviceName string, repoLib string, gomodName string, parser *domain.Parser) error {
	var (
		genServer  genServer
		f          = jen.NewFile("server")
		genCode    []jen.Code
		importName = map[string]string{
			gomodName + "/middleware":           "middleware",
			gomodName + "/repository":           "repository",
			gomodName + "/usecase":              "usecase",
			gomodName + "/transport/rest":       "rest",
			"github.com/gorilla/mux":            "mux",
			"go.mongodb.org/mongo-driver/mongo": "mongo",
			"github.com/jinzhu/gorm":            "gorm",
		}
	)

	usecaseFile, repoFile, handlerFile, err := genServer.getAllLayer(serviceName, gomodName, "rest")
	if err != nil {
		return err
	}

	libRepo, err := genServer.checkRepoLib(repoLib)
	if err != nil {
		return err
	}

	genCode = append(genCode, jen.Id("r").Op(":=").Qual("github.com/gorilla/mux", "NewRouter").Call())
	genCode = append(genCode, jen.Id("middl").Op(":=").Qual(gomodName+"/middleware", "InitGorillaMuxMiddleware").Call())
	genCode = append(genCode, jen.Id("r").Dot("Use").Call(jen.Qual("github.com/gorilla/mux", "CORSMethodMiddleware").Call(jen.Id("r"))))
	genCode = append(genCode, jen.Id("r").Dot("Use").Call(jen.Id("middl").Dot("MiddlewareLogging")))
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Id("timeoutContext").Op(":=").Qual("time", "Duration").Call(jen.Lit(2)).Op("*").Qual("time", "Second"))
	genCode = append(genCode, jen.Line())
	for _, i := range repoFile {
		genCode = append(genCode, i)
	}
	for _, i := range usecaseFile {
		genCode = append(genCode, i)
	}
	for _, i := range handlerFile {
		genCode = append(genCode, i)
	}
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Return(jen.Id("r")))

	f.ImportNames(importName)
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")
	f.Comment("GorillaMuxServer /")
	f.Func().Id("GorillaMuxServer").Params(libRepo).Op("*").Qual("github.com/gorilla/mux", "Router").Block(
		genCode[:]...,
	)

	fileDir := fmt.Sprintf("%s/gorilla_mux_server.go", dirName)
	err = f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenNetHTTPMuxServer(dirName string, serviceName string, repoLib string, gomodName string, parser *domain.Parser) error {
	var (
		genServer  genServer
		f          = jen.NewFile("server")
		genCode    []jen.Code
		importName = map[string]string{
			gomodName + "/middleware":           "middleware",
			gomodName + "/repository":           "repository",
			gomodName + "/usecase":              "usecase",
			gomodName + "/transport/rest":       "rest",
			"go.mongodb.org/mongo-driver/mongo": "mongo",
			"github.com/jinzhu/gorm":            "gorm",
		}
	)

	usecaseFile, repoFile, handlerFile, err := genServer.getAllLayer(serviceName, gomodName, "rest")
	if err != nil {
		return err
	}

	libRepo, err := genServer.checkRepoLib(repoLib)
	if err != nil {
		return err
	}

	genCode = append(genCode, jen.Id("r").Op(":=").Qual("net/http", "NewServeMux").Call())
	genCode = append(genCode, jen.Id("middl").Op(":=").Qual(gomodName+"/middleware", "InitNetHTTPMiddleware").Call())
	genCode = append(genCode, jen.Var().Id("handler").Qual("net/http", "Handler").Op("=").Id("r"))
	genCode = append(genCode, jen.Id("handler").Op("=").Id("middl").Dot("MiddlewareLogging").Call(jen.Id("handler")))
	genCode = append(genCode, jen.Id("handler").Op("=").Id("middl").Dot("CORS").Call(jen.Id("handler")))
	genCode = append(genCode, jen.Id("handler").Op("=").Id("parseURL").Call(jen.Id("handler")))
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Id("timeoutContext").Op(":=").Qual("time", "Duration").Call(jen.Lit(2)).Op("*").Qual("time", "Second"))
	genCode = append(genCode, jen.Line())
	for _, i := range repoFile {
		genCode = append(genCode, i)
	}
	for _, i := range usecaseFile {
		genCode = append(genCode, i)
	}
	for _, i := range handlerFile {
		genCode = append(genCode, i)
	}
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Return(jen.Id("handler")))

	f.ImportNames(importName)
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")
	f.Comment("MuxServer /")
	f.Func().Id("MuxServer").Params(libRepo).Qual("net/http", "Handler").Block(
		genCode[:]...,
	)
	f.Line()
	f.Func().Id("parseURL").Params(jen.Id("next").Qual("net/http", "Handler")).Qual("net/http", "Handler").Block(
		jen.Return(
			jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(
				jen.Id("w").Qual("net/http", "ResponseWriter"),
				jen.Id("r").Op("*").Qual("net/http", "Request"),
			).Block(
				jen.Id("ctx").Op(":=").Id("r").Dot("Context").Call(),
				jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
				jen.Id("ctx").Op("=").Qual("context", "WithValue").Call(jen.Id("ctx"), jen.Lit("uri"), jen.Id("r").Dot("RequestURI")),
				jen.Line(),
				jen.Id("s").Op(":=").Qual("strings", "Split").Call(jen.Id("r").Dot("RequestURI"), jen.Lit("/")),
				jen.If(jen.Id("s").Op("[1]==").Lit("example")).Block(
					jen.Id("r").Dot("URL").Op("=&").Qual("net/url", "URL").Values(jen.Dict{jen.Id("Path"): jen.Lit("/example")}),
				),
				jen.Line(),
				jen.Id("next").Dot("ServeHTTP").Call(jen.Id("w"), jen.Id("r").Dot("WithContext").Call(jen.Id("ctx"))),
			)),
		),
	)

	fileDir := fmt.Sprintf("%s/net_http_mux_server.go", dirName)
	err = f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGraphqlServer(dirName string, serviceName string, repoLib string, gomodName string, parser *domain.Parser) error {
	var (
		genServer  genServer
		f          = jen.NewFile("server")
		genCode    []jen.Code
		importName = map[string]string{
			gomodName + "/middleware":           "middleware",
			gomodName + "/repository":           "repository",
			gomodName + "/usecase":              "usecase",
			"go.mongodb.org/mongo-driver/mongo": "mongo",
			"github.com/jinzhu/gorm":            "gorm",
		}
	)

	usecaseFile, repoFile, handlerFile, err := genServer.getAllLayer(serviceName, gomodName, "graphql")
	if err != nil {
		return err
	}

	libRepo, err := genServer.checkRepoLib(repoLib)
	if err != nil {
		return err
	}

	genCode = append(genCode, jen.Id("r").Op(":=").Qual("net/http", "NewServeMux").Call())
	genCode = append(genCode, jen.Id("middl").Op(":=").Qual(gomodName+"/middleware", "InitNetHTTPMiddleware").Call())
	genCode = append(genCode, jen.Var().Id("handler").Qual("net/http", "Handler").Op("=").Id("r"))
	genCode = append(genCode, jen.Id("handler").Op("=").Id("middl").Dot("MiddlewareLogging").Call(jen.Id("handler")))
	genCode = append(genCode, jen.Id("handler").Op("=").Id("middl").Dot("CORS").Call(jen.Id("handler")))
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Id("timeoutContext").Op(":=").Qual("time", "Duration").Call(jen.Lit(2)).Op("*").Qual("time", "Second"))
	genCode = append(genCode, jen.Line())
	for _, i := range repoFile {
		genCode = append(genCode, i)
	}
	for _, i := range usecaseFile {
		genCode = append(genCode, i)
	}
	for _, i := range handlerFile {
		genCode = append(genCode, i)
	}
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Return(jen.Id("handler")))

	f.ImportNames(importName)
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")
	f.ImportAlias(gomodName+"/transport/graphql", "graphqlhandler")
	f.Comment("GraphQLServer /")
	f.Func().Id("GraphQLServer").Params(libRepo).Qual("net/http", "Handler").Block(
		genCode[:]...,
	)

	fileDir := fmt.Sprintf("%s/graphql_server.go", dirName)
	err = f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGrpcServer(dirName string, serviceName string, repoLib string, gomodName string, parser *domain.Parser) error {
	var (
		genServer  genServer
		f          = jen.NewFile("server")
		genCode    []jen.Code
		importName = map[string]string{
			gomodName + "/repository":           "repository",
			gomodName + "/usecase":              "usecase",
			"google.golang.org/grpc":            "grpc",
			"go.mongodb.org/mongo-driver/mongo": "mongo",
			"github.com/jinzhu/gorm":            "gorm",
		}
	)

	usecaseFile, repoFile, handlerFile, err := genServer.getAllLayer(serviceName, gomodName, "grpc")
	if err != nil {
		return err
	}

	libRepo, err := genServer.checkRepoLib(repoLib)
	if err != nil {
		return err
	}
	genCode = append(genCode, jen.Comment("slice of gRPC options"))
	genCode = append(genCode, jen.Comment("Here we can configure things like TLS"))
	genCode = append(genCode, jen.Id("opts").Op(":=[]").Qual("google.golang.org/grpc", "ServerOption").Values(jen.Dict{}))
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Id("s").Op(":=").Qual("google.golang.org/grpc", "NewServer").Call(jen.Id("opts").Op("...")))
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Id("timeoutContext").Op(":=").Qual("time", "Duration").Call(jen.Lit(2)).Op("*").Qual("time", "Second"))
	genCode = append(genCode, jen.Line())
	for _, i := range repoFile {
		genCode = append(genCode, i)
	}
	for _, i := range usecaseFile {
		genCode = append(genCode, i)
	}
	for _, i := range handlerFile {
		genCode = append(genCode, i)
	}
	genCode = append(genCode, jen.Line())
	genCode = append(genCode, jen.Return(jen.Id("s")))

	f.ImportNames(importName)
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")
	f.ImportAlias(gomodName+"/transport/grpc", "grpcHandler")
	f.Comment("GRPCServer /")
	f.Func().Id("GRPCServer").Params(libRepo).Op("*").Qual("google.golang.org/grpc", "Server").Block(
		genCode[:]...,
	)

	fileDir := fmt.Sprintf("%s/grpc_server.go", dirName)
	err = f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *genServer) getRepository(path string, gomodName string) (repo []jen.Code, err error) {
	newFs := fs.NewFsService()

	res, err := newFs.ReadDir("./" + path + "/repository")
	if err != nil {
		return repo, err
	}

	for i := range res {
		if filepath.Ext(res[i].Name()) == ".go" {
			p := parser.NewParserGeneral()
			par, err := p.GeneralParser(path + "/repository/" + res[i].Name())
			if err != nil {
				return repo, err
			}
			reg, err := regexp.Compile("[^a-zA-Z0-9]+")
			if err != nil {
				return repo, err
			}

			fileName := reg.ReplaceAllString(res[i].Name(), "")
			repoName := par.Repository.Method[0].Name
			repo = append(repo, jen.Id(fileName[:len(fileName)-2]).Op(":=").Qual(gomodName+"/repository", repoName).Call(jen.Id("db")))
		}
	}

	return repo, err
}

func (gen *genServer) getUsecase(path string, gomodName string) (usecase []jen.Code, err error) {
	newFs := fs.NewFsService()

	res, err := newFs.ReadDir("./" + path + "/usecase")
	if err != nil {
		return usecase, err
	}

	for i := range res {
		if filepath.Ext(res[i].Name()) == ".go" {
			p := parser.NewParserGeneral()
			par, err := p.GeneralParser(path + "/usecase/" + res[i].Name())
			if err != nil {
				return usecase, err
			}

			reg, err := regexp.Compile("[^a-zA-Z0-9]+")
			if err != nil {
				return usecase, err
			}

			fileName := reg.ReplaceAllString(res[i].Name(), "")
			usecaseName := par.Usecase.Method[0].Name
			repoName := fileName[:len(fileName)-9]
			usecase = append(usecase, jen.Id(fileName[:len(fileName)-2]).Op(":=").Qual(gomodName+"/usecase", usecaseName).Call(jen.Id(repoName+"repository"), jen.Id("timeoutContext")))
		}
	}

	return usecase, err
}

func (gen *genServer) getHandler(pathName string, gomodName string, transportType string) (handler []jen.Code, err error) {
	var usecaseName string
	pathName = pathName + "/transport/" + transportType

	newFs := fs.NewFsService()

	res, err := newFs.ReadDir("./" + pathName)
	if err != nil {
		return handler, err
	}

	for i := range res {
		if filepath.Ext(res[i].Name()) == ".go" {
			p := parser.NewParserGeneral()
			par, err := p.GeneralParser(pathName + "/" + res[i].Name())
			if err != nil {
				return handler, err
			}

			reg, err := regexp.Compile("[^a-zA-Z0-9]+")
			if err != nil {
				return handler, err
			}

			fileName := reg.ReplaceAllString(res[i].Name(), "")
			handlerName := par.Handler.Method[0].Name

			if transportType == "graphql" {
				res, err := newFs.ReadDir("./" + pathName + "/types")
				if err != nil {
					return handler, err
				}
				for i := range res {
					file := path.Base(res[i].Name())
					usecaseName = strings.TrimSuffix(file, filepath.Ext(file))
				}
			} else {
				usecaseName = fileName[:len(fileName)-9]
			}

			if transportType == "grpc" {
				handler = append(handler, jen.Qual(gomodName+"/transport/"+transportType, handlerName).Call(jen.Id("s"), jen.Id(usecaseName+"usecase")))
			} else {
				handler = append(handler, jen.Qual(gomodName+"/transport/"+transportType, handlerName).Call(jen.Id("r"), jen.Id(usecaseName+"usecase")))
			}
		}
	}

	return handler, err
}

func (gen *genServer) getAllLayer(serviceName string, gomodName string, transportType string) (usecase []jen.Code, repository []jen.Code, handler []jen.Code, err error) {
	var genServer genServer

	usecase, err = genServer.getUsecase(serviceName, gomodName)
	if err != nil {
		return usecase, repository, handler, err
	}

	repository, err = genServer.getRepository(serviceName, gomodName)
	if err != nil {
		return usecase, repository, handler, err
	}

	handler, err = genServer.getHandler(serviceName, gomodName, transportType)
	if err != nil {
		return usecase, repository, handler, err
	}

	return usecase, repository, handler, err
}

func (gen *genServer) checkRepoLib(repoLib string) (jen.Code, error) {
	if domain.GoPg == repoLib {
		return jen.Id("db").Op("*").Qual("github.com/go-pg/pg/v9", "DB"), nil
	}

	if domain.Gorm == repoLib {
		return jen.Id("db").Op("*").Qual("github.com/jinzhu/gorm", "DB"), nil
	}

	if domain.Sqlx == repoLib {
		return jen.Id("db").Op("*").Qual("github.com/jmoiron/sqlx", "DB"), nil
	}

	if domain.SQL == repoLib {
		return jen.Id("db").Op("*").Qual("database/sql", "DB"), nil
	}

	if domain.Mongod == repoLib {
		return jen.Id("db").Op("*").Qual("go.mongodb.org/mongo-driver/mongo", "Database"), nil
	}

	return jen.Err(), errors.New("Wrong repository library")
}
