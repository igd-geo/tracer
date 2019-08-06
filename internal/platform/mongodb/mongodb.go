package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	conn        *mongo.Client
	database    string
	collections Collections
}

type Collections struct {
	entity   string
	agent    string
	activity string
}

type Entity struct {
	UID  string
	ID   string
	Data json.RawMessage
}
type Agent struct {
	UID  string
	ID   string
	Data json.RawMessage
}
type Activity struct {
	UID  string
	ID   string
	Data json.RawMessage
}

func NewClient(mongoURL, mongoDatabase, collEntity, collAgent, collActivity string) *Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		conn:     client,
		database: mongoDatabase,
		collections: Collections{
			entity:   collEntity,
			agent:    collAgent,
			activity: collActivity,
		},
	}
}

func (client *Client) EntityUID(id string) (string, error) {
	collection := client.conn.Database(client.database).Collection(client.collections.entity)
	if collection == nil {
		return "", fmt.Errorf("No Collection")
	}

	var result Entity
	filter := bson.D{{Key: "id", Value: id}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	return result.UID, err
}

func (client *Client) AgentUID(id string) (string, error) {
	collection := client.conn.Database(client.database).Collection(client.collections.agent)
	if collection == nil {
		return "", fmt.Errorf("No Collection")
	}

	var result Agent
	filter := bson.D{{Key: "id", Value: id}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	return result.UID, err
}

func (client *Client) ActivitytUID(id string) (string, error) {
	collection := client.conn.Database(client.database).Collection(client.collections.activity)
	if collection == nil {
		return "", fmt.Errorf("No Collection")
	}

	var result Activity
	filter := bson.D{{Key: "id", Value: id}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	return result.UID, err
}

func (client *Client) InsertEntity(uid string, payload json.RawMessage) error {
	collection := client.conn.Database(client.database).Collection(client.collections.entity)
	if collection == nil {
		return fmt.Errorf("No Collection")
	}

	_, err := collection.InsertOne(context.TODO(), bson.Raw(payload))
	return err
}

func (client *Client) InsertAgent(uid string, payload json.RawMessage) error {
	collection := client.conn.Database(client.database).Collection(client.collections.agent)
	if collection == nil {
		return fmt.Errorf("No Collection")
	}

	_, err := collection.InsertOne(context.TODO(), bson.Raw(payload))
	return err
}

func (client *Client) InsertActivity(uid string, payload json.RawMessage) error {
	collection := client.conn.Database(client.database).Collection(client.collections.activity)
	if collection == nil {
		return fmt.Errorf("No Collection")
	}

	_, err := collection.InsertOne(context.TODO(), bson.Raw(payload))
	return err
}

func (client *Client) FetchEntity(id string) (*Entity, error) {
	/*
		collection := client.conn.Database(client.database).Collection(client.collections.entity)
		if collection == nil {
			return result, fmt.Errorf("No Collection")
		}

		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			return result, err
		}

		return result, nil
	*/
	var result Entity
	filter := bson.D{{Key: "id", Value: id}}
	return &result, client.fetch(client.collections.entity, filter, &result)
}

func (client *Client) FetchAgent(id string) (*Agent, error) {
	/*
		collection := client.conn.Database(client.database).Collection(client.collections.agent)
		if collection == nil {
			return result, fmt.Errorf("No Collection")
		}

		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			return result, err
		}
	*/
	var result Agent
	filter := bson.D{{Key: "id", Value: id}}
	return &result, client.fetch(client.collections.agent, filter, &result)
}

func (client *Client) FetchActivity(id string) (*Activity, error) {
	/*
		collection := client.conn.Database(client.database).Collection(client.collections.activity)
		if collection == nil {
			return result, fmt.Errorf("No Collection")
		}

		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			return result, err
		}
	*/
	var result Activity
	filter := bson.D{{Key: "id", Value: id}}
	return &result, client.fetch(client.collections.activity, filter, &result)

}

func (client *Client) fetch(collection string, filter bson.D, result interface{}) error {
	c := client.conn.Database(client.database).Collection(collection)
	if c == nil {
		return fmt.Errorf("No Collection")
	}

	return c.FindOne(context.TODO(), filter).Decode(&result)
}
