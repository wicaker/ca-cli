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
	expected_echo_example_transport = `package rest

import (
	"context"
	"github.com/example/exampletranposport/domain"
	"github.com/labstack/echo"
	"net/http"
)

type exampleHandler struct {
	ExampleUsecase domain.ExampleUsecase
}

// NewExampleHandler will initialize the example endpoint
func NewExampleHandler(e *echo.Echo, u domain.ExampleUsecase) {
	handler := &exampleHandler{ExampleUsecase: u}
	e.GET("/example/fetch", handler.FetchHandler)
	e.GET("/example/getbyid", handler.GetByIDHandler)
	e.GET("/example/store", handler.StoreHandler)
	e.GET("/example/update", handler.UpdateHandler)
	e.GET("/example/delete", handler.DeleteHandler)
}

func (eh *exampleHandler) FetchHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "FetchHandler"
	return c.JSON(http.StatusOK, domain.ResponseData)
}

func (eh *exampleHandler) GetByIDHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "GetByIDHandler"
	return c.JSON(http.StatusOK, domain.ResponseData)
}

func (eh *exampleHandler) StoreHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "StoreHandler"
	return c.JSON(http.StatusOK, domain.ResponseData)
}

func (eh *exampleHandler) UpdateHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "UpdateHandler"
	return c.JSON(http.StatusOK, domain.ResponseData)
}

func (eh *exampleHandler) DeleteHandler(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "DeleteHandler"
	return c.JSON(http.StatusOK, domain.ResponseData)
}
`

	expected_gin_example_transport = `package rest

import (
	"context"
	"github.com/example/exampletranposport/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type exampleHandler struct {
	ExampleUsecase domain.ExampleUsecase
}

// NewExampleHandler will initialize the example endpoint
func NewExampleHandler(r *gin.Engine, u domain.ExampleUsecase) {
	handler := &exampleHandler{ExampleUsecase: u}
	r.GET("/example/fetch", handler.FetchHandler)
	r.GET("/example/getbyid", handler.GetByIDHandler)
	r.GET("/example/store", handler.StoreHandler)
	r.GET("/example/update", handler.UpdateHandler)
	r.GET("/example/delete", handler.DeleteHandler)
}

func (eh *exampleHandler) FetchHandler(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "FetchHandler"
	c.JSON(http.StatusOK, domain.ResponseData)
	return
}

func (eh *exampleHandler) GetByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "GetByIDHandler"
	c.JSON(http.StatusOK, domain.ResponseData)
	return
}

func (eh *exampleHandler) StoreHandler(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "StoreHandler"
	c.JSON(http.StatusOK, domain.ResponseData)
	return
}

func (eh *exampleHandler) UpdateHandler(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "UpdateHandler"
	c.JSON(http.StatusOK, domain.ResponseData)
	return
}

func (eh *exampleHandler) DeleteHandler(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "DeleteHandler"
	c.JSON(http.StatusOK, domain.ResponseData)
	return
}
`

	expected_gorilla_mux_example_transport = `package rest

import (
	"context"
	"github.com/example/exampletranposport/domain"
	"github.com/gorilla/mux"
	json "github.com/json-iterator/go"
	"net/http"
)

type exampleHandler struct {
	ExampleUsecase domain.ExampleUsecase
}

// NewExampleHandler will initialize the example endpoint
func NewExampleHandler(r *mux.Router, u domain.ExampleUsecase) {
	handler := &exampleHandler{ExampleUsecase: u}
	r.HandleFunc("/example/fetch", handler.FetchHandler).Methods("GET")
	r.HandleFunc("/example/getbyid", handler.GetByIDHandler).Methods("GET")
	r.HandleFunc("/example/store", handler.StoreHandler).Methods("GET")
	r.HandleFunc("/example/update", handler.UpdateHandler).Methods("GET")
	r.HandleFunc("/example/delete", handler.DeleteHandler).Methods("GET")
}

func (eh *exampleHandler) FetchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "FetchHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}

func (eh *exampleHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "GetByIDHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}

func (eh *exampleHandler) StoreHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "StoreHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}

func (eh *exampleHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "UpdateHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}

func (eh *exampleHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "DeleteHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}
`

	expected_net_http_mux_example_transport = `package rest

import (
	"context"
	"github.com/example/exampletranposport/domain"
	json "github.com/json-iterator/go"
	"net/http"
	"path"
	"strings"
)

type exampleHandler struct {
	ExampleUsecase domain.ExampleUsecase
}

// NewExampleHandler will initialize the example endpoint
func NewExampleHandler(r *http.ServeMux, u domain.ExampleUsecase) {
	handler := &exampleHandler{ExampleUsecase: u}
	r.Handle("/example", handler)
}

func (eh *exampleHandler) FetchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "FetchHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}

func (eh *exampleHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "GetByIDHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}

func (eh *exampleHandler) StoreHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "StoreHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}

func (eh *exampleHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "UpdateHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}

func (eh *exampleHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	domain.ResponseData = make(map[string]interface{})
	domain.ResponseData["message"] = "DeleteHandler"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ResponseData)
	return
}

func (eh *exampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	uri := ctx.Value("uri")
	s := strings.Split(r.RequestURI, "/")

	if uri == "/example" {
		if r.Method == http.MethodGet {
			eh.FetchHandler(w, r)
			return
		}
	}

	// if the path contain params
	if len(s) == 3 {
		id, err := strconv.Atoi(path.Base(r.RequestURI))
		if err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(&domain.ResponseError{Message: "Method Not Allowed"})
			return
		}
		if r.Method == http.MethodGet {
			eh.FetchHandler(w, r, id)
			return
		}
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(&domain.ResponseError{Message: "Method Not Allowed"})
	return
}
`
	expected_graphql_example_types = `package types

import "github.com/graphql-go/graphql"

// ExampleType is the GraphQL schema for the example type.
var ExampleType = graphql.NewObject(graphql.ObjectConfig{
	Fields: graphql.Fields{"id": &graphql.Field{Type: graphql.ID}},
	Name:   "Example",
})
`
	expected_graphql_mutations = `package mutations

import (
	"github.com/example/exampletranposport/domain"
	"github.com/graphql-go/graphql"
)

// GraphQLMutation represent the graphQLMutation
type GraphQLMutation struct {
	ExampleUsecase domain.ExampleUsecase
}

// NewGraphQLMutation will initialize mutations
func NewGraphQLMutation(e domain.ExampleUsecase) *GraphQLMutation {
	return &GraphQLMutation{ExampleUsecase: e}
}

// GetRootMutationFields returns all the available mutations.
func (gm *GraphQLMutation) GetRootMutationFields() graphql.Fields {
	return graphql.Fields{
		"exampleDelete":  gm.DeleteExampleMutation(),
		"exampleFetch":   gm.FetchExampleMutation(),
		"exampleGetByID": gm.GetByIDExampleMutation(),
		"exampleStore":   gm.StoreExampleMutation(),
		"exampleUpdate":  gm.UpdateExampleMutation(),
	}
}
`
	expected_graphql_example_mutations = `package mutations

import (
	"context"
	"github.com/example/exampletranposport/transport/graphql/types"
	"github.com/graphql-go/graphql"
)

// FetchExampleMutation /.
func (gm *GraphQLMutation) FetchExampleMutation() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}

// GetByIDExampleMutation /.
func (gm *GraphQLMutation) GetByIDExampleMutation() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}

// StoreExampleMutation /.
func (gm *GraphQLMutation) StoreExampleMutation() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}

// UpdateExampleMutation /.
func (gm *GraphQLMutation) UpdateExampleMutation() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}

// DeleteExampleMutation /.
func (gm *GraphQLMutation) DeleteExampleMutation() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}
`
	expected_graphql_queries = `package queries

import (
	"github.com/example/exampletranposport/domain"
	"github.com/graphql-go/graphql"
)

// GraphQLQuery represent the GraphQLQuery
type GraphQLQuery struct {
	ExampleUsecase domain.ExampleUsecase
}

// NewGraphQLQuery will initialize queries
func NewGraphQLQuery(e domain.ExampleUsecase) *GraphQLQuery {
	return &GraphQLQuery{ExampleUsecase: e}
}

// GetRootQueryFields returns all the available queries.
func (gq *GraphQLQuery) GetRootQueryFields() graphql.Fields {
	return graphql.Fields{
		"exampleDelete":  gq.DeleteExampleQuery(),
		"exampleFetch":   gq.FetchExampleQuery(),
		"exampleGetByID": gq.GetByIDExampleQuery(),
		"exampleStore":   gq.StoreExampleQuery(),
		"exampleUpdate":  gq.UpdateExampleQuery(),
	}
}
`
	expected_graphql_example_queries = `package queries

import (
	"context"
	"github.com/example/exampletranposport/transport/graphql/types"
	"github.com/graphql-go/graphql"
)

// FetchExampleQuery /.
func (gq *GraphQLQuery) FetchExampleQuery() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.ID}},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}

// GetByIDExampleQuery /.
func (gq *GraphQLQuery) GetByIDExampleQuery() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.ID}},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}

// StoreExampleQuery /.
func (gq *GraphQLQuery) StoreExampleQuery() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.ID}},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}

// UpdateExampleQuery /.
func (gq *GraphQLQuery) UpdateExampleQuery() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.ID}},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}

// DeleteExampleQuery /.
func (gq *GraphQLQuery) DeleteExampleQuery() *graphql.Field {
	return &graphql.Field{
		Args:        graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.ID}},
		Description: "Example",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}
			return nil, nil
		},
		Type: types.ExampleType,
	}
}
`
	expected_graphql_index = `package graphqlhandler

import (
	"context"
	"github.com/example/exampletranposport/domain"
	"github.com/example/exampletranposport/transport/graphql/mutations"
	"github.com/example/exampletranposport/transport/graphql/queries"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// graphQLHandler represent the graphQLHandler
type graphQLHandler struct {
	ExampleUsecase domain.ExampleUsecase
}

// NewGraphQLHandler will initialize the graphql endpoint
func NewGraphQLHandler(r *http.ServeMux, e domain.ExampleUsecase) {
	handle := &graphQLHandler{ExampleUsecase: e}

	h := handler.New(&handler.Config{
		GraphiQL: false,
		Pretty:   true,
		Schema:   handle.schema(),
	})

	r.Handle("/graphql", httpHeaderMiddleware(h))
}

func (gh *graphQLHandler) schema() *graphql.Schema {
	rootMutation := mutations.NewGraphQLMutation(gh.ExampleUsecase)
	rootQuery := queries.NewGraphQLQuery(gh.ExampleUsecase)

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Fields: rootQuery.GetRootQueryFields(),
		Name:   "Query",
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Fields: rootMutation.GetRootMutationFields(),
		Name:   "Mutation",
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Mutation: mutationType,
		Query:    queryType,
	})
	if err != nil {
		log.Printf("errors: %v", err.Error())
	}

	return &schema
}

func httpHeaderMiddleware(next *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "example", r.Header.Get("example"))

		next.ContextHandler(ctx, w, r)
	})
}
`
	expected_grpc_example_transport = `package grpchandler

import (
	"context"
	"github.com/example/exampletranposport/domain"
	pb "github.com/example/exampletranposport/proto"
	"google.golang.org/grpc"
)

// GrpcExampleHandler represent the grpc handler for example
type GrpcExampleHandler struct {
	ExampleUsecase domain.ExampleUsecase
}

// NewGrpcExampleHandler will initialize the grpc endpoint for example entity
func NewGrpcExampleHandler(gs *grpc.Server, e domain.ExampleUsecase) {
	srv := &GrpcExampleHandler{ExampleUsecase: e}

	pb.RegisterExampleServiceServer(gs, srv)
}

// FetchExample will handle FetchExample request
func (gh *GrpcExampleHandler) FetchExample(ctx context.Context, req *pb.FetchExampleReq) (*pb.FetchExampleResp, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return &pb.FetchExampleResp{}, nil
}

// GetByIDExample will handle GetByIDExample request
func (gh *GrpcExampleHandler) GetByIDExample(ctx context.Context, req *pb.GetByIDExampleReq) (*pb.GetByIDExampleResp, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return &pb.GetByIDExampleResp{}, nil
}

// StoreExample will handle StoreExample request
func (gh *GrpcExampleHandler) StoreExample(ctx context.Context, req *pb.StoreExampleReq) (*pb.StoreExampleResp, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return &pb.StoreExampleResp{}, nil
}

// UpdateExample will handle UpdateExample request
func (gh *GrpcExampleHandler) UpdateExample(ctx context.Context, req *pb.UpdateExampleReq) (*pb.UpdateExampleResp, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return &pb.UpdateExampleResp{}, nil
}

// DeleteExample will handle DeleteExample request
func (gh *GrpcExampleHandler) DeleteExample(ctx context.Context, req *pb.DeleteExampleReq) (*pb.DeleteExampleResp, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return &pb.DeleteExampleResp{}, nil
}
`
)

func TestGenerateEchoTransport(t *testing.T) {
	var (
		serviceName = "test_echo_example_transport"
		dirLayer1   = "transport"
		dirLayer2   = "rest"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampletranposport"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirLayer1, dirLayer2)
		parser      = &domain.Parser{
			Usecase: domain.Usecase{
				Name: "ExampleUsecase",
				Method: []domain.Method{
					domain.Method{
						Name: "Fetch",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "[]*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "GetByID",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Store",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Update",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Delete",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "error"},
						},
					},
				},
			},
		}
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate an example_handler.go file", func(t *testing.T) {
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

		err = newFs.CreateDir(serviceName + "/" + dirLayer1 + "/" + dirLayer2)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate example_handler.go file
		gen := generator.NewGeneratorService()
		err = gen.GenEchoTransport(dirName, domainFile, gomodName, parser)
		resGopg, err := newFs.FindFile(dirName + "/example_handler.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/example_handler.go")
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
		assert.Equal(t, expected_echo_example_transport, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_handler file
		gen := generator.NewGeneratorService()
		err := gen.GenEchoTransport(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}

func TestGenerateGinTransport(t *testing.T) {
	var (
		serviceName = "test_gin_example_transport"
		dirLayer1   = "transport"
		dirLayer2   = "rest"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampletranposport"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirLayer1, dirLayer2)
		parser      = &domain.Parser{
			Usecase: domain.Usecase{
				Name: "ExampleUsecase",
				Method: []domain.Method{
					domain.Method{
						Name: "Fetch",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "[]*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "GetByID",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Store",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Update",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Delete",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "error"},
						},
					},
				},
			},
		}
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate an example_handler.go file", func(t *testing.T) {
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

		err = newFs.CreateDir(serviceName + "/" + dirLayer1 + "/" + dirLayer2)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate example_handler.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGinTransport(dirName, domainFile, gomodName, parser)
		resGopg, err := newFs.FindFile(dirName + "/example_handler.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/example_handler.go")
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
		assert.Equal(t, expected_gin_example_transport, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_handler file
		gen := generator.NewGeneratorService()
		err := gen.GenGinTransport(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}

func TestGenerateGorillaMuxTransport(t *testing.T) {
	var (
		serviceName = "test_gorilla_mux_example_transport"
		dirLayer1   = "transport"
		dirLayer2   = "rest"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampletranposport"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirLayer1, dirLayer2)
		parser      = &domain.Parser{
			Usecase: domain.Usecase{
				Name: "ExampleUsecase",
				Method: []domain.Method{
					domain.Method{
						Name: "Fetch",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "[]*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "GetByID",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Store",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Update",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Delete",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "error"},
						},
					},
				},
			},
		}
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate an example_handler.go file", func(t *testing.T) {
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

		err = newFs.CreateDir(serviceName + "/" + dirLayer1 + "/" + dirLayer2)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate example_handler.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGorillaMuxTransport(dirName, domainFile, gomodName, parser)
		resGopg, err := newFs.FindFile(dirName + "/example_handler.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/example_handler.go")
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
		assert.Equal(t, expected_gorilla_mux_example_transport, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_handler file
		gen := generator.NewGeneratorService()
		err := gen.GenGorillaMuxTransport(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}

func TestGenerateNetHTTPTransport(t *testing.T) {
	var (
		serviceName = "test_net_http_mux_example_transport"
		dirLayer1   = "transport"
		dirLayer2   = "rest"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampletranposport"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirLayer1, dirLayer2)
		parser      = &domain.Parser{
			Usecase: domain.Usecase{
				Name: "ExampleUsecase",
				Method: []domain.Method{
					domain.Method{
						Name: "Fetch",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "[]*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "GetByID",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Store",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Update",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Delete",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "error"},
						},
					},
				},
			},
		}
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate an example_handler.go file", func(t *testing.T) {
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

		err = newFs.CreateDir(serviceName + "/" + dirLayer1 + "/" + dirLayer2)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate example_handler.go file
		gen := generator.NewGeneratorService()
		err = gen.GenNetHTTPTransport(dirName, domainFile, gomodName, parser)
		resGopg, err := newFs.FindFile(dirName + "/example_handler.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/example_handler.go")
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
		assert.Equal(t, expected_net_http_mux_example_transport, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_handler file
		gen := generator.NewGeneratorService()
		err := gen.GenNetHTTPTransport(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}

func TestGenerateGraphqlTransport(t *testing.T) {
	var (
		serviceName = "test_graphql_example_transport"
		dirLayer1   = "transport"
		dirLayer2   = "graphql"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampletranposport"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirLayer1, dirLayer2)
		parser      = &domain.Parser{
			Usecase: domain.Usecase{
				Name: "ExampleUsecase",
				Method: []domain.Method{
					domain.Method{
						Name: "Fetch",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "[]*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "GetByID",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Store",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Update",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Delete",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "error"},
						},
					},
				},
			},
		}
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate graphql transport", func(t *testing.T) {
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

		err = newFs.CreateDir(serviceName + "/" + dirLayer1 + "/" + dirLayer2)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate graphql transport
		gen := generator.NewGeneratorService()
		err = gen.GenGraphqlTransport(dirName, domainFile, gomodName, parser)

		// types directory
		resGopg, err := newFs.FindFile(dirName + "/types/example.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/types/example.go")
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
		assert.Equal(t, expected_graphql_example_types, string(data))

		// mutations directory
		// mutations.go
		resMutations, err := newFs.FindFile(dirName + "/mutations/mutations.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resMutations)

		dataMutations, err := ioutil.ReadFile(dirName + "/mutations/mutations.go")
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
		assert.Equal(t, expected_graphql_mutations, string(dataMutations))

		// example.go
		resMutExample, err := newFs.FindFile(dirName + "/mutations/example.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resMutExample)

		dataMutExample, err := ioutil.ReadFile(dirName + "/mutations/example.go")
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
		assert.Equal(t, expected_graphql_example_mutations, string(dataMutExample))

		// queries directory
		// queries.go
		resQueries, err := newFs.FindFile(dirName + "/queries/queries.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resQueries)

		dataQueries, err := ioutil.ReadFile(dirName + "/queries/queries.go")
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
		assert.Equal(t, expected_graphql_queries, string(dataQueries))

		// example.go
		resQuExample, err := newFs.FindFile(dirName + "/queries/example.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resQuExample)

		dataQuExample, err := ioutil.ReadFile(dirName + "/queries/example.go")
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
		assert.Equal(t, expected_graphql_example_queries, string(dataQuExample))

		// index file
		resIndex, err := newFs.FindFile(dirName + "/index.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resIndex)

		dataIndex, err := ioutil.ReadFile(dirName + "/index.go")
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
		assert.Equal(t, expected_graphql_index, string(dataIndex))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_handler file
		gen := generator.NewGeneratorService()
		err := gen.GenGraphqlTransport(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}

func TestGenerateGrpcTransport(t *testing.T) {
	var (
		serviceName = "test_grpc_example_transport"
		dirLayer1   = "transport"
		dirLayer2   = "grpc"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampletranposport"
		dirName     = fmt.Sprintf("%s/%s/%s", serviceName, dirLayer1, dirLayer2)
		parser      = &domain.Parser{
			Usecase: domain.Usecase{
				Name: "ExampleUsecase",
				Method: []domain.Method{
					domain.Method{
						Name: "Fetch",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "[]*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "GetByID",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Store",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Update",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "exp", Type: "*domain.Example"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "*domain.Example"},
							domain.MethodValue{Type: "error"},
						},
					},
					domain.Method{
						Name: "Delete",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "ctx", Type: "context.Context"},
							domain.MethodValue{Name: "id", Type: "uint64"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "error"},
						},
					},
				},
			},
		}
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate an example_handler.go file", func(t *testing.T) {
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

		err = newFs.CreateDir(serviceName + "/" + dirLayer1 + "/" + dirLayer2)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate example_handler.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGrpcTransport(dirName, domainFile, gomodName, parser)
		resGopg, err := newFs.FindFile(dirName + "/example_handler.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/example_handler.go")
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
		assert.Equal(t, expected_grpc_example_transport, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_handler file
		gen := generator.NewGeneratorService()
		err := gen.GenGrpcTransport(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}
