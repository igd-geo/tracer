package rbmq

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

type Broker struct {
	conn *amqp.Connection
}

func NewBroker(url string) *Broker {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	return &Broker{
		conn: conn,
	}
}

// NewSession returns a new Session
func (b *Broker) NewSession(msgChan chan<- *util.Entity, exchange string, exchangeType string) *Session {
	return &Session{
		consumer: newConsumer(b.conn, msgChan),
		producer: newProducer(b.conn, exchange, exchangeType),
	}
}

// NewProducerOnly returns a new Session without a consumer
func (b *Broker) NewProducerOnlySession(exchange string, exchangeType string) *Session {
	return &Session{
		producer: newProducer(b.conn, exchange, exchangeType),
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

	return nil
}

func (b *Broker) Close() error {
	if err := b.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
