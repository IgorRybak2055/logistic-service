// Package ragger defines App work and functions to configure App.
package logistic

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/IgorRybak2055/logistic-service/internal/repository"
	"github.com/IgorRybak2055/logistic-service/internal/services"
	"github.com/IgorRybak2055/logistic-service/pkg/email"
)

var (
	errProjectNotSpecified = errors.New("project not specified in request")
	errTopicNotSpecified   = errors.New("topic not specified in request")
	errBadRequest          = errors.New("not enough data in request")
)

// HTTPConfig stores configs for ragger.
type HTTPConfig struct {
	Host     string `config:"HTTP_HOST,required"`
	Addr     string `config:"HTTP_ADDR,required"`
	LogLevel string `config:"LOG_LEVEL,required"`
}

// FullRestoreURL return restore url.
func (h HTTPConfig) FullRestoreURL() string {
	return h.Host + h.Addr + "/api/restore_password"
}

// App is a structure for ragger server.
type App struct {
	cfg    *HTTPConfig
	Logger *logrus.Logger
	Srv    *http.Server
	DBC    *sqlx.DB

	companyService   services.Company
	accountService services.Account
	projectService services.Project
	topicService   services.Topic

	sendToEmailCh chan email.MessageData
}

// New returns App for start App
func New(cfg *HTTPConfig, sendToEmailCg chan email.MessageData) *App {
	return &App{
		cfg:           cfg,
		Logger:        logrus.New(),
		sendToEmailCh: sendToEmailCg,
	}
}

// Start configures all needs App fields and start App
func (a *App) Start() error {
	projectRepo := repository.NewProjectRepository(a.DBC)

	a.companyService = services.NewCompanyService(repository.NewCompanyRepository(a.DBC), a.Logger)
	a.accountService = services.NewAccountService(repository.NewAccountRepository(a.DBC), a.Logger)
	a.projectService = services.NewProjectService(projectRepo, a.Logger)
	a.topicService = services.NewTopicService(repository.NewTopicRepository(a.DBC), projectRepo, a.Logger)

	router := mux.NewRouter()

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	router.Handle("/health", handle(a.handleHealth)).Methods(http.MethodGet)
	router.HandleFunc("/api/company/register", handle(a.handleCompanyRegister)).Methods(http.MethodPost)
	router.HandleFunc("/api/register", handle(a.handleRegister)).Methods(http.MethodPost)
	router.HandleFunc("/api/login", handle(a.handleLogin)).Methods(http.MethodPost)
	router.HandleFunc("/api/token", handle(a.handleGenerateToken)).Methods(http.MethodGet)
	router.HandleFunc("/api/restore_password", handle(a.handleRestorePassword))

	api := router.PathPrefix("/api").Subrouter()

	api.Use(JwtAuthentication)
	// api.HandleFunc("/new_password", handle(a.handleNewPassword)).Methods(http.MethodPost)

	// // actions with projects
	// api.HandleFunc("/projects", handle(a.handleGetAllProject)).Methods(http.MethodGet)
	// api.HandleFunc("/projects/{project_id}", handle(a.handleGetProject)).Methods(http.MethodGet)
	// api.HandleFunc("/projects", handle(a.handleNewProject)).Methods(http.MethodPost)
	// api.HandleFunc("/projects/{project_id}", handle(a.handleUpdateProject)).Methods(http.MethodPut)
	// api.HandleFunc("/projects/{project_id}", handle(a.handleDeleteProject)).Methods(http.MethodDelete)
	//
	// // actions with topic
	// api.HandleFunc("/topics", handle(a.handleGetAllTopics)).Methods(http.MethodGet)
	// api.HandleFunc("/topics/{topic_id}", handle(a.handleGetTopic)).Methods(http.MethodGet)
	// api.HandleFunc("/topics", handle(a.handleNewTopic)).Methods(http.MethodPost)
	// api.HandleFunc("/topics/{topic_id}", handle(a.handleUpdateTopic)).Methods(http.MethodPut)
	// api.HandleFunc("/topics/{topic_id}", handle(a.handleDeleteTopic)).Methods(http.MethodDelete)

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	a.Srv = &http.Server{
		Handler: router,
		Addr:    a.cfg.Addr,
	}

	a.Logger.Info("starting api server...")

	return a.Srv.ListenAndServe()
}
