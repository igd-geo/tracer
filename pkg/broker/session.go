package broker

import (
	"fmt"
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/provenance"
	"github.com/streadway/amqp"
)

// Config Configuration for the RabbitMQ client used as message broker
type Config struct {
	URL          string `yaml:"url"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	ExchangeName string `yaml:"exchangeName"`
	ExchangeType string `yaml:"exchangeType"`
	ConsumerTag  string `yaml:"consumerTag"`
	BindingKey   string `yaml:"bindingkey"`
}

// Session Wrapper for  RabbitMQ consumer and producer connections
type Session struct {
	conn     *amqp.Connection
	consumer *consumer
	producer *producer
	config   *Config
}

// New Connects to the RabbitMQ server and returns Session struct
func New(config *Config) (*Session, error) {
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s", config.User, config.Password, config.URL)
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	session := &Session{
		conn: conn,
	}
	return session, nil
}

// NewConsumer initializes a new consumer to receive provenance messages from an exchange
func (s *Session) NewConsumer(handlerFunc func(*provenance.Entity)) error {
	consumer, err := newConsumer(s.conn, s.config, handlerFunc)
	if err != nil {
		return err
	}
	s.consumer = consumer
	return nil
}

// NewProducer initializes a new producer to publish messages to an exchange
func (s *Session) NewProducer() error {
	producer, err := newProducer(s.conn, s.config)
	if err != nil {
		return err
	}
	s.producer = producer
	return nil
}

// Publish sends a new message to the defined routing key
func (s *Session) Publish(msg string, routingKey string) error {
	return s.producer.publish(msg, routingKey)
}

// Close Gracefully closes the Session
func (s *Session) Close() error {
	if err := s.consumer.close(); err != nil {
		return err
	}

	if err := s.producer.close(); err != nil {
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
