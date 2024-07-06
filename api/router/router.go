package router

import (
	"mysite/features/health"
	"mysite/features/login"
	"mysite/features/refresh"
	"mysite/features/register"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var defaultTimeout = time.Second * 30

func InitRouter() chi.Router {
	r := chi.NewRouter()
	defaultMiddleWare(r)
	return buildRoute(r)
}

func defaultMiddleWare(r chi.Router) {
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(defaultTimeout))
}

func buildRoute(r chi.Router) chi.Router {
	health.HandlerFromMux(health.NewHandler(), r)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		publicApi(r)
	})
	return r
}

func publicApi(r chi.Router) {
	r.Group(func(r chi.Router) {
		register.HandlerFromMux(register.NewHandler(), r)
		login.HandlerFromMux(login.NewHandler(), r)
		refresh.HandlerFromMux(refresh.NewHandler(), r)
	})
}
