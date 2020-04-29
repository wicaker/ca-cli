package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/wicaker/cacli/domain"
	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/generator"
	"github.com/wicaker/cacli/parser"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	serviceName, goModName, dbHelper, restServer string
	overWrite, grpcOpt, graphqlOpt               bool
	initCmd                                      = &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initiate first clean architecture code",
		Run:     runInit,
	}
)

func runInit(cmd *cobra.Command, args []string) {
	var (
		selectDBOpt = []domain.Option{
			domain.Option{Title: domain.GoPg, Description: "Will using package from https://github.com/go-pg/pg"},
			domain.Option{Title: domain.Gorm, Description: "Will using package from https://github.com/jinzhu/gorm"},
			domain.Option{Title: domain.Sqlx, Description: "Will using package from https://github.com/jmoiron/sqlx"},
			domain.Option{Title: domain.SQL, Description: "Will using package from https://golang.org/pkg/database/sql/"},
			domain.Option{Title: domain.Mongod, Description: "Will using package from go.mongodb.org/mongo-driver/mongo/"},
		}
		selectRestServerOpt = []domain.Option{
			domain.Option{Title: domain.Echo, Description: "Will using package from https://github.com/labstack/echo"},
			domain.Option{Title: domain.Gin, Description: "Will using package from https://github.com/gin-gonic/gin"},
			domain.Option{Title: domain.GorillaMux, Description: "Will using package from https://github.com/gorilla/mux"},
			domain.Option{Title: domain.NetHTTP, Description: "Will using package from https://golang.org/pkg/net/http/"},
			domain.Option{Title: "no", Description: "REST API transport will not added"},
		}
		selectGraphqlOpt = []domain.Option{
			domain.Option{Title: "yes", Description: "Will using package from github.com/graphql-go/graphql"},
			domain.Option{Title: "no", Description: "Graphql transport will not added"},
		}
		selectGrpcOpt = []domain.Option{
			domain.Option{Title: "yes", Description: "This transport wil using package from google.golang.org/grpc"},
			domain.Option{Title: "no", Description: "Grpc transport will not added"},
		}
		overwriteDirOpt = []domain.Option{
			domain.Option{Title: "yes", Description: "The existing directory will be overwriten"},
			domain.Option{Title: "no", Description: "Close the app"},
		}
		newFs = fs.NewFsService()
	)

	// input service name
	if serviceName == "" {
		resName, err := promptInit("Service name")
		serviceName = resName
		failOnInitError(err, `input service name `, serviceName)
	}

	// find existing directory which equal to service name
	res, err := newFs.FindDir(serviceName)
	failOnInitError(err, `find existing directory which equal to service name `, serviceName)

	// ask to overwrite already existing directory
	// will be removed if yes
	if res != nil {
		overwriteOpt, err := selectInit(overwriteDirOpt, "Directory name `"+serviceName+"` already exist, do you want to overwrite ? (be careful to remove existing directory)")
		if err != nil {
			log.Error(err, `ask to overwrite already directory`)
			os.Exit(1)
		}

		if overwriteOpt == "no" {
			log.Error("Directory name `" + serviceName + "` already exist")
			os.Exit(1)
		} else {
			err := newFs.RemoveDir(serviceName)
			failOnInitError(err, `overwrite already dir, remove existing directory `, serviceName)
		}
	}

	// input gomod name
	if goModName == "" {
		goModName, err = promptInit("Go module name")
		failOnInitError(err, `input gomod name `, serviceName)
	}

	// input dbHelper or ORM
	if dbHelper != domain.GoPg && dbHelper != domain.Gorm && dbHelper != domain.Sqlx && dbHelper != domain.SQL && dbHelper != domain.Mongod {
		dbHelper, err = selectInit(selectDBOpt, "DB Helper")
		failOnInitError(err, `input dbHelper or ORM `, serviceName)
	}

	// input http rest api server transport
	if restServer != domain.Echo && restServer != domain.Gin && restServer != domain.GorillaMux && restServer != domain.NetHTTP && restServer != "no" {
		restServer, err = selectInit(selectRestServerOpt, "Using REST API? , choose one if yes!")
		failOnInitError(err, `input http rest api server transport `, serviceName)
	}

	// input graphql transport
	if graphqlOpt == false {
		graphqlOp, err := selectInit(selectGraphqlOpt, "Using Graphql ?")
		failOnInitError(err, `input graphql transport `, serviceName)
		if graphqlOp == "yes" || graphqlOp == "y" {
			graphqlOpt = true
		}
	}

	// input htt2 gRPC transport
	if grpcOpt == false {
		grpcOp, err := selectInit(selectGrpcOpt, "Using gRPC ?")
		failOnInitError(err, `input htt2 gRPC transport `, serviceName)
		if grpcOp == "yes" || grpcOp == "y" {
			grpcOpt = true
		}
	}

	generateInit(
		newFs,
		serviceName,
		goModName,
		dbHelper,
		restServer,
		graphqlOpt,
		grpcOpt,
	)
}

func generateInit(
	newFs domain.FsService,
	serviceName string,
	goModName string,
	dbHelper string,
	restServer string,
	graphqlOpt bool,
	grpcOpt bool,
) {
	var (
		stdout, stderr bytes.Buffer
		transport      []string
	)

	// generator service
	newGen := generator.NewGeneratorService()

	// create project if no directory
	if serviceName != "" {
		// create directory service
		err := newFs.CreateDir(serviceName)
		failOnInitError(err, `create directory service `, serviceName)

		// create go module
		cmd := exec.Command("go", "mod", "init", goModName)
		cmd.Dir = "./" + serviceName
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		failOnInitError(err, `create go module `, serviceName)

		// create domain directory
		err = newFs.CreateDir("./" + serviceName + "/domain")
		failOnInitError(err, `create domain directory `, serviceName)

		// generate errors file inside domain
		err = newGen.GenDomainErrors(serviceName + "/domain")
		failOnInitError(err, `generate errors file inside domain `, serviceName)

		// generate status_code file inside domain
		err = newGen.GenDomainStatusCode(serviceName + "/domain")
		failOnInitError(err, `generate status_code file inside domain `, serviceName)

		// generate success file inside domain
		err = newGen.GenDomainSuccess(serviceName + "/domain")
		failOnInitError(err, `generate success file inside domain `, serviceName)

		// generate example file inside domain
		err = newGen.GenDomainExample(serviceName + "/domain")
		failOnInitError(err, `generate example file inside domain `, serviceName)

		// create usecase directory
		err = newFs.CreateDir("./" + serviceName + "/usecase")
		failOnInitError(err, `create usecase directory `, serviceName)

		// parse file in domain dir
		p := parser.NewParserDomain("example")
		par, err := p.DomainParser("./" + serviceName + "/domain/example.go")
		failOnInitError(err, `parse file in domain dir `, serviceName)

		// create example usecase based on interface in domain layer
		err = newGen.GenUsecase(serviceName+"/usecase", "example.go", goModName, par)
		failOnInitError(err, `create example usecase based on interface in domain layer `, serviceName)

		// create repository directory
		err = newFs.CreateDir("./" + serviceName + "/repository")
		failOnInitError(err, `create repository directory`, serviceName)

		// create repository layer
		if dbHelper == domain.GoPg {
			err = newGen.GenGopgRepository(serviceName+"/repository", "example.go", goModName, par)
			failOnInitError(err, `create gopg repository layer `, serviceName)
		} else if dbHelper == domain.Gorm {
			err = newGen.GenGormRepository(serviceName+"/repository", "example.go", goModName, par)
			failOnInitError(err, `create gorm repository layer `, serviceName)
		} else if dbHelper == domain.Sqlx {
			err = newGen.GenSqlxRepository(serviceName+"/repository", "example.go", goModName, par)
			failOnInitError(err, `create sqlx repository layer `, serviceName)
		} else if dbHelper == domain.SQL {
			err = newGen.GenSQLRepository(serviceName+"/repository", "example.go", goModName, par)
			failOnInitError(err, `create sql repository layer `, serviceName)
		}

		// create database directory
		err = newFs.CreateDir("./" + serviceName + "/database")
		failOnInitError(err, `create database directory `, serviceName)

		// create config directory
		err = newFs.CreateDir("./" + serviceName + "/database/config")
		failOnInitError(err, `create config directory `, serviceName)

		// create db config
		if dbHelper == domain.GoPg {
			err = newGen.GenGopgConfig(serviceName + "/database/config")
			failOnInitError(err, `create db config gopg `, serviceName)
		} else if dbHelper == domain.Gorm {
			err = newGen.GenGormConfig(serviceName + "/database/config")
			failOnInitError(err, `create db config gorm `, serviceName)
		} else if dbHelper == domain.Sqlx {
			err = newGen.GenSqlxConfig(serviceName + "/database/config")
			failOnInitError(err, `create db config sqlx `, serviceName)
		} else if dbHelper == domain.SQL {
			err = newGen.GenSQLConfig(serviceName + "/database/config")
			failOnInitError(err, `create db config sql `, serviceName)
		}

		// create transport directory
		err = newFs.CreateDir("./" + serviceName + "/transport")
		failOnInitError(err, `create transport directory `, serviceName)

		// create middleware directory
		err = newFs.CreateDir("./" + serviceName + "/middleware")
		failOnInitError(err, `create middleware directory `, serviceName)

		// create server directory
		err = newFs.CreateDir("./" + serviceName + "/server")
		failOnInitError(err, ` create server directory `, serviceName)

		// create rest directory
		if restServer != "no" {
			err = newFs.CreateDir("./" + serviceName + "/transport/rest")
			failOnInitError(err, `create transport rest directory `, serviceName)

		}

		// generate transport rest api
		if restServer == domain.Echo {
			err = newGen.GenEchoTransport(serviceName+"/transport/rest", "example.go", goModName, par)
			failOnInitError(err, `generate transport rest api echo `, serviceName)

			err = newGen.GenEchoMiddleware(serviceName + "/middleware")
			failOnInitError(err, `generate middleware echo `, serviceName)

			err = newGen.GenEchoServer(serviceName+"/server", serviceName, dbHelper, goModName, par)
			failOnInitError(err, `generate server echo `, serviceName)

			transport = append(transport, domain.Echo)
		} else if restServer == domain.Gin {
			err = newGen.GenGinTransport(serviceName+"/transport/rest", "example.go", goModName, par)
			failOnInitError(err, `generate transport rest api gin `, serviceName)

			err = newGen.GenGinMiddleware(serviceName + "/middleware")
			failOnInitError(err, `generate middleware gin `, serviceName)

			err = newGen.GenGinServer(serviceName+"/server", serviceName, dbHelper, goModName, par)
			failOnInitError(err, `generate server gin `, serviceName)

			transport = append(transport, domain.Gin)
		} else if restServer == domain.GorillaMux {
			err = newGen.GenGorillaMuxTransport(serviceName+"/transport/rest", "example.go", goModName, par)
			failOnInitError(err, `generate transport rest api gorilla mux `, serviceName)

			err = newGen.GenGorillaMuxMiddleware(serviceName + "/middleware")
			failOnInitError(err, `generate middleware gorilla mux `, serviceName)

			err = newGen.GenGorillaMuxServer(serviceName+"/server", serviceName, dbHelper, goModName, par)
			failOnInitError(err, `generate server gorilla mux `, serviceName)

			transport = append(transport, domain.GorillaMux)
		} else if restServer == domain.NetHTTP {
			err = newGen.GenNetHTTPTransport(serviceName+"/transport/rest", "example.go", goModName, par)
			failOnInitError(err, `generate transport rest api net/http `, serviceName)

			err = newGen.GenNetHTTPMiddleware(serviceName + "/middleware")
			failOnInitError(err, `generate middleware net/http `, serviceName)

			err = newGen.GenNetHTTPMuxServer(serviceName+"/server", serviceName, dbHelper, goModName, par)
			failOnInitError(err, `generate server net/http `, serviceName)

			transport = append(transport, domain.NetHTTP)
		}

		// generate tranport graphql
		if graphqlOpt {
			err = newFs.CreateDir("./" + serviceName + "/transport/graphql")
			failOnInitError(err, `create transport graphql directory `, serviceName)

			err = newGen.GenGraphqlTransport(serviceName+"/transport/graphql", "example.go", goModName, par)
			failOnInitError(err, `generate transport graphql `, serviceName)

			err = newGen.GenNetHTTPMiddleware(serviceName + "/middleware")
			failOnInitError(err, `generate middleware net/http  `, serviceName)

			err = newGen.GenGraphqlServer(serviceName+"/server", serviceName, dbHelper, goModName, par)
			failOnInitError(err, `generate server graphql`, serviceName)

			transport = append(transport, domain.Graphql)
		}

		// generate tranport grpc
		if grpcOpt {
			err = newFs.CreateDir("./" + serviceName + "/transport/grpc")
			failOnInitError(err, `create transport grpc directory `, serviceName)

			err = newFs.CreateDir("./" + serviceName + "/proto")
			failOnInitError(err, `create proto directory `, serviceName)

			err = newGen.GenProtobuf(serviceName+"/proto", "example.go", goModName, par)
			failOnInitError(err, `generate protobuf `, serviceName)

			// generate *.pb.go file
			// cmd := exec.Command("protoc", "--go_out=plugins=grpc:.", "proto/*.proto")
			// cmd.Dir = "./" + serviceName
			// cmd.Stdout = &stdout
			// cmd.Stderr = &stderr
			// err = cmd.Run()
			// failOnInitError(err, `generate *.pb.go file `, serviceName)
			// soda g config

			err = newGen.GenGrpcTransport(serviceName+"/transport/grpc", "example.go", goModName, par)
			failOnInitError(err, `generate transport grpc `, serviceName)

			err = newGen.GenGrpcServer(serviceName+"/server", serviceName, dbHelper, goModName, par)
			failOnInitError(err, `generate server grpc `, serviceName)

			transport = append(transport, domain.Grpc)
		}

		// generate main
		err = newGen.GenMain(serviceName, goModName, dbHelper, transport)
		failOnInitError(err, `generate main.go `, serviceName)

		// generate env
		err = newGen.GenEnv(serviceName)
		failOnInitError(err, `generate env `, serviceName)

		// generate Readme
		err = newGen.GenReadme(serviceName)
		failOnInitError(err, `generate README.md `, serviceName)

		// generate Dockerfile
		err = newGen.GenDockerfile(serviceName)
		failOnInitError(err, `generate Dockerfile `, serviceName)
	}

	log.Info("Congratulation, your `" + serviceName + "` service was successfully initiated !")
}

func failOnInitError(err error, msg string, svcName string) {
	if err != nil {
		newFs := fs.NewFsService()
		log.Errorf("%s: %s", msg, err)

		// remove directory of service
		err = newFs.RemoveDir(svcName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		os.Exit(1)
	}
}

func execConfigDb(serviceName string, stdout bytes.Buffer, stderr bytes.Buffer) {
	// create `database.yml` configuration
	cmd := exec.Command("soda", "g", "config")
	cmd.Dir = "./" + serviceName
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	failOnInitError(err, `create database.yml configuration error, you should install soda first. See : https://github.com/gobuffalo/pop`, serviceName)

	// create example migration
	cmd = exec.Command("soda", "generate", "-p", "./database/migrations", "fizz", "example")
	cmd.Dir = "./" + serviceName
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	failOnInitError(err, `create example migration`, serviceName)
}

func promptInit(label string) (string, error) {
	var result string

	selected := fmt.Sprintf(`{{ "✔" | green | bold }} {{ "%s" | bold }}: {{ "%s" | cyan }}`, label, result)
	templates := promptui.PromptTemplates{
		Success: selected,
	}
	prompt := promptui.Prompt{
		Label:     label,
		Templates: &templates,
		Validate: func(input string) error {
			if len(input) < 1 {
				return errors.New(label + " must have at least 1 characters")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func selectInit(items []domain.Option, label string) (string, error) {
	funcMap := promptui.FuncMap
	funcMap["truncate"] = func(size int, input string) string {
		if len(input) <= size {
			return input
		}
		return input[:size-3] + "..."
	}
	selected := fmt.Sprintf(`{{ "✔" | green | bold }} {{ "%s" | bold }}: {{ .Title | cyan }}`, label)
	templates := promptui.SelectTemplates{
		Active:   `🚩 {{ .Title | cyan | bold }}`,
		Inactive: `   {{ .Title | cyan }}`,
		Selected: selected,
		Details: `
		Description:
		{{ .Description | truncate 80 }}`,
	}

	list := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: &templates,
		Searcher: func(input string, idx int) bool {
			return false
		},
	}
	index, _, err := list.Run()
	if err != nil {
		return "", err
	}

	return items[index].Title, nil
}

func init() {
	initCmd.PersistentFlags().StringVar(&goModName, "gomod", "", "For initiate gomod name")
	initCmd.PersistentFlags().StringVar(&serviceName, "service", "", "Name of service")
	initCmd.PersistentFlags().StringVar(&dbHelper, "database", "", "Database helper library. Choose one of: gopg, gorm, sqlx, sql")
	initCmd.PersistentFlags().StringVar(&restServer, "rest", "", "Rest API server library. Choose one of: echo, gin, gorilla mux, net/http")

	initCmd.PersistentFlags().BoolVar(&overWrite, "overwrite", false, "True if will overwrite directory with the existing service name")
	initCmd.PersistentFlags().BoolVar(&grpcOpt, "grpc", false, "True if generate grpc server")
	initCmd.PersistentFlags().BoolVar(&graphqlOpt, "graphql", false, "True if will use graphql server")

	RootCmd.AddCommand(initCmd)
}
