package main

import (
	"context"
	"log/slog"
	"mysite/pkgs/database"
	"mysite/pkgs/env"
	"mysite/pkgs/logger"
	"mysite/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	srv := http.Server{
		Addr:    ":3000",
		Handler: router.InitRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("failed to listen and serve", logger.AttrError(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", logger.AttrError(err))
		return
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		slog.Info("timeout of 5 seconds.")
	}
	slog.Info("Server exiting")
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
