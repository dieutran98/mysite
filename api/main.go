package main

import (
	"log/slog"
	"mysite/router"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":3000", router.InitRouter()); err != nil {
		slog.Error("init router error", slog.String("err", err.Error()))
	}
}
