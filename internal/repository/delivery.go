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

// DeleteProject delete user project by ID.
func (d delivery) Create(ctx context.Context, delivery models.Delivery) (models.Delivery, error) {
	var query = ` INSERT INTO delivery
            (
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

	rows, err := sqlx.NamedQueryContext(ctx, d.dbc, query, delivery)
	if err != nil {
		return models.Delivery{}, errors.Wrap(err, "creating delivery")
	}

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
