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
	GenEchoServer(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGinServer(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGorillaMuxServer(dirName string, domainName string, gomodName string, parser *Parser) error
	GenNetHTTPServer(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGraphqlServer(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGrpcServer(dirName string, domainName string, gomodName string, parser *Parser) error
	GenEchoMiddleware(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGinMiddleware(dirName string, domainName string, gomodName string, parser *Parser) error
	GenGorillaMuxMiddleware(dirName string, domainName string, gomodName string, parser *Parser) error
	GenNetHTTPMiddleware(dirName string, domainName string, gomodName string, parser *Parser) error
	GenMain(dirName string, domainName string, gomodName string, parser *Parser) error
}
