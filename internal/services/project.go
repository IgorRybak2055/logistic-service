// Package services contains the basic logic of a Ragger.
package services

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/logistic-service/internal/models"
	"github.com/IgorRybak2055/logistic-service/internal/repository"
)

type projectService struct {
	projectRepo repository.Project
	log         *logrus.Logger
}

// NewProjectService will create new projectService object representation of Project interface
func NewProjectService(pr repository.Project, logger *logrus.Logger) Project {
	return &projectService{
		projectRepo: pr,
		log:         logger,
	}
}

func (s projectService) NewProject(ctx context.Context, project models.Project) (models.Project, error) {
	if err := project.Validate(); err != nil {
		s.log.Debugf("newproject: failed to %s", err)
		return models.Project{}, err
	}

	project, err := s.projectRepo.CreateProject(ctx, project)
	if err != nil {
		s.log.Warnf("newproject: failed to %s", err)
		return models.Project{}, errors.Wrap(err, "creating new project")
	}

	return project, err
}

func (s projectService) GetUserProjects(ctx context.Context, userID int64) ([]models.Project, error) {
	userProjects, err := s.projectRepo.GetProjects(ctx, userID)
	if err != nil {
		s.log.Warnf("getuserprojects: failed to %s", err)
		return nil, errors.Wrap(err, "getting user projects")
	}

	return userProjects, err
}

func (s projectService) GetProject(ctx context.Context, userID int64, projectID string) (models.Project, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, userID, projectID)
	if err != nil {
		s.log.Warnf("getproject: failed to %s", err)
		return models.Project{}, errors.Wrap(err, "getting user projects")
	}

	return project, err
}

func (s projectService) DeleteProject(ctx context.Context, userID int64, projectID string) error {
	if err := s.projectRepo.DeleteProject(ctx, userID, projectID); err != nil {
		s.log.Warnf("deleteproject: failed to %s", err)
		return errors.Wrap(err, "deleting project")
	}

	return nil
}

func (s projectService) UpdateProject(ctx context.Context, userID int64, projectID string,
	upds map[string]interface{}) (models.Project, error) {
	project, err := s.projectRepo.UpdateProject(ctx, userID, projectID, upds)
	if err != nil {
		s.log.Warnf("updateproject: failed to %s", err)
		return models.Project{}, errors.Wrap(err, "updating project")
	}

	return project, err
}
