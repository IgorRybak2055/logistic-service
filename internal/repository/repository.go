// Package repository defines ability to work with the database(PostgreSQL).
package repository

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // import for migration
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // ...
	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/logistic-service/internal/models"
)

// Company represents possible database actions with an account.
type Company interface {
	Create(ctx context.Context, company models.Company) (models.Company, error)
}

type Delivery interface {
	Create(ctx context.Context, delivery models.Delivery) (models.Delivery, error)
}

// Account represents possible database actions with an account.
type Account interface {
	CreateAccount(ctx context.Context, account models.Account) (models.Account, error)
	GetByEmail(ctx context.Context, email string) (models.Account, error)
	GetByID(ctx context.Context, accountID int64) (models.Account, error)
	SetNewPassword(ctx context.Context, newPassword string) error
}

// Project represents possible database actions with an project.
type Project interface {
	DeleteProject(ctx context.Context, userID int64, projectID string) error
	CreateProject(ctx context.Context, project models.Project) (models.Project, error)
	GetProjects(ctx context.Context, userID int64) ([]models.Project, error)
	GetProjectByID(ctx context.Context, userID int64, projectID string) (models.Project, error)
	UpdateProject(ctx context.Context, userID int64, projectID string, upds map[string]interface{}) (models.Project, error)
	UpdateProjectTime(ctx context.Context, updateTime time.Time, projectID, userID int64) error
}

// Topic represents possible database actions with a topic.
type Topic interface {
	DeleteTopic(ctx context.Context, userID int64, topicID string) error
	CreateTopic(ctx context.Context, topic models.Topic) (models.Topic, error)
	GetTopics(ctx context.Context, userID int64, projectID string) ([]models.Topic, error)
	GetTopicByID(ctx context.Context, userID int64, topicID string) (models.Topic, error)
	UpdateTopic(ctx context.Context, userID int64, topicID string, upds map[string]interface{}) (models.Topic, error)
}

// MakeMigrations provides an opportunity to work with migrations
func MakeMigrations(dsn string, logger *logrus.Logger) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logger.Fatal("creating new migrations:", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("upping migrations:", err)
		return err
	}

	return nil
}
