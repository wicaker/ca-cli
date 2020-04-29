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
	expected_echo_middleware = `package middleware

import (
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"time"
)

// EchoMiddleware represent the data-struct for middleware
type EchoMiddleware struct{}

// InitEchoMiddleware intialize the middleware
func InitEchoMiddleware() *EchoMiddleware {
	return &EchoMiddleware{}
}

// CORS will handle the CORS middleware
func (m *EchoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

// MiddlewareLogging for logging
func (m *EchoMiddleware) MiddlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		makeLogEntry(c).Info("incoming request")
		return next(c)
	}
}
func makeLogEntry(c echo.Context) *log.Entry {
	if c == nil {
		return log.WithFields(log.Fields{"at": time.Now().Format("2006-01-02 15:04:05")})
	}
	return log.WithFields(log.Fields{
		"at":     time.Now().Format("2006-01-02 15:04:05"),
		"ip":     c.Request().RemoteAddr,
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
	})
}
`
	expected_gin_middleware = `package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// GinMiddleware represent the data-struct for middleware
type GinMiddleware struct{}

// InitGinMiddleware intialize the middleware
func InitGinMiddleware() *GinMiddleware {
	return &GinMiddleware{}
}

// CORS will handle the CORS middleware
func (m *GinMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// MiddlewareLogging for logging
func (m *GinMiddleware) MiddlewareLogging() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n", param.ClientIP, param.TimeStamp.Format(time.RFC1123), param.Method, param.Path, param.Request.Proto, param.StatusCode, param.Latency, param.Request.UserAgent(), param.ErrorMessage)
	})
}
`
	expected_gorilla_mux_middleware = `package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// GorillaMuxMiddleware represent the data-struct for middleware
type GorillaMuxMiddleware struct{}

// InitGorillaMuxMiddleware intialize the middleware
func InitGorillaMuxMiddleware() *GorillaMuxMiddleware {
	return &GorillaMuxMiddleware{}
}

// CORS will handle the CORS middleware
func (m *GorillaMuxMiddleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

// MiddlewareLogging for logging
func (m *GorillaMuxMiddleware) MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"at":     time.Now().Format("2006-01-02 15:04:05"),
			"ip":     r.RemoteAddr,
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Info("incoming request")

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
`
	expected_net_http_middleware = `package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// NetHTTPMiddleware represent the data-struct for middleware
type NetHTTPMiddleware struct{}

// InitNetHTTPMiddleware intialize the middleware
func InitNetHTTPMiddleware() *NetHTTPMiddleware {
	return &NetHTTPMiddleware{}
}

// CORS will handle the CORS middleware
func (m *NetHTTPMiddleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

// MiddlewareLogging for logging
func (m *NetHTTPMiddleware) MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"at":     time.Now().Format("2006-01-02 15:04:05"),
			"ip":     r.RemoteAddr,
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Info("incoming request")

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
`
)

func TestGenerateEchoMiddleware(t *testing.T) {
	var (
		serviceName = "test_echo_middleware"
		dirLayer    = "middleware"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an echo_middleware.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(dirName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate echo_middleware.go file
		gen := generator.NewGeneratorService()
		err = gen.GenEchoMiddleware(dirName)
		resGopg, err := newFs.FindFile(dirName + "/echo_middleware.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/echo_middleware.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_echo_middleware, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate echo_middleware.go file
		gen := generator.NewGeneratorService()
		err := gen.GenEchoMiddleware(serviceName)

		assert.Error(t, err)
	})
}

func TestGenerateGinMiddleware(t *testing.T) {
	var (
		serviceName = "test_gin_middleware"
		dirLayer    = "middleware"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an gin_middleware.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(dirName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate gin_middleware.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGinMiddleware(dirName)
		resGopg, err := newFs.FindFile(dirName + "/gin_middleware.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/gin_middleware.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_gin_middleware, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate gin_middleware.go file
		gen := generator.NewGeneratorService()
		err := gen.GenGinMiddleware(serviceName)

		assert.Error(t, err)
	})
}

func TestGenerateGorillaMuxMiddleware(t *testing.T) {
	var (
		serviceName = "test_gorilla_mux_middleware"
		dirLayer    = "middleware"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an gorilla_mux_middleware.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(dirName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate gorilla_mux_middleware.go file
		gen := generator.NewGeneratorService()
		err = gen.GenGorillaMuxMiddleware(dirName)
		resGopg, err := newFs.FindFile(dirName + "/gorilla_mux_middleware.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/gorilla_mux_middleware.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_gorilla_mux_middleware, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate gorilla_mux_middleware.go file
		gen := generator.NewGeneratorService()
		err := gen.GenGorillaMuxMiddleware(serviceName)

		assert.Error(t, err)
	})
}

func TestGenerateNetHTTPMiddleware(t *testing.T) {
	var (
		serviceName = "test_net_http_middleware"
		dirLayer    = "middleware"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an net_http_middleware.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(dirName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate net_http_middleware.go file
		gen := generator.NewGeneratorService()
		err = gen.GenNetHTTPMiddleware(dirName)
		resGopg, err := newFs.FindFile(dirName + "/net_http_middleware.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/net_http_middleware.go")
		if err != nil {
			log.Error("File reading error", err)
			os.Exit(1)
		}
		assert.Equal(t, expected_net_http_middleware, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate net_http_middleware.go file
		gen := generator.NewGeneratorService()
		err := gen.GenNetHTTPMiddleware(serviceName)

		assert.Error(t, err)
	})
}
