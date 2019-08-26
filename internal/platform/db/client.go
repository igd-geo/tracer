package db

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

// Client is a wrapper for a Dgraph connection
type Client struct {
	conn *dgo.Dgraph
}

// Transaction bundles multiple entities to a mutation
type Transaction struct {
	Mutation  []*util.Entity
	Size      int
	StartTime time.Time
	InCommit  bool
}

// Result contains the results of a query
type Result struct {
	Entity     []*util.Entity   `json:"entity,omitempty"`
	Activity   []*util.Activity `json:"activity,omitempty"`
	Agent      []*util.Agent    `json:"agent,omitempty"`
	Supervisor []*util.Agent    `json:"supervisor,omitempty"`
	Graph      []*util.Graph    `json:"graph,omitempty"`
}

// NewClient returns a new Client
func NewClient(url string) *Client {
	d, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		conn: dgo.NewDgraphClient(
			api.NewDgraphClient(d),
		),
	}
}

// NewTransaction returns an empty Transaction
func (c *Client) NewTransaction() *Transaction {
	return &Transaction{
		Mutation:  []*util.Entity{},
		Size:      0,
		StartTime: time.Now(),
	}
}

// RunQuery runs the query and returns a result
func (c *Client) RunQuery(query string) (Result, error) {
	var res Result

	resp, err := c.conn.NewReadOnlyTxn().Query(context.TODO(), query)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(resp.GetJson(), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// RunQueryWithVars runs the query
func (c *Client) RunQueryWithVars(query *Query) (Result, error) {
	var res Result

	resp, err := c.conn.NewReadOnlyTxn().QueryWithVars(context.TODO(), query.queryString, query.variables)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(resp.GetJson(), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// RunMutation runs a mutation
func (c *Client) RunMutation(mutation *[]*util.Entity) (*api.Assigned, error) {
	txn := c.conn.NewTxn()
	defer txn.Discard(context.TODO())

	payload, err := json.Marshal(mutation)
	if err != nil {
		return nil, err
	}

	mu := &api.Mutation{
		CommitNow: true,
	}

	mu.SetJson = payload
	assigned, err := txn.Mutate(context.TODO(), mu)
	if err != nil {
		return nil, err
	}
	return assigned, nil
}
