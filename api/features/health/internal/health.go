package internal

import (
	"mysite/model"
	"mysite/utils/ptrconv"
	"net/http"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) HealthCheck() model.HealthResponse {
	return model.HealthResponse{
		Message: ptrconv.String(http.StatusText(http.StatusOK)),
	}
}
