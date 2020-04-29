package generator_test

import (
	"os"
	"testing"

	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/generator"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGenerateReadmeFile(t *testing.T) {
	serviceName := "testreadme"
	newFs := fs.NewFsService()

	t.Run("success, should generate a README.md file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate README.md file
		gen := generator.NewGeneratorService()
		err = gen.GenReadme(serviceName)
		resReadme, err := newFs.FindFile(serviceName + "/README.md")

		assert.NoError(t, err)
		assert.NotEqual(t, nil, resReadme)

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate README.md file
		gen := generator.NewGeneratorService()
		err := gen.GenReadme(serviceName)

		assert.Error(t, err)
	})
}
