package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"templify/pkg/app"
)

// Application metadata that is set at compile time.
// nolint
var (
	version     string
	buildDate   string
	description = "Templify"
)

// main just loads config and inits logger. Rest is done in app.Run.
func main() {
	appCfg, err := app.LoadConfig(
		version,
		buildDate,
		description,
	)
	if err != nil {
		fmt.Printf("could not load config: %s", err.Error())
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	err = app.Run(
		appCfg,
		quit,
	)

	if err != nil {
		slog.Error("error running app", "error", err)
	}
}
