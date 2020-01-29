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
	GenGopgRepository(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGormRepository(dirName string, domainName string, gomodName string, parser *Parser) error
	GenSQLRepository(dirName string, domainName string, gomodName string, parser *Parser) error
	GenSqlxRepository(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGopgConfig(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGormConfig(dirName string, domainName string, gomodName string, parser *Parser) error
	GenSQLConfig(dirName string, domainName string, gomodName string, parser *Parser) error
	GenSqlxConfig(dirName string, domainName string, gomodName string, parser *Parser) error
	GenEchoTransport(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGinTransport(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGorillaMuxTransport(dirName string, domainName string, gomodName string, parser *Parser) error
	GenNetHTTPTransport(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGraphqlTransport(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGrpcTransport(dirName string, domainName string, gomodName string, parser *Parser) error
}
