package config

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	DB struct {
		Host     string `envconfig:"HOST"     default:"localhost"`
		Port     string `envconfig:"PORT"     default:"5432"`
		User     string `envconfig:"USER"     default:"postgres"`
		Password string `envconfig:"PASSWORD" default:"postgres"`
		DBName   string `envconfig:"NAME"     default:"app_db"`
	} `envconfig:"DB"`

	GRPCconfig struct {
		Port int `envconfig:"PORT" default:"8084"`
	} `envconfig:"GRPC"`

	AccountService struct {
		Host string `envconfig:"HOST" default:"localhost"`
		Port int    `envconfig:"PORT" default:"8083"`
	} `envconfig:"ACCOUNT_SERVICE"`

	Kafka struct {
		Brokers string `envconfig:"BROKERS" default:"localhost:9092"`
		Topic   string `envconfig:"TOPIC"   default:"transactions.completed"`
	} `envconfig:"KAFKA"`
}
