package models

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Account a struct to user info.
type Company struct {
	ID         int64     `db:"id" json:"id"`
	Kind       string    `db:"kind" json:"kind"`
	Name       string    `db:"name" json:"name"`
	Phone      string    `db:"phone" json:"phone"`
	Email      string    `db:"email" json:"email"`
	BankDetail string    `db:"bank_detail" json:"bank_detail"`
	CreatedAt  time.Time `db:"created_at" json:"-"`
	UpdatedAt  time.Time `db:"updated_at" json:"-"`
}

// Validate checks incoming account details.
func (c *Company) Validate() error {
	if !strings.Contains(c.Email, "@") {
		return errors.New("email address is required")
	}

	if c.Kind == "" {
		return errors.New("company type is required")
	}

	if c.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
