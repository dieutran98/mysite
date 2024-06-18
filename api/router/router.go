package router

import (
	"mysite/features/health"
	"mysite/features/register"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter() chi.Router {
	r := chi.NewRouter()
	defaultMiddleWare(r)
	return buildRoute(r)
}

func defaultMiddleWare(r chi.Router) {
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
}

func buildRoute(r chi.Router) chi.Router {
	health.HandlerFromMux(health.NewHandler(), r)

	r.Route("/api/v1", func(r chi.Router) {
		publicApi(r)
	})
	return r
}

func publicApi(r chi.Router) {
	r.Group(func(r chi.Router) {
		register.HandlerFromMux(register.NewHandler(), r)
	})
}
