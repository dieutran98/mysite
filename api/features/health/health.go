package health

import (
	"fmt"
	"mysite/features/health/internal"
	"net/http"

	"github.com/go-chi/render"
)

type api struct {
	svc service
}

type service interface {
	HealthCheck() int
}

func NewHandler() *api {
	return &api{
		svc: internal.NewService(),
	}
}

func (a *api) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test")
	render.PlainText(w, r, http.StatusText(a.svc.HealthCheck()))
}
