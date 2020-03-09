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
	GenGopgConfig(dirName string) error
	GenGormConfig(dirName string) error
	GenSQLConfig(dirName string) error
	GenSqlxConfig(dirName string) error
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
	GenEchoMiddleware(dirName string) error
	GenGinMiddleware(dirName string) error
	GenGorillaMuxMiddleware(dirName string) error
	GenNetHTTPMiddleware(dirName string) error
	GenMain(dirName string, gomodName string) error
	GenEnv(dirName string) error
}
