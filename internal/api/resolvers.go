package api

import (
	"fmt"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"github.com/graphql-go/graphql"
)

func resolveQueryEntity(db provutil.InfoDB, p graphql.ResolveParams) (*provutil.Entity, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("Field \"id\" not set")
	}

	entity := db.FetchEntity(id)
	if entity == nil {
		return nil, fmt.Errorf("Entity %s not found!", id)
	}
	return entity, nil
}

func resolveQueryAgent(db provutil.InfoDB, p graphql.ResolveParams) (*provutil.Agent, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("Field \"id\" not set")
	}
	agent := db.FetchAgent(id)
	if agent == nil {
		return nil, fmt.Errorf("Agent %s not found!", id)
	}
	return agent, nil
}

func resolveQueryActivity(db provutil.InfoDB, p graphql.ResolveParams) (*provutil.Activity, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("Field \"id\" not set")
	}
	activity := db.FetchActivity(id)
	if activity == nil {
		return nil, fmt.Errorf("Activity %s not found!", id)
	}
	return activity, nil
}

func resolveProvInfo(idb provutil.InfoDB, gdb provutil.ProvDB, p graphql.ResolveParams) (*provutil.Entity, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("Field \"id\" not set")
	}

	entity := idb.FetchEntity(id)
	if entity == nil {
		return nil, fmt.Errorf("Entity %s not found!", id)
	}
	entity.Graph = gdb.FetchProvenanceGraph(entity.UID)

	return entity, nil
}
