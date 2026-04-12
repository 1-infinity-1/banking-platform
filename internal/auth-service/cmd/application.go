package cmd

import (
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/app"
	"github.com/1-infinity-1/banking-platform/internal/auth-service/internal/config"
	loadCfg "github.com/1-infinity-1/banking-platform/pkg/config"
	"github.com/1-infinity-1/banking-platform/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	envFilePath = `local.env`
	prefix      = `AUTH`
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

		application := app.NewApp(log, cfg)

		if err := application.GRPCSrv.Run(cmd.Context()); err != nil {
			panic(err)
		}
	},
}
