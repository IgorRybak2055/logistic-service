package models

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Truck struct {
	ID                int64   `db:"id" json:"id"`
	CompanyID         int64   `db:"company_id"`
	Name              string  `db:"name" json:"name"`
	TrailerType       string  `db:"trailer_type" json:"trailer_type"`
	TrailerParameters string  `db:"trailer_parameters" json:"trailer_parameters"`
	Carrying          float64 `db:"carrying" json:"carrying"`
	Year              int     `db:"year" json:"year"`
	CurrentLocation   string  `db:"current_location" json:"current_location"`
}

// Validate checks incoming account details.
func (t *Truck) Validate() error {
	switch {
	case t.Name == "":
		return errors.New("name is required")
	case t.TrailerType == "":
		return errors.New("trailer type is required")
	case len(strings.Split(t.TrailerParameters, "/")) != 3:
		return errors.New("trailer params (3) is required")
	case t.Carrying <= 0:
		return errors.New("carrying is required")
	case t.Year < 1980 || t.Year > time.Now().Year():
		return errors.New("year is required")
	default:
		// go on
	}

	return nil
}
