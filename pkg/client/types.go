// Package client represent client ragger API
package client

import (
	"encoding/json"
)

// User a struct to user info.
type User struct {
	ID        int64             `db:"id" json:"id"`
	CompanyId int64             `db:"company_id" json:"company_id"`
	Name      string            `db:"name" json:"name"`
	Email     string            `db:"email" json:"email"`
	Phone     string            `db:"phone" json:"phone,omitempty"`
	Password  string            `db:"password" json:"password,omitempty"`
	Token     map[string]string `json:"tokens"`
}

type Delivery struct {
	ID             int64   `db:"id" json:"id"`
	CompanyID      int64   `db:"company_id" json:"company_id"`
	ShipmentDate   string  `db:"shipment_date" json:"shipment_date"`
	ShipmentPlace  string  `db:"shipment_place" json:"shipment_place"`
	UnloadingPlace string  `db:"unloading_place" json:"unloading_place"`
	Cargo          string  `db:"cargo" json:"cargo"`
	WeightCargo    float64 `db:"weight_cargo" json:"weight_cargo"`
	VolumeCargo    float64 `db:"volume_cargo" json:"volume_cargo"`
	TrailerType    string  `db:"trailer_type" json:"trailer_type"`
	Price          float64 `db:"price" json:"price"`
	Status         string  `json:"status"`
	CreatedAt      string  `db:"created_at" json:"-"`
	UpdatedAt      string  `db:"updated_at" json:"-"`
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
		Status         string  `json:"status"`
	}

	var temp tmp

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// parseDate, err := time.Parse(time.RFC3339, temp.ShipmentDate)
	// if err != nil {
	// 	log.Println("err -", err)
	// 	return err
	// }

	if len(temp.ShipmentDate) > 10 {
		d.ShipmentDate = temp.ShipmentDate[:10]
	} else {
		d.ShipmentDate = temp.ShipmentDate
	}

	// d.ShipmentDate = parseDate
	d.ID = temp.ID
	d.CompanyID = temp.CompanyID
	d.ShipmentPlace = temp.ShipmentPlace
	d.UnloadingPlace = temp.UnloadingPlace
	d.Cargo = temp.Cargo
	d.WeightCargo = temp.WeightCargo
	d.VolumeCargo = temp.VolumeCargo
	d.TrailerType = temp.TrailerType
	d.Price = temp.Price
	d.Status = temp.Status

	return nil
}

// Project a struct to project info.
type Project struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Topic a struct to topic info.
type Topic struct {
	ID          int64  `json:"id"`
	ParentID    int64  `json:"parent_id,omitempty"`
	ProjectID   int64  `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
