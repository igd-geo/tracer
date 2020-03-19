package api

import (
	"encoding/json"
	"fmt"
	"reflect"

	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/database"
	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/provenance"
	"github.com/graphql-go/graphql"
)

const (
	nodeTypeRoot              = "root"
	nodeTypeEntity            = "entity"
	nodeTypeAgent             = "agent"
	nodeTypeActivity          = "activity"
	edgeTypeWasGeneratedBy    = "wasGeneratedBy"
	edgeTypeWasAssociatedWith = "wasAssociatedWith"
	edgeTypeActedOnBehalfOf   = "actedOnBehalfOf"
	edgeTypeWasAttributedTo   = "wasAttributedTo"
	edgeTypeWasDerivedFrom    = "wasDerivedFrom"
	edgeTypeUsed              = "used"
)

func resolveQueryEntity(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf(`field "id" not set`)
	}

	query := database.NewQuery(database.QueryEntityFullByID)
	query.SetVariable(database.VariableEntityID, id)

	res, err := p.Context.Value(dbConnectionKey("db")).(*database.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Entity) != 1 {
		return nil, fmt.Errorf("entity %s not found", id)
	}

	return res.Entity[0], nil
}

func resolveQueryAgent(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf(`field "id" not set`)
	}

	query := database.NewQuery(database.QueryAgentFullByID)
	query.SetVariable(database.VariableAgentID, id)

	res, err := p.Context.Value(dbConnectionKey("db")).(*database.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Agent) != 1 {
		return nil, fmt.Errorf("agent %s not found", id)
	}

	return res.Agent[0], nil
}

func resolveQueryActivity(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf(`field "id" not set`)
	}

	query := database.NewQuery(database.QueryActivityFullByID)
	query.SetVariable(database.VariableActivityID, id)

	res, err := p.Context.Value(dbConnectionKey("db")).(*database.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Activity) != 1 {
		return nil, fmt.Errorf("activity %s not found", id)
	}

	return res.Activity[0], nil
}

func resolveWasGeneratedBy(p graphql.ResolveParams) (interface{}, error) {
	entity, ok := p.Source.(*provenance.Entity)
	if !ok {
		return nil, fmt.Errorf(`nested entity "id" not set`)
	}

	query := database.NewQuery(database.QueryWasGeneratedBy)
	query.SetVariable(database.VariableEntityID, entity.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*database.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Activity) != 1 {
		return nil, fmt.Errorf("WasGeneratedBy not found for %s", entity.UID)
	}

	return res.Activity[0], nil
}

func resolveWasDerivedFrom(p graphql.ResolveParams) (interface{}, error) {
	entity, ok := p.Source.(*provenance.Entity)
	if !ok {
		return nil, fmt.Errorf(`nested entity "id" not set`)
	}

	query := database.NewQuery(database.QueryWasDerivedFrom)
	query.SetVariable(database.VariableEntityID, entity.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*database.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Entity) < 1 {
		return nil, nil
	}

	return res.Entity, nil
}

func resolveWasAssociatedWith(p graphql.ResolveParams) (interface{}, error) {
	activity, ok := p.Source.(*provenance.Activity)
	if !ok {
		return nil, fmt.Errorf(`nested activity "id" not set`)
	}

	query := database.NewQuery(database.QueryWasAssociatedWith)
	query.SetVariable(database.VariableActivityID, activity.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*database.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Agent) != 1 {
		return nil, nil
	}

	return res.Agent[0], nil
}

func resolveUsed(p graphql.ResolveParams) (interface{}, error) {
	activity, ok := p.Source.(*provenance.Activity)
	if !ok {
		return nil, fmt.Errorf(`nested activity "id" not set`)
	}

	query := database.NewQuery(database.QueryUsed)
	query.SetVariable(database.VariableActivityID, activity.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*database.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Entity) < 1 {
		return nil, nil
	}

	return res.Entity, nil
}

func resolveActedOnBehalfOf(p graphql.ResolveParams) (interface{}, error) {
	agent, ok := p.Source.(*provenance.Agent)
	if !ok {
		return nil, fmt.Errorf(`nested agent "id" not set`)
	}
	query := database.NewQuery(database.QueryActedOnBehalfOf)
	query.SetVariable(database.VariableAgentID, agent.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*database.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Agent) != 1 {
		return nil, nil
	}

	return res.Agent[0], nil
}

func resolveProvGraph(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf(`field "id" not set`)
	}

	query := database.NewQuery(database.QueryProvenanceGraph)
	query.SetVariable(database.VariableGraphRootID, id)

	res, err := p.Context.Value(dbConnectionKey("db")).(*database.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Graph) != 1 {
		return nil, fmt.Errorf("%s is no valid graph root", id)
	}

	return res.Graph[0], nil
}

func parseGraph(graph *provenance.Graph) error {
	var nodes []provenance.Node
	var edges []provenance.Edge

	var v map[string]interface{}
	err := json.Unmarshal(graph.RawMessage, &v)
	if err != nil {
		return err
	}

	reflectFields(&nodes, &edges, nodeTypeRoot, v)

	graph.Nodes, err = json.Marshal(nodes)
	if err != nil {
		return err
	}

	graph.Edges, err = json.Marshal(edges)
	if err != nil {
		return err
	}

	return nil
}

func reflectFields(nodes *[]provenance.Node, edges *[]provenance.Edge, nodeType string, v map[string]interface{}) {
	values := reflect.ValueOf(v)
	node := make(provenance.Node)
	node["nodeType"] = nodeType

	id, ok := v["id"].(string)
	if !ok {
		return
	}

	iter := values.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		if value.IsZero() {
			continue
		}

		switch key.String() {
		case edgeTypeWasGeneratedBy:
			activity := value.Interface().([]interface{})[0].(map[string]interface{})

			activityID, ok := activity["id"].(string)
			if !ok {
				return
			}

			*edges = append(*edges, drawEdge(id, activityID, edgeTypeWasGeneratedBy))
			reflectFields(nodes, edges, nodeTypeActivity, activity)

		case edgeTypeWasAssociatedWith:
			agent := value.Interface().([]interface{})[0].(map[string]interface{})

			agentID, ok := agent["id"].(string)
			if !ok {
				return
			}

			*edges = append(*edges, drawEdge(id, agentID, edgeTypeWasAssociatedWith))
			reflectFields(nodes, edges, nodeTypeAgent, agent)

		case edgeTypeActedOnBehalfOf:
			agent := value.Interface().([]interface{})[0].(map[string]interface{})

			agentID, ok := agent["id"].(string)
			if !ok {
				return
			}

			*edges = append(*edges, drawEdge(id, agentID, edgeTypeActedOnBehalfOf))
			reflectFields(nodes, edges, nodeTypeAgent, agent)

		case edgeTypeUsed:
			entities := value.Interface().([]interface{})

			for _, entity := range entities {
				entityID, ok := entity.(map[string]interface{})["id"].(string)
				if !ok {
					return
				}

				*edges = append(*edges, drawEdge(id, entityID, edgeTypeUsed))
				reflectFields(nodes, edges, nodeTypeEntity, entity.(map[string]interface{}))
			}

		case edgeTypeWasDerivedFrom:
			entities := value.Interface().([]interface{})

			for _, entity := range entities {
				entityID, ok := entity.(map[string]interface{})["id"].(string)
				if !ok {
					return
				}

				*edges = append(*edges, drawEdge(id, entityID, edgeTypeWasDerivedFrom))
				reflectFields(nodes, edges, nodeTypeEntity, entity.(map[string]interface{}))
			}

		default:
			node[key.String()] = value.Interface().(string)
		}
	}
	*nodes = append(*nodes, node)
}

func drawEdge(source string, target string, edgeType string) provenance.Edge {
	return provenance.Edge{
		Source:   source,
		Target:   target,
		EdgeType: edgeType,
	}
}
