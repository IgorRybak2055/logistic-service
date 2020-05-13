// Package client represent client ragger API
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const contentType = "application/json"

var (
	templatesDir    = filepath.Join(".", "pkg", "client", "templates")
	mainTmpl        = template.Must(template.ParseFiles(filepath.Join(templatesDir, "main.html")))
	loginTmpl       = template.Must(template.ParseFiles(filepath.Join(templatesDir, "login.html")))
	registerTmpl    = template.Must(template.ParseFiles(filepath.Join(templatesDir, "register.html")))
	addDeliveryTmpl = template.Must(template.ParseFiles(filepath.Join(templatesDir, "addDelivery.html")))
	activeTenders = template.Must(template.ParseFiles(filepath.Join(templatesDir, "myTenders.html")))
)

// Client struct to client ragger API
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiToken   string
}

// NewClient return new client for ragger
func NewClient(url string, timeout time.Duration) *Client {
	return &Client{
		baseURL:  url,
		apiToken: "token",
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Run() {
	simpleMux := http.NewServeMux()
	simpleMux.HandleFunc("/", c.handleLogin)
	simpleMux.HandleFunc("/main", c.handleMain)
	simpleMux.HandleFunc("/add_delivery", c.handleAddDelivery)
	simpleMux.HandleFunc("/active_tenders", c.handleActiveTenders)
	simpleMux.HandleFunc("/delivery/info", c.handleDeliveryInfo)
	// simpleMux.HandleFunc("/registration", srv.handleRegistration)
	// simpleMux.Handle("/m/", loginHandler)

	log.Printf("Start server at :%v ...", ":8383")
	err := http.ListenAndServe(":8383", simpleMux)
	if err != nil {
		log.Fatal("Error happened:", err.Error())
	}
}

func (c *Client) request(method, path string, data url.Values, body, entity interface{}) error {
	fullPath := c.baseURL + path
	if data != nil {
		fullPath += "?" + data.Encode()
	}

	var buf io.ReadWriter

	if body != nil {
		buf = new(bytes.Buffer)

		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, fullPath, buf)
	if err != nil {
		return err
	}

	if err = c.call(req, entity); err != nil {
		return err
	}

	return nil
}

func (c *Client) get(path string, data url.Values, entity interface{}) error {
	return c.request(http.MethodGet, path, data, nil, entity)
}

func (c *Client) post(path string, data url.Values, body, entity interface{}) error {
	return c.request(http.MethodPost, path, data, body, entity)
}

func (c *Client) put(path string, data url.Values, body, entity interface{}) error {
	return c.request(http.MethodPut, path, data, body, entity)
}

func (c *Client) delete(path string) error {
	fullPath := c.baseURL + path

	req, err := http.NewRequest(http.MethodDelete, fullPath, nil)
	if err != nil {
		return err
	}

	if err = c.call(req, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) call(req *http.Request, entity interface{}) error {
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)

	if len(c.apiToken) > 0 {
		req.Header.Set("Authorization", c.apiToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "calling http")
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "reading http body")
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.Errorf("http request %d: %s", resp.StatusCode, string(body))
	}

	if entity == nil {
		return nil
	}

	log.Println("BODY --", string(body))

	err = json.Unmarshal(body, entity)

	return err
}

// Registration doing request to registration new user in ragger
func (c *Client) Registration(user User) (User, error) {
	var u User

	err := c.post("/api/register", nil, user, &u)

	return u, err
}

// Login doing request to login in ragger
func (c *Client) Login(user User) (User, error) {
	var u User

	err := c.post("/api/login", nil, user, &u)

	c.apiToken = u.Token["access_token"]

	return u, err
}

// Projects doing request for getting all user projects from ragger
func (c *Client) Deliveries() ([]Delivery, error) {
	var deliveries []Delivery

	err := c.get("/api/deliveries", nil, &deliveries)

	return deliveries, err
}

// CreateProject doing request for creating new user project and returns new project
func (c *Client) CreateDelivery(dlv Delivery) (Delivery, error) {
	var dl Delivery

	err := c.post("/api/deliveries", nil, dlv, &dl)

	return dl, err
}

// Project doing request for getting user projects from ragger by project id
func (c *Client) ActiveDeliveries() ([]Delivery, error) {
	var deliveries []Delivery

	err := c.get("/api/active_tenders", nil, &deliveries)

	return deliveries, err
}

// Project doing request for getting user projects from ragger by project id
func (c *Client) Delivery(deliveryID int64) (Delivery, error) {
	var dl Delivery

	err := c.get(fmt.Sprintf("/api/deliveries/%v", deliveryID), nil, &dl)

	return dl, err
}

// ----------------------------------------------------------------------------------------------------------------------


// Project doing request for getting user projects from ragger by project id
func (c *Client) Project(projectID int64) (Project, error) {
	var project Project

	err := c.get(fmt.Sprintf("/api/projects/%v", projectID), nil, &project)

	return project, err
}

// CreateProject doing request for creating new user project and returns new project
func (c *Client) CreateProject(project Project) (Project, error) {
	var pr Project

	err := c.post("/api/projects", nil, project, &pr)

	return pr, err
}

// UpdateProject doing a request to update an existing user project and returns the updated project
func (c *Client) UpdateProject(project Project) (Project, error) {
	var pr Project

	err := c.put(fmt.Sprintf("/api/projects/%v", project.ID), nil, project, &pr)

	return pr, err
}

// DeleteProject doing a request to delete an existing user project by project id
func (c *Client) DeleteProject(pr int64) error {
	return c.delete(fmt.Sprintf("/api/projects/%v", pr))
}

// CreateTopic doing request for creating new project topic and returns new topic
func (c *Client) CreateTopic(topic Topic) (Topic, error) {
	var t Topic

	data := url.Values{}
	data.Set("project", strconv.FormatInt(topic.ProjectID, 10))
	err := c.post("/api/topics", data, topic, &t)

	return t, err
}

// Topics doing request for getting all topics for project from ragger
func (c *Client) Topics(projectID int64) ([]Topic, error) {
	var topics []Topic

	data := url.Values{}

	data.Set("project", strconv.FormatInt(projectID, 10))

	err := c.get("/api/topics", data, &topics)

	return topics, err
}

// Topic doing request for getting topic from ragger
func (c *Client) Topic(topicID int64) (Topic, error) {
	var topic Topic

	err := c.get(fmt.Sprintf("/api/topics/%v", topicID), nil, &topic)

	return topic, err
}

// UpdateTopic doing a request to update an existing topic and returns the updated topic
func (c *Client) UpdateTopic(topic Topic) (Topic, error) {
	var t Topic

	err := c.put(fmt.Sprintf("/api/topics/%v", topic.ID), nil, topic, &t)

	return t, err
}

// DeleteTopic doing a request to delete an existing topic by id
func (c *Client) DeleteTopic(t int64) error {
	return c.delete(fmt.Sprintf("/api/topics/%v", t))
}
