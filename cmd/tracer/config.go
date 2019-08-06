package main

import (
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
	flag "github.com/spf13/pflag"
)

const (
	// Dgraph
	defaultDgraphURL = "localhost:9080"

	// MongoDB
	defaultMongoURL                = "mongodb://root:example@localhost:27017"
	defaultMongoDatabase           = "tracer"
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
