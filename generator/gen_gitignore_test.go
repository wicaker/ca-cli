package generator_test

import (
	"os"
	"testing"

	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/generator"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGenerateGitIgnoreFile(t *testing.T) {
	serviceName := "testgitignore"
	newFs := fs.NewFsService()

	t.Run("success, should generate an .gitignore file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate gitignore file
		gen := generator.NewGeneratorService()
		err = gen.GenGitIgnore(serviceName)
		resgitignore, err := newFs.FindFile(serviceName + "/.gitignore")

		assert.NoError(t, err)
		assert.NotEqual(t, nil, resgitignore)

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate gitignore file
		gen := generator.NewGeneratorService()
		err := gen.GenGitIgnore(serviceName)

		assert.Error(t, err)
	})
}
