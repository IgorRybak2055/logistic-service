package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/IgorRybak2055/logistic-service/internal/models"
)

type company struct {
	dbc *sqlx.DB
}

// NewAccountRepository will create an object that represent the Account interface
func NewCompanyRepository(dbc *sqlx.DB) Company {
	return &company{dbc}
}

func (c company) Create (ctx context.Context, company models.Company) (models.Company, error) {
	var query = `INSERT INTO company (
                        company_type,
						name,
						phone,
                        email,
                        bank_detail,
                        created_at,
                        updated_at
            )
            VALUES (
                        :kind,
						:name,
						:phone,
                        :email,
                        :bank_detail,
                        :created_at,
                        :updated_at
            )
            returning id`

	var (
		rows = &sqlx.Rows{}
		err  error
	)

	rows, err = sqlx.NamedQueryContext(ctx, c.dbc, query, company)
	if err != nil {
		return models.Company{}, errors.Wrap(err, "executing query")
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println("error close rows:", err)
		}
	}()

	if rows.Next() {
		err = rows.StructScan(&company)
		if err != nil {
			return models.Company{}, errors.Wrap(err, "scanning rows:")
		}
	}

	return company, nil
}
