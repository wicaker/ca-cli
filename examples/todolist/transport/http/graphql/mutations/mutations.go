package mutations

import (
	"todolist/domain"

	"github.com/graphql-go/graphql"
)

// GraphQLMutation represent the graphQLMutation
type GraphQLMutation struct {
	UserUsecase domain.UserUsecase
	TaskUsecase domain.TaskUsecase
}

// NewGraphQLMutation will initialize the user endpoint
func NewGraphQLMutation(u domain.UserUsecase, t domain.TaskUsecase) *GraphQLMutation {
	return &GraphQLMutation{
		UserUsecase: u,
		TaskUsecase: t,
	}
}

// GetRootMutationFields returns all the available mutations.
func (gm *GraphQLMutation) GetRootMutationFields() graphql.Fields {
	return graphql.Fields{
		"userLogin":    gm.GetLoginUserMutation(),
		"userRegister": gm.GetRegisterUserMutation(),
		"taskStore":    gm.GetStoreTaskMutation(),
		"taskUpdate":   gm.GetUpdateTaskMutation(),
		"taskDelete":   gm.GetDeleteTaskMutation(),
	}
}
