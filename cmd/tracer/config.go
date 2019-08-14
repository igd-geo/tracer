package main

import (
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
	flag "github.com/spf13/pflag"
)

const (
	defaultProvDBURL   = "localhost:9080"
	defaultInfoDBURL   = "mongodb://root:example@localhost:27017"
	defaultRabbitMQ    = "amqp://guest:guest@localhost:5672/"
	defaultConsumerTag = "tracer_consumer"
)

func installFlags(config *config.Config) {
	flag.StringVar(&config.ProvDB, "provdb", defaultProvDBURL, "")
	flag.StringVar(&config.InfoDB, "infodb", defaultInfoDBURL, "")
	flag.StringVar(&config.RabbitURL, "rabbit", defaultRabbitMQ, "")
	flag.StringVar(&config.ConsumerTag, "ctag", defaultConsumerTag, "")

	flag.Parse()
}
