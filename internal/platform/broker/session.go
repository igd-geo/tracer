package broker

import (
	"fmt"
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
	"github.com/streadway/amqp"
)

// Session is a wrapper for a RabbitMQ consumer and producer
type Session struct {
	consumer *consumer
	producer *producer
	conn     *amqp.Connection
}

// New returns a new Session
func New(url string, msgChan chan<- *util.Entity, ctag string, exchange string, routingKey string) *Session {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	return &Session{
		consumer: newConsumer(conn, url, msgChan, ctag),
		producer: newProducer(conn, url, exchange, routingKey),
		conn:     conn,
	}
}

// Publish sends a new message to the defined routing key
func (s *Session) Publish(msg string, routingKey string) error {
	return s.producer.publish(msg, routingKey)
}

// Shutdown gracefully closes the Session
func (s *Session) Shutdown() error {
	if err := s.consumer.shutdown(); err != nil {
		return err
	}

	if err := s.producer.shutdown(); err != nil {
		return err
	}

	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
