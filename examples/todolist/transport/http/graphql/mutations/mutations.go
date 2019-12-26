package mutations

import (
	"todolist/domain"

	"github.com/graphql-go/graphql"
)

// GraphQLMutation represent the graphQLMutation
type GraphQLMutation struct {
	UserUsecase domain.UserUsecase
}

// NewGraphQLMutation will initialize the user endpoint
func NewGraphQLMutation(u domain.UserUsecase) *GraphQLMutation {
	return &GraphQLMutation{
		UserUsecase: u,
	}
}

// GetRootMutationFields returns all the available mutations.
func (gm *GraphQLMutation) GetRootMutationFields() graphql.Fields {
	return graphql.Fields{
		"userLogin": gm.GetLoginUserMutation(),
	}
}
