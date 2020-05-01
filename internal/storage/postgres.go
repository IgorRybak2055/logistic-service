// Package storage responsible for connecting with the database and send connection object to the ragger.
package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Config stores configs for postgres storage.
type Config struct {
	Host    string `config:"DATABASE_HOST,required"`
	Name    string `config:"DATABASE_NAME,required"`
	User    string `config:"DATABASE_USER,required"`
	Pass    string `config:"DATABASE_PASSWORD,required"`
	SSLMode string `config:"DATABASE_SSLMODE,required"`
}

// Postgres returns postgres dsn.
func (c Config) Postgres() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Name, c.User, c.Pass, c.SSLMode)
}

// Connect return database connection and error if postgres connection don't open
func Connect(dsn string, logger *logrus.Logger) (*sqlx.DB, error) {
	var err error

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("opening db connection:", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Fatal("pinging db connection:", err)
		return nil, err
	}

	return db, nil
}
