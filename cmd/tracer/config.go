package main

import (
	"flag"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

const (
	defaultDgraph      = ""
	defaultMongoDB     = ""
	defaultRabbitMQ    = "amqp://guest:guest@localhost:5672/"
	defaultConsumerTag = "tracer_consumer"
)

func installFlags(config *config.Config) {
	flag.StringVar(&config.DgraphURL, "dgraph", defaultDgraph, "")
	flag.StringVar(&config.MongoURL, "mongo", defaultMongoDB, "")
	flag.StringVar(&config.RabbitURL, "rabbit", defaultRabbitMQ, "")
	flag.StringVar(&config.ConsumerTag, "ctag", defaultConsumerTag, "")

	flag.Parse()
}
