package config

import (
	flag "github.com/spf13/pflag"
)

const (
	defaultDB          = "localhost:9080"
	defaultBroker      = "amqp://guest:guest@localhost:5672/"
	defaultConsumerTag = "tracer_consumer"
)

// Config contains configuration information
type Config struct {
	DB          string
	Broker      string
	ConsumerTag string
	BatchSize   int
}

// New returns an empty Config struct
func New() *Config {
	return &Config{}
}

// InstallFlags fills a config with values passed by flags
func (config *Config) InstallFlags() {
	flag.StringVar(&config.DB, "db", defaultDB, "")
	flag.StringVar(&config.Broker, "broker", defaultBroker, "")
	flag.StringVar(&config.ConsumerTag, "ctag", defaultConsumerTag, "")

	flag.Parse()
}
