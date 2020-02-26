package generator_test

import (
	"os"
	"testing"

	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/generator"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGenerateEnvFile(t *testing.T) {
	serviceName := "testenv"
	newFs := fs.NewFsService()

	t.Run("success, should generate an .env and env.example file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate env file
		gen := generator.NewGeneratorService()
		err = gen.GenEnv(serviceName)
		resEnv, err := newFs.FindFile(serviceName + "/.env")
		resExample, err := newFs.FindFile(serviceName + "/.env.example")

		assert.NoError(t, err)
		assert.NotEqual(t, nil, resEnv)
		assert.NotEqual(t, nil, resExample)

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate env file
		gen := generator.NewGeneratorService()
		err := gen.GenEnv(serviceName)

		assert.Error(t, err)
	})
}
