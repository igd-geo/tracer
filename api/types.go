package api

import (
	"encoding/json"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var entityType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Entity",
		Description: "Provenance Entity Object",
		Fields: graphql.Fields{
			"uid": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"uri": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"creationDate": &graphql.Field{
				Type: graphql.String,
			},
			"data": &graphql.Field{
				Type: jsonRaw,
			},
			"graph": &graphql.Field{
				Type: jsonRaw,
			},
		},
	},
)

var agentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Agent",
		Description: "Provenance Agent Object",
		Fields: graphql.Fields{
			"uid": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"data": &graphql.Field{
				Type: jsonRaw,
			},
			"graph": &graphql.Field{
				Type: jsonRaw,
			},
		},
	},
)

var activityType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Activity",
		Description: "Provenance Activity Object",
		Fields: graphql.Fields{
			"uid": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"startDate": &graphql.Field{
				Type: graphql.String,
			},
			"endDate": &graphql.Field{
				Type: graphql.String,
			},
			"data": &graphql.Field{
				Type: jsonRaw,
			},
			"graph": &graphql.Field{
				Type: jsonRaw,
			},
		},
	},
)

var graphType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Graph",
		Description: "Provenance Graph Object",
		Fields: graphql.Fields{
			"json": &graphql.Field{
				Type: jsonRaw,
			},
			"nodes": &graphql.Field{
				Type: jsonRaw,
			},
			"edges": &graphql.Field{
				Type: jsonRaw,
			},
		},
	},
)

/*
var nodeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Node Type",
		Description: "Provenance Graph Node",
		Fields: graphql.Fields{
			"nodeType": &graphql.Field{
				Type: graphql.String,
			},
			"uid": &graphql.Field{
				Type: graphql.String,
			},
			"ID": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"startDate": &graphql.Field{
				Type: graphql.String,
			},
			"endDate": &graphql.Field{
				Type: graphql.String,
			},
			"uri": &graphql.Field{
				Type: graphql.String,
			},
			"creationDate": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var edgeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Edge Type",
		Description: "Provenance Graph Edge",
		Fields: graphql.Fields{
			"edgeType": &graphql.Field{
				Type: graphql.String,
			},
			"source": &graphql.Field{
				Type: graphql.String,
			},
			"target": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
*/
var jsonRaw = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "JSON",
	Description: "Raw JSON Byte Array",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case json.RawMessage:
			return value
		case *json.RawMessage:
			return *value
		default:
			return nil
		}
	},
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			return json.RawMessage(value)
		case *string:
			return json.RawMessage(*value)
		default:
			return nil
		}
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return json.RawMessage(valueAST.Value)
		default:
			return nil
		}
	},
})
