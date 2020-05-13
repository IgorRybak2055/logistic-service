package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/IgorRybak2055/logistic-service/internal/models"
)

type delivery struct {
	dbc *sqlx.DB
}

// NewProjectRepository will create an object that represent the Project interface
func NewDeliveryRepository(dbc *sqlx.DB) Delivery {
	return &delivery{dbc}
}

// Create delete user project by ID.
func (d delivery) Create(ctx context.Context, delivery models.Delivery) (models.Delivery, error) {
	var query = ` INSERT INTO delivery
            (
						company_id,
                        shipment_date,
                        shipment_place,
                        unloading_place,
                        cargo,
                        weight_cargo,
                        volume_cargo,
                        trailer_type,
                        price
            )
            VALUES
            (
						:company_id,
                        :shipment_date,
                        :shipment_place,
                        :unloading_place,
                        :cargo,
                        :weight_cargo,
                        :volume_cargo,
                        :trailer_type,
                        :price
            )
            returning id`

	log.Println("create dlv: ", delivery)

	rows, err := sqlx.NamedQueryContext(ctx, d.dbc, query, delivery)
	if err != nil {
		return models.Delivery{}, errors.Wrap(err, "executing query")
	}

	log.Println("create dlv: ", delivery)

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println("closing rows:", err)
		}
	}()

	if rows.Next() {
		err = rows.StructScan(&delivery)
		if err != nil {
			return models.Delivery{}, errors.Wrap(err, "scanning result")
		}
	}

	return delivery, nil
}

// GetProjects returns all users projects.
func (d delivery) Deliveries(ctx context.Context) ([]models.Delivery, error) {
	var (
		query = `SELECT *
				 FROM   delivery LIMIT 20`
		dl  []models.Delivery
		err error
	)

	err = sqlx.SelectContext(ctx, d.dbc, &dl, query)
	if err != nil {
		return nil, errors.Wrap(err, "getting trucks")
	}

	return dl, nil
}

// GetProjects returns all users projects.
func (d delivery) InterestingDeliveries(ctx context.Context, companyID int64) ([]models.Delivery, error) {
	var (
		query = `SELECT *
				 FROM delivery WHERE delivery.company_id = ANY(SELECT * FROM subscribers WHERE company_id = $1)`
		dl  []models.Delivery
		err error
	)

	err = sqlx.SelectContext(ctx, d.dbc, &dl, query, companyID)
	if err != nil {
		return nil, errors.Wrap(err, "getting deliveries")
	}

	return dl, nil
}

// GetProjects returns all users projects.
func (d delivery) Delivery(ctx context.Context, id string) (models.Delivery, error) {
	var (
		query = `SELECT *
				 FROM delivery WHERE id = $1`
		dl  models.Delivery
		err error
	)

	err = sqlx.GetContext(ctx, d.dbc, &dl, query, id)
	if err != nil {
		return models.Delivery{}, errors.Wrap(err, "getting trucks")
	}

	return dl, nil
}
