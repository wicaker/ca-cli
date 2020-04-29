package generator_test

import (
	"os"
	"testing"

	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/generator"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGenDockerFile(t *testing.T) {
	serviceName := "testdocker"
	newFs := fs.NewFsService()

	t.Run("success, should generate a Dockerfile file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate Dockerfile file
		gen := generator.NewGeneratorService()
		err = gen.GenDockerfile(serviceName)
		resDocker, err := newFs.FindFile(serviceName + "/Dockerfile")

		assert.NoError(t, err)
		assert.NotEqual(t, nil, resDocker)

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate Dockerfile file
		gen := generator.NewGeneratorService()
		err := gen.GenDockerfile(serviceName)

		assert.Error(t, err)
	})
}
