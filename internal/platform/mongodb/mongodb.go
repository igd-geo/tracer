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
	collections collections
}

type collections struct {
	entity   string
	agent    string
	activity string
}

type Entity struct {
	UID          string          `bson:"uid,omitempty"`
	ID           string          `bson:"id,omitempty"`
	URI          string          `bson:"uri,omitempty"`
	Name         string          `bson:"name,omitempty"`
	CreationDate string          `bson:"creationDate,omitempty"`
	Data         json.RawMessage `bson:"data,omitempty"`
}
type Agent struct {
	UID  string          `bson:"uid,omitempty"`
	ID   string          `bson:"id,omitempty"`
	Name string          `bson:"name,omitempty"`
	Data json.RawMessage `bson:"data,omitempty"`
}
type Activity struct {
	UID       string          `bson:"uid,omitempty"`
	ID        string          `bson:"id,omitempty"`
	StartDate string          `bson:"startDate,omitempty"`
	EndDate   string          `bson:"endDate,omitempty"`
	Data      json.RawMessage `bson:"data,omitempty"`
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
		collections: collections{
			entity:   collEntity,
			agent:    collAgent,
			activity: collActivity,
		},
	}
}

func NewEntity() *Entity {
	return &Entity{}
}
func NewAgent() *Agent {
	return &Agent{}
}
func NewActivity() *Activity {
	return &Activity{}
}

func (client *Client) EntityUID(id string) string {
	return client.FetchEntity(id).UID
}

func (client *Client) AgentUID(id string) string {
	return client.FetchAgent(id).UID
}

func (client *Client) ActivitytUID(id string) string {
	return client.FetchActivity(id).UID
}

func (client *Client) InsertEntity(entity *Entity) error {
	collection := client.conn.Database(client.database).Collection(client.collections.entity)
	if collection == nil {
		return fmt.Errorf("No Collection")
	}

	payload, _ := bson.Marshal(entity)

	_, err := collection.InsertOne(context.TODO(), payload)
	return err
}

func (client *Client) InsertAgent(agent *Agent) error {
	collection := client.conn.Database(client.database).Collection(client.collections.agent)
	if collection == nil {
		return fmt.Errorf("No Collection")
	}

	payload, _ := bson.Marshal(agent)

	_, err := collection.InsertOne(context.TODO(), payload)
	return err
}

func (client *Client) InsertActivity(activity *Activity) error {
	collection := client.conn.Database(client.database).Collection(client.collections.activity)
	if collection == nil {
		return fmt.Errorf("No Collection")
	}

	payload, _ := bson.Marshal(activity)

	_, err := collection.InsertOne(context.TODO(), payload)
	return err
}

func (client *Client) FetchEntity(id string) *Entity {
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
	client.fetch(client.collections.entity, filter, &result)
	return &result
}

func (client *Client) FetchAgent(id string) *Agent {
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
	client.fetch(client.collections.agent, filter, &result)
	return &result
}

func (client *Client) FetchActivity(id string) *Activity {
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
	client.fetch(client.collections.activity, filter, &result)
	return &result

}

func (client *Client) fetch(collection string, filter bson.D, result interface{}) error {
	c := client.conn.Database(client.database).Collection(collection)
	if c == nil {
		return fmt.Errorf("No Collection")
	}

	return c.FindOne(context.TODO(), filter).Decode(&result)
}
