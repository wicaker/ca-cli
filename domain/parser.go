package domain

// Parser /
type Parser struct {
	Repository
	Usecase
	Handler
}

// ParserDomain /
type ParserDomain interface {
	DomainParser(filePath string) (*Parser, error)
}

// ParserGeneral /
type ParserGeneral interface {
	GeneralParser(filePath string) (*Parser, error)
}
