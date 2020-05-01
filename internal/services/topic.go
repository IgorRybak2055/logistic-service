// Package services contains the basic logic of a Ragger.
package services

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/logistic-service/internal/models"
	"github.com/IgorRybak2055/logistic-service/internal/repository"
)

type topicService struct {
	topicRepo   repository.Topic
	projectRepo repository.Project
	log         *logrus.Logger
}

// NewTopicService will create new projectService object representation of Topic interface
func NewTopicService(tr repository.Topic, pr repository.Project, logger *logrus.Logger) Topic {
	return &topicService{
		topicRepo:   tr,
		projectRepo: pr,
		log:         logger,
	}
}

func (s topicService) NewTopic(ctx context.Context, topic models.Topic, userID int64) (models.Topic, error) {
	if err := topic.Validate(); err != nil {
		s.log.Debugf("newptopic: failed to %s", err)
		return models.Topic{}, err
	}

	updTime := time.Now()

	err := s.projectRepo.UpdateProjectTime(ctx, updTime, topic.ProjectID, userID)
	if err != nil {
		s.log.Warnf("newtopic: cannot update project: %s", err)

		return models.Topic{}, err
	}

	topic, err = s.topicRepo.CreateTopic(ctx, topic)
	if err != nil {
		s.log.Warnf("newtopic: failed to %s", err)
		return models.Topic{}, errors.Wrap(err, "creating new topic")
	}

	return topic, err
}

func (s topicService) GetTopics(ctx context.Context, userID int64, projectID string) ([]models.Topic, error) {
	userTopics, err := s.topicRepo.GetTopics(ctx, userID, projectID)
	if err != nil {
		s.log.Warnf("gettopics: failed to %s", err)
		return nil, errors.Wrap(err, "getting topics by project")
	}

	return userTopics, err
}

func (s topicService) GetTopic(ctx context.Context, userID int64, topicID string) (models.Topic, error) {
	topic, err := s.topicRepo.GetTopicByID(ctx, userID, topicID)
	if err != nil {
		s.log.Warnf("gettopic: failed to %s", err)
		return models.Topic{}, errors.Wrap(err, "getting user projects")
	}

	return topic, err
}

func (s topicService) DeleteTopic(ctx context.Context, userID int64, topicID string) error {
	if err := s.topicRepo.DeleteTopic(ctx, userID, topicID); err != nil {
		s.log.Warnf("deleteproject: failed to %s", err)
		return errors.Wrap(err, "deleting project")
	}

	return nil
}

func (s topicService) UpdateTopic(ctx context.Context, userID int64, topicID string,
	upds map[string]interface{}) (models.Topic, error) {
	topic, err := s.topicRepo.UpdateTopic(ctx, userID, topicID, upds)
	if err != nil {
		s.log.Warnf("updatetopic: failed to %s", err)
		return models.Topic{}, errors.Wrap(err, "updating topic")
	}

	return topic, err
}
