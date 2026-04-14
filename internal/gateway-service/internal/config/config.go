package config

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	HTTPConfig struct {
		Port string `envconfig:"PORT" default:"8081"`
	} `envconfig:"HTTP"`
}
