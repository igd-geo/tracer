package rabbitmq

import (
	"log"
)

type Session struct {
	consumer *consumer
	producer *producer
}

func New(url string, msgChan chan<- Delivery, ctag string, exchange string, routingKey string) *Session {
	return &Session{
		consumer: newConsumer(url, msgChan, ctag),
		producer: newProducer(url, exchange, routingKey),
	}
}

func (s *Session) Shutdown() error {
	if err := s.consumer.shutdown(); err != nil {
		return err
	}

	if err := s.producer.shutdown(); err != nil {
		return err
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
