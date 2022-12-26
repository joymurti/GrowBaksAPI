package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/repository"
)

type (
	// IHealthCheckService is an interface that has all the function to be implemented inside health check service
	IHealthCheckService interface {
		HealthCheck() (bool, error)
	}

	// HealthCheckService is an app health check struct that consists of all the dependencies needed for health check service
	HealthCheckService struct {
		Context         context.Context
		Config          *config.Configuration
		Logger          *logrus.Logger
		HealthCheckRepo repository.IHealthCheckRepository
	}
)

// HealthCheck service layer to checking database is ok or not
func (hcs *HealthCheckService) HealthCheck() (bool, error) {
	ok, err := hcs.HealthCheckRepo.HealthCheck()
	if err != nil || !ok {
		return false, err
	}

	return ok, nil
}
