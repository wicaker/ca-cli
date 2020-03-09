package parser_test

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wicaker/cacli/domain"
	"github.com/wicaker/cacli/parser"

	"github.com/stretchr/testify/assert"
)

func TestParserDomainLayerSuccess(t *testing.T) {
	var (
		pathFile string              = "./mocks/domain/example.go"
		file                         = path.Base(pathFile)                          //example.go
		fileName                     = strings.TrimSuffix(file, filepath.Ext(file)) //example
		pars     domain.ParserDomain = parser.NewParserDomain(fileName)
		expected *domain.Parser      = &domain.Parser{
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
			Repository: domain.Repository{
				Name: "ExampleRepository",
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
	)

	t.Run("success, get domain layer information", func(t *testing.T) {
		res, err := pars.DomainParser(pathFile)
		fmt.Println("file", file)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestParserDomainLayerFailed(t *testing.T) {
	var (
		path2 string              = "./mocks/domain/example2.go"
		path3 string              = "./mocks/domain/example3.go"
		path4 string              = "./mocks/domain/example4.go"
		pars  domain.ParserDomain = parser.NewParserDomain("example")
	)

	t.Run("failed, no file or directory ", func(t *testing.T) {
		_, err := pars.DomainParser("./d")
		assert.Error(t, err)
	})

	t.Run("failed, no file or just directory ", func(t *testing.T) {
		_, err := pars.DomainParser("mocks/domain")
		assert.Error(t, err)
	})

	t.Run("failed, not appropriate structure of file", func(t *testing.T) {
		_, err := pars.DomainParser(path2)
		assert.Error(t, err)
	})

	t.Run("failed, not appropriate structure of file", func(t *testing.T) {
		_, err := pars.DomainParser(path3)
		assert.Error(t, err)
	})

	t.Run("failed, filename and interface name not same", func(t *testing.T) {
		pars = parser.NewParserDomain("example4")
		_, err := pars.DomainParser(path4)
		assert.Error(t, err)
	})
}
