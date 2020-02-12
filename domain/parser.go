package domain

// Parser /
type Parser struct {
	Repository
	Usecase
	Handler
}

// ParserService /
type ParserService interface {
	Parser(filePath string) (*Parser, error)
}

// RepoParserService /
type RepoParserService interface {
	RepoParser(filePath string) (*Parser, error)
}
