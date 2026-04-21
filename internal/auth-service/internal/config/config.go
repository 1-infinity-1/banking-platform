package config

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	DB struct {
		Host     string `envconfig:"HOST" default:"localhost"`
		Port     string `envconfig:"PORT" default:"5432"`
		User     string `envconfig:"USER" default:"postgres"`
		Password string `envconfig:"PASSWORD" default:"postgres"`
		DBName   string `envconfig:"NAME" default:"app_db"`
	} `envconfig:"DB"`

	GRPCconfig struct {
		Port int `envconfig:"PORT" default:"8082"`
	} `envconfig:"GRPC"`
}
