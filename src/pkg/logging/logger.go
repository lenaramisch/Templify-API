package logging

import (
	"log/slog"
	"os"
	"strings"
)

func SetLogger() *slog.Logger {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slogLevel := slog.LevelDebug
	logLevel, ok := os.LookupEnv("LOGGER_LEVEL")
	if ok && logLevel != "" {
		logLevel = strings.ToLower(logLevel)
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
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel,
		}))
	} else if ok && loggerVar == "prettyjson" {
		logger = slog.New(NewPrettyHandler(&slog.HandlerOptions{
			Level:     slogLevel,
			AddSource: true,
		}))
	} else {
		slog.Warn("LOGGER environment variable not set to json, using default console logger")
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel,
		}))
	}

	logger.With("Loglevel", logLevel, "Logformat", loggerVar).Info("Logger set")
	logger.Debug("DEBUG MODE IS ACTIVE")
	return logger
}
