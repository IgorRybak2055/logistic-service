// Package models defines used entities in Ragger.
package models

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// Token (JWT) claims struct.
type Token struct {
	UserID int64 `json:"id"`
	jwt.StandardClaims
}

// Account a struct to user info.
type Account struct {
	ID        int64             `db:"id" json:"id"`
	Name      string            `db:"name" json:"name"`
	Email     string            `db:"email" json:"email"`
	Password  string            `db:"password" json:"password,omitempty"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt time.Time         `db:"updated_at"`
	Token     map[string]string `json:"tokens"`
}

// Validate checks incoming account details.
func (a *Account) Validate() error {
	if !strings.Contains(a.Email, "@") {
		return errors.New("email address is required")
	}

	if len(a.Password) < 6 {
		return errors.New("password is required")
	}

	if a.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
