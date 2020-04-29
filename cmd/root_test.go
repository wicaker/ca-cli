package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/wicaker/cacli/cmd"
	"github.com/wicaker/cacli/fs"
)

func TestInitCommand(t *testing.T) {
	var (
		stdout, stderr bytes.Buffer
		b              = bytes.NewBufferString("")
		newFs          = fs.NewFsService()
		serviceName    = "test"
	)

	t.Run("success, should initiate a go project structure", func(t *testing.T) {
		cmd.RootCmd.SetOut(b)
		cmd.RootCmd.SetArgs([]string{
			`init`,
			`--service=` + serviceName,
			`--gomod=github.com/wicaker/test`,
			`--database=gopg`,
			`--rest=echo`,
			`--grpc=true`,
			`--graphql=true`,
		})
		cmd.RootCmd.Execute()

		_, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		cmdR := exec.Command("protoc", "--go_out=plugins=grpc:.", "proto/example.proto")
		cmdR.Dir = "./" + serviceName
		cmdR.Stdout = &stdout
		cmdR.Stderr = &stderr
		err = cmdR.Run()

		if err != nil {
			fmt.Println(err)
			t.Fatal(err)
		}

		cmdR = exec.Command("go", "build")
		cmdR.Dir = "./" + serviceName
		cmdR.Stdout = &stdout
		cmdR.Stderr = &stderr
		err = cmdR.Run()

		if err != nil {
			fmt.Println(err)
			t.Fatal(err)
		}

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})
}
