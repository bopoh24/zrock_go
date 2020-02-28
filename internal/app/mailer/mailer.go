package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"time"

	"github.com/bopoh24/zrock_go/internal/app/settings"
)

// EmailTypes email type
type EmailTypes string

// email  types
const (
	EmailRegistration     EmailTypes = "registration"
	EmailPasswordRecovery            = "recovery"
)

var subjects map[EmailTypes]string = map[EmailTypes]string{
	EmailRegistration: "Confirm Email address",
}

// Mailer ...
type Mailer struct {
	auth smtp.Auth
}

// NewMailer ...
func NewMailer() *Mailer {
	return &Mailer{
		auth: smtp.PlainAuth(
			"",
			settings.App.SMTPUser,
			settings.App.SMTPPassword,
			settings.App.SMTPHost,
		),
	}
}

// SendEmail sends email type data
func (m *Mailer) SendEmail(mt EmailTypes, data interface{}, recipient string) error {
	htmlBody, err := m.processTemplate(mt, data)
	if err != nil {
		return err
	}
	emailBody := mailBody(recipient, subjects[mt], string(htmlBody))
	// dump email
	if settings.App.SMTPHost == "" {
		return ioutil.WriteFile(
			fmt.Sprintf("dumps/%s_%d.eml", mt, time.Now().Unix()), emailBody, 0644)
	}
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", settings.App.SMTPHost, settings.App.SMTPPort),
		m.auth,
		settings.App.SMTPFrom,
		[]string{recipient},
		emailBody,
	)
}

func (m *Mailer) processTemplate(mt EmailTypes, data interface{}) ([]byte, error) {

	templateName := filepath.Join("../../../assets/email_templates", fmt.Sprintf("%s.html", mt))
	t, err := template.ParseFiles(templateName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func mailBody(to string, subject string, body string) []byte {
	return []byte(
		fmt.Sprintf("From: ZRock <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"utf-8\"\r\n"+
			"%s\r\n", settings.App.SMTPFrom, to, subject, body))
}
