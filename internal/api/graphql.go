package api

import (
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"github.com/graphql-go/graphql"
)

func initGraphQL(db *db.Client) graphql.Schema {
	entityType.AddFieldConfig(
		"wasDerivedFrom",
		&graphql.Field{
			Type:    graphql.NewList(entityType),
			Resolve: resolveWasDerivedFrom,
		},
	)

	entityType.AddFieldConfig(
		"wasGeneratedBy",
		&graphql.Field{
			Type:    activityType,
			Resolve: resolveWasGeneratedBy,
		},
	)

	agentType.AddFieldConfig(
		"actedOnBehalfOf",
		&graphql.Field{
			Type:    agentType,
			Resolve: resolveActedOnBehalfOf,
		},
	)

	activityType.AddFieldConfig(
		"wasAssociatedWith",
		&graphql.Field{
			Type:    agentType,
			Resolve: resolveWasAssociatedWith,
		},
	)

	activityType.AddFieldConfig(
		"used",
		&graphql.Field{
			Type:    graphql.NewList(entityType),
			Resolve: resolveUsed,
		},
	)

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"entity": &graphql.Field{
					Type:        entityType,
					Description: "Get entity by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolveQueryEntity,
					//Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					//return resolveQueryEntity(db, p)
					//},
				},
				"agent": &graphql.Field{
					Type:        agentType,
					Description: "Get agent by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolveQueryAgent,
					// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// 	return resolveQueryAgent(db, p)
					// },
				},
				"activity": &graphql.Field{
					Type:        activityType,
					Description: "Get activity by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolveQueryActivity,
					// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// 	return resolveQueryActivity(db, p)
					// },
				},
				"graph": &graphql.Field{
					Type:        graphType,
					Description: "Get provenance graph of object",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolveProvGraph,
					// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// 	return resolveProvGraph(db, p)
					// },
				},
			},
		})
	schema, _ := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: nil,
		})
	return schema
}
