package config

import "time"

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

	AccessTokenTTL  time.Duration `envconfig:"AUTH_ACCESS_TOKEN_TTL" default:"5m"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL" default:"1h"`

	SecretKeyForToken string `envconfig:"SECRET_KEY_FOR_TOKEN" default:"secret"`
}
