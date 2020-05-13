package models

import "time"

type Tender struct {
	ID         int64     `db:"id" json:"id"`
	CompanyID  int64     `db:"company_id" json:"company_id"`
	DeliveryID int64     `db:"delivery_id" json:"delivery_id"`
	Start      time.Time `db:"start" json:"start"`
}
