package services

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/logistic-service/internal/models"
	"github.com/IgorRybak2055/logistic-service/internal/repository"
)

type deliveryService struct {
	deliveryRepo repository.Delivery
	log         *logrus.Logger
}

// NewProjectService will create new projectService object representation of Project interface
func NewDeliveryService(dr repository.Delivery, logger *logrus.Logger) Delivery {
	return &deliveryService{
		deliveryRepo: dr,
		log:         logger,
	}
}

func (d deliveryService) CreateDelivery(ctx context.Context, delivery models.Delivery) (models.Delivery, error) {
	if err := delivery.Validate(); err != nil {
		d.log.Debugf("createdelivery: %s", err)
		return models.Delivery{}, err
	}

	delivery, err := d.deliveryRepo.Create(ctx, delivery)
	if err != nil {
		d.log.Warnf("createdelivery: failed to %s", err)
		return models.Delivery{}, errors.Wrap(err, "creating new delivery")
	}

	return delivery, err
}
