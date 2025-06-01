package logger

import (
	"canary/internal/types"
	"log/slog"
	"os"
)

type Logger types.Logger

var log *Logger

func New(level, format string) *Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var handler slog.Handler
	if format == "text" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	logger := Logger{
		Logger: slog.New(handler),
	}
	return &logger
}

func Initialize(level, format string) {
	log = New(level, format)

}

func Get() *Logger {
	if log == nil {
		Initialize("info", "json")
	}
	return log
}

func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}
