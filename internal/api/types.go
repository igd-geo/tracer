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
			"Attributes": &graphql.Field{
				Type: attributesType,
			},
			"Data": &graphql.Field{
				Type: jsonRaw,
			},
			"Graph": &graphql.Field{
				Type: jsonRaw,
			},
		},
	},
)

var agentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Attributes",
		Description: "Provenance Agent Object",
		Fields: graphql.Fields{
			"Attributes": &graphql.Field{
				Type: attributesType,
			},
			"Data": &graphql.Field{
				Type: jsonRaw,
			},
			"Graph": &graphql.Field{
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
			"Attributes": &graphql.Field{
				Type: attributesType,
			},
			"Data": &graphql.Field{
				Type: jsonRaw,
			},
			"Graph": &graphql.Field{
				Type: jsonRaw,
			},
		},
	},
)

var attributesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Attributes",
		Description: "Attribute collection",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: graphql.String,
			},
			"URI": &graphql.Field{
				Type: graphql.String,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"CreationDate": &graphql.Field{
				Type: graphql.String,
			},
			"StartDate": &graphql.Field{
				Type: graphql.String,
			},
			"EndDate": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var edgesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Edges",
		Description: "Edge collection",
		Fields: graphql.Fields{
			"WasGeneratedBy": &graphql.Field{},
		},
	},
)

var objectType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Object",
		Fields: graphql.Fields{
			"Attributes": &graphql.Field{
				Type: attributesType,
			},
		},
	},
)

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
