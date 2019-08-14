package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName       = "tracer"
	collectionEntity   = "entity"
	collectionAgent    = "agent"
	collectionActivity = "activity"
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

func NewClient(url string) *Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		conn:     client,
		database: databaseName,
		collections: collections{
			entity:   collectionEntity,
			agent:    collectionAgent,
			activity: collectionActivity,
		},
	}
}

func (client *Client) EntityUID(id string) string {
	entity := client.FetchEntity(id)
	if entity == nil {
		return ""
	}
	return entity.UID
}

func (client *Client) AgentUID(id string) string {
	agent := client.FetchAgent(id)
	if agent == nil {
		return ""
	}
	return agent.UID
}

func (client *Client) ActivitytUID(id string) string {
	activity := client.FetchActivity(id)
	if activity == nil {
		return ""
	}
	return activity.UID
}

func (client *Client) InsertEntity(entity *provutil.Entity) error {
	collection := client.conn.Database(client.database).Collection(client.collections.entity)
	if collection == nil {
		return fmt.Errorf("No Collection")
	}

	payload, _ := bson.Marshal(entity)
	_, err := collection.InsertOne(context.TODO(), payload)
	return err
}

func (client *Client) InsertAgent(agent *provutil.Agent) error {
	collection := client.conn.Database(client.database).Collection(client.collections.agent)
	if collection == nil {
		return fmt.Errorf("No Collection")
	}

	payload, _ := bson.Marshal(agent)
	_, err := collection.InsertOne(context.TODO(), payload)
	return err
}

func (client *Client) InsertActivity(activity *provutil.Activity) error {
	collection := client.conn.Database(client.database).Collection(client.collections.activity)
	if collection == nil {
		return fmt.Errorf("No Collection")
	}

	payload, _ := bson.Marshal(activity)

	_, err := collection.InsertOne(context.TODO(), payload)
	return err
}

func (client *Client) FetchEntity(id string) *provutil.Entity {
	var entity provutil.Entity
	filter := bson.D{
		{Key: "id", Value: id},
	}
	result := client.fetch(client.collections.entity, filter)
	err := result.Decode(&entity)
	if err != nil {
		return nil
	}

	return &entity
}

func (client *Client) FetchAgent(id string) *provutil.Agent {
	var agent provutil.Agent
	filter := bson.D{
		{Key: "id", Value: id},
	}

	result := client.fetch(client.collections.agent, filter)
	err := result.Decode(&agent)
	if err != nil {
		return nil
	}

	return &agent
}

func (client *Client) FetchActivity(id string) *provutil.Activity {
	var activity provutil.Activity
	filter := bson.D{
		{Key: "id", Value: id},
	}
	result := client.fetch(client.collections.activity, filter)
	err := result.Decode(&activity)
	if err != nil {
		return nil
	}

	return &activity
}

func (client *Client) fetch(collection string, filter bson.D) *mongo.SingleResult {
	c := client.conn.Database(client.database).Collection(collection)
	if c == nil {
		return nil
	}
	return c.FindOne(context.TODO(), filter)
}
