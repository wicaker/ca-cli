package domain

// Generator domain /
type Generator struct {
}

// GeneratorService /
type GeneratorService interface {
	GenDomainErrors(dirName string) error
	GenDomainStatusCode(dirName string) error
}
