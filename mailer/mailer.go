package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

const (
	mime = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
)

type Mailer interface {
	SendEmail() error
	ParseTemplate(filename string, data interface{}) error
}

type mailer struct {
	Auth    smtp.Auth
	Request *request
}

func NewMailer(Auth smtp.Auth, Request *request) Mailer {
	return &mailer{
		Auth:    Auth,
		Request: Request,
	}
}

func (m *mailer) SendEmail() error {
	subject := "Subject: " + m.Request.Subject + "!\n"
	msg := []byte(subject + mime + "\n" + m.Request.Body)
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", m.Request.Host, m.Request.Port),
		m.Auth,
		m.Request.From,
		m.Request.To,
		msg,
	)
}

func (m *mailer) ParseTemplate(filename string, data interface{}) error {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	m.Request.Body = buf.String()
	return nil
}
