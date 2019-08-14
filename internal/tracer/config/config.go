package config

type Config struct {
	InfoDB      string
	ProvDB      string
	RabbitURL   string
	ConsumerTag string
}

func New() *Config {
	return &Config{}
}
