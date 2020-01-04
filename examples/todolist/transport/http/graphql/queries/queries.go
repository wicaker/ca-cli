package queries

import (
	"todolist/domain"

	"github.com/graphql-go/graphql"
)

// GraphQLQuery represent the GraphQLQuery
type GraphQLQuery struct {
	UserUsecase domain.UserUsecase
	TaskUsecase domain.TaskUsecase
}

// NewGraphQLQuery will initialize the user endpoint
func NewGraphQLQuery(u domain.UserUsecase, t domain.TaskUsecase) *GraphQLQuery {
	return &GraphQLQuery{
		UserUsecase: u,
		TaskUsecase: t,
	}
}

// GetRootQueryFields returns all the available queries.
func (gq *GraphQLQuery) GetRootQueryFields() graphql.Fields {
	return graphql.Fields{
		// "user": gq.GetUserQuery(),
		"task": gq.GetTaskQuery(),
	}
}
