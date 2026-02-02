package logger

import (
	"log/slog"
	"os"
)

// New консструктор создания логгера
func New() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
