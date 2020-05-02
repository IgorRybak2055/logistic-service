package services

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/logistic-service/internal/models"
	"github.com/IgorRybak2055/logistic-service/internal/repository"
)

type companyService struct {
	companyRepo repository.Company
	log         *logrus.Logger
}

// NewCompanyService will create new accountService object representation of Account interface
func NewCompanyService(cr repository.Company, logger *logrus.Logger) Company {
	return &companyService{
		companyRepo: cr,
		log:         logger,
	}
}

// Create allows to create a new user in the system,
// checks the incoming of this account according to the requirements.
func (s companyService) Create(ctx context.Context, company models.Company) (models.Company, error) {
	var err error

	if err = company.Validate(); err != nil {
		s.log.Debugf("create account: failed to validate account: %s", err)
		return models.Company{}, err
	}

	var createTime = time.Now().UTC()

	company.CreatedAt, company.UpdatedAt = createTime, createTime

	s.log.Debug("create account: successfully created")

	return s.companyRepo.Create(ctx, company)
}