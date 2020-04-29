package generator

import (
	"github.com/dave/jennifer/jen"
)

func (gen *caGen) GenGopgConfig(dirName string) error {
	f := jen.NewFile("config")
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")

	f.Comment("GopgInit will connecting service to databsase using go-pg orm")
	f.Func().Id("GopgInit").Params().Op("*").Qual("github.com/go-pg/pg/v9", "DB").Block(
		jen.Id("db").Op(":=").Qual("github.com/go-pg/pg/v9", "Connect").Call(jen.Op("&").Qual("github.com/go-pg/pg/v9", "Options").Values(jen.Dict{
			jen.Id("Addr"):     jen.Qual("os", "Getenv").Call(jen.Lit("DATABASE_HOST")).Op("+").Lit(":").Op("+").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PORT")),
			jen.Id("Database"): jen.Qual("os", "Getenv").Call(jen.Lit("DATABASE_NAME")),
			jen.Id("Password"): jen.Qual("os", "Getenv").Call(jen.Lit("DATABASE_PASSWORD")),
			jen.Id("User"):     jen.Qual("os", "Getenv").Call(jen.Lit("DATABASE_USER")),
		})),
		jen.Line(),
		jen.Return(jen.Id("db")),
	)

	err := f.Save(dirName + "/gopg_config.go")
	if err != nil {
		return err
	}
	return nil
}

func (gen *caGen) GenGormConfig(dirName string) error {
	f := jen.NewFile("config")
	f.ImportName("github.com/jinzhu/gorm", "gorm")

	f.Comment("GormInit will connecting service to databsase using GORM orm")
	f.Func().Id("GormInit").Params().Op("*").Qual("github.com/jinzhu/gorm", "DB").Block(
		jen.Id("dbHost").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_HOST")),
		jen.Id("dbUser").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_USER")),
		jen.Id("dbPassword").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PASSWORD")),
		jen.Id("dbName").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_NAME")),
		jen.Line(),
		jen.Id("dsn").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("host=%s user=%s dbname=%s sslmode=disable password=%s"), jen.Id("dbHost"), jen.Id("dbUser"), jen.Id("dbName"), jen.Id("dbPassword")),
		jen.Line(),
		jen.List(jen.Id("db"), jen.Err()).Op(":=").Qual("github.com/jinzhu/gorm", "Open").Call(jen.Lit("postgres"), jen.Id("dsn")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Panic(jen.Lit("failed to connect to database")),
		),
		jen.Line(),
		jen.Return(jen.Id("db")),
	)

	err := f.Save(dirName + "/gorm_config.go")
	if err != nil {
		return err
	}
	return nil
}

func (gen *caGen) GenSQLConfig(dirName string) error {
	f := jen.NewFile("config")

	f.Comment("SQLInit will connecting service to databsase using database/sql standard library")
	f.Func().Id("SQLInit").Params().Op("*").Qual("database/sql", "DB").Block(
		jen.Id("dbHost").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_HOST")),
		jen.Id("dbPort").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PORT")),
		jen.Id("dbUser").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_USER")),
		jen.Id("dbPass").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PASSWORD")),
		jen.Id("dbName").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_NAME")),
		jen.Line(),
		jen.Id("connection").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s:%s@tcp(%s:%s)/%s"), jen.Id("dbUser"), jen.Id("dbPass"), jen.Id("dbHost"), jen.Id("dbPort"), jen.Id("dbName")),
		jen.Id("val").Op(":=").Qual("net/url", "Values").Op("{}"),
		jen.Id("val").Dot("Add").Call(jen.Lit("parseTime"), jen.Lit("1")),
		jen.Id("val").Dot("Add").Call(jen.Lit("loc"), jen.Lit("Asia/Jakarta")),
		jen.Id("dsn").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s?%s"), jen.Id("connection"), jen.Id("val").Dot("Encode").Op("()")),
		jen.List(jen.Id("dbConn"), jen.Err()).Op(":=").Qual("database/sql", "Open").Call(jen.Lit("mysql"), jen.Id("dsn")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Panic(jen.Err()),
		),
		jen.Line(),
		jen.Return(jen.Id("dbConn")),
	)

	err := f.Save(dirName + "/sql_config.go")
	if err != nil {
		return err
	}
	return nil
}

func (gen *caGen) GenSqlxConfig(dirName string) error {
	f := jen.NewFile("config")
	f.ImportName("github.com/jmoiron/sqlx", "sqlx")

	f.Comment("SqlxInit will connecting service to databsase using sqlx library")
	f.Func().Id("SqlxInit").Params().Op("*").Qual("github.com/jmoiron/sqlx", "DB").Block(
		jen.Id("dbHost").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_HOST")),
		jen.Id("dbPort").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PORT")),
		jen.Id("dbUser").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_USER")),
		jen.Id("dbPass").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PASSWORD")),
		jen.Id("dbName").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_NAME")),
		jen.Line(),
		jen.Id("connection").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s:%s@tcp(%s:%s)/%s"), jen.Id("dbUser"), jen.Id("dbPass"), jen.Id("dbHost"), jen.Id("dbPort"), jen.Id("dbName")),
		jen.Id("val").Op(":=").Qual("net/url", "Values").Op("{}"),
		jen.Id("val").Dot("Add").Call(jen.Lit("parseTime"), jen.Lit("1")),
		jen.Id("val").Dot("Add").Call(jen.Lit("loc"), jen.Lit("Asia/Jakarta")),
		jen.Id("dsn").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s?%s"), jen.Id("connection"), jen.Id("val").Dot("Encode").Op("()")),
		jen.List(jen.Id("dbConn"), jen.Err()).Op(":=").Qual("github.com/jmoiron/sqlx", "Connect").Call(jen.Lit("mysql"), jen.Id("dsn")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Panic(jen.Err()),
		),
		jen.Line(),
		jen.Return(jen.Id("dbConn")),
	)

	err := f.Save(dirName + "/sqlx_config.go")
	if err != nil {
		return err
	}
	return nil
}

func (gen *caGen) GenMongodConfig(dirName string) error {
	var (
		importName = map[string]string{
			"go.mongodb.org/mongo-driver/mongo":          "mongo",
			"go.mongodb.org/mongo-driver/mongo/options":  "options",
			"go.mongodb.org/mongo-driver/mongo/readpref": "readpref",
		}
		f = jen.NewFile("config")
	)
	f.ImportNames(importName)

	f.Comment("MongodInit will connecting service to databsase using mongo-driver")
	f.Func().Id("MongodInit").Params().Op("*").Qual("go.mongodb.org/mongo-driver/mongo", "Database").Block(
		jen.Id("clientOptions").Op(":=").Qual("go.mongodb.org/mongo-driver/mongo/options", "Client").Call().Dot("ApplyURI").Call(jen.Qual("os", "Getenv").Call(jen.Lit("DATABASE_MONGO_URL"))),
		jen.Id("client").Op(",").Err().Op(":=").Qual("go.mongodb.org/mongo-driver/mongo", "NewClient").Call(jen.Id("clientOptions")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Qual("log", "Fatal").Call(jen.Err()),
		),
		jen.Line(),
		jen.Id("ctx").Op(",").Id("cancel").Op(":=").Qual("context", "WithTimeout").Call(jen.Qual("context", "Background").Call(), jen.Lit(10).Op("*").Qual("time", "Second")),
		jen.Err().Op("=").Id("client").Dot("Connect").Call(jen.Id("ctx")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Qual("log", "Fatal").Call(jen.Err()),
		),
		jen.Line(),
		jen.Defer().Id("cancel").Call(),
		jen.Line(),
		jen.Err().Op("=").Id("client").Dot("Ping").Call(jen.Qual("context", "Background").Call(), jen.Qual("go.mongodb.org/mongo-driver/mongo/readpref", "Primary").Call()),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Qual("log", "Fatal").Call(jen.Lit("Couldn't connect to the database "), jen.Err()),
		).Else().Block(
			jen.Qual("log", "Println").Call(jen.Lit("Connected!")),
		),
		jen.Line(),
		jen.Id("db").Op(":=").Id("client").Dot("Database").Call(jen.Qual("os", "Getenv").Call(jen.Lit("DATABASE_MONGO_NAME"))),
		jen.Return(jen.Id("db")),
	)

	err := f.Save(dirName + "/mongod_config.go")
	if err != nil {
		return err
	}
	return nil
}
