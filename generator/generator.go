/*
Package generator used to generate codes and create a file based on the layer
*/
package generator

import (
	"github.com/wicaker/cacli/domain"
)

type caGen struct {
	gen domain.Generator
}

// NewGeneratorService will create new a caGen object representation of domain.Generator interface
func NewGeneratorService() domain.GeneratorService {
	return &caGen{}
}
