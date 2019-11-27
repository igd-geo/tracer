package rbmq

import (
	"encoding/json"
	"fmt"
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
	"github.com/streadway/amqp"
)

const consumerTag = "tracer_consumer"

// Delivery is a custom wrapper for RabbitMQ deliveries
type Delivery struct {
	Entity *util.Entity `json:"entity,omitempty"`
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
		"notifications",	// name
		"topic",         	// type
		true,          	// durable
		false,      	// auto-deleted
		false,         	// internal
		false,         	// no-wait
		nil,      		// arguments
	)
	failOnError(err, "Failed to declare a exchange")

	queue, err := c.channel.QueueDeclare(
		"",    		// name
		false, 		// durable
		false, 	// delete when unused
		true,  	// exclusive
		false, 		// no-wait
		nil,   		// arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = c.channel.QueueBind(
		queue.Name,      			// queue name
		"#.tracer.#",    		// routing key
		"notifications",	// exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	deliveries, err := c.channel.Consume(
		queue.Name,  		// queue
		consumerTag, 		// consumer
		true,    	// auto ack
		false,	// exclusive
		false,     	// no local
		false,     	// no wait
		nil,       	// args
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
		delivery := Delivery{
			Entity: util.NewEntity(),
		}
		err := json.Unmarshal(d.Body, &delivery)
		if err != nil {
			log.Println(err)
			continue
		}
		ch <- delivery.Entity
	}
	done <- nil
}
