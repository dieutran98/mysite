package router

import (
	"mysite/features/health"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()
	defaultMiddleWare(r)
	return buildRoute(r)
}

func defaultMiddleWare(r *chi.Mux) {
	r.Use(middleware.Recoverer)
}

func buildRoute(r *chi.Mux) *chi.Mux {
	health.HandlerFromMux(health.NewHandler(), r)
	return r
}
