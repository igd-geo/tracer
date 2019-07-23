package config

type Config struct {
	DgraphURL               string
	MongoURL                string
	MongoDatabase           string
	MongoCollectionEntity   string
	MongoCollectionAgent    string
	MongoCollectionActivity string
	RabbitURL               string
	ConsumerTag             string
}

func New() *Config {
	return &Config{}
}
