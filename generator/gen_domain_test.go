package generator_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/generator"
)

const (
	expected_domain_errors = `package domain

import "errors"

// ResponseError represent the response error struct
type ResponseError struct {
	Message string ` + "`json:" + `"` + "message" + `"` + "`" + `
}

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Given Param is not valid")
	// ErrUnauthorized will throw if the given request-header token is not valid
	ErrUnauthorized = errors.New("Unauthorized")
)
`
	expected_domain_status_code = `package domain

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetStatusCode will return status code based on type of error
func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	log.Error(err)
	switch err {
	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
`
	expected_domain_success = `package domain

// ResponseSuccess represent the reseponse success struct
type ResponseSuccess struct {
	Message string      ` + "`json:" + `"` + "message" + `"` + "`" + `
	Data    interface{} ` + "`json:" + `"` + "data" + `"` + "`" + `
}

var (
	// ResponseData with type map used to response json if no error
	ResponseData map[string]interface{}
)
`
)

func TestGenerateDomainErrors(t *testing.T) {
	var (
		serviceName = "test_domain_errors"
		dirLayer    = "domain"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an errors.go file", func(t *testing.T) {
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

		// generate errors.go file
		gen := generator.NewGeneratorService()
		err = gen.GenDomainErrors(dirName)
		resGopg, err := newFs.FindFile(dirName + "/errors.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/errors.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_domain_errors, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate errors.go file
		gen := generator.NewGeneratorService()
		err := gen.GenDomainErrors(serviceName)

		assert.Error(t, err)
	})
}

func TestGenerateDomainStatusCode(t *testing.T) {
	var (
		serviceName = "test_domain_status_code"
		dirLayer    = "domain"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an status_code.go file", func(t *testing.T) {
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

		// generate status_code.go file
		gen := generator.NewGeneratorService()
		err = gen.GenDomainStatusCode(dirName)
		resGopg, err := newFs.FindFile(dirName + "/status_code.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/status_code.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_domain_status_code, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate status_code.go file
		gen := generator.NewGeneratorService()
		err := gen.GenDomainStatusCode(serviceName)

		assert.Error(t, err)
	})
}

func TestGenerateDomainSuccess(t *testing.T) {
	var (
		serviceName = "test_domain_success"
		dirLayer    = "domain"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an success.go file", func(t *testing.T) {
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

		// generate success.go file
		gen := generator.NewGeneratorService()
		err = gen.GenDomainSuccess(dirName)
		resGopg, err := newFs.FindFile(dirName + "/success.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/success.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_domain_success, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate success.go file
		gen := generator.NewGeneratorService()
		err := gen.GenDomainSuccess(serviceName)

		assert.Error(t, err)
	})
}

func TestGenerateDomainExample(t *testing.T) {
	var (
		serviceName = "test_domain_example"
		dirLayer    = "domain"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an example.go file", func(t *testing.T) {
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

		// generate example.go file
		gen := generator.NewGeneratorService()
		err = gen.GenDomainExample(dirName)
		resGopg, err := newFs.FindFile(dirName + "/example.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example.go file
		gen := generator.NewGeneratorService()
		err := gen.GenDomainExample(serviceName)

		assert.Error(t, err)
	})
}
