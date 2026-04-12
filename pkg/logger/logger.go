package logger

import (
	"log/slog"
	"os"
)

const (
	localLogLevel = "local"
	devLogLevel   = "dev"
	prodLogLevel  = "prod"
)

func NewLogger(cfgLevel string) *slog.Logger {
	var log *slog.Logger

	switch cfgLevel {
	case localLogLevel:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case devLogLevel:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case prodLogLevel:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
