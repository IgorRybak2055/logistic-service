// Package email provides opportunity to send email message.
package email

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"

	"github.com/pkg/errors"
)

// Sender describes the function of sending a message
type Sender interface {
	send(MessageData) error
}

// Config stores configs for email.
type Config struct {
	URL            string `config:"SMTP_URL,required"`
	Port           string `config:"SMTP_PORT,required"`
	Sender         string `config:"SMTP_SENDER,required"`
	SenderPassword string `config:"SMTP_SENDER_PASS,required"`
	TemplatePath   string `config:"SMTP_TEMPLATE,required"`
	RestoreURL     string `config:"-"`
}

// MessageData store data necessary for send email.
type MessageData struct {
	RecvEmail string
	UserID    string
	UserToken string
	message   string
}

// Service provides methods to send email message.
type Service struct {
	cfg       *Config
	messageCh <-chan MessageData
}

// New returns Service for start email service.
func New(cfg *Config, ch <-chan MessageData) *Service {
	return &Service{
		cfg:       cfg,
		messageCh: ch,
	}
}

// Run start mail service
func (s Service) Run() {
	var ctx = context.TODO()

LOOP:
	for {
		select {
		case msg := <-s.messageCh:
			go handleSend(s, msg, s.cfg.RestoreURL, s.cfg.TemplatePath)
		case <-ctx.Done():
			break LOOP
		}
	}
}

func handleSend(s Sender, e MessageData, restoreURL, templatePath string) {
	var err error

	type TemplateData struct {
		Address  string
		UserHash string
	}

	var data = TemplateData{
		Address:  restoreURL,
		UserHash: e.UserToken,
	}

	e.message, err = parseTemplate(filepath.Join(templatePath, "template.html"), data)
	if err != nil {
		log.Println("parsing letter template: ", err.Error())
		return
	}

	if err = s.send(e); err != nil {
		return
	}
}

func (s Service) send(e MessageData) error {
	var emailAuth = smtp.PlainAuth("", s.cfg.Sender, s.cfg.SenderPassword, s.cfg.URL)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "Ragger: Restore account password.\n"
	msg := []byte(subject + mime + "\n" + e.message)

	err := smtp.SendMail(s.cfg.URL+":"+s.cfg.Port, emailAuth, s.cfg.Sender, []string{e.RecvEmail}, msg)
	if err != nil {
		return errors.Wrap(err, "sending email")
	}

	return nil
}

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	var t, err = template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}

	var buf = new(bytes.Buffer)

	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
