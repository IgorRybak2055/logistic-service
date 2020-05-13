package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/IgorRybak2055/logistic-service/internal/models"
)

type truck struct {
	dbc *sqlx.DB
}

// NewTruckRepository will create an object that represent the Truck interface
func NewTruckRepository(dbc *sqlx.DB) Truck {
	return &truck{dbc}
}

// Create delete user project by ID.
func (t truck) New(ctx context.Context, truck models.Truck) (models.Truck, error) {
	var query = ` INSERT INTO truck
            (
						company_id,
                        name,
                        carrying,
                        year,
                        current_location,
                        trailer_parameters,
                        trailer_type
            )
            VALUES
            (
						:company_id,
                        :name,
                        :carrying,
                        :year,
                        :current_location,
                        :trailer_parameters,
                        :trailer_type
            )
            returning id`

	rows, err := sqlx.NamedQueryContext(ctx, t.dbc, query, truck)
	if err != nil {
		return models.Truck{}, errors.Wrap(err, "executing query")
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println("closing rows:", err)
		}
	}()

	if rows.Next() {
		err = rows.StructScan(&truck)
		if err != nil {
			return models.Truck{}, errors.Wrap(err, "scanning result")
		}
	}

	return truck, nil
}

// GetProjects returns all users projects.
func (t truck) Trucks(ctx context.Context, companyID int64) ([]models.Truck, error) {
	var (
		query = `SELECT *
				 FROM   truck
				 WHERE  company_id = $1`
		trucks []models.Truck
		err    error
	)

	err = sqlx.SelectContext(ctx, t.dbc, &trucks, query, companyID)
	if err != nil {
		return nil, errors.Wrap(err, "getting trucks")
	}

	return trucks, nil
}
