package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	Entity   []Entity   `json:"entity,omitempty"`
	Agent    []Agent    `json:"agent,omitempty"`
	Activity []Activity `json:"activity,omitempty"`
}

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

func NewEntity() *Entity {
	return &Entity{
		UID: "_:entity",
		WasGeneratedBy: &Activity{
			UID: "_:activity",
			WasAssociatedWith: &Agent{
				UID: "_:agent",
				ActedOnBehalfOf: &Agent{
					UID: "_:supervisor",
				},
			},
		},
	}
}

func NewAgent(uid string) *Agent {
	return &Agent{UID: uid}
}

func NewActivity(uid string) *Activity {
	return &Activity{}
}

func (c *Client) AddDerivate(derivate *Entity) (map[string]string, error) {
	payload, err := json.Marshal(derivate)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(payload))

	assigned, err := c.runMutation(payload)
	if err != nil {
		return nil, err
	}

	return assigned.GetUids(), nil
}

func (c *Client) QueryParentEntity(id string, revision string) *Entity {
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
	fmt.Println(decode.Entity)
	return &decode.Entity[0]
}

func (c *Client) QueryAgentByName(name string) *Agent {
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

func (c *Client) QueryEntityByID(id string) *Entity {
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
		fmt.Println(err)
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
