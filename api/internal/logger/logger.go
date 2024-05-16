package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	Log = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func New() *slog.Logger {
	return Log
}
