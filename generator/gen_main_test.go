package generator_test

import (
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
	expected_main_file = `package main

import (
	"github.com/example/examplemain/database/config"
	"github.com/example/examplemain/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}
}

func main() {
	dbgopg := config.GopgInit()

	errChan := make(chan error)

	go func() {
		eServer := server.EchoServer(dbgopg)
		srv := &http.Server{
			Addr:         ":" + os.Getenv("SERVER_ECHO_PORT"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		eServer.HideBanner = true
		eServer.HidePort = true
		log.WithFields(log.Fields{"at": time.Now().Format("2006-01-02 15:04:05")}).Printf("Starting echo server on port :%s...", os.Getenv("SERVER_ECHO_PORT"))
		log.Fatal(eServer.StartServer(srv))
	}()

	go func() {
		gServer := server.GinServer(dbgopg)
		srv := &http.Server{
			Addr:         ":" + os.Getenv("SERVER_GIN_PORT"),
			Handler:      gServer,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		log.WithFields(log.Fields{"at": time.Now().Format("2006-01-02 15:04:05")}).Printf("Starting gin server on port :%s...", os.Getenv("SERVER_GIN_PORT"))
		log.Fatal(srv.ListenAndServe())
	}()

	go func() {
		gmServer := server.GorillaMuxServer(dbgopg)
		srv := &http.Server{
			Addr:         ":" + os.Getenv("SERVER_GORILLA_MUX_PORT"),
			Handler:      gmServer,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		log.WithFields(log.Fields{"at": time.Now().Format("2006-01-02 15:04:05")}).Printf("Starting gorilla mux server on port :%s...", os.Getenv("SERVER_GORILLA_MUX_PORT"))
		log.Fatal(srv.ListenAndServe())
	}()

	go func() {
		httpMuxServer := server.MuxServer(dbgopg)
		srv := &http.Server{
			Addr:         ":" + os.Getenv("SERVER_NET_HTTP_SERVER_MUX_PORT"),
			Handler:      httpMuxServer,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		log.WithFields(log.Fields{"at": time.Now().Format("2006-01-02 15:04:05")}).Printf("Starting net/http ServerMux on port :%s...", os.Getenv("SERVER_NET_HTTP_SERVER_MUX_PORT"))
		log.Fatal(srv.ListenAndServe())
	}()

	go func() {
		httpMuxServer := server.GraphQLServer(dbgopg)
		srv := &http.Server{
			Addr:         ":" + os.Getenv("SERVER_GRAPHQL_SERVER_MUX_PORT"),
			Handler:      httpMuxServer,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		}
		log.WithFields(log.Fields{"at": time.Now().Format("2006-01-02 15:04:05")}).Printf("Starting Graphql Server on port :%s...", os.Getenv("SERVER_GRAPHQL_SERVER_MUX_PORT"))
		log.Fatal(srv.ListenAndServe())
	}()

	go func() {
		// 50051 is the default port for gRPC
		listener, err := net.Listen("tcp", ":"+os.Getenv("SERVER_GRPC_PORT"))

		if err != nil {
			log.Fatalf("Unable to listen on port :%v : %v", os.Getenv("SERVER_GRPC_PORT"), err)
		}

		s := server.GRPCServer(dbgopg)

		// Start the server
		log.WithFields(log.Fields{"at": time.Now().Format("2006-01-02 15:04:05")}).Printf("Starting GRPC server on port :%s...", os.Getenv("SERVER_GRPC_PORT"))
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Fatalln(<-errChan)
}
`
)

func TestGenerateMain(t *testing.T) {
	var (
		serviceName = "test_main"
		gomodName   = "github.com/example/examplemain"
		newFs       = fs.NewFsService()
		transport   = []string{domain.Echo, domain.Gin, domain.GorillaMux, domain.NetHTTP, domain.Graphql, domain.Grpc}
	)

	t.Run("success, should generate an main.go file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate main.go file
		gen := generator.NewGeneratorService()
		err = gen.GenMain(serviceName, gomodName, domain.GoPg, transport)
		resMain, err := newFs.FindFile(serviceName + "/main.go")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resMain)

		data, err := ioutil.ReadFile(serviceName + "/main.go")
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
		assert.Equal(t, expected_main_file, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate main.go file
		gen := generator.NewGeneratorService()
		err := gen.GenMain(serviceName, gomodName, domain.GoPg, transport)

		assert.Error(t, err)
	})
}
