package api

import (
	"encoding/json"
	"fmt"
	"reflect"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
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

	query := db.NewQuery(db.QueryEntityFullByID)
	query.SetVariable(db.VariableEntityID, id)

	res, err := p.Context.Value(dbConnectionKey("db")).(*db.Client).RunQueryWithVars(query)
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

	query := db.NewQuery(db.QueryAgentFullByID)
	query.SetVariable(db.VariableAgentID, id)

	res, err := p.Context.Value(dbConnectionKey("db")).(*db.Client).RunQueryWithVars(query)
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

	query := db.NewQuery(db.QueryActivityFullByID)
	query.SetVariable(db.VariableActivityID, id)

	res, err := p.Context.Value(dbConnectionKey("db")).(*db.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Activity) != 1 {
		return nil, fmt.Errorf("activity %s not found", id)
	}

	return res.Activity[0], nil
}

func resolveWasGeneratedBy(p graphql.ResolveParams) (interface{}, error) {
	entity, ok := p.Source.(*util.Entity)
	if !ok {
		return nil, fmt.Errorf(`nested entity "id" not set`)
	}

	query := db.NewQuery(db.QueryWasGeneratedBy)
	query.SetVariable(db.VariableEntityID, entity.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*db.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Activity) != 1 {
		return nil, fmt.Errorf("WasGeneratedBy not found for %s", entity.UID)
	}

	return res.Activity[0], nil
}

func resolveWasDerivedFrom(p graphql.ResolveParams) (interface{}, error) {
	entity, ok := p.Source.(*util.Entity)
	if !ok {
		return nil, fmt.Errorf(`nested entity "id" not set`)
	}

	query := db.NewQuery(db.QueryWasDerivedFrom)
	query.SetVariable(db.VariableEntityID, entity.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*db.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Entity) < 1 {
		return nil, nil
	}

	return res.Entity, nil
}

func resolveWasAssociatedWith(p graphql.ResolveParams) (interface{}, error) {
	activity, ok := p.Source.(*util.Activity)
	if !ok {
		return nil, fmt.Errorf(`nested activity "id" not set`)
	}

	query := db.NewQuery(db.QueryWasAssociatedWith)
	query.SetVariable(db.VariableActivityID, activity.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*db.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Agent) != 1 {
		return nil, nil
	}

	return res.Agent[0], nil
}

func resolveUsed(p graphql.ResolveParams) (interface{}, error) {
	activity, ok := p.Source.(*util.Activity)
	if !ok {
		return nil, fmt.Errorf(`nested activity "id" not set`)
	}

	query := db.NewQuery(db.QueryUsed)
	query.SetVariable(db.VariableActivityID, activity.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*db.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Entity) < 1 {
		return nil, nil
	}

	return res.Entity, nil
}

func resolveActedOnBehalfOf(p graphql.ResolveParams) (interface{}, error) {
	agent, ok := p.Source.(*util.Agent)
	if !ok {
		return nil, fmt.Errorf(`nested agent "id" not set`)
	}
	query := db.NewQuery(db.QueryActedOnBehalfOf)
	query.SetVariable(db.VariableAgentID, agent.UID)
	res, err := p.Context.Value(dbConnectionKey("db")).(*db.Client).RunQueryWithVars(query)
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

	query := db.NewQuery(db.QueryProvenanceGraph)
	query.SetVariable(db.VariableGraphRootID, id)

	res, err := p.Context.Value(dbConnectionKey("db")).(*db.Client).RunQueryWithVars(query)
	if err != nil {
		return nil, err
	}

	if len(res.Graph) != 1 {
		return nil, fmt.Errorf("%s is no valid graph root", id)
	}

	/*
		err = parseGraph(res.Graph[0])
		if err != nil {
			return nil, err
		}
	*/

	return res.Graph[0], nil
}

func parseGraph(graph *util.Graph) error {
	var nodes []util.Node
	var edges []util.Edge

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

func reflectFields(nodes *[]util.Node, edges *[]util.Edge, nodeType string, v map[string]interface{}) {
	values := reflect.ValueOf(v)
	node := make(util.Node)
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

func drawEdge(source string, target string, edgeType string) util.Edge {
	return util.Edge{
		Source:   source,
		Target:   target,
		EdgeType: edgeType,
	}
}
