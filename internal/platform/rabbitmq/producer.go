package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// TODO: Confirms

type Producer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	tag        string
	exchange   string
	routingKey string
	done       chan error
}

func NewProducer(url string, exchange string, routingKey string, ctag string, reliable bool) *Producer {
	p := &Producer{
		conn:       nil,
		channel:    nil,
		tag:        ctag,
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
	if reliable {
		log.Printf("enabling publishing confirms.")
		err := p.channel.Confirm(false)
		failOnError(err, "Channel could not be put into confirm mode")

		confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

	err = p.channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte("Hi"),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	)
	failOnError(err, "Failed to declare a queue")

	return p
}

func (p *Producer) Publish(body string) error {
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

func (p *Producer) Shutdown() error {
	log.Println("\nShutting Down...")
	if err := p.channel.Cancel(p.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
