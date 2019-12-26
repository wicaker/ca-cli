package queries

import (
	"todolist/domain"

	"github.com/graphql-go/graphql"
)

// GraphQLQuery represent the GraphQLQuery
type GraphQLQuery struct {
	UserUsecase domain.UserUsecase
}

// NewGraphQLQuery will initialize the user endpoint
func NewGraphQLQuery(u domain.UserUsecase) *GraphQLQuery {
	return &GraphQLQuery{
		UserUsecase: u,
	}
}

// GetRootQueryFields returns all the available queries.
func (gq *GraphQLQuery) GetRootQueryFields() graphql.Fields {
	return graphql.Fields{
		"user": gq.GetUserQuery(),
	}
}
