package health

import (
	"mysite/features/health/internal"
	"mysite/model"
	"net/http"

	"github.com/go-chi/render"
)

type api struct {
	svc service
}

type service interface {
	HealthCheck() model.HealthResponse
}

func NewHandler() *api {
	return &api{
		svc: internal.NewService(),
	}
}

func (a *api) Health(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, a.svc.HealthCheck())
}
