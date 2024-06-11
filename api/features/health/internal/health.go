package internal

import (
	"net/http"

	"mysite/models/model"
	"mysite/utils/ptrconv"
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
