package broker

import (
	"log"

	"github.com/streadway/amqp"
)

type producer struct {
	channel  *amqp.Channel
	exchange string
}

func newProducer(conn *amqp.Connection, config *Config) (*producer, error) {
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

	producer := &producer{
		channel:  channel,
		exchange: config.ExchangeName,
	}

	return producer, nil
}

func (p *producer) publish(body string, routingKey string) error {
	err := p.channel.Publish(
		p.exchange, // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	)
	return err
}

func (p *producer) close() error {
	log.Println("shutting Down Producer...")
	return p.channel.Close()
}
