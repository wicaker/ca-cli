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
	expected_example_usecase_no_repo = `package usecase

import (
	"context"
	domain "github.com/example/exampleusecase/domain"
	"time"
)

type exampleUsecase struct {
	contextTimeout time.Duration
}

// NewExampleUsecase will create new an exampleUsecase object representation of domain.ExampleUsecase interface
func NewExampleUsecase(timeout time.Duration) domain.ExampleUsecase {
	return &exampleUsecase{contextTimeout: timeout}
}

func (eu *exampleUsecase) Fetch(ctx context.Context) ([]*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}

func (eu *exampleUsecase) GetByID(ctx context.Context, id uint64) (*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}

func (eu *exampleUsecase) Store(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}

func (eu *exampleUsecase) Update(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}

func (eu *exampleUsecase) Delete(ctx context.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil
}
`

	expected_example_usecase_with_repo = `package usecase

import (
	"context"
	domain "github.com/example/exampleusecase/domain"
	"time"
)

type exampleUsecase struct {
	exampleRepo    domain.ExampleRepository
	contextTimeout time.Duration
}

// NewExampleUsecase will create new an exampleUsecase object representation of domain.ExampleUsecase interface
func NewExampleUsecase(er domain.ExampleRepository, timeout time.Duration) domain.ExampleUsecase {
	return &exampleUsecase{
		contextTimeout: timeout,
		exampleRepo:    er,
	}
}

func (eu *exampleUsecase) Fetch(ctx context.Context) ([]*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}

func (eu *exampleUsecase) GetByID(ctx context.Context, id uint64) (*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}

func (eu *exampleUsecase) Store(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}

func (eu *exampleUsecase) Update(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil, nil
}

func (eu *exampleUsecase) Delete(ctx context.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(ctx, eu.contextTimeout)
	defer cancel()
	return nil
}
`
)

func TestGenerateUsecaseLayer(t *testing.T) {
	var (
		serviceName = "test_example_usecase"
		dirLayer    = "usecase"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampleusecase"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		parser      = &domain.Parser{}
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate a file in usecase layer and without repository layer, with appropriate file name", func(t *testing.T) {
		parser = &domain.Parser{
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

		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirLayer)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate example_usecase.go file
		gen := generator.NewGeneratorService()
		err = gen.GenUsecase(dirName, domainFile, gomodName, parser)
		resGopg, err := newFs.FindFile(dirName + "/example_usecase.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/example_usecase.go")
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
		assert.Equal(t, expected_example_usecase_no_repo, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("success, should generate a file in usecase layer and with repository layer, with appropriate file name", func(t *testing.T) {
		parser = &domain.Parser{
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

		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirLayer)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate example_usecase.go file
		gen := generator.NewGeneratorService()
		err = gen.GenUsecase(dirName, domainFile, gomodName, parser)
		resGopg, err := newFs.FindFile(dirName + "/example_usecase.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/example_usecase.go")
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
		assert.Equal(t, expected_example_usecase_with_repo, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because parser not contain appropriate value", func(t *testing.T) {
		// generate status_code.go file
		gen := generator.NewGeneratorService()
		err := gen.GenUsecase(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}
