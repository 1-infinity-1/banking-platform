package config

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	GRPCconfig struct {
		Port int `envconfig:"PORT" default:"8082"`
	} `envconfig:"GRPC"`
}
