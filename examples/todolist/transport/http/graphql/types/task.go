package types

import (
	"github.com/graphql-go/graphql"
)

// TaskType is the GraphQL schema for the task type.
var TaskType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Task",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"due_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"completed": &graphql.Field{
				Type: graphql.Boolean,
			},
			"user_id": &graphql.Field{
				Type: graphql.ID,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
