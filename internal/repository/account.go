// Package repository defines ability to work with the database(PostgreSQL).
package repository

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/IgorRybak2055/logistic-service/internal/models"
)

type account struct {
	dbc *sqlx.DB
}

// NewAccountRepository will create an object that represent the Account interface
func NewAccountRepository(dbc *sqlx.DB) Account {
	return &account{dbc}
}

// CreateAccount use account data for registration new account in database.
func (a account) CreateAccount(ctx context.Context, account models.Account) (models.Account, error) {
	var query = `INSERT INTO account (
                        NAME,
                        email,
                        password,
						phone,
						company_id,
                        created_at,
                        updated_at
            )
            VALUES (
                        :name,
                        :email,
                        :password,
						:phone,
						:company_id,
                        :created_at,
                        :updated_at
            )
            returning id`

	var (
		rows = &sqlx.Rows{}
		err  error
	)

	rows, err = sqlx.NamedQueryContext(ctx, a.dbc, query, account)
	if err != nil {
		return models.Account{}, errors.Wrap(err, "executing query")
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println("error close rows:", err)
		}
	}()

	if rows.Next() {
		err = rows.StructScan(&account)
		if err != nil {
			return models.Account{}, errors.Wrap(err, "scanning rows:")
		}
	}

	return account, nil
}

// Login checks duplicate emails. Email must be unique.
func (a account) GetByEmail(ctx context.Context, email string) (models.Account, error) {
	var (
		query = `SELECT *
				 FROM   account
				 WHERE  email = $1`
		acc = models.Account{}
		err error
	)

	err = sqlx.GetContext(ctx, a.dbc, &acc, query, email)
	if err != nil {
		return models.Account{}, errors.Wrap(err, "getting data")
	}

	return acc, nil
}

// GetByID returns account data from database by accountID.
func (a account) GetByID(ctx context.Context, accountID int64) (models.Account, error) {
	var (
		query = `SELECT email,
       					updated_at
				 FROM   account
				 WHERE  id = $1`
		acc = models.Account{}
		err error
	)

	err = sqlx.GetContext(ctx, a.dbc, &acc, query, accountID)
	if err != nil {
		return models.Account{}, errors.Wrap(err, "getting data")
	}

	return acc, nil
}

func (a account) SetNewPassword(ctx context.Context, newPassword string) error {
	var (
		query = `UPDATE account
			      SET 	 password = $1,
       					 updated_at = $2
				  WHERE  id = $3`
		err    error
		userID = ctx.Value("user")
	)

	_, err = a.dbc.QueryxContext(ctx, query, newPassword, time.Now(), userID)

	return err
}
