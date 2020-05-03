// Package services contains the basic logic of a Ragger.
package services

import (
	"context"

	"github.com/IgorRybak2055/logistic-service/internal/models"
	"github.com/IgorRybak2055/logistic-service/pkg/email"
)

// Company represents possible actions with a company.
type Company interface {
	Create(ctx context.Context, account models.Company) (models.Company, error)
}

// Token represents possible actions with a token.
type Token interface {
	GenerateToken(ctx context.Context, refreshToken string) (map[string]string, error)
}

// Account represents possible actions with an account.
type Account interface {
	CreateAccount(ctx context.Context, account models.Account) (models.Account, error)
	Login(ctx context.Context, email, password string) (models.Account, error)
	RestorePassword(ctx context.Context, ch chan email.MessageData, emailAddress string) error
	SetNewPassword(ctx context.Context, newPassword string) error
	Token
}

// Delivery represents possible actions with a deliveries.
type Delivery interface {
	CreateDelivery(ctx context.Context, delivery models.Delivery) (models.Delivery, error)
	// GetUserProjects(ctx context.Context, userID int64) ([]models.Project, error)
	// GetProject(ctx context.Context, userID int64, projectID string) (models.Project, error)
	// DeleteProject(ctx context.Context, userID int64, projectID string) error
	// UpdateProject(ctx context.Context, userID int64, projectID string, upds map[string]interface{}) (models.Project, error)
}

// Project represents possible actions with a project.
type Project interface {
	NewProject(ctx context.Context, project models.Project) (models.Project, error)
	GetUserProjects(ctx context.Context, userID int64) ([]models.Project, error)
	GetProject(ctx context.Context, userID int64, projectID string) (models.Project, error)
	DeleteProject(ctx context.Context, userID int64, projectID string) error
	UpdateProject(ctx context.Context, userID int64, projectID string, upds map[string]interface{}) (models.Project, error)
}

// Topic represents possible actions with a topic.
type Topic interface {
	NewTopic(ctx context.Context, topic models.Topic, userID int64) (models.Topic, error)
	GetTopics(ctx context.Context, userID int64, projectID string) ([]models.Topic, error)
	GetTopic(ctx context.Context, userID int64, topicID string) (models.Topic, error)
	DeleteTopic(ctx context.Context, userID int64, topicID string) error
	UpdateTopic(ctx context.Context, userID int64, topicID string, upds map[string]interface{}) (models.Topic, error)
}
