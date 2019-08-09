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

type result struct {
	Entity   []provutil.Entity   `json:"entity,omitempty"`
	Agent    []provutil.Agent    `json:"agent,omitempty"`
	Activity []provutil.Activity `json:"activity,omitempty"`
}

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

func (c *Client) FetchProvenanceGraph(uid string) *json.RawMessage {
	query := `
		query entity($id: string) {
  			entity(func: uid($id)) {
    			expand(_all_) {
      				expand(_all_) {
        				expand(_all_) {
        					expand(_all_)
      					}
      				}
    			}
  			}
		}`
	variables := map[string]string{"$id": uid}

	res, err := c.runQuery(query, variables)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &res
}

func (c *Client) runQuery(query string, variables map[string]string) (json.RawMessage, error) {
	txn := c.conn.NewTxn()
	defer txn.Discard(context.Background())

	resp, err := txn.QueryWithVars(context.Background(), query, variables)
	if err != nil {
		return nil, err
	}

	return resp.Json, nil
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
