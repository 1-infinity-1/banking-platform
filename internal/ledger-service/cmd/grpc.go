package cmd

import (
	"log/slog"

	grpcapp "github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/app/grpc"
	"github.com/1-infinity-1/banking-platform/internal/ledger-service/internal/config"
	loadCfg "github.com/1-infinity-1/banking-platform/pkg/config"
	"github.com/1-infinity-1/banking-platform/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	envFilePath = `local.env`
	prefix      = `LEDGER`
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start gRPC server",
	Run: func(cmd *cobra.Command, _ []string) {
		loader := loadCfg.NewLoaderConfig(envFilePath, prefix)

		var cfg config.Config
		if err := loader.Load(&cfg); err != nil {
			panic(err)
		}

		log := logger.NewLogger(cfg.LogLevel)

		log.Info("grpc server starting", slog.Int("port", cfg.GRPCConfig.Port))

		application, err := grpcapp.NewApp(cmd.Context(), log, cfg)
		if err != nil {
			panic(err)
		}

		if err = application.Run(cmd.Context()); err != nil {
			panic(err)
		}
	},
}
