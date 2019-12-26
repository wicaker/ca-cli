package mutations

import (
	"context"
	"todolist/domain"
	"todolist/transport/http/graphql/types"

	"github.com/graphql-go/graphql"
)

// GetLoginUserMutation creates a new user and returns it.
func (gm *GraphQLMutation) GetLoginUserMutation() *graphql.Field {
	return &graphql.Field{
		Type:        types.UserType,
		Description: "Register new user",
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"email": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"username": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// decode process
			user := &domain.User{
				Password: params.Args["password"].(string),
			}

			if email, emailOk := params.Args["email"].(string); emailOk {
				user.Email = email
			}
			if username, usernameOk := params.Args["username"].(string); usernameOk {
				user.Username = username
			}

			// login usecase process
			ctx := context.Background()

			type data struct {
				Token string `json:"token"`
			}
			d := &data{}
			token, err := gm.UserUsecase.Login(ctx, user)
			if err != nil {
				return nil, err
			}

			// encode and return response
			d.Token = token
			return d, nil
		},
	}
}
