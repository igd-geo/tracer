package api

import (
	"encoding/json"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"github.com/graphql-go/graphql"
)

type infoDB interface {
	FetchEntity(id string) *provutil.Entity
	FetchAgent(id string) *provutil.Agent
	FetchActivity(id string) *provutil.Activity
}

type graphDB interface {
	FetchProvenanceGraph(uid string) *json.RawMessage
}

func initGraphQL(graphDB graphDB, infoDB infoDB) graphql.Schema {
	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"entity": &graphql.Field{
					Type:        entityType,
					Description: "Get entity by id",
					Args: graphql.FieldConfigArgument{
						"ID": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return resolveQueryEntity(infoDB, p)
					},
				},
				"agent": &graphql.Field{
					Type:        objectType,
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
					Type:        objectType,
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
						"ID": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return resolveProvInfo(infoDB, graphDB, p)
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
