package api

import (
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"github.com/graphql-go/graphql"
)

func initGraphQL(infoDB provutil.InfoDB, provDB provutil.ProvDB) graphql.Schema {
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
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return resolveQueryEntity(infoDB, p)
					},
				},
				"agent": &graphql.Field{
					Type:        agentType,
					Description: "Get agent by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return resolveQueryAgent(infoDB, p)
					},
				},
				"activity": &graphql.Field{
					Type:        activityType,
					Description: "Get activity by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return resolveQueryActivity(infoDB, p)
					},
				},
				"prov": &graphql.Field{
					Type:        entityType,
					Description: "Get provenance info of object",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return resolveProvInfo(infoDB, provDB, p)
					},
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
