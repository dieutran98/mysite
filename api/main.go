package main

import (
	"log/slog"
	"mysite/pkgs/database"
	"mysite/pkgs/env"
	"mysite/pkgs/logger"
	"mysite/router"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func main() {
	if err := setup(); err != nil {
		slog.Error("failed to setup application", logger.AttrError(err))
		return
	}

	defer func() { // all defer functions will running here
		if err := database.Close(); err != nil {
			slog.Error("failed to close database", logger.AttrError(err))
		}
	}()

	if err := http.ListenAndServe(":3000", router.InitRouter()); err != nil {
		slog.Error("init router error", logger.AttrError(err))
		return
	}
}

func setup() error {
	logger.SetLogger(os.Stdout)
	logger.SetLogLevel(slog.LevelDebug)

	if err := env.ReadEnv(); err != nil {
		return errors.Wrap(err, "failed to readEnv")
	}

	if err := database.SetupDatabase(); err != nil {
		return errors.Wrap(err, "failed to setup database")
	}

	return nil
}
