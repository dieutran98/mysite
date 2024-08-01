package internal

import (
	"net/http"

	"mysite/dtos"
	"mysite/utils/ptrconv"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) HealthCheck() dtos.HealthResponse {
	return dtos.HealthResponse{
		Message: ptrconv.String(http.StatusText(http.StatusOK)),
	}
}
