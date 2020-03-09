package parser_test

import (
	"testing"

	"github.com/wicaker/cacli/domain"
	"github.com/wicaker/cacli/parser"

	"github.com/stretchr/testify/assert"
)

func TestParserGeneralLayerSuccess(t *testing.T) {
	var (
		pars      domain.ParserGeneral = parser.NewParserGeneral()
		expecRepo *domain.Parser          = &domain.Parser{
			Repository: domain.Repository{
				Name: "",
				Method: []domain.Method{
					domain.Method{
						Name: "NewGopgExampleRepository",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "Conn", Type: "*pg.DB"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "domain.ExampleRepository"},
						},
					},
				},
			},
		}
		expecUsecase *domain.Parser = &domain.Parser{
			Usecase: domain.Usecase{
				Name: "",
				Method: []domain.Method{
					domain.Method{
						Name: "NewExampleUsecase",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "er", Type: "domain.ExampleRepository"},
							domain.MethodValue{Name: "timeout", Type: "time.Duration"},
						},
						ResultList: []domain.MethodValue{
							domain.MethodValue{Type: "domain.ExampleUsecase"},
						},
					},
				},
			},
		}
		expecHand *domain.Parser = &domain.Parser{
			Handler: domain.Handler{
				Name: "",
				Method: []domain.Method{
					domain.Method{
						Name: "NewExampleHandler",
						ParameterList: []domain.MethodValue{
							domain.MethodValue{Name: "e", Type: "*echo.Echo"},
							domain.MethodValue{Name: "u", Type: "domain.ExampleUsecase"},
						},
						ResultList: []domain.MethodValue{},
					},
				},
			},
		}
	)

	t.Run("success, get repository layer information", func(t *testing.T) {
		res, err := pars.GeneralParser("./mocks/repository/example_repository.go")
		assert.NoError(t, err)
		assert.Equal(t, expecRepo.Repository, res.Repository)
	})

	t.Run("success, get usecase layer information", func(t *testing.T) {
		res, err := pars.GeneralParser("./mocks/usecase/example_usecase.go")
		assert.NoError(t, err)
		assert.Equal(t, expecUsecase.Usecase, res.Usecase)
	})

	t.Run("success, get transport layer information", func(t *testing.T) {
		res, err := pars.GeneralParser("./mocks/transport/rest/example_handler.go")
		assert.NoError(t, err)
		assert.Equal(t, expecHand.Handler, res.Handler)
	})
}

func TestParserGeneralLayerFailed(t *testing.T) {
	var (
		pars domain.ParserGeneral = parser.NewParserGeneral()
	)

	t.Run("failed, no file or directory ", func(t *testing.T) {
		_, err := pars.GeneralParser("./d")
		assert.Error(t, err)
	})

	t.Run("failed, no file or just directory ", func(t *testing.T) {
		_, err := pars.GeneralParser("mocks/domain")
		assert.Error(t, err)
	})
}
