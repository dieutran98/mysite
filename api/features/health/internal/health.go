package internal

import (
	"net/http"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) HealthCheck() int {
	return http.StatusOK
}
