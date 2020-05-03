// Package ragger defines App work and functions to configure App.
package logistic

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/IgorRybak2055/logistic-service/internal/models"
)

// Health check godoc
// @Summary Ragger health check
// @Description Health check ragger service
// @Produce  json
// @Success 200 {string} string "response structure: {status:"UP"}"
// @Router /health [get]
func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var answer = map[string]string{"status": "UP"}

	if err := json.NewEncoder(w).Encode(answer); err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	return nil
}

// Company registration godoc
// @Summary Registration company of logistic-service service
// @Description Create a new company.
// @Tags Company
// @Produce  json
// @Param company body models.Company true "company"
// @Success 200 {object} models.Account "response structure: {message:"answer"}"
// @Failure 400 {string} string "response structure: {error:"error message"}"
// @Failure 500 {string} string "response structure: {error:"error message"}"
// @Router /api/company/register [post]
func (a *App) handleCompanyRegister(w http.ResponseWriter, r *http.Request) error {
	var newCompany models.Company

	if err := json.NewDecoder(r.Body).Decode(&newCompany); err != nil {
		a.Logger.Errorf("failed to decode company: %s", err)
		return err
	}

	var account, err = a.companyService.Create(r.Context(), newCompany)
	if err != nil {
		a.Logger.Warn("creating account:", err)

		return newError(http.StatusBadRequest, err)
	}

	Respond(w, http.StatusOK, Message(account))

	return nil
}

// Registration godoc
// @Summary Registration user of contact service
// @Description Create a new user with the input name, email & password.
// @Tags Account
// @Produce  json
// @Param account body models.Account true "account"
// @Success 200 {object} models.Account "response structure: {message:"answer"}"
// @Failure 400 {string} string "response structure: {error:"error message"}"
// @Router /api/register [post]
func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) error {
	var newAccount models.Account

	if err := json.NewDecoder(r.Body).Decode(&newAccount); err != nil {
		a.Logger.Errorf("failed to decode account: %s", err)
		return err
	}

	var account, err = a.accountService.CreateAccount(r.Context(), newAccount)
	if err != nil {
		a.Logger.Warn("creating account:", err)

		return newError(http.StatusBadRequest, err)
	}

	Respond(w, http.StatusOK, Message(account))

	return nil
}

// Set new password godoc
// @Summary Set new password
// @Description Setting new password for account
// @Tags Account
// @Produce json
// @Param password query string true "new password for account"
// @Param confirm_password query string true "confirm_password new password"
// @Success 200
// @Failure 400 {string} string "response structure: {error:"error message"}"
// @Router /api/new_password [post]
// @Security ApiKeyAuth
func (a *App) handleNewPassword(w http.ResponseWriter, r *http.Request) error {
	var (
		password        = r.FormValue("password")
		confirmPassword = r.FormValue("confirm_password")
	)

	if password != confirmPassword {
		a.Logger.Warnf("creating account: %s", "mismatched passwords")

		return newError(http.StatusBadRequest, errors.New("password mismatch"))
	}

	var err = a.accountService.SetNewPassword(r.Context(), password)
	if err != nil {
		return newError(http.StatusBadRequest, err)
	}

	Respond(w, http.StatusOK, nil)

	return nil
}

// Login godoc
// @Summary Login in ragger
// @Description Login in ragger with email and password
// @Tags Account
// @Produce  json
// @Param email query string true "account email"
// @Param password query string true "account password len(password) > 6"
// @Success 200 {object} models.Account "response structure: {message:"answer"}"
// @Failure 400 {string} string "response structure: {error:"error message"}"
// @Router /api/login [post]
func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var (
		email    = r.FormValue("email")
		password = r.FormValue("password")
	)

	var account, err = a.accountService.Login(r.Context(), email, password)
	if err != nil {
		a.Logger.Warn("login:", err)

		return newError(http.StatusBadRequest, err)
	}

	Respond(w, http.StatusOK, Message(account))

	return nil
}

// Generate Token godoc
// @Summary Generate Token in ragger
// @Description Generate Token for access to ragger
// @Tags Token
// @Produce json
// @Param refresh_token query string true "last generated refresh_token"
// @Success 200 {object} string "response structure: {message:"access_token:access, refresh_token:refresh"}"
// @Failure 401 {string} string "response structure: {error:"error message"}"
// @Router /api/token [get]
func (a *App) handleGenerateToken(w http.ResponseWriter, r *http.Request) error {
	var refreshToken = r.FormValue("refresh_token")

	newTokenPair, err := a.accountService.GenerateToken(r.Context(), refreshToken)
	if err != nil {
		a.Logger.Warn("generating token:", err)

		return newError(http.StatusUnauthorized, err)
	}

	Respond(w, http.StatusOK, Message(newTokenPair))

	return nil
}

// Restore password godoc
// @Summary Restore password
// @Description Returned html page for setting new password
// @Tags Account
// @Produce html
// @Success 200
// @Failure 500 {string} string "response structure: {error:"error message"}"
// @Router /api/restore_password [get]
func (a *App) handleRestorePassword(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		var body, err = ioutil.ReadFile("./assets/new_password.html")
		if err != nil {
			a.Logger.Warn("reading file:", err)

			return newError(http.StatusInternalServerError, err)
		}

		fmt.Fprint(w, string(body))
	case http.MethodPost:
		var email = r.FormValue("email")

		var err = a.accountService.RestorePassword(r.Context(), a.sendToEmailCh, email)
		if err != nil {
			a.Logger.Warn("checking email:", err)

			return newError(http.StatusBadRequest, err)
		}

		Respond(w, http.StatusOK, Message("check your email please"))
	default:
		return newError(http.StatusBadRequest, errors.New("undefined http method for this request"))
	}

	return nil
}

func contextAccountID(ctx context.Context) int64 { return ctx.Value("user").(int64) }

// New delivery godoc
// @Summary Creating new delivery.
// @Description Creating new delivery.
// @Tags Delivery
// @Produce json
// @Param delivery body models.Delivery true "delivery"
// @Success 201 {object} models.Delivery "response structure: {message:delivery}"
// @Failure 400 {string} string "response structure: {error:"error message"}"
// @Router /api/deliveries [post]
// @Security ApiKeyAuth
func (a *App) handleCreateDelivery(w http.ResponseWriter, r *http.Request) error {
	var (
		ctx      = r.Context()
		delivery models.Delivery
	)

	if err := json.NewDecoder(r.Body).Decode(&delivery); err != nil {
		a.Logger.Errorf("failed to decode delivery: %s", err)
		return err
	}

	var account, err = a.deliveryService.CreateDelivery(ctx, delivery)
	if err != nil {
		return newError(http.StatusBadRequest, err)
	}

	Respond(w, http.StatusCreated, Message(account))

	return nil
}

// // Get all project godoc
// // @Summary Get all projects of user.
// // @Description Get all projects of authorized user.
// // @Tags Project
// // @Produce json
// // @Success 200 {array} models.Project "response structure: {message:[]project}"
// // @Failure 400 {string} string "response structure: {error:"error message"}"
// // @Router /api/projects [get]
// // @Security ApiKeyAuth
// func (a *App) handleGetAllProject(w http.ResponseWriter, r *http.Request) error {
// 	ctx := r.Context()
// 	userID := contextUserID(ctx)
//
// 	usersProjects, err := a.projectService.GetUserProjects(ctx, userID)
// 	if err != nil {
// 		return newError(http.StatusInternalServerError, err)
// 	}
//
// 	Respond(w, http.StatusOK, Message(usersProjects))
//
// 	return nil
// }
//
// // Get project godoc
// // @Summary Get projects by ID.
// // @Description Get project of authorized user by ID.
// // @Tags Project
// // @Produce json
// // @Param project_id path string true "project title"
// // @Success 200 {object} models.Project "response structure: {message:project}"
// // @Failure 400 {string} string "response structure: {error:"error message"}"
// // @Router /api/projects/{project_id} [get]
// // @Security ApiKeyAuth
// func (a *App) handleGetProject(w http.ResponseWriter, r *http.Request) error {
// 	var (
// 		vars      = mux.Vars(r)
// 		projectID string
// 		ok        bool
// 		ctx       context.Context
// 	)
//
// 	if projectID, ok = vars["project_id"]; !ok {
// 		return newError(http.StatusBadRequest, errProjectNotSpecified)
// 	}
//
// 	userID := contextUserID(ctx)
//
// 	project, err := a.projectService.GetProject(ctx, userID, projectID)
// 	if err != nil {
// 		return newError(http.StatusInternalServerError, err)
// 	}
//
// 	Respond(w, http.StatusOK, Message(project))
//
// 	return nil
// }
//
// // Delete project godoc
// // @Summary Delete projects by ID.
// // @Description Delete project by ID of authorized user.
// // @Tags Project
// // @Produce json
// // @Param project_id path string true "project title"
// // @Success 200
// // @Failure 400 {string} string "response structure: {error:"error message"}"
// // @Router /api/projects/{project_id} [delete]
// // @Security ApiKeyAuth
// func (a *App) handleDeleteProject(w http.ResponseWriter, r *http.Request) error {
// 	var (
// 		vars      = mux.Vars(r)
// 		projectID string
// 		ok        bool
// 		ctx       = r.Context()
// 	)
//
// 	if projectID, ok = vars["project_id"]; !ok {
// 		return newError(http.StatusBadRequest, errProjectNotSpecified)
// 	}
//
// 	userID := contextUserID(ctx)
//
// 	err := a.projectService.DeleteProject(ctx, userID, projectID)
// 	if err != nil {
// 		return newError(http.StatusInternalServerError, err)
// 	}
//
// 	Respond(w, http.StatusOK, nil)
//
// 	return nil
// }
//
// // Update project godoc
// // @Summary Update projects by ID.
// // @Description Update project by ID of authorized user.
// // @Tags Project
// // @Produce json
// // @Param project_id path string true "project title"
// // @Param title query string false "project title"
// // @Param description query string false "project description"
// // @Success 200 {object} models.Project "response structure: {message:project}"
// // @Failure 400 {string} string "response structure: {error:"error message"}"
// // @Router /api/projects/{project_id} [put]
// // @Security ApiKeyAuth
// func (a *App) handleUpdateProject(w http.ResponseWriter, r *http.Request) error {
// 	var (
// 		vars      = mux.Vars(r)
// 		projectID string
// 		ok        bool
// 		ctx       = r.Context()
// 	)
//
// 	if projectID, ok = vars["project_id"]; !ok {
// 		return newError(http.StatusBadRequest, errProjectNotSpecified)
// 	}
//
// 	var upds = make(map[string]interface{}, 1)
//
// 	if title := r.FormValue("title"); title != "" {
// 		upds["title"] = title
// 	}
//
// 	if description := r.FormValue("description"); description != "" {
// 		upds["description"] = description
// 	}
//
// 	if len(upds) == 0 {
// 		a.Logger.Warnf("handleUpdateProject: failed to parse query: %s", errBadRequest)
// 		return newError(http.StatusBadRequest, errBadRequest)
// 	}
//
// 	userID := contextUserID(ctx)
//
// 	project, err := a.projectService.UpdateProject(r.Context(), userID, projectID, upds)
// 	if err != nil {
// 		return newError(http.StatusInternalServerError, err)
// 	}
//
// 	Respond(w, http.StatusOK, Message(project))
//
// 	return nil
// }
//
// // New topic godoc
// // @Summary Creating new topic.
// // @Description  Creating new topic with title and description.
// // @Tags Topic
// // @Produce json
// // @Param project query string true "project id"
// // @Param parent query string false "parent topic id"
// // @Param title query string true "topic title"
// // @Param description query string true "topic description"
// // @Success 201 {object} models.Topic "response structure: {message:topic}"
// // @Failure 400 {string} string "response structure: {error:"error message"}"
// // @Router /api/topics [post]
// // @Security ApiKeyAuth
// func (a *App) handleNewTopic(w http.ResponseWriter, r *http.Request) error {
// 	var ctx = r.Context()
//
// 	projectID, err := strconv.ParseInt(r.FormValue("project"), 10, 64)
// 	if err != nil {
// 		a.Logger.Warnf("parsing project id: %s", err)
// 		return newError(http.StatusBadRequest, errors.Wrap(err, "failed to parse project id"))
// 	}
//
// 	var (
// 		newTopic = models.Topic{
// 			Title:       r.FormValue("title"),
// 			Description: r.FormValue("description"),
// 			ProjectID:   projectID,
// 		}
// 		parentID int64
// 	)
//
// 	if parent := r.FormValue("parent"); parent != "" {
// 		parentID, err = strconv.ParseInt(parent, 10, 64)
// 		if err != nil {
// 			a.Logger.Warnf("parsing project id: %s", err)
// 		}
//
// 		newTopic.ParentID.Int64 = parentID
// 	}
//
// 	topic, err := a.topicService.NewTopic(ctx, newTopic, contextUserID(ctx))
// 	if err != nil {
// 		return newError(http.StatusBadRequest, err)
// 	}
//
// 	Respond(w, http.StatusCreated, Message(topic))
//
// 	return nil
// }
//
// // Get all topics by project godoc
// // @Summary Get all topics by project.
// // @Description Get all topics by project of authorized user .
// // @Tags Topic
// // @Produce json
// // @Param project query string true "project id"
// // @Success 200 {array} models.Topic "response structure: {message:[]topics}"
// // @Failure 400 {string} string "response structure: {error:"error message"}"
// // @Router /api/topics [get]
// // @Security ApiKeyAuth
// func (a *App) handleGetAllTopics(w http.ResponseWriter, r *http.Request) error {
// 	var (
// 		ctx       = r.Context()
// 		userID    = contextUserID(ctx)
// 		projectID = r.FormValue("project")
// 	)
//
// 	usersProjects, err := a.topicService.GetTopics(ctx, userID, projectID)
// 	if err != nil {
// 		return newError(http.StatusInternalServerError, err)
// 	}
//
// 	Respond(w, http.StatusOK, Message(usersProjects))
//
// 	return nil
// }
//
// // Get topic godoc
// // @Summary Get topic by ID.
// // @Description Get topic by ID of authorized user .
// // @Tags Topic
// // @Produce json
// // @Param topic_id path string true "topic ID"
// // @Success 200 {object} models.Topic "response structure: {message:topic}"
// // @Failure 400 {string} string "response structure: {error:"error message"}"
// // @Router /api/topics/{topic_id} [get]
// // @Security ApiKeyAuth
// func (a *App) handleGetTopic(w http.ResponseWriter, r *http.Request) error {
// 	var (
// 		vars    = mux.Vars(r)
// 		topicID string
// 		ok      bool
// 		ctx     = r.Context()
// 		userID  = contextUserID(ctx)
// 	)
//
// 	if topicID, ok = vars["topic_id"]; !ok {
// 		return newError(http.StatusBadRequest, errTopicNotSpecified)
// 	}
//
// 	topic, err := a.topicService.GetTopic(ctx, userID, topicID)
// 	if err != nil {
// 		return newError(http.StatusInternalServerError, err)
// 	}
//
// 	Respond(w, http.StatusOK, Message(topic))
//
// 	return nil
// }
//
// // Delete topic godoc
// // @Summary Delete topic by ID.
// // @Description Delete topic by ID of authorized user.
// // @Tags Topic
// // @Produce json
// // @Param topic_id path string true "topic ID"
// // @Success 200
// // @Failure 400 {string} string "response structure: {error:"error message"}"
// // @Router /api/projects/{project_id} [delete]
// // @Security ApiKeyAuth
// func (a *App) handleDeleteTopic(w http.ResponseWriter, r *http.Request) error {
// 	var (
// 		vars    = mux.Vars(r)
// 		topicID string
// 		ok      bool
// 	)
//
// 	if topicID, ok = vars["topic_id"]; !ok {
// 		return newError(http.StatusBadRequest, errProjectNotSpecified)
// 	}
//
// 	var ctx = r.Context()
//
// 	userID := contextUserID(ctx)
//
// 	err := a.topicService.DeleteTopic(ctx, userID, topicID)
// 	if err != nil {
// 		return newError(http.StatusInternalServerError, err)
// 	}
//
// 	Respond(w, http.StatusOK, nil)
//
// 	return nil
// }
//
// // Update topic godoc
// // @Summary Update topic by ID.
// // @Description Update topic by ID of authorized user.
// // @Tags Topic
// // @Produce json
// // @Param topic_id path string true "topic id"
// // @Param title query string false "topic title"
// // @Param description query string false "topic description"
// // @Success 200 {object} models.Topic "response structure: {message:topic}"
// // @Failure 400 {string} string "response structure: {error:"error message"}"
// // @Router /api/topics/{topic_id} [put]
// // @Security ApiKeyAuth
// func (a *App) handleUpdateTopic(w http.ResponseWriter, r *http.Request) error {
// 	var (
// 		vars    = mux.Vars(r)
// 		topicID string
// 		ok      bool
// 		ctx     = r.Context()
// 	)
//
// 	if topicID, ok = vars["topic_id"]; !ok {
// 		return newError(http.StatusBadRequest, errProjectNotSpecified)
// 	}
//
// 	var upds = make(map[string]interface{}, 1)
//
// 	if title := r.FormValue("title"); title != "" {
// 		upds["title"] = title
// 	}
//
// 	if description := r.FormValue("description"); description != "" {
// 		upds["description"] = description
// 	}
//
// 	if len(upds) == 0 {
// 		a.Logger.Warnf("handleUpdateProject: failed to parse query: %s", errBadRequest)
// 		return newError(http.StatusBadRequest, errBadRequest)
// 	}
//
// 	userID := contextUserID(ctx)
//
// 	topic, err := a.topicService.UpdateTopic(r.Context(), userID, topicID, upds)
// 	if err != nil {
// 		return newError(http.StatusInternalServerError, err)
// 	}
//
// 	Respond(w, http.StatusOK, Message(topic))
//
// 	return nil
// }
