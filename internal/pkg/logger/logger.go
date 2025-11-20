package pkg

import (
	"log/slog"
	"os"
)

func SetUpLogger() *slog.Logger {
	var log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)
	return log
}
