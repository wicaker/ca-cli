package generator

import (
	"fmt"
	"strings"

	"github.com/wicaker/cacli/domain"

	"github.com/dave/jennifer/jen"
)

type goServer struct{}

func (gen *caGen) GenMain(dirName string, gomodName string, repoLib string, transport []string) error {
	var (
		g          goServer
		importName = map[string]string{
			gomodName + "/server":          "server",
			gomodName + "/database/config": "config",
			"github.com/joho/godotenv":     "godotenv",
		}
		server   []jen.Code
		dbConfig jen.Code
		dbConf   = "db" + strings.ToLower(repoLib)
	)

	if repoLib == domain.GoPg {
		dbConfig = jen.Id(dbConf).Op(":=").Qual(gomodName+"/database/config", "GopgInit").Call()
	}
	if repoLib == domain.Gorm {
		dbConfig = jen.Id(dbConf).Op(":=").Qual(gomodName+"/database/config", "GormInit").Call()
	}
	if repoLib == domain.SQL {
		dbConfig = jen.Id(dbConf).Op(":=").Qual(gomodName+"/database/config", "SQLInit").Call()
	}
	if repoLib == domain.Sqlx {
		dbConfig = jen.Id(dbConf).Op(":=").Qual(gomodName+"/database/config", "SqlxInit").Call()
	}
	if repoLib == domain.Mongod {
		dbConfig = jen.Id(dbConf).Op(":=").Qual(gomodName+"/database/config", "MongodInit").Call()
	}

	server = append(server, dbConfig)
	server = append(server, jen.Line())
	server = append(server, jen.Id("errChan").Op(":=").Make(jen.Chan().Error()))
	server = append(server, jen.Line())
	for i := range transport {
		if transport[i] == domain.Echo {
			server = append(server, g.EchoServer(gomodName, dbConf))
			server = append(server, jen.Line())
		}
		if transport[i] == domain.Gin {
			server = append(server, g.GinServer(gomodName, dbConf))
			server = append(server, jen.Line())
		}
		if transport[i] == domain.GorillaMux {
			server = append(server, g.GorillaMuxServer(gomodName, dbConf))
			server = append(server, jen.Line())
		}
		if transport[i] == domain.NetHTTP {
			server = append(server, g.NetHTTPMuxServer(gomodName, dbConf))
			server = append(server, jen.Line())
		}
		if transport[i] == domain.Graphql {
			server = append(server, g.GraphQLServer(gomodName, dbConf))
			server = append(server, jen.Line())
		}
		if transport[i] == domain.Grpc {
			server = append(server, g.GRPCServer(gomodName, dbConf))
			server = append(server, jen.Line())
		}
	}
	server = append(server, jen.Qual("github.com/sirupsen/logrus", "Fatalln").Call(jen.Op("<-").Id("errChan")))

	f := jen.NewFile("main")
	f.ImportAlias("github.com/sirupsen/logrus", "log")
	f.Anon("github.com/go-sql-driver/mysql")
	f.Anon("github.com/jinzhu/gorm/dialects/postgres")
	f.ImportNames(importName)

	f.Func().Id("init").Params().Block(
		jen.Err().Op(":=").Qual("github.com/joho/godotenv", "Load").Call(),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Qual("github.com/sirupsen/logrus", "Print").Call(jen.Err()),
		),
	)

	f.Line()
	f.Func().Id("main").Params().Block(
		server...,
	)

	fileDir := fmt.Sprintf("%s/main.go", dirName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (g *goServer) EchoServer(gomodName string, dbConfig string) (code jen.Code) {
	code = (jen.Go().Func().Params().Block(
		jen.Id("eServer").Op(":=").Qual(gomodName+"/server", "EchoServer").Call(jen.Id(dbConfig)),
		jen.Id("srv").Op(":=").Op("&").Qual("net/http", "Server").Values(jen.Dict{
			jen.Id("Addr"):         jen.Lit(":").Op("+").Qual("os", "Getenv").Call(jen.Lit("SERVER_ECHO_PORT")),
			jen.Id("WriteTimeout"): jen.Lit(15).Op("*").Qual("time", "Second"),
			jen.Id("ReadTimeout"):  jen.Lit(15).Op("*").Qual("time", "Second"),
		}),
		jen.Id("eServer").Dot("HideBanner").Op("=").True(),
		jen.Id("eServer").Dot("HidePort").Op("=").True(),
		jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
			jen.Lit("at"): jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
		})).Dot("Printf").Call(jen.Lit("Starting echo server on port :%s..."), jen.Qual("os", "Getenv").Call(jen.Lit("SERVER_ECHO_PORT"))),
		jen.Qual("github.com/sirupsen/logrus", "Fatal").Call(jen.Id("eServer").Dot("StartServer").Call(jen.Id("srv"))),
	).Call())
	return
}

func (g *goServer) GinServer(gomodName string, dbConfig string) (code jen.Code) {
	code = (jen.Go().Func().Params().Block(
		jen.Id("gServer").Op(":=").Qual(gomodName+"/server", "GinServer").Call(jen.Id(dbConfig)),
		jen.Id("srv").Op(":=").Op("&").Qual("net/http", "Server").Values(jen.Dict{
			jen.Id("Addr"):         jen.Lit(":").Op("+").Qual("os", "Getenv").Call(jen.Lit("SERVER_GIN_PORT")),
			jen.Id("Handler"):      jen.Id("gServer"),
			jen.Id("WriteTimeout"): jen.Lit(15).Op("*").Qual("time", "Second"),
			jen.Id("ReadTimeout"):  jen.Lit(15).Op("*").Qual("time", "Second"),
		}),
		jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
			jen.Lit("at"): jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
		})).Dot("Printf").Call(jen.Lit("Starting gin server on port :%s..."), jen.Qual("os", "Getenv").Call(jen.Lit("SERVER_GIN_PORT"))),
		jen.Qual("github.com/sirupsen/logrus", "Fatal").Call(jen.Id("srv").Dot("ListenAndServe").Call()),
	).Call())
	return
}

func (g *goServer) GorillaMuxServer(gomodName string, dbConfig string) (code jen.Code) {
	code = (jen.Go().Func().Params().Block(
		jen.Id("gmServer").Op(":=").Qual(gomodName+"/server", "GorillaMuxServer").Call(jen.Id(dbConfig)),
		jen.Id("srv").Op(":=").Op("&").Qual("net/http", "Server").Values(jen.Dict{
			jen.Id("Addr"):         jen.Lit(":").Op("+").Qual("os", "Getenv").Call(jen.Lit("SERVER_GORILLA_MUX_PORT")),
			jen.Id("Handler"):      jen.Id("gmServer"),
			jen.Id("WriteTimeout"): jen.Lit(15).Op("*").Qual("time", "Second"),
			jen.Id("ReadTimeout"):  jen.Lit(15).Op("*").Qual("time", "Second"),
		}),
		jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
			jen.Lit("at"): jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
		})).Dot("Printf").Call(jen.Lit("Starting gorilla mux server on port :%s..."), jen.Qual("os", "Getenv").Call(jen.Lit("SERVER_GORILLA_MUX_PORT"))),
		jen.Qual("github.com/sirupsen/logrus", "Fatal").Call(jen.Id("srv").Dot("ListenAndServe").Call()),
	).Call())
	return
}

func (g *goServer) NetHTTPMuxServer(gomodName string, dbConfig string) (code jen.Code) {
	code = (jen.Go().Func().Params().Block(
		jen.Id("httpMuxServer").Op(":=").Qual(gomodName+"/server", "MuxServer").Call(jen.Id(dbConfig)),
		jen.Id("srv").Op(":=").Op("&").Qual("net/http", "Server").Values(jen.Dict{
			jen.Id("Addr"):         jen.Lit(":").Op("+").Qual("os", "Getenv").Call(jen.Lit("SERVER_NET_HTTP_SERVER_MUX_PORT")),
			jen.Id("Handler"):      jen.Id("httpMuxServer"),
			jen.Id("WriteTimeout"): jen.Lit(15).Op("*").Qual("time", "Second"),
			jen.Id("ReadTimeout"):  jen.Lit(15).Op("*").Qual("time", "Second"),
		}),
		jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
			jen.Lit("at"): jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
		})).Dot("Printf").Call(jen.Lit("Starting net/http ServerMux on port :%s..."), jen.Qual("os", "Getenv").Call(jen.Lit("SERVER_NET_HTTP_SERVER_MUX_PORT"))),
		jen.Qual("github.com/sirupsen/logrus", "Fatal").Call(jen.Id("srv").Dot("ListenAndServe").Call()),
	).Call())
	return
}

func (g *goServer) GraphQLServer(gomodName string, dbConfig string) (code jen.Code) {
	code = (jen.Go().Func().Params().Block(
		jen.Id("httpMuxServer").Op(":=").Qual(gomodName+"/server", "GraphQLServer").Call(jen.Id(dbConfig)),
		jen.Id("srv").Op(":=").Op("&").Qual("net/http", "Server").Values(jen.Dict{
			jen.Id("Addr"):         jen.Lit(":").Op("+").Qual("os", "Getenv").Call(jen.Lit("SERVER_GRAPHQL_SERVER_MUX_PORT")),
			jen.Id("Handler"):      jen.Id("httpMuxServer"),
			jen.Id("WriteTimeout"): jen.Lit(15).Op("*").Qual("time", "Second"),
			jen.Id("ReadTimeout"):  jen.Lit(15).Op("*").Qual("time", "Second"),
		}),
		jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
			jen.Lit("at"): jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
		})).Dot("Printf").Call(jen.Lit("Starting Graphql Server on port :%s..."), jen.Qual("os", "Getenv").Call(jen.Lit("SERVER_GRAPHQL_SERVER_MUX_PORT"))),
		jen.Qual("github.com/sirupsen/logrus", "Fatal").Call(jen.Id("srv").Dot("ListenAndServe").Call()),
	).Call())
	return
}

func (g *goServer) GRPCServer(gomodName string, dbConfig string) (code jen.Code) {
	code = (jen.Go().Func().Params().Block(
		jen.Comment("50051 is the default port for gRPC"),
		jen.Id("listener").Op(",").Err().Op(":=").Qual("net", "Listen").Call(jen.Lit("tcp"), jen.Lit(":").Op("+").Qual("os", "Getenv").Call(jen.Lit("SERVER_GRPC_PORT"))),
		jen.Line(),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Qual("github.com/sirupsen/logrus", "Fatalf").Call(jen.Lit("Unable to listen on port :%v : %v"), jen.Qual("os", "Getenv").Call(jen.Lit("SERVER_GRPC_PORT")), jen.Err()),
		),
		jen.Line(),
		jen.Id("s").Op(":=").Qual(gomodName+"/server", "GRPCServer").Call(jen.Id(dbConfig)),
		jen.Line(),
		jen.Comment("Start the server"),
		jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
			jen.Lit("at"): jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
		})).Dot("Printf").Call(jen.Lit("Starting GRPC server on port :%s..."), jen.Qual("os", "Getenv").Call(jen.Lit("SERVER_GRPC_PORT"))),
		jen.If(jen.Err().Op(":=").Id("s").Dot("Serve").Call(jen.Id("listener")).Op(";").Err().Op("!=").Nil()).Block(
			jen.Qual("github.com/sirupsen/logrus", "Fatalf").Call(jen.Lit("Failed to serve: %v"), jen.Err()),
		),
	).Call())
	return
}
