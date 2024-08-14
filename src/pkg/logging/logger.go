package logging

import (
	"log/slog"
	"os"
	"strings"
)

func SetLogger() *slog.Logger {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	var slogLevel slog.Level
	logLevel, ok := os.LookupEnv("LOGGER_LEVEL")
	if ok && logLevel != "" {
		switch logLevel {
		case "debug":
			slogLevel = slog.LevelDebug
		case "info":
			slogLevel = slog.LevelInfo
		case "warn":
			slogLevel = slog.LevelWarn
		case "error":
			slogLevel = slog.LevelError
		default:
			slog.Warn("Invalid LOGGER_LEVEL environment variable, using default logger level DEBUG")
		}
	}

	var logger *slog.Logger
	loggerVar, ok := os.LookupEnv("LOGGER")

	loggerVar = strings.ToLower(loggerVar)
	// set slog to json or console
	if ok && loggerVar == "json" {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     slogLevel,
		})
		logger = slog.New(handler)
	} else {
		slog.Warn("LOGGER environment variable not set to json, using default console logger")
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     slogLevel,
		})
		logger = slog.New(handler)
	}
	return logger
}
