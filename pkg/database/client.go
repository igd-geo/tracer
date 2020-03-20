package database

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/provenance"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

// Config contains the database configuration
type Config struct {
	URL string `yaml:"url"`
}

// Client is a wrapper for a Dgraph connection
type Client struct {
	dgraph *dgo.Dgraph
}

// Transaction bundles multiple entities to a mutation
type Transaction struct {
	Mutation []*provenance.Entity
	Size     int
}

// Result contains the results of a query
type Result struct {
	Entity     []*provenance.Entity   `json:"entity,omitempty"`
	Activity   []*provenance.Activity `json:"activity,omitempty"`
	Agent      []*provenance.Agent    `json:"agent,omitempty"`
	Supervisor []*provenance.Agent    `json:"supervisor,omitempty"`
	Graph      []*provenance.Graph    `json:"graph,omitempty"`
}

// New returns a new database client
func New(config *Config) *Client {
	grpcClient, err := grpc.Dial(config.URL, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		dgraph: dgo.NewDgraphClient(
			api.NewDgraphClient(grpcClient),
		),
	}
}

// NewTransaction returns an empty Transaction
func (c *Client) NewTransaction() *Transaction {
	transaction := &Transaction{
		Mutation: []*provenance.Entity{},
		Size:     0,
	}
	return transaction
}

// RunQuery runs the query and returns a result
func (c *Client) RunQuery(query string) (Result, error) {
	var res Result

	resp, err := c.dgraph.NewReadOnlyTxn().Query(context.TODO(), query)
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

	resp, err := c.dgraph.NewReadOnlyTxn().QueryWithVars(context.TODO(), query.queryString, query.variables)
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
func (c *Client) RunMutation(mutation []*provenance.Entity) (*api.Assigned, error) {
	txn := c.dgraph.NewTxn()
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

// Add appends new entities to an exisiting transaction
func (t *Transaction) Add(entity *provenance.Entity) {
	t.Mutation = append(t.Mutation, entity)
	t.Size++
}
