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

	AccountGRPC struct {
		Host string `envconfig:"HOST" default:"localhost"`
		Port string `envconfig:"PORT" default:"8083"`
	} `envconfig:"ACCOUNT_GRPC"`

	TransactionGRPC struct {
		Host string `envconfig:"HOST" default:"localhost"`
		Port string `envconfig:"PORT" default:"8084"`
	} `envconfig:"TRANSACTION_GRPC"`

	LedgerGRPC struct {
		Host string `envconfig:"HOST" default:"localhost"`
		Port string `envconfig:"PORT" default:"8085"`
	} `envconfig:"LEDGER_GRPC"`

	JWT struct {
		Secret string `envconfig:"SECRET" required:"true"`
	} `envconfig:"JWT"`
}
