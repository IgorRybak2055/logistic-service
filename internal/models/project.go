// Package models defines used entities in Ragger.
package models

import (
	"time"

	"github.com/pkg/errors"
)

// Project a struct to project info.
type Project struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	UserID      int64     `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Validate checks incoming profect details.
func (p *Project) Validate() error {
	if p.Title == "" {
		return errors.New("email address is required")
	}

	if p.Description == "" {
		return errors.New("email address is required")
	}

	return nil
}
