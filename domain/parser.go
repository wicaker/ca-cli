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

var (
	// MockParser used for mock data testing
	MockParser = &Parser{
		Usecase: Usecase{
			Name: "ExampleUsecase",
			Method: []Method{
				Method{
					Name: "Fetch",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "[]*domain.Example"},
						MethodValue{Type: "error"},
					},
				},
				Method{
					Name: "GetByID",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
						MethodValue{Name: "id", Type: "uint64"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "*domain.Example"},
						MethodValue{Type: "error"},
					},
				},
				Method{
					Name: "Store",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
						MethodValue{Name: "exp", Type: "*domain.Example"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "*domain.Example"},
						MethodValue{Type: "error"},
					},
				},
				Method{
					Name: "Update",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
						MethodValue{Name: "exp", Type: "*domain.Example"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "*domain.Example"},
						MethodValue{Type: "error"},
					},
				},
				Method{
					Name: "Delete",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
						MethodValue{Name: "id", Type: "uint64"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "error"},
					},
				},
			},
		},
		Repository: Repository{
			Name: "ExampleRepository",
			Method: []Method{
				Method{
					Name: "Fetch",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "[]*domain.Example"},
						MethodValue{Type: "error"},
					},
				},
				Method{
					Name: "GetByID",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
						MethodValue{Name: "id", Type: "uint64"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "*domain.Example"},
						MethodValue{Type: "error"},
					},
				},
				Method{
					Name: "Store",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
						MethodValue{Name: "exp", Type: "*domain.Example"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "*domain.Example"},
						MethodValue{Type: "error"},
					},
				},
				Method{
					Name: "Update",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
						MethodValue{Name: "exp", Type: "*domain.Example"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "*domain.Example"},
						MethodValue{Type: "error"},
					},
				},
				Method{
					Name: "Delete",
					ParameterList: []MethodValue{
						MethodValue{Name: "ctx", Type: "context.Context"},
						MethodValue{Name: "id", Type: "uint64"},
					},
					ResultList: []MethodValue{
						MethodValue{Type: "error"},
					},
				},
			},
		},
	}
)
