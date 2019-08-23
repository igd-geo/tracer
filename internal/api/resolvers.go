package api

import (
	"fmt"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
	"github.com/graphql-go/graphql"
)

func resolveQueryEntity(dbClient *db.Client, p graphql.ResolveParams) (*util.Entity, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf(`Field "id" not set`)
	}

	query := db.NewQuery(db.QueryEntityUIDByID)
	query.SetVariable(db.VariableEntityID, id)
	fmt.Printf("%+v", query)

	res, err := dbClient.RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Entity) != 1 {
		return nil, fmt.Errorf("Entity %s not found", id)
	}

	return res.Entity[0], nil
}

func resolveQueryAgent(dbClient *db.Client, p graphql.ResolveParams) (*util.Agent, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf(`Field "id" not set`)
	}

	query := db.NewQuery(db.QueryAgentUIDByID)
	query.SetVariable(db.VariableAgentID, id)

	res, err := dbClient.RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Agent) != 1 {
		return nil, fmt.Errorf("Agent %s not found", id)
	}

	return res.Agent[0], nil
}

func resolveQueryActivity(dbClient *db.Client, p graphql.ResolveParams) (*util.Activity, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf(`Field "id" not set`)
	}

	query := db.NewQuery(db.QueryActivityUIDByID)
	query.SetVariable(db.VariableActivityID, id)

	res, err := dbClient.RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Activity) != 1 {
		return nil, fmt.Errorf("Activity %s not found", id)
	}

	return res.Activity[0], nil
}

func resolveProvGraph(dbClient *db.Client, p graphql.ResolveParams) (*util.Graph, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf(`Field "id" not set`)
	}

	query := db.NewQuery(db.QueryProvenanceGraph)
	query.SetVariable(db.VariableGraphRootID, id)

	res, err := dbClient.RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Graph) != 1 {
		return nil, fmt.Errorf("%s is no valid graph root", id)
	}

	return res.Graph[0], nil
}
