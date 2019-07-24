package tracer

import (
	"fmt"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/dgraph"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/mongodb"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rabbitmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

type Tracer struct {
	msgs     <-chan rabbitmq.Delivery
	consumer *rabbitmq.Consumer
	mongoDB  *mongodb.Client
	dgraph   *dgraph.Client
	config   *config.Config
}

func New(config *config.Config) *Tracer {
	msgChan := make(chan rabbitmq.Delivery)
	tracer := Tracer{
		msgs:     msgChan,
		consumer: rabbitmq.NewConsumer(config.RabbitURL, msgChan, config.ConsumerTag),
		mongoDB: mongodb.NewClient(
			config.MongoURL,
			config.MongoDatabase,
			config.MongoCollectionEntity,
			config.MongoCollectionAgent,
			config.MongoCollectionActivity,
		),
		dgraph: dgraph.NewClient(config.DgraphURL),
	}
	return &tracer
}

func (tracer *Tracer) Listen() {
	go func() {
		for msg := range tracer.msgs {
			fmt.Println(msg)
		}
	}()
}

func (tracer *Tracer) Cleanup() error {
	return tracer.consumer.Shutdown()
}
