package tracer

import (
	"fmt"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rabbitmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

type Tracer struct {
	msgs     <-chan []byte
	consumer *rabbitmq.Consumer
	config   *config.Config
}

func New(config *config.Config) *Tracer {
	ch := make(chan []byte)
	tracer := Tracer{
		msgs: ch,
	}
	tracer.consumer = rabbitmq.NewConsumer(config.RabbitURL, ch, config.ConsumerTag)
	return &tracer
}

func (tracer *Tracer) Listen() {
	go func() {
		for {
			select {
			case msg := <-tracer.msgs:
				fmt.Println(msg)
			}
		}
	}()
}

func (tracer *Tracer) Cleanup() error {
	return tracer.consumer.Shutdown()
}
