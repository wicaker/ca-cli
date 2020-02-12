package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
)

func (gen *caGen) GenMain(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	f := jen.NewFile("main")
	f.ImportAlias("github.com/sirupsen/logrus", "log")

	f.Func().Id("init").Params().Block(
		jen.Comment("load .env file"),
		jen.Id("err").Op(":=").Qual("github.com/joho/godotenv", "Load").Call(),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Qual("github.com/sirupsen/logrus", "Print").Call(jen.Id("err")),
		),

		jen.Comment("viper configuration"),
		jen.Qual("github.com/spf13/viper", "SetConfigFile").Call(jen.Lit("config.json")),
		jen.Id("err").Op("=").Qual("github.com/spf13/viper", "ReadInConfig").Call(),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Id("panic").Call(jen.Id("err")),
		),
	)

	f.Func().Id("main").Params().Block(
		jen.Id("dbGopg").Op(":=").Qual(gomodName+"/database/config", "GopgInit").Call(),
		jen.Id("errChan").Op(":=").Make(jen.Chan().Error()),

		jen.Go().Func().Params().Block(
			jen.Id("eServer").Op(":=").Qual(gomodName+"/server", "EchoServer").Call(jen.Id("dbGopg")),
			jen.Id("srv").Op(":=").Op("&").Qual("net/http", "Server").Values(jen.Dict{
				jen.Id("Addr"):         jen.Qual("github.com/spf13/viper", "GetString").Call(jen.Lit("server.echo.address")),
				jen.Id("WriteTimeout"): jen.Lit(15).Op("*").Qual("time", "Second"),
				jen.Id("ReadTimeout"):  jen.Lit(15).Op("*").Qual("time", "Second"),
			}),
			jen.Qual("github.com/sirupsen/logrus", "Fatal").Call(jen.Id("eServer").Dot("StartServer").Call(jen.Id("srv"))),
		).Call(),
		jen.Qual("github.com/sirupsen/logrus", "Fatalln").Call(jen.Op("<-").Id("errChan")),
	)

	fileDir := fmt.Sprintf("%s/main.go", dirName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}
