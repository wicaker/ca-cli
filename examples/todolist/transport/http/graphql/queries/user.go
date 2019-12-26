package queries

import (
	"todolist/domain"
	"todolist/transport/http/graphql/types"

	"github.com/graphql-go/graphql"
	log "github.com/sirupsen/logrus"
)

// GetUserQuery returns the queries available against user type.
func (gq *GraphQLQuery) GetUserQuery() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(types.UserType),
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			log.Printf("[query] user\n")
			// var users []domain.User
			users := []domain.User{
				{
					ID:   1,
					Name: "test",
				},
			}
			return users, nil
		},
	}
}
