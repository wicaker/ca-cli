package main

import (
	"net"
	"net/http"
	"time"
	"todolist/database/config"
	"todolist/server"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	// load .env file
	e := godotenv.Load()
	if e != nil {
		log.Print(e)
	}

	// viper configuration
	// viper.SetConfigFile(`config.json`)
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	panic(err)
	// }
}

func main() {
	dbGopg := config.GopgInit()

	// conccurrent
	errChan := make(chan error)

	// echo server
	go func() {
		eServer := server.Echo(dbGopg)
		log.Fatal(eServer.Start(":" + "9090"))
	}()

	// gin server
	go func() {
		gServer := server.Gin(dbGopg)
		log.Fatal(gServer.Run(":" + "8090"))
	}()

	// gorilla mux server
	go func() {
		gmServer := server.GorillaMux(dbGopg)
		srv := &http.Server{
			Handler: gmServer,
			Addr:    "127.0.0.1:7090",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())
	}()

	// net/http ServeMux
	go func() {
		httpMuxServer := server.ServeMux(dbGopg)
		srv := &http.Server{
			Handler: httpMuxServer,
			Addr:    "127.0.0.1:6090",
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
			Addr:    "127.0.0.1:5090",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())
	}()

	// gRPC server
	go func() {
		// 50051 is the default port for gRPC
		listener, err := net.Listen("tcp", ":50051")

		if err != nil {
			log.Fatalf("Unable to listen on port :50051: %v", err)
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
