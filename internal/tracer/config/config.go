package config

type Config struct {
	DgraphURL   string
	MongoURL    string
	RabbitURL   string
	ConsumerTag string
}

func New() *Config {
	return &Config{}
}
