package api

import (
	"fmt"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"github.com/graphql-go/graphql"
)

func resolveQueryEntity(db infoDB, p graphql.ResolveParams) (*provutil.Entity, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("TODO")
	}
	return db.FetchEntity(id), nil
}

func resolveQueryAgent(db infoDB, p graphql.ResolveParams) (*provutil.Agent, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("TODO")
	}
	return db.FetchAgent(id), nil
}

func resolveQueryActivity(db infoDB, p graphql.ResolveParams) (*provutil.Activity, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("TODO")
	}
	return db.FetchActivity(id), nil
}

func resolveProvInfo(idb infoDB, gdb graphDB, p graphql.ResolveParams) (*provutil.Entity, error) {
	id, ok := p.Args["ID"].(string)
	if !ok {
		return nil, fmt.Errorf("TODO")
	}

	entity := idb.FetchEntity(id)
	entity.Graph = gdb.FetchProvenanceGraph(entity.UID)

	return entity, nil
}
