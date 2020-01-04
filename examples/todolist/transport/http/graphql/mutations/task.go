package mutations

import (
	"context"
	"time"
	"todolist/domain"
	"todolist/middleware"
	"todolist/transport/http/graphql/types"

	"github.com/graphql-go/graphql"
)

// GetStoreTaskMutation store task.
func (gm *GraphQLMutation) GetStoreTaskMutation() *graphql.Field {
	return &graphql.Field{
		Type:        types.TaskType,
		Description: "Create new task",
		Args: graphql.FieldConfigArgument{
			"title": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"due_date": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"completed": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// decode process
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}

			task := domain.Task{
				Title:       params.Args["title"].(string),
				Description: params.Args["description"].(string),
				DueDate:     params.Args["due_date"].(time.Time),
			}
			if completed, completedOk := params.Args["completed"].(bool); completedOk {
				task.Completed = completed
			}

			// get token
			tokenHeader := ctx.Value("token").(string)
			token, err := middleware.JwtVerify(tokenHeader)
			if err != nil {
				return nil, err
			}

			// validation request
			if ok, err := middleware.IsRequestValid(&task); !ok {
				return nil, err
			}

			// store task
			resp, err := gm.TaskUsecase.Store(ctx, token.ID, &task)
			if err != nil {
				return nil, err
			}

			return resp, nil
		},
	}
}

// GetUpdateTaskMutation for update task.
func (gm *GraphQLMutation) GetUpdateTaskMutation() *graphql.Field {
	return &graphql.Field{
		Type:        types.TaskType,
		Description: "Update a task",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"title": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"due_date": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"completed": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// decode process
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}

			task := domain.Task{
				ID:          uint64(params.Args["id"].(int)),
				Title:       params.Args["title"].(string),
				Description: params.Args["description"].(string),
				DueDate:     params.Args["due_date"].(time.Time),
				Completed:   params.Args["completed"].(bool),
			}

			// get token
			tokenHeader := ctx.Value("token").(string)
			token, err := middleware.JwtVerify(tokenHeader)
			if err != nil {
				return nil, err
			}
			task.UserID = token.ID

			// validation request
			if ok, err := middleware.IsRequestValid(&task); !ok {
				return nil, err
			}

			// update task
			err = gm.TaskUsecase.Update(ctx, &task)
			if err != nil {
				return nil, err
			}

			return task, nil
		},
	}
}

// GetDeleteTaskMutation for delete task.
func (gm *GraphQLMutation) GetDeleteTaskMutation() *graphql.Field {
	return &graphql.Field{
		Type:        types.TaskType,
		Description: "Delete a task",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// decode process
			ctx := params.Context
			if ctx == nil {
				ctx = context.Background()
			}

			task := &domain.Task{
				ID: uint64(params.Args["id"].(int)),
			}

			// get token
			tokenHeader := ctx.Value("token").(string)
			_, err := middleware.JwtVerify(tokenHeader)
			if err != nil {
				return nil, err
			}

			// delete task
			err = gm.TaskUsecase.Delete(ctx, uint64(params.Args["id"].(int)))
			if err != nil {
				return nil, err
			}

			return task, nil
		},
	}
}
