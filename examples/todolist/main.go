package main

import (
	"net"
	"net/http"
	"time"
	"todolist/database/config"
	"todolist/server"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// load .env file
	e := godotenv.Load()
	if e != nil {
		log.Print(e)
	}

	// viper configuration
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
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
			Addr: viper.GetString("server.echo.address"),
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Fatal(eServer.StartServer(srv))
	}()

	// gin server
	go func() {
		gServer := server.Gin(dbGopg)
		srv := &http.Server{
			Handler: gServer,
			Addr:    viper.GetString("server.gin.address"),
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())
	}()

	// gorilla mux server
	go func() {
		gmServer := server.GorillaMux(dbSQL)
		srv := &http.Server{
			Handler: gmServer,
			Addr:    viper.GetString("server.gorilla-mux.address"),
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())
	}()

	// net/http ServeMux
	go func() {
		httpMuxServer := server.ServeMux(dbGorm)
		srv := &http.Server{
			Handler: httpMuxServer,
			Addr:    viper.GetString("server.net-http-server-mux.address"),
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())
	}()

	// net/http ServeMux with Graphql server
	go func() {
		httpMuxServer := server.GraphQLServer(dbGopg)
		srv := &http.Server{
			Handler: httpMuxServer,
			Addr:    viper.GetString("server.graphql-server-mux.address"),
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())
	}()

	// gRPC server
	go func() {
		// 50051 is the default port for gRPC
		listener, err := net.Listen("tcp", viper.GetString("server.grpc.address"))

		if err != nil {
			log.Fatalf("Unable to listen on port :%v : %v", viper.GetString("server.grpc.address"), err)
		}

		s := server.GRPCServer(dbGopg)

		// Start the server
		log.Println("Starting grpc server on port :50051...")
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Fatalln(<-errChan)
}
