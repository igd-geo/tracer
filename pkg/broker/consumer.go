package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/provenance"
	"github.com/streadway/amqp"
)

// ProvenanceMessage Wrapper for consumed provenance messages
type ProvenanceMessage struct {
	StartDate string   `json:"startDate,omitempty"`
	EndDate   string   `json:"endDate,omitempty"`
	Input     []string `json:"input"`
	Output    string   `json:"output"`
}

type consumer struct {
	channel     *amqp.Channel
	consumerTag string
	done        chan error
}

func newConsumer(conn *amqp.Connection, config *Config, handlerFunc func(*provenance.Entity)) (*consumer, error) {
	done := make(chan error)

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(
		config.ExchangeName, // name
		config.ExchangeType, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = channel.QueueBind(
		queue.Name,          // queue name
		config.BindingKey,   // routing key
		config.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	deliveries, err := channel.Consume(
		queue.Name,         // queue
		config.ConsumerTag, // consumer
		true,               // auto ack
		false,              // exclusive
		false,              // no local
		false,              // no wait
		nil,                // args
	)
	if err != nil {
		return nil, err
	}

	go handle(deliveries, handlerFunc, done)

	consumer := &consumer{
		channel:     channel,
		consumerTag: config.ConsumerTag,
		done:        done,
	}

	return consumer, nil
}

func (c *consumer) close() error {
	log.Println("shutting Down Consumer...")
	if err := c.channel.Cancel(c.consumerTag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, handlerFunc func(*provenance.Entity), done chan error) {
	for delivery := range deliveries {
		routingkey := delivery.RoutingKey
		sections := strings.Split(routingkey, ".")

		serviceID := sections[0]
		processID := sections[1]

		msg := ProvenanceMessage{}

		err := json.Unmarshal(delivery.Body, &msg)
		if err != nil {
			log.Println(err)
			continue
		}

		usedEntities := []*provenance.Entity{}
		for _, usedID := range msg.Input {
			usedEntities = append(usedEntities, &provenance.Entity{ID: usedID})
		}
		agent := &provenance.Agent{
			ID:              serviceID,
			ActedOnBehalfOf: []*provenance.Agent{&provenance.Agent{}},
		}

		activity := &provenance.Activity{
			ID:                processID,
			WasAssociatedWith: []*provenance.Agent{agent},
			Used:              usedEntities,
			StartDate:         msg.StartDate,
			EndDate:           msg.EndDate,
		}

		entity := &provenance.Entity{
			UID:            fmt.Sprintf("_:%s", msg.Output),
			ID:             msg.Output,
			CreationDate:   msg.EndDate,
			WasGeneratedBy: []*provenance.Activity{activity},
			WasDerivedFrom: usedEntities,
		}

		go handlerFunc(entity)
	}
	done <- nil
}
