package queries

import (
	"context"
	"strconv"
	"todolist/domain"
	"todolist/middleware"
	"todolist/transport/http/graphql/types"

	"github.com/graphql-go/graphql"
	log "github.com/sirupsen/logrus"
)

// GetTaskQuery returns the queries available against task type.
func (gq *GraphQLQuery) GetTaskQuery() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(types.TaskType),
		Description: "Get task data",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.ID,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			log.Printf("[query] task\n")
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}

			tasks := []*domain.Task{}
			tokenHeader := ctx.Value("token").(string)
			token, err := middleware.JwtVerify(tokenHeader)
			if err != nil {
				return nil, err
			}

			// if query contain arguments `id`
			if idTask, idTaskOk := params.Args["id"].(string); idTaskOk {
				taskID, err := strconv.Atoi(idTask)
				if err != nil {
					return nil, err
				}

				res, err := gq.TaskUsecase.GetByID(ctx, uint64(taskID))
				if err != nil {
					return nil, err
				}
				tasks = append(tasks, res)
				return tasks, nil
			}

			// if not, just return array
			res, err := gq.TaskUsecase.Fetch(ctx, token.ID)
			if err != nil {
				return nil, err
			}

			if len(res) != 0 {
				tasks = res
			}
			return tasks, nil
		},
	}
}
