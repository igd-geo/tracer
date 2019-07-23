package tracer

import (
	"fmt"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/dgraph"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/mongodb"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rabbitmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

type Tracer struct {
	msgs     <-chan []byte
	consumer *rabbitmq.Consumer
	mongoDB  *mongodb.Client
	dgraph   *dgraph.Client
	config   *config.Config
}

func New(config *config.Config) *Tracer {
	ch := make(chan []byte)
	tracer := Tracer{
		msgs:     make(chan []byte),
		consumer: rabbitmq.NewConsumer(config.RabbitURL, ch, config.ConsumerTag),
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
		for {
			select {
			case msg := <-tracer.msgs:
				fmt.Println(msg)
			}
		}
	}()
}

func (tracer *Tracer) Cleanup() error {
	return tracer.consumer.Shutdown()
}
