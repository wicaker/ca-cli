package generator

import (
	"github.com/wicaker/cacli/domain"
)

type caGen struct {
	gen domain.Generator
}

// NewGeneratorService /
func NewGeneratorService() domain.GeneratorService {
	return &caGen{}
}
