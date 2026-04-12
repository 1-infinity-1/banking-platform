package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type LoaderConfig struct {
	envFilePath string
	prefix      string
}

func NewLoaderConfig(envFilePath string, prefix string) *LoaderConfig {
	return &LoaderConfig{
		envFilePath: envFilePath,
		prefix:      prefix,
	}
}

func (l *LoaderConfig) Load(cfg interface{}) error {
	if l.envFilePath != "" {
		_ = godotenv.Load(l.envFilePath)
	}

	if err := envconfig.Process(l.prefix, cfg); err != nil {
		return fmt.Errorf("failed to process env vars: %w", err)
	}

	return nil
}
