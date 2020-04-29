package generator

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
	"github.com/wicaker/cacli/fs"
)

func (gen *caGen) GenEchoTransport(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		useCase    = parser.Usecase.Name
		newD       = fmt.Sprintf("New%sHandler", strings.ToUpper(string(domainName[0]))+domainName[1:])
		comment    = fmt.Sprintf("%s will initialize the %s endpoint", newD, domainName)
		handler    []jen.Code
		f          = jen.NewFile("rest")
		importName = map[string]string{
			gomodName + "/domain":      "domain",
			"github.com/labstack/echo": "echo",
		}
	)

	handler = append(handler, jen.Id("handler").Op(":=").Op("&").Id(domainName+"Handler").Values(jen.Dict{
		jen.Id(useCase): jen.Id("u"),
	}))

	for _, i := range parser.Usecase.Method {
		path := fmt.Sprintf("/%s/%s", domainName, strings.ToLower(i.Name))
		handler = append(handler, jen.Id("e").Dot("GET").Call(jen.Lit(path), jen.Id("handler").Dot(i.Name+"Handler")))
	}

	f.ImportNames(importName)

	f.Type().Id(domainName + "Handler").Struct(
		jen.Id(useCase).Qual(gomodName+"/domain", useCase),
	)

	f.Comment(comment)
	f.Func().Id(newD).Params(
		jen.Id("e").Op("*").Qual("github.com/labstack/echo", "Echo"),
		jen.Id("u").Qual(gomodName+"/domain", useCase),
	).Block(handler[:]...)

	for _, i := range parser.Usecase.Method {
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"h").Op("*").Id(domainName+"Handler")).
			Id(i.Name+"Handler").Params(jen.Id("c").Qual("github.com/labstack/echo", "Context")).Call(jen.Error()).Block(
			jen.Id("ctx").Op(":=").Id("c").Dot("Request").Call().Dot("Context").Call(),
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),

			jen.Qual(gomodName+"/domain", "ResponseData").Op("=").Make(jen.Map(jen.String()).Interface()),
			jen.Qual(gomodName+"/domain", "ResponseData").Op("[").Lit("message").Op("]=").Lit(i.Name+"Handler"),
			jen.Return(jen.Id("c").Dot("JSON").Call(jen.Qual("net/http", "StatusOK"), jen.Qual(gomodName+"/domain", "ResponseData"))),
		)
	}

	fileDir := fmt.Sprintf("%s/%s_handler.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGinTransport(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		useCase    = parser.Usecase.Name
		newD       = fmt.Sprintf("New%sHandler", strings.ToUpper(string(domainName[0]))+domainName[1:])
		comment    = fmt.Sprintf("%s will initialize the %s endpoint", newD, domainName)
		handler    []jen.Code
		f          = jen.NewFile("rest")
		importName = map[string]string{
			gomodName + "/domain":      "domain",
			"github.com/gin-gonic/gin": "gin",
		}
	)

	handler = append(handler, jen.Id("handler").Op(":=").Op("&").Id(domainName+"Handler").Values(jen.Dict{
		jen.Id(useCase): jen.Id("u"),
	}))

	for _, i := range parser.Usecase.Method {
		path := fmt.Sprintf("/%s/%s", domainName, strings.ToLower(i.Name))
		handler = append(handler, jen.Id("r").Dot("GET").Call(jen.Lit(path), jen.Id("handler").Dot(i.Name+"Handler")))
	}

	f.ImportNames(importName)

	f.Type().Id(domainName + "Handler").Struct(
		jen.Id(useCase).Qual(gomodName+"/domain", useCase),
	)

	f.Comment(comment)
	f.Func().Id(newD).Params(
		jen.Id("r").Op("*").Qual("github.com/gin-gonic/gin", "Engine"),
		jen.Id("u").Qual(gomodName+"/domain", useCase),
	).Block(handler[:]...)

	for _, i := range parser.Usecase.Method {
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"h").Op("*").Id(domainName+"Handler")).
			Id(i.Name+"Handler").Params(jen.Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context")).Block(
			jen.Id("ctx").Op(":=").Id("c").Dot("Request").Dot("Context").Call(),
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),

			jen.Qual(gomodName+"/domain", "ResponseData").Op("=").Make(jen.Map(jen.String()).Interface()),
			jen.Qual(gomodName+"/domain", "ResponseData").Op("[").Lit("message").Op("]=").Lit(i.Name+"Handler"),
			jen.Id("c").Dot("JSON").Call(jen.Qual("net/http", "StatusOK"), jen.Qual(gomodName+"/domain", "ResponseData")),
			jen.Return(),
		)
	}

	fileDir := fmt.Sprintf("%s/%s_handler.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGorillaMuxTransport(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		useCase    = parser.Usecase.Name
		newD       = fmt.Sprintf("New%sHandler", strings.ToUpper(string(domainName[0]))+domainName[1:])
		comment    = fmt.Sprintf("%s will initialize the %s endpoint", newD, domainName)
		handler    []jen.Code
		f          = jen.NewFile("rest")
		importName = map[string]string{
			gomodName + "/domain":    "domain",
			"github.com/gorilla/mux": "mux",
		}
	)

	handler = append(handler, jen.Id("handler").Op(":=").Op("&").Id(domainName+"Handler").Values(jen.Dict{
		jen.Id(useCase): jen.Id("u"),
	}))

	for _, i := range parser.Usecase.Method {
		path := fmt.Sprintf("/%s/%s", domainName, strings.ToLower(i.Name))
		handler = append(handler, jen.Id("r").Dot("HandleFunc").Call(jen.Lit(path), jen.Id("handler").Dot(i.Name+"Handler")).Dot("Methods").Call(jen.Lit("GET")))
	}

	f.ImportNames(importName)
	f.ImportAlias("github.com/json-iterator/go", "json")

	f.Type().Id(domainName + "Handler").Struct(
		jen.Id(useCase).Qual(gomodName+"/domain", useCase),
	)

	f.Comment(comment)
	f.Func().Id(newD).Params(
		jen.Id("r").Op("*").Qual("github.com/gorilla/mux", "Router"),
		jen.Id("u").Qual(gomodName+"/domain", useCase),
	).Block(handler[:]...)

	for _, i := range parser.Usecase.Method {
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"h").Op("*").Id(domainName+"Handler")).
			Id(i.Name+"Handler").Params(jen.Id("w").Qual("net/http", "ResponseWriter"), jen.Id("r").Op("*").Qual("net/http", "Request")).Block(
			jen.Id("ctx").Op(":=").Id("r").Dot("Context").Call(),
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),

			jen.Qual(gomodName+"/domain", "ResponseData").Op("=").Make(jen.Map(jen.String()).Interface()),
			jen.Qual(gomodName+"/domain", "ResponseData").Op("[").Lit("message").Op("]=").Lit(i.Name+"Handler"),
			jen.Id("w").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
			jen.Qual("github.com/json-iterator/go", "NewEncoder").Call(jen.Id("w")).Dot("Encode").Call(jen.Qual(gomodName+"/domain", "ResponseData")),
			jen.Return(),
		)
	}

	fileDir := fmt.Sprintf("%s/%s_handler.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenNetHTTPTransport(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file        = path.Base(domainFile)
		domainName  = strings.TrimSuffix(file, filepath.Ext(file))
		useCase     = parser.Usecase.Name
		newD        = fmt.Sprintf("New%sHandler", strings.ToUpper(string(domainName[0]))+domainName[1:])
		comment     = fmt.Sprintf("%s will initialize the %s endpoint", newD, domainName)
		handlerName []string
		f           = jen.NewFile("rest")
		path        = fmt.Sprintf("/%s", domainName)
		importName  = map[string]string{
			gomodName + "/domain": "domain",
		}
	)

	f.ImportNames(importName)
	f.ImportAlias("github.com/json-iterator/go", "json")

	f.Type().Id(domainName + "Handler").Struct(
		jen.Id(useCase).Qual(gomodName+"/domain", useCase),
	)

	f.Comment(comment)
	f.Func().Id(newD).Params(
		jen.Id("r").Op("*").Qual("net/http", "ServeMux"),
		jen.Id("u").Qual(gomodName+"/domain", useCase),
	).Block(
		jen.Id("handler").Op(":=").Op("&").Id(domainName+"Handler").Values(jen.Dict{
			jen.Id(useCase): jen.Id("u"),
		}),
		jen.Id("r").Dot("Handle").Call(jen.Lit(path), jen.Id("handler")),
	)

	for _, i := range parser.Usecase.Method {
		handlerName = append(handlerName, i.Name+"Handler")
		f.Line()
		f.Func().
			Params(jen.Id(string(domainName[0])+"h").Op("*").Id(domainName+"Handler")).
			Id(i.Name+"Handler").Params(jen.Id("w").Qual("net/http", "ResponseWriter"), jen.Id("r").Op("*").Qual("net/http", "Request")).Block(
			jen.Id("ctx").Op(":=").Id("r").Dot("Context").Call(),
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),

			jen.Qual(gomodName+"/domain", "ResponseData").Op("=").Make(jen.Map(jen.String()).Interface()),
			jen.Qual(gomodName+"/domain", "ResponseData").Op("[").Lit("message").Op("]=").Lit(i.Name+"Handler"),
			jen.Id("w").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
			jen.Qual("github.com/json-iterator/go", "NewEncoder").Call(jen.Id("w")).Dot("Encode").Call(jen.Qual(gomodName+"/domain", "ResponseData")),
			jen.Return(),
		)
	}

	f.Line()
	f.Func().
		Params(jen.Id(string(domainName[0])+"h").Op("*").Id(domainName+"Handler")).
		Id("ServeHTTP").Params(jen.Id("w").Qual("net/http", "ResponseWriter"), jen.Id("r").Op("*").Qual("net/http", "Request")).Block(
		jen.Id("ctx").Op(":=").Id("r").Dot("Context").Call(),
		jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
		jen.Line(),

		jen.Id("uri").Op(":=").Id("ctx").Dot("Value").Call(jen.Lit("uri")),
		jen.Id("s").Op(":=").Qual("strings", "Split").Call(jen.Id("r").Dot("RequestURI"), jen.Lit("/")),
		jen.Line(),

		jen.If(jen.Id("uri").Op("==").Lit(path)).Block(
			jen.If(jen.Id("r").Dot("Method").Op("==").Qual("net/http", "MethodGet")).Block(
				jen.Id(string(domainName[0])+"h").Dot(handlerName[0]).Call(jen.Id("w"), jen.Id("r")),
				jen.Return(),
			),
		),
		jen.Line(),

		jen.Comment("// if the path contain params"),
		jen.If(jen.Len(jen.Id("s")).Op("==").Lit(3)).Block(
			jen.Id("id").Op(",").Err().Op(":=").Id("strconv").Dot("Atoi").Call(jen.Qual("path", "Base").Call(jen.Id("r").Dot("RequestURI"))),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Id("w").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusMethodNotAllowed")),
				jen.Qual("github.com/json-iterator/go", "NewEncoder").Call(jen.Id("w")).Dot("Encode").Call(jen.Op("&").Qual(gomodName+"/domain", "ResponseError").Values(jen.Dict{
					jen.Id("Message"): jen.Lit("Method Not Allowed"),
				})),
				jen.Return(),
			),
			jen.If(jen.Id("r").Dot("Method").Op("==").Qual("net/http", "MethodGet")).Block(
				jen.Id(string(domainName[0])+"h").Dot(handlerName[0]).Call(jen.Id("w"), jen.Id("r"), jen.Id("id")),
				jen.Return(),
			),
		),
		jen.Line(),

		jen.Id("w").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusMethodNotAllowed")),
		jen.Qual("github.com/json-iterator/go", "NewEncoder").Call(jen.Id("w")).Dot("Encode").Call(jen.Op("&").Qual(gomodName+"/domain", "ResponseError").Values(jen.Dict{
			jen.Id("Message"): jen.Lit("Method Not Allowed"),
		})),
		jen.Return(),
	)

	fileDir := fmt.Sprintf("%s/%s_handler.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGraphqlTransport(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		useCase    = parser.Usecase.Name
		newFs      = fs.NewFsService()
		importName = map[string]string{
			gomodName + "/domain":                      "domain",
			"github.com/graphql-go/graphql":            "graphql",
			gomodName + "/transport/graphql/types":     "types",
			gomodName + "/transport/graphql/mutations": "mutations",
			gomodName + "/transport/graphql/queries":   "queries",
			"github.com/graphql-go/handler":            "handler",
		}
	)
	// types
	// =>>example.go
	typeGen := func(dirName string, domainName string) error {
		typeName := fmt.Sprintf("%sType", strings.ToUpper(string(domainName[0]))+domainName[1:])
		f := jen.NewFile("types")
		f.ImportNames(importName)

		f.Comment(fmt.Sprintf("%s is the GraphQL schema for the %s type.", typeName, domainName))
		f.Var().Id(typeName).Op("=").Qual("github.com/graphql-go/graphql", "NewObject").Call(
			jen.Qual("github.com/graphql-go/graphql", "ObjectConfig").Values(jen.Dict{
				jen.Id("Name"): jen.Lit(strings.ToUpper(string(domainName[0])) + domainName[1:]),
				jen.Id("Fields"): jen.Qual("github.com/graphql-go/graphql", "Fields").Values(jen.Dict{
					jen.Lit("id"): jen.Op("&").Qual("github.com/graphql-go/graphql", "Field").Values(jen.Dict{
						jen.Id("Type"): jen.Qual("github.com/graphql-go/graphql", "ID"),
					}),
				}),
			}),
		)
		// create types directory
		err := newFs.CreateDir(dirName + "/types")
		if err != nil {
			return err
		}

		// save file
		fileDir := fmt.Sprintf("%s/types/%s.go", dirName, domainName)
		err = f.Save(fileDir)
		if err != nil {
			return err
		}

		return nil
	}
	err := typeGen(dirName, domainName)
	if err != nil {
		return err
	}

	// mutations
	// =>>mutations.go
	mutationGen := func(dirName string, domainName string) error {
		graphFields := jen.Dict{}
		f := jen.NewFile("mutations")
		f.ImportNames(importName)

		f.Comment("GraphQLMutation represent the graphQLMutation")
		f.Type().Id("GraphQLMutation").Struct(
			jen.Id(useCase).Qual(gomodName+"/domain", useCase),
		)

		f.Comment("NewGraphQLMutation will initialize mutations")
		f.Func().Id("NewGraphQLMutation").Params(
			jen.Id(string(domainName[0])).Qual(gomodName+"/domain", useCase),
		).Op("*").Id("GraphQLMutation").Block(
			jen.Return(jen.Op("&").Id("GraphQLMutation").Values(jen.Dict{
				jen.Id(useCase): jen.Id(string(domainName[0])),
			})),
		)

		for _, i := range parser.Usecase.Method {
			graphFields[jen.Lit(domainName+i.Name)] = jen.Id("gm").Dot(i.Name + strings.ToUpper(string(domainName[0])) + domainName[1:] + "Mutation").Call()
		}

		f.Comment("GetRootMutationFields returns all the available mutations.")
		f.Func().
			Params(jen.Id("gm").Op("*").Id("GraphQLMutation")).Id("GetRootMutationFields").Params().Qual("github.com/graphql-go/graphql", "Fields").Block(
			jen.Return(jen.Qual("github.com/graphql-go/graphql", "Fields").Values(graphFields)),
		)

		// create mutations directory
		err := newFs.CreateDir(dirName + "/mutations")
		if err != nil {
			return err
		}

		// save file
		fileDir := fmt.Sprintf("%s/mutations/mutations.go", dirName)
		err = f.Save(fileDir)
		if err != nil {
			return err
		}

		return nil
	}

	// =>>example.go
	mutationFieldGen := func(dirName string, domainName string) error {
		typeName := fmt.Sprintf("%sType", strings.ToUpper(string(domainName[0]))+domainName[1:])
		f := jen.NewFile("mutations")
		f.ImportNames(importName)

		for _, i := range parser.Usecase.Method {
			f.Line()
			f.Comment(i.Name + strings.ToUpper(string(domainName[0])) + domainName[1:] + "Mutation /.")
			f.Func().
				Params(jen.Id("gm").Op("*").Id("GraphQLMutation")).
				Id(i.Name+strings.ToUpper(string(domainName[0]))+domainName[1:]+"Mutation").Params().Op("*").Qual("github.com/graphql-go/graphql", "Field").Block(
				jen.Return(jen.Op("&").Qual("github.com/graphql-go/graphql", "Field").Values(jen.Dict{
					jen.Id("Type"):        jen.Qual(gomodName+"/transport/graphql/types", typeName),
					jen.Id("Description"): jen.Lit(strings.ToUpper(string(domainName[0])) + domainName[1:]),
					jen.Id("Args"):        jen.Qual("github.com/graphql-go/graphql", "FieldConfigArgument").Values(jen.Dict{}),
					jen.Id("Resolve"): jen.Func().Params(jen.Id("params").Qual("github.com/graphql-go/graphql", "ResolveParams")).Call(jen.Interface(), jen.Error()).Block(
						jen.Id("ctx").Op(":=").Id("params").Dot("Context"),
						jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
						jen.Return(jen.Nil(), jen.Nil()),
					),
				})),
			)
		}

		// save file
		fileDir := fmt.Sprintf("%s/mutations/%s.go", dirName, domainName)
		err = f.Save(fileDir)
		if err != nil {
			return err
		}

		return nil
	}

	err = mutationGen(dirName, domainName)
	if err != nil {
		return err
	}

	err = mutationFieldGen(dirName, domainName)
	if err != nil {
		return err
	}

	// queries
	// =>>queries.go
	queryGen := func(dirName string, domainName string) error {
		graphFields := jen.Dict{}
		f := jen.NewFile("queries")
		f.ImportNames(importName)

		f.Comment("GraphQLQuery represent the GraphQLQuery")
		f.Type().Id("GraphQLQuery").Struct(
			jen.Id(useCase).Qual(gomodName+"/domain", useCase),
		)

		f.Comment("NewGraphQLQuery will initialize queries")
		f.Func().Id("NewGraphQLQuery").Params(
			jen.Id(string(domainName[0])).Qual(gomodName+"/domain", useCase),
		).Op("*").Id("GraphQLQuery").Block(
			jen.Return(jen.Op("&").Id("GraphQLQuery").Values(jen.Dict{
				jen.Id(useCase): jen.Id(string(domainName[0])),
			})),
		)

		for _, i := range parser.Usecase.Method {
			graphFields[jen.Lit(domainName+i.Name)] = jen.Id("gq").Dot(i.Name + strings.ToUpper(string(domainName[0])) + domainName[1:] + "Query").Call()
		}

		f.Comment("GetRootQueryFields returns all the available queries.")
		f.Func().
			Params(jen.Id("gq").Op("*").Id("GraphQLQuery")).Id("GetRootQueryFields").Params().Qual("github.com/graphql-go/graphql", "Fields").Block(
			jen.Return(jen.Qual("github.com/graphql-go/graphql", "Fields").Values(graphFields)),
		)

		// create queries directory
		err := newFs.CreateDir(dirName + "/queries")
		if err != nil {
			return err
		}

		// save file
		fileDir := fmt.Sprintf("%s/queries/queries.go", dirName)
		err = f.Save(fileDir)
		if err != nil {
			return err
		}

		return nil
	}

	// =>>example.go
	queryFieldGen := func(dirName string, domainName string) error {
		typeName := fmt.Sprintf("%sType", strings.ToUpper(string(domainName[0]))+domainName[1:])
		f := jen.NewFile("queries")
		f.ImportNames(importName)

		for _, i := range parser.Usecase.Method {
			f.Line()
			f.Comment(i.Name + strings.ToUpper(string(domainName[0])) + domainName[1:] + "Query /.")
			f.Func().
				Params(jen.Id("gq").Op("*").Id("GraphQLQuery")).
				Id(i.Name+strings.ToUpper(string(domainName[0]))+domainName[1:]+"Query").Params().Op("*").Qual("github.com/graphql-go/graphql", "Field").Block(
				jen.Return(jen.Op("&").Qual("github.com/graphql-go/graphql", "Field").Values(jen.Dict{
					jen.Id("Type"):        jen.Qual(gomodName+"/transport/graphql/types", typeName),
					jen.Id("Description"): jen.Lit(strings.ToUpper(string(domainName[0])) + domainName[1:]),
					jen.Id("Args"): jen.Qual("github.com/graphql-go/graphql", "FieldConfigArgument").Values(jen.Dict{
						jen.Lit("id"): jen.Op("&").Qual("github.com/graphql-go/graphql", "ArgumentConfig").Values(jen.Dict{
							jen.Id("Type"): jen.Qual("github.com/graphql-go/graphql", "ID"),
						}),
					}),
					jen.Id("Resolve"): jen.Func().Params(jen.Id("params").Qual("github.com/graphql-go/graphql", "ResolveParams")).Call(jen.Interface(), jen.Error()).Block(
						jen.Id("ctx").Op(":=").Id("params").Dot("Context"),
						jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
						jen.Return(jen.Nil(), jen.Nil()),
					),
				})),
			)
		}

		// save file
		fileDir := fmt.Sprintf("%s/queries/%s.go", dirName, domainName)
		err = f.Save(fileDir)
		if err != nil {
			return err
		}

		return nil
	}

	err = queryGen(dirName, domainName)
	if err != nil {
		return err
	}

	err = queryFieldGen(dirName, domainName)
	if err != nil {
		return err
	}

	// index.go
	indexGen := func() error {
		f := jen.NewFile("graphqlhandler")
		f.ImportNames(importName)
		f.ImportAlias("github.com/sirupsen/logrus", "log")

		f.Comment("graphQLHandler represent the graphQLHandler")
		f.Type().Id("graphQLHandler").Struct(
			jen.Id(useCase).Qual(gomodName+"/domain", useCase),
		)

		f.Comment("NewGraphQLHandler will initialize the graphql endpoint")
		f.Func().Id("NewGraphQLHandler").Params(
			jen.Id("r").Op("*").Qual("net/http", "ServeMux"),
			jen.Id(string(domainName[0])).Qual(gomodName+"/domain", useCase),
		).Block(
			jen.Id("handle").Op(":=&").Id("graphQLHandler").Values(jen.Dict{
				jen.Id(useCase): jen.Id(string(domainName[0])),
			}),
			jen.Line(),
			jen.Id("h").Op(":=").Qual("github.com/graphql-go/handler", "New").Call(jen.Op("&").Qual("github.com/graphql-go/handler", "Config").Values(jen.Dict{
				jen.Id("Schema"):   jen.Id("handle").Dot("schema").Call(),
				jen.Id("Pretty"):   jen.True(),
				jen.Id("GraphiQL"): jen.False(),
			})),
			jen.Line(),
			jen.Id("r").Dot("Handle").Call(jen.Lit("/graphql"), jen.Id("httpHeaderMiddleware").Call(jen.Id("h"))),
		)
		f.Line()
		f.Func().Params(jen.Id("gh").Op("*").Id("graphQLHandler")).Id("schema").Params().Op("*").Qual("github.com/graphql-go/graphql", "Schema").Block(
			jen.Id("rootMutation").Op(":=").Qual(gomodName+"/transport/graphql/mutations", "NewGraphQLMutation").Call(jen.Id("gh").Dot(useCase)),
			jen.Id("rootQuery").Op(":=").Qual(gomodName+"/transport/graphql/queries", "NewGraphQLQuery").Call(jen.Id("gh").Dot(useCase)),
			jen.Line(),
			jen.Id("queryType").Op(":=").Qual("github.com/graphql-go/graphql", "NewObject").Call(
				jen.Qual("github.com/graphql-go/graphql", "ObjectConfig").Values(jen.Dict{
					jen.Id("Name"):   jen.Lit("Query"),
					jen.Id("Fields"): jen.Id("rootQuery").Dot("GetRootQueryFields").Call(),
				}),
			),
			jen.Line(),
			jen.Id("mutationType").Op(":=").Qual("github.com/graphql-go/graphql", "NewObject").Call(
				jen.Qual("github.com/graphql-go/graphql", "ObjectConfig").Values(jen.Dict{
					jen.Id("Name"):   jen.Lit("Mutation"),
					jen.Id("Fields"): jen.Id("rootMutation").Dot("GetRootMutationFields").Call(),
				}),
			),
			jen.Line(),
			jen.Id("schema").Op(",").Err().Op(":=").Qual("github.com/graphql-go/graphql", "NewSchema").Call(
				jen.Qual("github.com/graphql-go/graphql", "SchemaConfig").Values(jen.Dict{
					jen.Id("Query"):    jen.Id("queryType"),
					jen.Id("Mutation"): jen.Id("mutationType"),
				}),
			),
			jen.If(jen.Err().Op("!=").Nil().Block(
				jen.Qual("github.com/sirupsen/logrus", "Printf").Call(jen.Lit("errors: %v"), jen.Err().Dot("Error").Call()),
			)),
			jen.Line(),
			jen.Return(jen.Op("&").Id("schema")),
		)
		f.Line()
		f.Func().Id("httpHeaderMiddleware").Params(jen.Id("next").Op("*").Qual("github.com/graphql-go/handler", "Handler")).Qual("net/http", "Handler").Block(
			jen.Return(
				jen.Qual("net/http", "HandlerFunc").Call(
					jen.Func().Params(jen.Id("w").Qual("net/http", "ResponseWriter"), jen.Id("r").Op("*").Qual("net/http", "Request")).Block(
						jen.Id("ctx").Op(":=").Qual("context", "WithValue").Call(jen.Id("r").Dot("Context").Call(), jen.Lit("example"), jen.Id("r").Dot("Header").Dot("Get").Call(jen.Lit("example"))),
						jen.Line(),
						jen.Id("next").Dot("ContextHandler").Call(jen.Id("ctx"), jen.Id("w"), jen.Id("r")),
					),
				),
			),
		)

		// save file
		fileDir := fmt.Sprintf("%s/index.go", dirName)
		err = f.Save(fileDir)
		if err != nil {
			return err
		}

		return nil
	}
	err = indexGen()
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGrpcTransport(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file       = path.Base(domainFile)
		domainName = strings.TrimSuffix(file, filepath.Ext(file))
		useCase    = parser.Usecase.Name
		handler    = fmt.Sprintf("%sHandler", strings.ToUpper(string(domainName[0]))+domainName[1:])
		f          = jen.NewFile("grpchandler")
		importName = map[string]string{
			gomodName + "/domain":    "domain",
			"google.golang.org/grpc": "grpc",
		}
	)

	f.ImportNames(importName)
	f.ImportAlias(gomodName+"/proto", "pb")

	f.Comment("Grpc" + handler + " represent the grpc handler for " + domainName)
	f.Type().Id("Grpc" + handler).Struct(
		jen.Id(useCase).Qual(gomodName+"/domain", useCase),
	)
	f.Line()
	f.Comment("NewGrpc" + handler + " will initialize the grpc endpoint for " + domainName + " entity")
	f.Func().Id("NewGrpc"+handler).Params(
		jen.Id("gs").Op("*").Qual("google.golang.org/grpc", "Server"),
		jen.Id(string(domainName[0])).Qual(gomodName+"/domain", useCase),
	).Block(
		jen.Id("srv").Op(":=&").Id("Grpc"+handler).Values(jen.Dict{
			jen.Id(useCase): jen.Id(string(domainName[0])),
		}),
		jen.Line(),
		jen.Qual(gomodName+"/proto", "Register"+strings.ToUpper(string(domainName[0]))+domainName[1:]+"ServiceServer").Call(jen.Id("gs"), jen.Id("srv")),
	)
	f.Line()
	for _, i := range parser.Usecase.Method {
		funcName := i.Name + strings.ToUpper(string(domainName[0])) + domainName[1:]
		f.Line()
		f.Comment(funcName + " will handle " + funcName + " request")
		f.Func().
			Params(jen.Id("gh").Op("*").Id("Grpc"+handler)).
			Id(funcName).Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("req").Op("*").Qual(gomodName+"/proto", funcName+"Req"),
		).Call(
			jen.Op("*").Qual(gomodName+"/proto", funcName+"Resp"),
			jen.Error(),
		).Block(
			jen.If(jen.Id("ctx").Op("==").Nil()).Block(jen.Id("ctx").Op("=").Qual("context", "Background").Call()),
			jen.Return(jen.Op("&").Qual(gomodName+"/proto", funcName+"Resp").Values(jen.Dict{}), jen.Nil()),
		)
	}
	fileDir := fmt.Sprintf("%s/%s_handler.go", dirName, domainName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}
	return nil
}
