package config

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	DB struct {
		Host     string `envconfig:"HOST"     default:"localhost"`
		Port     string `envconfig:"PORT"     default:"5432"`
		User     string `envconfig:"USER"     default:"postgres"`
		Password string `envconfig:"PASSWORD" default:"postgres"`
		DBName   string `envconfig:"NAME"     default:"ledger_db"`
	} `envconfig:"DB"`

	GRPCConfig struct {
		Port int `envconfig:"PORT" default:"8083"`
	} `envconfig:"GRPC"`

	Kafka struct {
		Brokers []string `envconfig:"BROKERS" default:"localhost:9092"`
		Topic   string   `envconfig:"TOPIC"   default:"transactions.completed"`
		GroupID string   `envconfig:"GROUP"   default:"ledger-service"`
	} `envconfig:"KAFKA"`
}
