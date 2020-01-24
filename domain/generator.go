package domain

// Generator domain /
type Generator struct {
}

// GeneratorService /
type GeneratorService interface {
	GenDomainErrors(dirName string) error
	GenDomainStatusCode(dirName string) error
	GenDomainSuccess(dirName string) error
	GenDomainExample(dirName string) error
	GenUsecase(dirName string, domainName string, gomodName string, parser *Parser) error
}
