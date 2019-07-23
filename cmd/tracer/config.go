package main

import (
	"flag"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

const (
	// Dgraph
	defaultDgraphURL = "http://localhost:9080"

	// MongoDB
	defaultMongoURL                = "mongodb://localhost:27017"
	defaultMongoDatabase           = "provenance"
	defaultMongoCollectionEntity   = "entity"
	defaultMongoCollectionAgent    = "agent"
	defaultMongoCollectionActivity = "activity"

	// RabbitMQ
	defaultRabbitMQ    = "amqp://guest:guest@localhost:5672/"
	defaultConsumerTag = "tracer_consumer"
)

func installFlags(config *config.Config) {
	flag.StringVar(&config.DgraphURL, "dgraph", defaultDgraphURL, "")

	flag.StringVar(&config.MongoURL, "mongoURL", defaultMongoURL, "")
	flag.StringVar(&config.MongoDatabase, "mongoDatabase", defaultMongoDatabase, "")
	flag.StringVar(&config.MongoCollectionEntity, "mongoCollectionEntity", defaultMongoCollectionEntity, "")
	flag.StringVar(&config.MongoCollectionAgent, "mongoCollectionAgent", defaultMongoCollectionAgent, "")
	flag.StringVar(&config.MongoCollectionActivity, "mongoCollectionActivity", defaultMongoCollectionActivity, "")

	flag.StringVar(&config.RabbitURL, "rabbit", defaultRabbitMQ, "")
	flag.StringVar(&config.ConsumerTag, "ctag", defaultConsumerTag, "")

	flag.Parse()
}
