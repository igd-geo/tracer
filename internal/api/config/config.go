package config

import (
	flag "github.com/spf13/pflag"
)

const (
	defaultDB          = "localhost:9080"
	defaultBroker      = "amqp://guest:guest@localhost:5672/"
	defaultConsumerTag = "tracer_consumer"
	defaultPort        = ":1234"
)

// Config contains configuration information
type Config struct {
	DB          string
	Broker      string
	ConsumerTag string
	Port        string
}

// New returns an empty Config struct
func New() *Config {
	return &Config{}
}

// InstallFlags fills a config with values passed by flags
func (config *Config) InstallFlags() {
	flag.StringVar(&config.DB, "db", defaultDB, "database grpc url")
	flag.StringVar(&config.Broker, "broker", defaultBroker, "rabbitmq url")
	flag.StringVar(&config.ConsumerTag, "ctag", defaultConsumerTag, "consumer tag")
	flag.StringVar(&config.Port, "port", defaultPort, "tracer api port")

	flag.Parse()
}
