package mongodb

import (
	"context"
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
