package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

func (gen *caGen) GenEchoMiddleware(dirName string) error {
	f := jen.NewFile("middleware")
	f.ImportAlias("github.com/sirupsen/logrus", "log")

	f.Comment("EchoMiddleware represent the data-struct for middleware")
	f.Type().Id("EchoMiddleware").Struct()

	f.Comment("InitEchoMiddleware intialize the middleware")
	f.Func().Id("InitEchoMiddleware").Params().Op("*").Id("EchoMiddleware").Block(
		jen.Return(jen.Op("&").Id("EchoMiddleware").Values(jen.Dict{})),
	)

	f.Comment("CORS will handle the CORS middleware")
	f.Func().Params(jen.Id("m").Op("*").Id("EchoMiddleware")).Id("CORS").Params(jen.Id("next").Qual("github.com/labstack/echo", "HandlerFunc")).Qual("github.com/labstack/echo", "HandlerFunc").Block(
		jen.Return(jen.Func().Params(jen.Id("c").Qual("github.com/labstack/echo", "Context")).Error().Block(
			jen.Id("c").Dot("Response").Call().Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Origin"), jen.Lit("*")),
			jen.Return(jen.Id("next").Call(jen.Id("c"))),
		)),
	)

	f.Comment("MiddlewareLogging for logging")
	f.Func().Params(jen.Id("m").Op("*").Id("EchoMiddleware")).Id("MiddlewareLogging").Params(jen.Id("next").Qual("github.com/labstack/echo", "HandlerFunc")).Qual("github.com/labstack/echo", "HandlerFunc").Block(
		jen.Return(jen.Func().Params(jen.Id("c").Qual("github.com/labstack/echo", "Context")).Error().Block(
			jen.Id("makeLogEntry").Call(jen.Id("c")).Dot("Info").Call(jen.Lit("incoming request")),
			jen.Return(jen.Id("next").Call(jen.Id("c"))),
		)),
	)

	f.Func().Id("makeLogEntry").Params(jen.Id("c").Qual("github.com/labstack/echo", "Context")).Op("*").Qual("github.com/sirupsen/logrus", "Entry").Block(
		jen.If(jen.Id("c").Op("==").Nil().Block(
			jen.Return(jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
				jen.Lit("at"): jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
			}))),
		)),
		jen.Return(jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Qual("github.com/sirupsen/logrus", "Fields").Values(jen.Dict{
			jen.Lit("at"):     jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
			jen.Lit("method"): jen.Id("c").Dot("Request").Call().Dot("Method"),
			jen.Lit("uri"):    jen.Id("c").Dot("Request").Call().Dot("URL").Dot("String").Call(),
			jen.Lit("ip"):     jen.Id("c").Dot("Request").Call().Dot("RemoteAddr"),
		}))),
	)

	fileDir := fmt.Sprintf("%s/echo_middleware.go", dirName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGinMiddleware(dirName string) error {
	f := jen.NewFile("middleware")
	f.ImportAlias("github.com/sirupsen/logrus", "log")

	f.Comment("GinMiddleware represent the data-struct for middleware")
	f.Type().Id("GinMiddleware").Struct()

	f.Comment("InitGinMiddleware intialize the middleware")
	f.Func().Id("InitGinMiddleware").Params().Op("*").Id("GinMiddleware").Block(
		jen.Return(jen.Op("&").Id("GinMiddleware").Values(jen.Dict{})),
	)

	f.Comment("CORS will handle the CORS middleware")
	f.Func().Params(jen.Id("m").Op("*").Id("GinMiddleware")).Id("CORS").Params().Qual("github.com/gin-gonic/gin", "HandlerFunc").Block(
		jen.Return(jen.Func().Params(jen.Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context")).Block(
			jen.Id("c").Dot("Writer").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Origin"), jen.Lit("*")),
			jen.Id("c").Dot("Writer").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Credentials"), jen.Lit("true")),
			jen.Id("c").Dot("Writer").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Headers"), jen.Lit("Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")),
			jen.Id("c").Dot("Writer").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Methods"), jen.Lit("POST, OPTIONS, GET, PUT")),
			jen.Line(),
			jen.If(jen.Id("c").Dot("Request").Dot("Method").Op("==").Lit("OPTIONS")).Block(
				jen.Id("c").Dot("AbortWithStatus").Call(jen.Lit(204)),
				jen.Return(),
			),
			f.Line(),
			jen.Id("c").Dot("Next").Call(),
		)),
	)

	f.Comment("MiddlewareLogging for logging")
	f.Func().Params(jen.Id("m").Op("*").Id("GinMiddleware")).Id("MiddlewareLogging").Params().Qual("github.com/gin-gonic/gin", "HandlerFunc").Block(
		jen.Return(
			jen.Qual("github.com/gin-gonic/gin", "LoggerWithFormatter").Call(
				jen.Func().Params(
					jen.Id("param").Qual("github.com/gin-gonic/gin", "LogFormatterParams"),
				).String().Block(
					jen.Return(
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n"),
							jen.Id("param").Dot("ClientIP"),
							jen.Id("param").Dot("TimeStamp").Dot("Format").Call(jen.Qual("time", "RFC1123")),
							jen.Id("param").Dot("Method"),
							jen.Id("param").Dot("Path"),
							jen.Id("param").Dot("Request").Dot("Proto"),
							jen.Id("param").Dot("StatusCode"),
							jen.Id("param").Dot("Latency"),
							jen.Id("param").Dot("Request").Dot("UserAgent").Call(),
							jen.Id("param").Dot("ErrorMessage"),
						),
					),
				),
			),
		),
	)

	fileDir := fmt.Sprintf("%s/gin_middleware.go", dirName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenGorillaMuxMiddleware(dirName string) error {
	f := jen.NewFile("middleware")
	f.ImportAlias("github.com/sirupsen/logrus", "log")

	f.Comment("GorillaMuxMiddleware represent the data-struct for middleware")
	f.Type().Id("GorillaMuxMiddleware").Struct()

	f.Comment("InitGorillaMuxMiddleware intialize the middleware")
	f.Func().Id("InitGorillaMuxMiddleware").Params().Op("*").Id("GorillaMuxMiddleware").Block(
		jen.Return(jen.Op("&").Id("GorillaMuxMiddleware").Values(jen.Dict{})),
	)

	f.Comment("CORS will handle the CORS middleware")
	f.Func().Params(jen.Id("m").Op("*").Id("GorillaMuxMiddleware")).Id("CORS").Params(jen.Id("next").Qual("net/http", "Handler")).Qual("net/http", "Handler").Block(
		jen.Return(jen.Qual("net/http", "HandlerFunc").Call(
			jen.Func().Params(jen.Id("w").Qual("net/http", "ResponseWriter"), jen.Id("r").Op("*").Qual("net/http", "Request")).Block(
				jen.Id("w").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Origin"), jen.Lit("*")),
				jen.Id("w").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Methods"), jen.Lit("POST, GET, OPTIONS, PUT, DELETE, PATCH")),
				jen.Id("w").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Headers"), jen.Lit("Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")),
				jen.Id("next").Dot("ServeHTTP").Call(jen.Id("w"), jen.Id("r")),
			),
		)),
	)

	f.Comment("MiddlewareLogging for logging")
	f.Func().Params(jen.Id("m").Op("*").Id("GorillaMuxMiddleware")).Id("MiddlewareLogging").Params(jen.Id("next").Qual("net/http", "Handler")).Qual("net/http", "Handler").Block(
		jen.Return(jen.Qual("net/http", "HandlerFunc").Call(
			jen.Func().Params(jen.Id("w").Qual("net/http", "ResponseWriter"), jen.Id("r").Op("*").Qual("net/http", "Request")).Block(
				jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Id("log").Dot("Fields").Values(jen.Dict{
					jen.Lit("at"):     jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
					jen.Lit("method"): jen.Id("r").Dot("Method"),
					jen.Lit("uri"):    jen.Id("r").Dot("RequestURI"),
					jen.Lit("ip"):     jen.Id("r").Dot("RemoteAddr"),
				})).Dot("Info").Call(jen.Lit("incoming request")),
				jen.Line(),
				jen.Comment("// Call the next handler, which can be another middleware in the chain, or the final handler."),
				jen.Id("next").Dot("ServeHTTP").Call(jen.Id("w"), jen.Id("r")),
			),
		)),
	)

	fileDir := fmt.Sprintf("%s/gorilla_mux_middleware.go", dirName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func (gen *caGen) GenNetHTTPMiddleware(dirName string) error {
	f := jen.NewFile("middleware")
	f.ImportAlias("github.com/sirupsen/logrus", "log")

	f.Comment("NetHTTPMiddleware represent the data-struct for middleware")
	f.Type().Id("NetHTTPMiddleware").Struct()

	f.Comment("InitNetHTTPMiddleware intialize the middleware")
	f.Func().Id("InitNetHTTPMiddleware").Params().Op("*").Id("NetHTTPMiddleware").Block(
		jen.Return(jen.Op("&").Id("NetHTTPMiddleware").Values(jen.Dict{})),
	)

	f.Comment("CORS will handle the CORS middleware")
	f.Func().Params(jen.Id("m").Op("*").Id("NetHTTPMiddleware")).Id("CORS").Params(jen.Id("next").Qual("net/http", "Handler")).Qual("net/http", "Handler").Block(
		jen.Return(jen.Qual("net/http", "HandlerFunc").Call(
			jen.Func().Params(jen.Id("w").Qual("net/http", "ResponseWriter"), jen.Id("r").Op("*").Qual("net/http", "Request")).Block(
				jen.Id("w").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Origin"), jen.Lit("*")),
				jen.Id("w").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Methods"), jen.Lit("POST, GET, OPTIONS, PUT, DELETE, PATCH")),
				jen.Id("w").Dot("Header").Call().Dot("Set").Call(jen.Lit("Access-Control-Allow-Headers"), jen.Lit("Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")),
				jen.Id("next").Dot("ServeHTTP").Call(jen.Id("w"), jen.Id("r")),
			),
		)),
	)

	f.Comment("MiddlewareLogging for logging")
	f.Func().Params(jen.Id("m").Op("*").Id("NetHTTPMiddleware")).Id("MiddlewareLogging").Params(jen.Id("next").Qual("net/http", "Handler")).Qual("net/http", "Handler").Block(
		jen.Return(jen.Qual("net/http", "HandlerFunc").Call(
			jen.Func().Params(jen.Id("w").Qual("net/http", "ResponseWriter"), jen.Id("r").Op("*").Qual("net/http", "Request")).Block(
				jen.Qual("github.com/sirupsen/logrus", "WithFields").Call(jen.Id("log").Dot("Fields").Values(jen.Dict{
					jen.Lit("at"):     jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Lit("2006-01-02 15:04:05")),
					jen.Lit("method"): jen.Id("r").Dot("Method"),
					jen.Lit("uri"):    jen.Id("r").Dot("RequestURI"),
					jen.Lit("ip"):     jen.Id("r").Dot("RemoteAddr"),
				})).Dot("Info").Call(jen.Lit("incoming request")),
				jen.Line(),
				jen.Comment("// Call the next handler, which can be another middleware in the chain, or the final handler."),
				jen.Id("next").Dot("ServeHTTP").Call(jen.Id("w"), jen.Id("r")),
			),
		)),
	)

	fileDir := fmt.Sprintf("%s/net_http_middleware.go", dirName)
	err := f.Save(fileDir)
	if err != nil {
		return err
	}

	return nil
}
