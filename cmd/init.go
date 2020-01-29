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

// initCmd represents the new command
var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initiate first clean architecture code",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			goModName, dbHelper, restServer, graphqlOpt, grpcOpt string
			selectDBOpt                                          = []domain.Option{
				domain.Option{Title: domain.GoPg, Description: "Will using package from https://github.com/go-pg/pg"},
				domain.Option{Title: domain.Gorm, Description: "Will using package from https://github.com/jinzhu/gorm"},
				domain.Option{Title: domain.Sqlx, Description: "Will using package from https://github.com/jmoiron/sqlx"},
				domain.Option{Title: domain.SQL, Description: "Will using package from https://golang.org/pkg/database/sql/"},
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
			// isAgain        string = "yes"
			transport      []string
			stdout, stderr bytes.Buffer
		)

		if len(args) == 0 {
			log.Error("You must provide a name for the service")
			return
		}

		if len(args) > 1 {
			log.Error("Input command properly. Exp: cacli init service_name")
			return
		}

		// input gomod name
		goModName, err := promptInit("Go module name")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// input dbHelper or ORM
		dbHelper, err = selectInit(selectDBOpt, "DB Helper")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// input http rest api server transport
		restServer, err = selectInit(selectRestServerOpt, "Using REST API? , choose if yes!")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		// input graphql transport
		graphqlOpt, err = selectInit(selectGraphqlOpt, "Using Graphql ?")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// input htt2 gRPC transport
		grpcOpt, err = selectInit(selectGrpcOpt, "Using gRPC ?")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// fs service
		newFs := fs.NewFsService()
		res, err := newFs.FindDir(args[0])
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generator service
		newGen := generator.NewGeneratorService()

		// overwrite already dir ?
		if res != nil {
			log.Error("Directory name already exist, use another name")
			os.Exit(1)
		}

		// create project if no directory
		if res == nil {
			// create directory service
			err := newFs.CreateDir(args[0])
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// create go module
			cmd := exec.Command("go", "mod", "init", goModName)
			cmd.Dir = "./" + args[0]
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// create domain directory
			err = newFs.CreateDir("./" + args[0] + "/domain")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// generate errors file inside domain
			err = newGen.GenDomainErrors(args[0])
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// generate status_code file inside domain
			err = newGen.GenDomainStatusCode(args[0])
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// generate success file inside domain
			err = newGen.GenDomainSuccess(args[0])
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// generate example file inside domain
			err = newGen.GenDomainExample(args[0])
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// create usecase directory
			err = newFs.CreateDir("./" + args[0] + "/usecase")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// parse file in domain dir
			p := parser.NewParserService("example")
			par, err := p.Parser("./" + args[0] + "/domain/example.go")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// create example usecase based on interface in domain layer
			err = newGen.GenUsecase(args[0], "example", goModName, par)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// create repository directory
			err = newFs.CreateDir("./" + args[0] + "/repository")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// create repository layer
			if dbHelper == domain.GoPg {
				err = newGen.GenGopgRepository(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else if dbHelper == domain.Gorm {
				err = newGen.GenGormRepository(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else if dbHelper == domain.Sqlx {
				err = newGen.GenSqlxRepository(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else if dbHelper == domain.SQL {
				err = newGen.GenSQLRepository(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

			// create database directory
			err = newFs.CreateDir("./" + args[0] + "/database")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// create config directory
			err = newFs.CreateDir("./" + args[0] + "/database/config")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// create db config
			if dbHelper == domain.GoPg {
				err = newGen.GenGopgConfig(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else if dbHelper == domain.Gorm {
				err = newGen.GenGormConfig(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else if dbHelper == domain.Sqlx {
				err = newGen.GenSqlxConfig(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else if dbHelper == domain.SQL {
				err = newGen.GenSQLConfig(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

			// create transport directory
			err = newFs.CreateDir("./" + args[0] + "/transport")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// create rest directory
			if restServer != "no" {
				err = newFs.CreateDir("./" + args[0] + "/transport/rest")
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

			// generate transport rest api
			if restServer == domain.Echo {
				err = newGen.GenEchoTransport(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else if restServer == domain.Gin {
				err = newGen.GenEchoTransport(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else if restServer == domain.GorillaMux {
				err = newGen.GenGorillaMuxTransport(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else if restServer == domain.NetHTTP {
				err = newGen.GenNetHTTPTransport(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

			// generate tranport graphql
			if graphqlOpt == "yes" || graphqlOpt == "y" {
				err = newFs.CreateDir("./" + args[0] + "/transport/graphql")
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}

				err = newGen.GenGraphqlTransport(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

			// generate tranport grpc
			if grpcOpt == "yes" || grpcOpt == "y" {
				err = newFs.CreateDir("./" + args[0] + "/transport/grpc")
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}

				err = newGen.GenGrpcTransport(args[0], "example", goModName, par)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

		}
		fmt.Println(dbHelper, transport)
	},
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
	rootCmd.AddCommand(initCmd)
}
