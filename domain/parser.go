package domain

// Parser /
type Parser struct {
	Repository
	Usecase
}

// ParserService /
type ParserService interface {
	Parser(filePath string) (*Parser, error)
}
