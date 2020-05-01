// Package repository defines ability to work with the database(PostgreSQL).
package repository

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/IgorRybak2055/logistic-service/internal/models"
)

type project struct {
	dbc *sqlx.DB
}

// NewProjectRepository will create an object that represent the Project interface
func NewProjectRepository(dbc *sqlx.DB) Project {
	return &project{dbc}
}

// DeleteProject delete user project by ID.
func (p project) CreateProject(ctx context.Context, project models.Project) (models.Project, error) {
	var query = ` INSERT INTO project
            (
                        title,
                        description,
                        user_id
            )
            VALUES
            (
                        :title,
                        :description,
                        :user_id
            )
            returning id`

	rows, err := sqlx.NamedQueryContext(ctx, p.dbc, query, project)
	if err != nil {
		return models.Project{}, errors.Wrap(err, "creating project")
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println("closing rows:", err)
		}
	}()

	if rows.Next() {
		err = rows.StructScan(&project)
		if err != nil {
			return models.Project{}, errors.Wrap(err, "scanning result")
		}
	}

	return project, nil
}

// GetProjects returns all users projects.
func (p project) GetProjects(ctx context.Context, userID int64) ([]models.Project, error) {
	var (
		query = `SELECT *
				 FROM   project
				 WHERE  user_id = $1`
		projects []models.Project
		err      error
	)

	err = sqlx.SelectContext(ctx, p.dbc, &projects, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, "getting projects")
	}

	return projects, nil
}

// GetProjectByID returns user project by ID.
func (p project) GetProjectByID(ctx context.Context, userID int64, projectID string) (models.Project, error) {
	var (
		query = `SELECT *
				 FROM   project
				 WHERE  user_id = $1
				 AND id = $2  `
		project models.Project
		err     error
	)

	if err = sqlx.GetContext(ctx, p.dbc, &project, query, userID, projectID); err != nil {
		return models.Project{}, errors.Wrap(err, "getting project")
	}

	return project, nil
}

// DeleteProject delete user project by ID.
func (p project) DeleteProject(ctx context.Context, userID int64, projectID string) error {
	var (
		query = `DELETE FROM project WHERE user_id = $1 AND id = $2`
		err   error
	)

	if _, err = p.dbc.ExecContext(ctx, query, userID, projectID); err != nil {
		return errors.Wrap(err, "deleting project")
	}

	return nil
}

// UpdateProject update user project.
func (p project) UpdateProject(ctx context.Context, userID int64, projectID string,
	upds map[string]interface{}) (models.Project, error) {
	var (
		query   = `UPDATE project SET`
		err     error
		project models.Project
		comma   = ""
		params  []interface{}
	)

	if title, ok := upds["title"]; ok {
		query += " title = $1"
		comma = ","

		params = append(params, title)
	}

	if description, ok := upds["description"]; ok {
		query += comma + " description = $2"

		params = append(params, description)
	}

	params = append(params, time.Now(), userID, projectID)

	query += `, updated_at = $3 WHERE user_id = $4 AND id = $5 RETURNING *`
	row := p.dbc.QueryRowxContext(ctx, query, params...)

	if err = row.StructScan(&project); err != nil {
		return models.Project{}, errors.Wrap(err, "updating project")
	}

	return project, nil
}

// UpdateProjectTime updates field updated_at when create or update topic of project.
func (p project) UpdateProjectTime(ctx context.Context, updateTime time.Time, projectID, userID int64) error {
	var query = `UPDATE project set updated_at = $1 where id = $2 AND user_id = $3`

	result, err := p.dbc.ExecContext(ctx, query, updateTime, projectID, userID)
	if err != nil {
		return errors.Wrap(err, "updating project")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("permission denied")
	}

	return nil
}
