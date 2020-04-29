package generator_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/wicaker/cacli/domain"
	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/generator"
)

const (
	expected_echo_server = `package server

import (
	"github.com/example/exampleserver/middleware"
	"github.com/example/exampleserver/repository"
	"github.com/example/exampleserver/transport/rest"
	"github.com/example/exampleserver/usecase"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// EchoServer /
func EchoServer(db *mongo.Database) *echo.Echo {
	r := echo.New()
	middl := middleware.InitEchoMiddleware()
	r.Use(middl.MiddlewareLogging)
	r.Use(middl.CORS)

	timeoutContext := time.Duration(2) * time.Second

	examplerepository := repository.NewMongodExampleRepository(db)
	exampleusecase := usecase.NewExampleUsecase(examplerepository, timeoutContext)
	rest.NewExampleHandler(r, exampleusecase)

	return r
}
`
	expected_gin_server = `package server

import (
	"github.com/example/exampleserver/middleware"
	"github.com/example/exampleserver/repository"
	"github.com/example/exampleserver/transport/rest"
	"github.com/example/exampleserver/usecase"
	"github.com/gin-gonic/gin"
	pg "github.com/go-pg/pg/v9"
	"os"
	"time"
)

// GinServer /
func GinServer(db *pg.DB) *gin.Engine {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	middl := middleware.InitGinMiddleware()
	r.Use(middl.MiddlewareLogging())
	r.Use(middl.CORS())
	r.Use(gin.Recovery())

	timeoutContext := time.Duration(2) * time.Second

	examplerepository := repository.NewGopgExampleRepository(db)
	exampleusecase := usecase.NewExampleUsecase(examplerepository, timeoutContext)
	rest.NewExampleHandler(r, exampleusecase)

	return r
}
`
	expected_gorilla_mux_server = `package server

import (
	"database/sql"
	"github.com/example/exampleserver/middleware"
	"github.com/example/exampleserver/repository"
	"github.com/example/exampleserver/transport/rest"
	"github.com/example/exampleserver/usecase"
	"github.com/gorilla/mux"
	"time"
)

// GorillaMuxServer /
func GorillaMuxServer(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	middl := middleware.InitGorillaMuxMiddleware()
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middl.MiddlewareLogging)

	timeoutContext := time.Duration(2) * time.Second

	examplerepository := repository.NewSQLExampleRepository(db)
	exampleusecase := usecase.NewExampleUsecase(examplerepository, timeoutContext)
	rest.NewExampleHandler(r, exampleusecase)

	return r
}
`
	expected_net_http_mux_server = `package server

import (
	"context"
	"github.com/example/exampleserver/middleware"
	"github.com/example/exampleserver/repository"
	"github.com/example/exampleserver/transport/rest"
	"github.com/example/exampleserver/usecase"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MuxServer /
func MuxServer(db *gorm.DB) http.Handler {
	r := http.NewServeMux()
	middl := middleware.InitNetHTTPMiddleware()
	var handler http.Handler = r
	handler = middl.MiddlewareLogging(handler)
	handler = middl.CORS(handler)
	handler = parseURL(handler)

	timeoutContext := time.Duration(2) * time.Second

	examplerepository := repository.NewGormExampleRepository(db)
	exampleusecase := usecase.NewExampleUsecase(examplerepository, timeoutContext)
	rest.NewExampleHandler(r, exampleusecase)

	return handler
}

func parseURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		ctx = context.WithValue(ctx, "uri", r.RequestURI)

		s := strings.Split(r.RequestURI, "/")
		if s[1] == "example" {
			r.URL = &url.URL{Path: "/example"}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
`
	expected_graphql_server = `package server

import (
	"database/sql"
	"github.com/example/exampleserver/middleware"
	"github.com/example/exampleserver/repository"
	graphqlhandler "github.com/example/exampleserver/transport/graphql"
	"github.com/example/exampleserver/usecase"
	"net/http"
	"time"
)

// GraphQLServer /
func GraphQLServer(db *sql.DB) http.Handler {
	r := http.NewServeMux()
	middl := middleware.InitNetHTTPMiddleware()
	var handler http.Handler = r
	handler = middl.MiddlewareLogging(handler)
	handler = middl.CORS(handler)

	timeoutContext := time.Duration(2) * time.Second

	examplerepository := repository.NewSQLExampleRepository(db)
	exampleusecase := usecase.NewExampleUsecase(examplerepository, timeoutContext)
	graphqlhandler.NewGraphQLHandler(r, exampleusecase)

	return handler
}
`
	expected_grpc_server = `package server

import (
	"github.com/example/exampleserver/repository"
	grpcHandler "github.com/example/exampleserver/transport/grpc"
	"github.com/example/exampleserver/usecase"
	pg "github.com/go-pg/pg/v9"
	"google.golang.org/grpc"
	"time"
)

// GRPCServer /
func GRPCServer(db *pg.DB) *grpc.Server {
	// slice of gRPC options
	// Here we can configure things like TLS
	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	timeoutContext := time.Duration(2) * time.Second

	examplerepository := repository.NewGopgExampleRepository(db)
	exampleusecase := usecase.NewExampleUsecase(examplerepository, timeoutContext)
	grpcHandler.NewGrpcExampleHandler(s, exampleusecase)

	return s
}
`
)

func TestGenerateEchoServer(t *testing.T) {
	var (
		serviceName = "test_echo_example_server"
		dirLayer1   = "server"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampleserver"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer1)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an echo_server.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirLayer1)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/usecase")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/repository")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport/rest")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate echo_server.go file
		gen := generator.NewGeneratorService()
		err = gen.GenMongodRepository(serviceName+"/repository", domainFile, gomodName, domain.MockParser)
		err = gen.GenUsecase(serviceName+"/usecase", domainFile, gomodName, domain.MockParser)
		err = gen.GenEchoTransport(serviceName+"/transport/rest", domainFile, gomodName, domain.MockParser)
		err = gen.GenEchoServer(dirName, serviceName, domain.Mongod, gomodName, domain.MockParser)
		resEcho, err := newFs.FindFile(dirName + "/echo_server.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resEcho)

		data, err := ioutil.ReadFile(dirName + "/echo_server.go")
		if err != nil {
			log.Error("File reading error", err)
			// remove directory of service
			err = newFs.RemoveDir(serviceName)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			os.Exit(1)
		}
		assert.Equal(t, expected_echo_server, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate echo_server file
		gen := generator.NewGeneratorService()
		err := gen.GenEchoServer(dirName, serviceName, domain.Sqlx, gomodName, domain.MockParser)
		assert.Error(t, err)
	})
}
func TestGenerateGinServer(t *testing.T) {
	var (
		serviceName = "test_gin_example_server"
		dirLayer1   = "server"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampleserver"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer1)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an gin_server.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirLayer1)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/usecase")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/repository")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport/rest")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate gin_server.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGopgRepository(serviceName+"/repository", domainFile, gomodName, domain.MockParser)
		err = gen.GenUsecase(serviceName+"/usecase", domainFile, gomodName, domain.MockParser)
		err = gen.GenGinTransport(serviceName+"/transport/rest", domainFile, gomodName, domain.MockParser)
		err = gen.GenGinServer(dirName, serviceName, domain.GoPg, gomodName, domain.MockParser)
		resGin, err := newFs.FindFile(dirName + "/gin_server.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGin)

		data, err := ioutil.ReadFile(dirName + "/gin_server.go")
		if err != nil {
			log.Error("File reading error", err)
			// remove directory of service
			err = newFs.RemoveDir(serviceName)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			os.Exit(1)
		}
		assert.Equal(t, expected_gin_server, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate gin_server file
		gen := generator.NewGeneratorService()
		err := gen.GenGinServer(dirName, serviceName, domain.Sqlx, gomodName, domain.MockParser)
		assert.Error(t, err)
	})
}
func TestGenerateGorillaMuxServer(t *testing.T) {
	var (
		serviceName = "test_gorilla_mux_example_server"
		dirLayer1   = "server"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampleserver"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer1)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an gorilla_mux_server.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirLayer1)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/usecase")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/repository")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport/rest")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate gorilla_mux_server.go file
		gen := generator.NewGeneratorService()
		err = gen.GenSQLRepository(serviceName+"/repository", domainFile, gomodName, domain.MockParser)
		err = gen.GenUsecase(serviceName+"/usecase", domainFile, gomodName, domain.MockParser)
		err = gen.GenGorillaMuxTransport(serviceName+"/transport/rest", domainFile, gomodName, domain.MockParser)
		err = gen.GenGorillaMuxServer(dirName, serviceName, domain.SQL, gomodName, domain.MockParser)
		resGorillaMux, err := newFs.FindFile(dirName + "/gorilla_mux_server.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGorillaMux)

		data, err := ioutil.ReadFile(dirName + "/gorilla_mux_server.go")
		if err != nil {
			log.Error("File reading error", err)
			// remove directory of service
			err = newFs.RemoveDir(serviceName)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			os.Exit(1)
		}
		assert.Equal(t, expected_gorilla_mux_server, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate gorilla_mux_server file
		gen := generator.NewGeneratorService()
		err := gen.GenGinServer(dirName, serviceName, domain.Sqlx, gomodName, domain.MockParser)
		assert.Error(t, err)
	})
}
func TestGenNetHTTPMuxServer(t *testing.T) {
	var (
		serviceName = "test_net_http_mux_example_server"
		dirLayer1   = "server"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampleserver"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer1)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an net_http_mux_server.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirLayer1)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/usecase")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/repository")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport/rest")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate net_http_mux_server.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGormRepository(serviceName+"/repository", domainFile, gomodName, domain.MockParser)
		err = gen.GenUsecase(serviceName+"/usecase", domainFile, gomodName, domain.MockParser)
		err = gen.GenNetHTTPTransport(serviceName+"/transport/rest", domainFile, gomodName, domain.MockParser)
		err = gen.GenNetHTTPMuxServer(dirName, serviceName, domain.Gorm, gomodName, domain.MockParser)
		resNetHTTP, err := newFs.FindFile(dirName + "/net_http_mux_server.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resNetHTTP)

		data, err := ioutil.ReadFile(dirName + "/net_http_mux_server.go")
		if err != nil {
			log.Error("File reading error", err)
			// remove directory of service
			err = newFs.RemoveDir(serviceName)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			os.Exit(1)
		}
		assert.Equal(t, expected_net_http_mux_server, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate net_http_mux_server file
		gen := generator.NewGeneratorService()
		err := gen.GenGinServer(dirName, serviceName, domain.Sqlx, gomodName, domain.MockParser)
		assert.Error(t, err)
	})
}
func TestGenerateGraphqlServer(t *testing.T) {
	var (
		serviceName = "test_graphql_example_server"
		dirLayer1   = "server"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampleserver"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer1)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an graphql_server.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirLayer1)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/usecase")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/repository")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport/graphql")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate graphql_server.go file
		gen := generator.NewGeneratorService()
		err = gen.GenSQLRepository(serviceName+"/repository", domainFile, gomodName, domain.MockParser)
		err = gen.GenUsecase(serviceName+"/usecase", domainFile, gomodName, domain.MockParser)
		err = gen.GenGraphqlTransport(serviceName+"/transport/graphql", domainFile, gomodName, domain.MockParser)
		err = gen.GenGraphqlServer(dirName, serviceName, domain.SQL, gomodName, domain.MockParser)
		resGraphql, err := newFs.FindFile(dirName + "/graphql_server.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGraphql)

		data, err := ioutil.ReadFile(dirName + "/graphql_server.go")
		if err != nil {
			log.Error("File reading error", err)
			// remove directory of service
			err = newFs.RemoveDir(serviceName)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			os.Exit(1)
		}
		assert.Equal(t, expected_graphql_server, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate graphql_server file
		gen := generator.NewGeneratorService()
		err := gen.GenGinServer(dirName, serviceName, domain.Sqlx, gomodName, domain.MockParser)
		assert.Error(t, err)
	})
}
func TestGenerateGrpcServer(t *testing.T) {
	var (
		serviceName = "test_grpc_example_server"
		dirLayer1   = "server"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampleserver"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer1)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an grpc_server.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirLayer1)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/usecase")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/repository")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/transport/grpc")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate grpc_server.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGopgRepository(serviceName+"/repository", domainFile, gomodName, domain.MockParser)
		err = gen.GenUsecase(serviceName+"/usecase", domainFile, gomodName, domain.MockParser)
		err = gen.GenGrpcTransport(serviceName+"/transport/grpc", domainFile, gomodName, domain.MockParser)
		err = gen.GenGrpcServer(dirName, serviceName, domain.GoPg, gomodName, domain.MockParser)
		resGrpc, err := newFs.FindFile(dirName + "/grpc_server.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGrpc)

		data, err := ioutil.ReadFile(dirName + "/grpc_server.go")
		if err != nil {
			log.Error("File reading error", err)
			// remove directory of service
			err = newFs.RemoveDir(serviceName)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			os.Exit(1)
		}
		assert.Equal(t, expected_grpc_server, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate grpc_server file
		gen := generator.NewGeneratorService()
		err := gen.GenGinServer(dirName, serviceName, domain.Sqlx, gomodName, domain.MockParser)
		assert.Error(t, err)
	})
}
