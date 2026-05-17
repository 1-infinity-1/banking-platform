package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "ledger-service",
}

func Execute() {
	rootCmd.AddCommand(grpcCmd)
	rootCmd.AddCommand(consumerCmd)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
	cancel()
}
