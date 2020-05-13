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

func (d deliveryService) Deliveries(ctx context.Context) ([]models.Delivery, error) {
	dlvs, err := d.deliveryRepo.Deliveries(ctx)
	if err != nil {
		d.log.Warnf("gettrucks: failed to %s", err)
		return nil, errors.Wrap(err, "getting company trucks")
	}

	return dlvs, err
}

func (d deliveryService) InterestingDeliveries(ctx context.Context, companyID int64) ([]models.Delivery, error) {
	dlvs, err := d.deliveryRepo.InterestingDeliveries(ctx, companyID)
	if err != nil {
		d.log.Warnf("gettrucks: failed to %s", err)
		return nil, errors.Wrap(err, "getting company trucks")
	}

	return dlvs, err
}

func (d deliveryService) Delivery(ctx context.Context, id string) (models.Delivery, error) {
	dlvs, err := d.deliveryRepo.Delivery(ctx, id)
	if err != nil {
		d.log.Warnf("gettrucks: failed to %s", err)
		return models.Delivery{}, errors.Wrap(err, "getting company trucks")
	}

	return dlvs, err
}