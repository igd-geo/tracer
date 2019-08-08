package dgraph

import (
	"context"
	"encoding/json"
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

type Client struct {
	conn *dgo.Dgraph
}

type query struct {
	text      string
	variables map[string]string
}

type decode struct {
	Entity   []provutil.Entity   `json:"entity,omitempty"`
	Agent    []provutil.Agent    `json:"agent,omitempty"`
	Activity []provutil.Activity `json:"activity,omitempty"`
}

/*
type Entity struct {
	UID            string    `json:"uid,omitempty"`
	ID             string    `json:"id,omitempty"`
	URI            string    `json:"uri,omitempty"`
	Name           string    `json:"name,omitempty"`
	CreationDate   string    `json:"creationDate,omitempty"`
	WasDerivedFrom []*Entity `json:"wasDerivedFrom,omitempty"`
	WasGeneratedBy *Activity `json:"wasGeneratedBy,omitempty"`
}

type Activity struct {
	UID               string    `json:"uid,omitempty"`
	ID                string    `json:"id,omitempty"`
	StartDate         string    `json:"startDate,omitempty"`
	EndDate           string    `json:"endDate,omitempty"`
	WasAssociatedWith *Agent    `json:"wasAssociatedWith,omitempty"`
	Used              []*Entity `json:"used,omitempty"`
}

type Agent struct {
	UID             string `json:"uid,omitempty"`
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	ActedOnBehalfOf *Agent `json:"actedOnBehalfOf,omitempty"`
}
*/

func NewClient(dgraphURL string) *Client {
	d, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		conn: dgo.NewDgraphClient(
			api.NewDgraphClient(d),
		),
	}
}

func (c *Client) AddDerivate(derivate *provutil.Entity) (map[string]string, error) {
	activity := derivate.WasGeneratedBy
	agent := derivate.WasGeneratedBy.WasAssociatedWith
	supervisor := derivate.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf

	entityData, activityData, agentData, supervisorData := derivate.Data, activity.Data, agent.Data, supervisor.Data
	derivate.Data, activity.Data, agent.Data, supervisor.Data = nil, nil, nil, nil

	payload, err := json.Marshal(derivate)
	if err != nil {
		return nil, err
	}

	assigned, err := c.runMutation(payload)
	if err != nil {
		return nil, err
	}
	derivate.Data, activity.Data, agent.Data, supervisor.Data = entityData, activityData, agentData, supervisorData

	return assigned.GetUids(), nil
}

func (c *Client) QueryParentEntity(id string, revision string) *provutil.Entity {
	query := query{
		text: `
		query entity($id: string, $revision: int){
			entity(func:eq(id, $id)) @filter(eq(revision, $revision)) {
				uid
			}
		}`,
		variables: map[string]string{"$id": id, "$revision": revision},
	}

	var decode decode
	c.runQuery(&decode, query)

	if len(decode.Entity) == 0 {
		return nil
	}
	return &decode.Entity[0]
}

func (c *Client) QueryAgentByName(name string) *provutil.Agent {
	query := query{
		text: `
		query agent($name: string){
			agent(func:allofterms(name, $name)) {
				uid
			}
		}`,
		variables: map[string]string{"$name": name},
	}

	var decode decode
	c.runQuery(&decode, query)

	if len(decode.Agent) == 0 {
		return nil
	}
	return &decode.Agent[0]
}

func (c *Client) QueryEntityByID(id string) *provutil.Entity {
	query := query{
		text: `
		query entity($id: string){
			entity(func:eq(id, $id)) {
				uid
			}
		}`,
		variables: map[string]string{"$id": id},
	}

	var decode decode
	c.runQuery(&decode, query)

	if len(decode.Entity) == 0 {
		return nil
	}
	return &decode.Entity[0]
}

func (c *Client) runQuery(decode *decode, query query) error {
	txn := c.conn.NewTxn()
	defer txn.Discard(context.Background())

	resp, err := txn.QueryWithVars(context.Background(), query.text, query.variables)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp.Json, decode)
	return nil
}

func (c *Client) runMutation(payload []byte) (*api.Assigned, error) {
	txn := c.conn.NewTxn()
	defer txn.Discard(context.Background())

	mu := &api.Mutation{
		CommitNow: true,
	}

	mu.SetJson = payload
	assigned, err := txn.Mutate(context.Background(), mu)
	if err != nil {
		return nil, err
	}
	return assigned, nil
}
