package cmd

import (
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/app"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/config"
	loadCfg "github.com/1-infinity-1/banking-platform/pkg/config"
	"github.com/1-infinity-1/banking-platform/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	envFilePath = `local.env`
	prefix      = `GATEWAY`
)

var runApplicationCmd = &cobra.Command{
	Use: "application",
	RunE: func(cmd *cobra.Command, _ []string) error {
		loader := loadCfg.NewLoaderConfig(envFilePath, prefix)

		var cfg config.Config
		if err := loader.Load(&cfg); err != nil {
			return fmt.Errorf("loader.Load: %w", err)
		}

		log := logger.NewLogger(cfg.LogLevel)

		application, err := app.NewApp(log, cfg)
		if err != nil {
			return fmt.Errorf("app.NewApp: %w", err)
		}

		return application.Run(cmd.Context())
	},
}
