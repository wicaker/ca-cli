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
	"github.com/thoas/go-funk"
)

// initCmd represents the new command
var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initiate first clean architecture code",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			goModName, dbHelper string
			selectDBOpt         = []domain.Option{
				domain.Option{Title: domain.GoPg, Description: "https://github.com/go-pg/pg"},
				domain.Option{Title: domain.Gorm, Description: "https://github.com/jinzhu/gorm"},
				domain.Option{Title: domain.Sqlx, Description: "https://github.com/jmoiron/sqlx"},
				domain.Option{Title: domain.SQL, Description: "https://golang.org/pkg/database/sql/"},
			}
			selectTransportOpt = []domain.Option{
				domain.Option{Title: domain.Echo, Description: "https://github.com/labstack/echo"},
				domain.Option{Title: domain.Gin, Description: "https://github.com/gin-gonic/gin"},
				domain.Option{Title: domain.GorillaMux, Description: "https://github.com/gorilla/mux"},
				domain.Option{Title: domain.NetHTTP, Description: "https://golang.org/pkg/net/http/"},
				domain.Option{Title: domain.Graphql, Description: "github.com/graphql-go/graphql"},
				domain.Option{Title: domain.Grpc, Description: "google.golang.org/grpc"},
			}
			isAgain        string = "yes"
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

		// input transport used
		numTransport := 1
		for isAgain == "yes" || isAgain == "y" && len(selectTransportOpt) != 0 {
			// give transport option to user
			label := fmt.Sprintf("Transport %d", numTransport)
			t, err := selectInit(selectTransportOpt, label)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			// append new transport and remove option which already used
			transport = append(transport, t)
			r := funk.Filter(selectTransportOpt, func(x domain.Option) bool {
				return x.Title != t
			})
			selectTransportOpt = r.([]domain.Option)

			// ask if user want add another transport method
			if len(selectTransportOpt) > 0 {
				isAgain, err = promptInit("Add new transport ? (yes/y)")
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

			numTransport++
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
