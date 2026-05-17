package cmd

import (
	"log/slog"

	consumerapp "github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/app/consumer"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/config"
	loadCfg "github.com/1-infinity-1/banking-platform/pkg/config"
	"github.com/1-infinity-1/banking-platform/pkg/logger"
	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Start Kafka consumer",
	Run: func(cmd *cobra.Command, _ []string) {
		loader := loadCfg.NewLoaderConfig(envFilePath, prefix)

		var cfg config.Config
		if err := loader.Load(&cfg); err != nil {
			panic(err)
		}

		log := logger.NewLogger(cfg.LogLevel)

		log.Info("kafka consumer starting",
			slog.Any("brokers", cfg.Kafka.Brokers),
			slog.String("topic", cfg.Kafka.Topic),
		)

		application, err := consumerapp.NewApp(cmd.Context(), log, cfg)
		if err != nil {
			panic(err)
		}

		if err = application.Run(cmd.Context()); err != nil {
			panic(err)
		}
	},
}
