package services

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/logistic-service/internal/models"
	"github.com/IgorRybak2055/logistic-service/internal/repository"
)

type truckService struct {
	truckRepo repository.Truck
	log       *logrus.Logger
}

// NewProjectService will create new projectService object representation of Project interface
func NewTruckService(tr repository.Truck, logger *logrus.Logger) Truck {
	return &truckService{
		truckRepo: tr,
		log:       logger,
	}
}

func (t truckService) Create(ctx context.Context, truck models.Truck) (models.Truck, error) {
	if err := truck.Validate(); err != nil {
		t.log.Debugf("createtruck: %s", err)
		return models.Truck{}, err
	}

	truck, err := t.truckRepo.New(ctx, truck)
	if err != nil {
		t.log.Warnf("createtruck: failed to %s", err)
		return models.Truck{}, errors.Wrap(err, "creating new delivery")
	}

	return truck, err
}

func (t truckService) Trucks(ctx context.Context, companyID int64) ([]models.Truck, error) {
	userProjects, err := t.truckRepo.Trucks(ctx, companyID)
	if err != nil {
		t.log.Warnf("gettrucks: failed to %s", err)
		return nil, errors.Wrap(err, "getting company trucks")
	}

	return userProjects, err
}