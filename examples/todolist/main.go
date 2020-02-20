package main

import (
	"net"
	"net/http"
	"os"
	"time"
	"todolist/database/config"
	"todolist/server"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	// load .env file
	e := godotenv.Load()
	if e != nil {
		log.Print(e)
	}
}

func main() {
	dbGopg := config.GopgInit()
	dbSqlx := config.SqlxInit()
	dbSQL := config.SQLInit()
	dbGorm := config.GormInit()

	// conccurrent
	errChan := make(chan error)

	// echo server
	go func() {
		eServer := server.Echo(dbSqlx)
		srv := &http.Server{
			Addr:         ":" + os.Getenv("SERVER_ECHO_PORT"),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		eServer.HideBanner = true
		eServer.HidePort = true
		log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		}).Printf("Starting echo server on port :%s...", os.Getenv("SERVER_ECHO_PORT"))
		log.Fatal(eServer.StartServer(srv))
	}()

	// gin server
	go func() {
		gServer := server.Gin(dbGopg)
		srv := &http.Server{
			Handler:      gServer,
			Addr:         ":" + os.Getenv("SERVER_GIN_PORT"),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		}).Printf("Starting gin server on port :%s...", os.Getenv("SERVER_GIN_PORT"))
		log.Fatal("GIN SERVER", srv.ListenAndServe())
	}()

	// // gorilla mux server
	go func() {
		gmServer := server.GorillaMux(dbSQL)
		srv := &http.Server{
			Handler:      gmServer,
			Addr:         ":" + os.Getenv("SERVER_GORILLA_MUX_PORT"),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		}).Printf("Starting gorilla mux server on port :%s...", os.Getenv("SERVER_GORILLA_MUX_PORT"))
		log.Fatal(srv.ListenAndServe())
	}()

	// // net/http ServeMux
	go func() {
		httpMuxServer := server.ServeMux(dbGorm)
		srv := &http.Server{
			Handler:      httpMuxServer,
			Addr:         ":" + os.Getenv("SERVER_NET_HTTP_SERVER_MUX_PORT"),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		}).Printf("Starting net/http ServerMux on port :%s...", os.Getenv("SERVER_NET_HTTP_SERVER_MUX_PORT"))
		log.Fatal(srv.ListenAndServe())
	}()

	// net/http ServeMux with Graphql server
	go func() {
		httpMuxServer := server.GraphQLServer(dbGopg)
		srv := &http.Server{
			Handler:      httpMuxServer,
			Addr:         ":" + os.Getenv("SERVER_GRAPHQL_SERVER_MUX_PORT"),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		}).Printf("Starting Grahlql Server on port :%s...", os.Getenv("SERVER_GRAPHQL_SERVER_MUX_PORT"))
		log.Fatal(srv.ListenAndServe())
	}()

	// gRPC server
	go func() {
		// 50051 is the default port for gRPC
		listener, err := net.Listen("tcp", ":"+os.Getenv("SERVER_GRPC_PORT"))

		if err != nil {
			log.Fatalf("Unable to listen on port :%v : %v", os.Getenv("SERVER_GRPC_PORT"), err)
		}

		s := server.GRPCServer(dbGopg)

		// Start the server
		log.WithFields(log.Fields{
			"at": time.Now().Format("2006-01-02 15:04:05"),
		}).Println("Starting grpc server on port :50051...")
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Fatalln(<-errChan)
}
