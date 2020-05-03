package models

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Delivery struct {
	ID             int64     `db:"id" json:"id"`
	CompanyID      int64     `db:"company_id" json:"company_id"`
	ShipmentDate   time.Time `db:"shipment_date" json:"shipment_date"`
	ShipmentPlace  string    `db:"shipment_place" json:"shipment_place"`
	UnloadingPlace string    `db:"unloading_place" json:"unloading_place"`
	Cargo          string    `db:"cargo" json:"cargo"`
	WeightCargo    float64   `db:"weight_cargo" json:"weight_cargo"`
	VolumeCargo    float64   `db:"volume_cargo" json:"volume_cargo"`
	TrailerType    string    `db:"trailer_type" json:"trailer_type"`
	Price          float64   `db:"price" json:"price"`
	CreatedAt      time.Time `db:"created_at" json:"-"`
	UpdatedAt      time.Time `db:"updated_at" json:"-"`
}

// Validate checks incoming account details.
func (d *Delivery) Validate() error {
	switch {
	case d.ShipmentDate.Before(time.Now()):
		return errors.New("shipment date is required")
	case d.ShipmentPlace == "":
		return errors.New("shipment place is required")
	case d.UnloadingPlace == "":
		return errors.New("unloading place is required")
	case d.Cargo == "":
		return errors.New("cargo describe is required")
	case d.WeightCargo <= 0:
		return errors.New("cargo weight is required")
	case d.VolumeCargo <= 0:
		return errors.New("cargo volume is required")
	case d.TrailerType == "":
		return errors.New("trailer type is required")
	case d.Price == 0:
		return errors.New("start price is required")
	default:
		// go on
	}

	return nil
}

func (d *Delivery) UnmarshalJSON(data []byte) error {
	type tmp struct {
		ID             int64   `json:"id"`
		CompanyID      int64   `json:"company_id"`
		ShipmentDate   string  `json:"shipment_date"`
		ShipmentPlace  string  `json:"shipment_place"`
		UnloadingPlace string  `json:"unloading_place"`
		Cargo          string  `json:"cargo"`
		WeightCargo    float64 `json:"weight_cargo"`
		VolumeCargo    float64 `json:"volume_cargo"`
		TrailerType    string  `json:"trailer_type"`
		Price          float64 `json:"price"`
	}

	var temp tmp

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	parseDate, err := time.Parse("2006-01-02", temp.ShipmentDate)
	if err != nil {
		return err
	}

	d.ShipmentDate = parseDate
	d.ID = temp.ID
	d.CompanyID = temp.CompanyID
	d.ShipmentPlace = temp.ShipmentPlace
	d.UnloadingPlace = temp.UnloadingPlace
	d.Cargo = temp.Cargo
	d.WeightCargo = temp.WeightCargo
	d.VolumeCargo = temp.VolumeCargo
	d.TrailerType = temp.TrailerType
	d.Price = temp.Price

	return nil
}
