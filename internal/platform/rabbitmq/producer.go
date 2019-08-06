package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// TODO: Confirms

type producer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	routingKey string
	done       chan error
}

func newProducer(url string, exchange string, routingKey string) *producer {
	p := &producer{
		conn:       nil,
		channel:    nil,
		exchange:   exchange,
		routingKey: routingKey,
		done:       make(chan error),
	}
	var err error

	p.conn, err = amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	go func() {
		err := <-p.conn.NotifyClose(make(chan *amqp.Error))
		if err != nil {
			log.Printf("closing: %s", err)
		}
	}()

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

	err = p.publish("Hi")
	failOnError(err, "Failed to initialize producer")

	return p
}

func (p *producer) publish(body string) error {
	err := p.channel.Publish(
		p.exchange,   // publish to an exchange
		p.routingKey, // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
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
	log.Println("\nShutting Down...")
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	return <-p.done
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
