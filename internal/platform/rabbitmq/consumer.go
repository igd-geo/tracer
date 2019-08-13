package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"github.com/streadway/amqp"
)

type consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func newConsumer(conn *amqp.Connection, url string, msgChan chan<- *provutil.Entity, ctag string) *consumer {
	c := &consumer{
		conn:    conn,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}
	var err error

	go func() {
		err := <-c.conn.NotifyClose(make(chan *amqp.Error))
		if err != nil {
			log.Printf("closing: %s", err)
		}
	}()

	c.channel, err = c.conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = c.channel.ExchangeDeclare(
		"notifications", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	queue, err := c.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = c.channel.QueueBind(
		queue.Name,      // queue name
		"#.tracer.#",    // routing key
		"notifications", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	deliveries, err := c.channel.Consume(
		queue.Name, // queue
		c.tag,      // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	go handle(deliveries, msgChan, c.done)

	return c
}

func (c *consumer) shutdown() error {
	log.Println("Shutting Down Consumer...")
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, ch chan<- *provutil.Entity, done chan error) {
	for d := range deliveries {
		delivery := struct {
			Derivate *provutil.Entity `json:"entity,omitempty"`
		}{
			Derivate: provutil.NewEntity(),
		}
		err := json.Unmarshal(d.Body, &delivery)
		if err != nil {
			log.Println(err)
			return
		}
		ch <- delivery.Derivate
	}
	done <- nil
}
