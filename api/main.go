package main

import (
	"log/slog"
	"mysite/pkgs/env"
	"mysite/pkgs/logger"
	"mysite/router"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func main() {
	logger.SetLogger(os.Stdout)

	if err := env.ReadEnv(viper.New()); err != nil {
		slog.Error("read env error", slog.String("err", err.Error()))
		return
	}

	if err := http.ListenAndServe(":3000", router.InitRouter()); err != nil {
		slog.Error("init router error", slog.String("err", err.Error()))
		return
	}
}
