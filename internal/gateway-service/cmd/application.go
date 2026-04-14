package cmd

import (
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
	Run: func(cmd *cobra.Command, _ []string) {
		loadCfg := loadCfg.NewLoaderConfig(envFilePath, prefix)

		var cfg config.Config
		if err := loadCfg.Load(&cfg); err != nil {
			panic(err)
		}

		log := logger.NewLogger(cfg.LogLevel)

		application, err := app.NewApp(log, cfg)
		if err != nil {
			panic(err)
		}

		if err := application.HTTPSrv.Run(cmd.Context()); err != nil {
			panic(err)
		}
	},
}
