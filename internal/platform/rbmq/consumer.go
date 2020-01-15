package rbmq

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
	"github.com/streadway/amqp"
)

const (
	consumerTag = "tracer_consumer"
	scenarioID  = "5ce8a0c4-9947-4ad8-bf9e-f95aa7134c37"
	bindingKey  = "#.#.provenance"
)

// Delivery is a custom wrapper for RabbitMQ deliveries
type Delivery struct {
	Entity *util.Entity `json:"entity,omitempty"`
}

type ProvMessage struct {
	StartDate string   `json:"startDate,omitempty"`
	EndDate   string   `json:"endDate,omitempty"`
	Input     []string `json:"input"`
	Output    string   `json:"output"`
}

type consumer struct {
	channel *amqp.Channel
	done    chan error
}

func newConsumer(conn *amqp.Connection, msgChan chan<- *util.Entity) *consumer {
	c := &consumer{
		channel: nil,
		done:    make(chan error),
	}
	var err error

	go func() {
		err := <-conn.NotifyClose(make(chan *amqp.Error))
		if err != nil {
			log.Printf("closing: %s", err)
		}
	}()

	c.channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = c.channel.ExchangeDeclare(
		scenarioID, // name
		"topic",    // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	queue, err := c.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = c.channel.QueueBind(
		queue.Name, // queue name
		bindingKey, // routing key
		scenarioID, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	deliveries, err := c.channel.Consume(
		queue.Name,  // queue
		consumerTag, // consumer
		true,        // auto ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	go handle(deliveries, msgChan, c.done)

	return c
}

func (c *consumer) shutdown() error {
	log.Println("Shutting Down Consumer...")
	if err := c.channel.Cancel(consumerTag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, ch chan<- *util.Entity, done chan error) {
	for d := range deliveries {
		routingkey := d.RoutingKey
		rk := strings.Split(routingkey, ".")

		serviceID := rk[0]
		processID := rk[1]

		provmsg := ProvMessage{}

		err := json.Unmarshal(d.Body, &provmsg)
		if err != nil {
			log.Println(err)
			continue
		}

		used := []*util.Entity{}
		for _, usedID := range provmsg.Input {
			used = append(used, &util.Entity{ID: usedID})
		}
		agent := &util.Agent{
			ID:              serviceID,
			ActedOnBehalfOf: []*util.Agent{&util.Agent{}},
		}

		activity := &util.Activity{
			ID:                processID,
			WasAssociatedWith: []*util.Agent{agent},
			Used:              used,
			StartDate:         provmsg.StartDate,
			EndDate:           provmsg.EndDate,
		}

		entity := &util.Entity{
			UID:            fmt.Sprintf("_:%s", provmsg.Output),
			ID:             provmsg.Output,
			CreationDate:   provmsg.EndDate,
			WasGeneratedBy: []*util.Activity{activity},
			WasDerivedFrom: used,
		}

		ch <- entity
	}
	done <- nil
}
