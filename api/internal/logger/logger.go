package logger

import (
	"context"
	"log/slog"
	"os"
)

type ContextValues struct {
	RequestID string
}

var Log *slog.Logger
var contextValues ContextValues

func init() {
	Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	contextValues = ContextValues{
		RequestID: "request-id",
	}
}

func New() *slog.Logger {
	return Log
}

func putContextValuesToArgs(args []interface{}, ctx context.Context) []interface{} {
	requestId := ctx.Value(contextValues.RequestID)
	args = append(args, contextValues.RequestID, requestId)
	return args
}

func Info(ctx context.Context, message string, args ...interface{}) {
	args = putContextValuesToArgs(args, ctx)
	Log.InfoContext(ctx, message, args...)
}

func Error(ctx context.Context, message string, args ...interface{}) {
	args = putContextValuesToArgs(args, ctx)
	Log.ErrorContext(ctx, message, args...)
}
