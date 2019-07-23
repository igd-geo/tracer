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
	Title          string    `json:"title,omitempty"`
	WasDerivedFrom *Entity   `json:"wasDerivedFrom,omitempty"`
	WasGeneratedBy *Activity `json:"wasGeneratedBy,omitempty"`
	/*
		Revision       int       `json:"revision,omitempty"`
		Type           string    `json:"type,omitempty"`
		Date           string    `json:"date,omitempty"`
	*/
}

type Activity struct {
	UID               string    `json:"uid,omitempty"`
	ID                string    `json:"id,omitempty"`
	Title             string    `json:"title,omitempty"`
	WasAssociatedWith *Agent    `json:"wasAssociatedWith,omitempty"`
	Used              []*Entity `json:"used,omitempty"`
	/*
		Type              string    `json:"type,omitempty"`
		StartTime         string    `json:"startTime,omitempty"`
		EndTime           string    `json:"endTime,omitempty"`
	*/
}

type Agent struct {
	UID             string `json:"uid,omitempty"`
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	ActedOnBehalfOf *Agent `json:"actedOnBehalfOf,omitempty"`
	/*
		Telephone          string `json:"phone,omitempty"`
		Facsimile          string `json:"fax,omitempty"`
		Address            string `json:"address,omitempty"`
		City               string `json:"city,omitempty"`
		AdministrativeArea string `json:"area,omitempty"`
		PostalCode         string `json:"postal,omitempty"`
		Country            string `json:"country,omitempty"`
		Email              string `json:"email,omitempty"`
		Role               string `json:"role,omitempty"`
	*/
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

func (c *Client) AddDerivate(derivate *Entity) error {
	payload, err := json.MarshalIndent(derivate, "", " ")
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = c.runMutation(payload)
	if err != nil {
		log.Fatal(err)
	}

	return nil
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
