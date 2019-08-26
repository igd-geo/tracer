package broker

import (
	"log"

	"github.com/streadway/amqp"
)

// TODO: Confirms

type producer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

func newProducer(conn *amqp.Connection, url string, exchange string, routingKey string) *producer {
	var err error
	p := &producer{
		conn:       conn,
		channel:    nil,
		exchange:   exchange,
		routingKey: routingKey,
	}

	p.channel, err = p.conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = p.channel.ExchangeDeclare(
		exchange,   // name
		routingKey, // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	return p
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
			// a bunch of application/implementation-specific fields
		},
	)
	return err
}

func (p *producer) shutdown() error {
	log.Println("Shutting Down Producer...")
	p.channel.Close()
	return nil
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
