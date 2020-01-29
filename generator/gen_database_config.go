package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
)

func (gen *caGen) GenGopgConfig(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	f := jen.NewFile("config")
	f.ImportAlias("github.com/go-pg/pg/v9", "pg")

	f.Comment("GopgInit will connecting service to databsase using go-pg orm")
	f.Func().Id("GopgInit").Params().Op("*").Qual("github.com/go-pg/pg/v9", "DB").Block(
		jen.Id("db").Op(":=").Qual("github.com/go-pg/pg/v9", "Connect").Call(jen.Op("&").Qual("github.com/go-pg/pg/v9", "Options").Values(jen.Dict{
			jen.Id("User"):     jen.Qual("os", "Getenv").Call(jen.Lit("DATABASE_USER")),
			jen.Id("Database"): jen.Qual("os", "Getenv").Call(jen.Lit("DATABASE_NAME")),
			jen.Id("Password"): jen.Qual("os", "Getenv").Call(jen.Lit("DATABASE_PASSWORD")),
		})),
		jen.Return(jen.Id("db")),
	)

	err := f.Save(dirName + "/database/config/gopg_config.go")
	if err != nil {
		return err
	}
	return nil
}

func (gen *caGen) GenGormConfig(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	f := jen.NewFile("config")

	f.Comment("GormInit will connecting service to databsase using GORM orm")
	f.Func().Id("GormInit").Params().Op("*").Qual("github.com/jinzhu/gorm", "DB").Block(
		jen.Id("dbHost").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_HOST")),
		jen.Id("dbUser").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_USER")),
		jen.Id("dbName").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_NAME")),
		jen.Id("dbPassword").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PASSWORD")),

		jen.Id("dsn").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("host=%s user=%s dbname=%s sslmode=disable password=%s"), jen.Id("dbHost"), jen.Id("dbUser"), jen.Id("dbName"), jen.Id("dbPassword")),

		jen.List(jen.Id("db"), jen.Err()).Op(":=").Qual("github.com/jinzhu/gorm", "Open").Call(jen.Lit("postgres"), jen.Id("dsn")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Panic(jen.Lit("failed to connect to database")),
		),

		jen.Return(jen.Id("db")),
	)

	err := f.Save(dirName + "/database/config/gorm_config.go")
	if err != nil {
		return err
	}
	return nil
}

func (gen *caGen) GenSQLConfig(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	f := jen.NewFile("config")

	f.Comment("SQLInit will connecting service to databsase using database/sql standard library")
	f.Func().Id("SQLInit").Params().Op("*").Qual("database/sql", "DB").Block(
		jen.Id("dbHost").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_HOST")),
		jen.Id("dbPort").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PORT")),
		jen.Id("dbUser").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_USER")),
		jen.Id("dbPass").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PASSWORD")),
		jen.Id("dbName").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_NAME")),

		jen.Id("connection").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s:%s@tcp(%s:%s)/%s"), jen.Id("dbUser"), jen.Id("dbPass"), jen.Id("dbHost"), jen.Id("dbPort"), jen.Id("dbName")),

		jen.Id("val").Op(":=").Qual("net/url", "Values").Op("{}"),
		jen.Qual("val", "Add").Call(jen.Lit("parseTime"), jen.Lit("1")),
		jen.Qual("val", "Add").Call(jen.Lit("loc"), jen.Lit("Asia/Jakarta")),

		jen.Id("dsn").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s?%s"), jen.Id("connection"), jen.Qual("val", "Encode").Op("()")),
		// dbConn, err := sql.Open(`mysql`, dsn)
		jen.List(jen.Id("dbConn"), jen.Err()).Op(":=").Qual("database/sql", "Open").Call(jen.Lit("mysql"), jen.Id("dsn")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Panic(jen.Lit("failed to connect to database")),
		),

		jen.Return(jen.Id("dbConn")),
	)

	err := f.Save(dirName + "/database/config/sql_config.go")
	if err != nil {
		return err
	}
	return nil
}

func (gen *caGen) GenSqlxConfig(dirName string, domainName string, gomodName string, parser *domain.Parser) error {
	f := jen.NewFile("config")

	f.Comment("SqlxInit will connecting service to databsase using sqlx library")
	f.Func().Id("SqlxInit").Params().Op("*").Qual("github.com/jmoiron/sqlx", "DB").Block(
		jen.Id("dbHost").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_HOST")),
		jen.Id("dbPort").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PORT")),
		jen.Id("dbUser").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_USER")),
		jen.Id("dbPass").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_PASSWORD")),
		jen.Id("dbName").Op(":=").Qual("os", "Getenv").Call(jen.Lit("DATABASE_NAME")),

		jen.Id("connection").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s:%s@tcp(%s:%s)/%s"), jen.Id("dbUser"), jen.Id("dbPass"), jen.Id("dbHost"), jen.Id("dbPort"), jen.Id("dbName")),

		jen.Id("val").Op(":=").Qual("net/url", "Values").Op("{}"),
		jen.Qual("val", "Add").Call(jen.Lit("parseTime"), jen.Lit("1")),
		jen.Qual("val", "Add").Call(jen.Lit("loc"), jen.Lit("Asia/Jakarta")),

		jen.Id("dsn").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%s?%s"), jen.Id("connection"), jen.Qual("val", "Encode").Op("()")),
		// dbConn, err := sql.Open(`mysql`, dsn)
		jen.List(jen.Id("dbConn"), jen.Err()).Op(":=").Qual("github.com/jmoiron/sqlx", "Connect").Call(jen.Lit("mysql"), jen.Id("dsn")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Panic(jen.Lit("failed to connect to database")),
		),

		jen.Return(jen.Id("dbConn")),
	)

	err := f.Save(dirName + "/database/config/sqlx_config.go")
	if err != nil {
		return err
	}
	return nil
}
