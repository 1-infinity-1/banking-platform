package cmd

import (
	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/app"
	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/config"
	loadCfg "github.com/1-infinity-1/banking-platform/pkg/config"
	"github.com/1-infinity-1/banking-platform/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	envFilePath = `local.env`
	prefix      = `ACCOUNT`
)

var runApplicationCmd = &cobra.Command{
	Use: "application",
	Run: func(cmd *cobra.Command, _ []string) {
		loader := loadCfg.NewLoaderConfig(envFilePath, prefix)

		var cfg config.Config
		if err := loader.Load(&cfg); err != nil {
			panic(err)
		}

		log := logger.NewLogger(cfg.LogLevel)

		application, err := app.NewApp(cmd.Context(), log, cfg)
		if err != nil {
			panic(err)
		}

		err = application.GRPCSrv.Run(cmd.Context())
		if err != nil {
			panic(err)
		}
	},
}
