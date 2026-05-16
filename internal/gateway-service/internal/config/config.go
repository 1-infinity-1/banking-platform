package config

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	HTTPConfig struct {
		Port string `envconfig:"PORT" default:"8081"`
	} `envconfig:"HTTP"`

	AuthGRPC struct {
		Host string `envconfig:"HOST" default:"localhost"`
		Port string `envconfig:"PORT" default:"8082"`
	} `envconfig:"AUTH_GRPC"`
}
