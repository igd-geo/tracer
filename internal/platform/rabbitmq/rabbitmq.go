package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Session struct {
	consumer *consumer
	producer *producer
	conn     *amqp.Connection
}

func New(url string, msgChan chan<- Delivery, ctag string, exchange string, routingKey string) *Session {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	return &Session{
		consumer: newConsumer(conn, url, msgChan, ctag),
		producer: newProducer(conn, url, exchange, routingKey),
		conn:     conn,
	}
}

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
