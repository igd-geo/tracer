package mongodb

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	conn        *mongo.Client
	collections Collections
}

type Collections struct {
	entity   string
	agent    string
	activity string
}

func NewClient(mongoURL, mongoDatabase, collEntity, collAgent, collActivity string) *Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		conn: client,
		collections: Collections{
			entity:   collEntity,
			agent:    collAgent,
			activity: collActivity,
		},
	}
}

func (client *Client) EntityUID(id string) (string, error) {
	return "", nil
}

func (client *Client) AgentUID(id string) (string, error) {
	return "", nil
}

func (client *Client) ActivitytUID(id string) (string, error) {
	return "", nil
}

func (client *Client) InsertEntity(uid string, payload json.RawMessage) error {
	return nil
}

func (client *Client) InsertAgent(uid string, payload json.RawMessage) error {
	return nil
}

func (client *Client) InsertActivity(uid string, payload json.RawMessage) error {
	return nil
}

func (client *Client) FetchEntity() error {
	return nil
}

func (client *Client) FetchAgent() error {
	return nil
}

func (client *Client) FetchActivity() error {
	return nil
}
