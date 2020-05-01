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

type topic struct {
	dbc *sqlx.DB
}

// NewTopicRepository will create an object that represent the Topic interface
func NewTopicRepository(dbc *sqlx.DB) Topic {
	return &topic{dbc}
}

// func (t topic) tmp(ctx context.Context, updateTime time.Time, projectID int64) error {
// 	var query = `UPDATE project set updated_at = $1 where id = $2 AND user_id = $3`
//
// 	result, err := t.dbc.ExecContext(ctx, query, updateTime, projectID, ctx.Value("user"))
// 	if err != nil {
// 		return errors.Wrap(err, "updating project")
// 	}
//
// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
//
// 	if rows != 1 {
// 		return errors.New("no changes")
// 		// log.Fatalf("expected to affect 1 row, affected %d", rows)
// 	}
//
// 	return nil
// }

// CreateAccount use account data for registration new account in database.
func (t topic) CreateTopic(ctx context.Context, topic models.Topic) (models.Topic, error) {
	var (
		query = `INSERT INTO topic `
		cols  = ` project_id,
				 title,
				 description`

		values = ` :project_id,
				   :title,
				   :description`
	)

	if topic.ParentID.Int64 != 0 {
		cols = "parent_id,\n" + cols
		values = ":parent_id,\n" + values
	}

	query += `(` +
		cols +
		`) 
			VALUES 
			(` +
		values +
		`)
			returning id`

	rows, err := sqlx.NamedQueryContext(ctx, t.dbc, query, topic)
	if err != nil {
		return models.Topic{}, errors.Wrap(err, "creating project")
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println("closing rows:", err)
		}
	}()

	if rows.Next() {
		err = rows.StructScan(&topic)
		if err != nil {
			return models.Topic{}, errors.Wrap(err, "scanning result")
		}
	}

	return topic, nil
}

// GetProjects returns all users projects.
func (t topic) GetTopics(ctx context.Context, userID int64, projectID string) ([]models.Topic, error) {
	var (
		query = `SELECT topic.*
				 FROM   topic
				 	    JOIN project p
						  ON topic.project_id = p.id
					    JOIN account a
						  ON p.user_id = a.id
				 WHERE  a.id = $1
					    AND p.id = $2`
		topics []models.Topic
		err    error
	)

	err = sqlx.SelectContext(ctx, t.dbc, &topics, query, userID, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "getting topics")
	}

	return topics, nil
}

// GetProjectByID returns user project by ID.
func (t topic) GetTopicByID(ctx context.Context, userID int64, topicID string) (models.Topic, error) {
	var (
		query = `SELECT topic.*
				 FROM   topic
				 	    JOIN project p
						  ON topic.project_id = p.id
					    JOIN account a
						  ON p.user_id = a.id
				 WHERE  a.id = $1
					    AND topic.id = $2`
		topic models.Topic
		err   error
	)

	err = sqlx.GetContext(ctx, t.dbc, &topic, query, userID, topicID)
	if err != nil {
		return models.Topic{}, errors.Wrap(err, "getting topic")
	}

	return topic, nil
}

// DeleteProject delete user project by ID.
func (t topic) DeleteTopic(ctx context.Context, userID int64, topicID string) error {
	var (
		query = `DELETE FROM topic t
				 USING  account AS a
				 WHERE  a.id = $1
       					AND t.id = $2`
	)

	result, err := t.dbc.ExecContext(ctx, query, userID, topicID)
	if err != nil {
		return errors.Wrap(err, "deleting topic")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("permission for deleting denied")
	}

	return nil
}

// UpdateProject update user project.
func (t topic) UpdateTopic(ctx context.Context, userID int64, topicID string,
	upds map[string]interface{}) (models.Topic, error) {

	// Update topic SET title='dsdsds' WHERE project_id = ANY(SELECT project.id from project join account a on project.user_id = a.id where a.id = 1) AND id = 3;

	var (
		query  = `	UPDATE topic
					SET    title = 'dsdsds'
					WHERE  project_id = ANY (SELECT project.id
                         					 FROM   project
                                			 		join account a
                                  			 			ON project.user_id = a.id
                         					 WHERE  a.id = 1)
       						AND id = 3;  `
		err    error
		topic  models.Topic
		comma  = ""
		params []interface{}
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

	params = append(params, time.Now(), userID, topicID)

	query += `, updated_at = $3 WHERE user_id = $4 AND id = $5 RETURNING *`
	row := t.dbc.QueryRowxContext(ctx, query, params...)

	if err = row.StructScan(&topic); err != nil {
		return models.Topic{}, errors.Wrap(err, "updating topic")
	}

	return topic, nil
}
