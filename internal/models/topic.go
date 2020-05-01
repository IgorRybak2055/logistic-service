// Package models defines used entities in Ragger.
package models

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// Topic a struct to topic info.
type Topic struct {
	ID          int64         `db:"id" json:"id"`
	ParentID    sql.NullInt64 `db:"parent_id" json:"parent_id,omitempty"`
	ProjectID   int64         `db:"project_id" json:"project_id"`
	Title       string        `db:"title" json:"title"`
	Description string        `db:"description" json:"description"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}

// Validate checks incoming profect details.
func (t *Topic) Validate() error {
	if t.Title == "" {
		return errors.New("email address is required")
	}

	if t.Description == "" {
		return errors.New("email address is required")
	}

	return nil
}
