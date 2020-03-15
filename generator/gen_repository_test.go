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
	expected_gopg_example_repository = `package repository

import (
	"context"
	domain "github.com/example/examplerepository/domain"
	pg "github.com/go-pg/pg/v9"
)

type gopgExampleRepository struct {
	Conn *pg.DB
}

// NewGopgExampleRepository will create new an gopgExampleRepository object representation of domain.ExampleRepository interface
func NewGopgExampleRepository(Conn *pg.DB) domain.ExampleRepository {
	return &gopgExampleRepository{Conn: Conn}
}

func (er *gopgExampleRepository) Fetch(ctx context.Context) ([]*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *gopgExampleRepository) GetByID(ctx context.Context, id uint64) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *gopgExampleRepository) Store(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *gopgExampleRepository) Update(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *gopgExampleRepository) Delete(ctx context.Context, id uint64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil
}
`

	expected_gorm_example_repository = `package repository

import (
	"context"
	domain "github.com/example/examplerepository/domain"
	gorm "github.com/jinzhu/gorm"
)

type gormExampleRepository struct {
	Conn *gorm.DB
}

// NewGormExampleRepository will create new an gormExampleRepository object representation of domain.ExampleRepository interface
func NewGormExampleRepository(Conn *gorm.DB) domain.ExampleRepository {
	return &gormExampleRepository{Conn: Conn}
}

func (er *gormExampleRepository) Fetch(ctx context.Context) ([]*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *gormExampleRepository) GetByID(ctx context.Context, id uint64) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *gormExampleRepository) Store(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *gormExampleRepository) Update(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *gormExampleRepository) Delete(ctx context.Context, id uint64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil
}
`

	expected_sql_example_repository = `package repository

import (
	"context"
	"database/sql"
	domain "github.com/example/examplerepository/domain"
)

type sqlExampleRepository struct {
	Conn *sql.DB
}

// NewSQLExampleRepository will create new an sqlExampleRepository object representation of domain.ExampleRepository interface
func NewSQLExampleRepository(Conn *sql.DB) domain.ExampleRepository {
	return &sqlExampleRepository{Conn: Conn}
}

func (er *sqlExampleRepository) Fetch(ctx context.Context) ([]*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *sqlExampleRepository) GetByID(ctx context.Context, id uint64) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *sqlExampleRepository) Store(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *sqlExampleRepository) Update(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *sqlExampleRepository) Delete(ctx context.Context, id uint64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil
}
`

	expected_sqlx_example_repository = `package repository

import (
	"context"
	domain "github.com/example/examplerepository/domain"
	sqlx "github.com/jmoiron/sqlx"
)

type sqlxExampleRepository struct {
	Conn *sqlx.DB
}

// NewSqlxExampleRepository will create new an sqlxExampleRepository object representation of domain.ExampleRepository interface
func NewSqlxExampleRepository(Conn *sqlx.DB) domain.ExampleRepository {
	return &sqlxExampleRepository{Conn: Conn}
}

func (er *sqlxExampleRepository) Fetch(ctx context.Context) ([]*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *sqlxExampleRepository) GetByID(ctx context.Context, id uint64) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *sqlxExampleRepository) Store(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *sqlxExampleRepository) Update(ctx context.Context, exp *domain.Example) (*domain.Example, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil, nil
}

func (er *sqlxExampleRepository) Delete(ctx context.Context, id uint64) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return nil
}
`
)

func TestGenerateGopgRepository(t *testing.T) {
	var (
		serviceName = "test_gppg_example_repository"
		dirLayer    = "repository"
		domainFile  = "example.go"
		gomodName   = "github.com/example/examplerepository"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		parser      = &domain.Parser{
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
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate an example_repository.go file", func(t *testing.T) {
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

		// generate example_repository.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGopgRepository(dirName, domainFile, gomodName, parser)
		resGopg, err := newFs.FindFile(dirName + "/example_repository.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/example_repository.go")
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
		assert.Equal(t, expected_gopg_example_repository, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_repository file
		gen := generator.NewGeneratorService()
		err := gen.GenGopgRepository(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}

func TestGenerateGormRepository(t *testing.T) {
	var (
		serviceName = "test_gorm_example_repository"
		dirLayer    = "repository"
		domainFile  = "example.go"
		gomodName   = "github.com/example/examplerepository"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		parser      = &domain.Parser{
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
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate an example_repository.go file", func(t *testing.T) {
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

		// generate example_repository.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGormRepository(dirName, domainFile, gomodName, parser)
		resGorm, err := newFs.FindFile(dirName + "/example_repository.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGorm)

		data, err := ioutil.ReadFile(dirName + "/example_repository.go")
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
		assert.Equal(t, expected_gorm_example_repository, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_repository file
		gen := generator.NewGeneratorService()
		err := gen.GenGopgRepository(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}

func TestGenerateSQLRepository(t *testing.T) {
	var (
		serviceName = "test_sql_example_repository"
		dirLayer    = "repository"
		domainFile  = "example.go"
		gomodName   = "github.com/example/examplerepository"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		parser      = &domain.Parser{
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
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate an example_repository.go file", func(t *testing.T) {
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

		// generate example_repository.go file
		gen := generator.NewGeneratorService()
		err = gen.GenSQLRepository(dirName, domainFile, gomodName, parser)
		resSQL, err := newFs.FindFile(dirName + "/example_repository.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resSQL)

		data, err := ioutil.ReadFile(dirName + "/example_repository.go")
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
		assert.Equal(t, expected_sql_example_repository, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_repository file
		gen := generator.NewGeneratorService()
		err := gen.GenGopgRepository(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}

func TestGenerateSqlxRepository(t *testing.T) {
	var (
		serviceName = "test_example_repository"
		dirLayer    = "repository"
		domainFile  = "example.go"
		gomodName   = "github.com/example/examplerepository"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		parser      = &domain.Parser{
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
		newFs = fs.NewFsService()
	)

	t.Run("success, should generate an example_repository.go file", func(t *testing.T) {
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

		// generate example_repository.go file
		gen := generator.NewGeneratorService()
		err = gen.GenSqlxRepository(dirName, domainFile, gomodName, parser)
		resSqlx, err := newFs.FindFile(dirName + "/example_repository.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resSqlx)

		data, err := ioutil.ReadFile(dirName + "/example_repository.go")
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
		assert.Equal(t, expected_sqlx_example_repository, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example_repository file
		gen := generator.NewGeneratorService()
		err := gen.GenGopgRepository(dirName, domainFile, gomodName, parser)

		assert.Error(t, err)
	})
}
