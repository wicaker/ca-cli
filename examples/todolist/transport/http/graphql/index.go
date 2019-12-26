package graphqlhandler

import (
	"net/http"

	"todolist/domain"
	"todolist/transport/http/graphql/mutations"
	"todolist/transport/http/graphql/queries"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	log "github.com/sirupsen/logrus"
)

// graphQLHandler represent the graphQLHandler
type graphQLHandler struct {
	UserUsecase domain.UserUsecase
}

// NewGraphQLHandler will initialize the graphql endpoint
func NewGraphQLHandler(r *http.ServeMux, u domain.UserUsecase) {
	handle := &graphQLHandler{
		UserUsecase: u,
	}

	h := handler.New(&handler.Config{
		Schema:   handle.schema(),
		Pretty:   true,
		GraphiQL: false,
	})

	r.Handle("/graphql", h)
}

func (gh *graphQLHandler) schema() *graphql.Schema {
	rootMutation := mutations.NewGraphQLMutation(gh.UserUsecase)
	rootQuery := queries.NewGraphQLQuery(gh.UserUsecase)

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: rootQuery.GetRootQueryFields(),
		})

	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: rootMutation.GetRootMutationFields(),
	})

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)
	if err != nil {
		log.Printf("errors: %v", err.Error())
	}
	return &schema
}
