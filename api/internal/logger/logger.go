package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func New() *slog.Logger {
	Log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	return Log
}
