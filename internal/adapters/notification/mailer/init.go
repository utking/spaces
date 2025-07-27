// Package mailer provides functionality to send email notifications using SMTP.
package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"strconv"

	"github.com/utking/spaces/internal/adapters/notification/mailer/templates"
	"github.com/utking/spaces/internal/application/domain"
	gomail "gopkg.in/mail.v2"
)

type Sensitive string

func (s Sensitive) String() string {
	return "(redacted)"
}

type Mailer struct {
	Host     string
	From     string
	Username Sensitive
	Password Sensitive
	Port     int32
	UseTLS   bool
}

// New creates a new instance of Mailer with the provided configuration.
func New(host string, port int32, username, password, from string, useTLS bool) *Mailer {
	return &Mailer{
		Host:     host,
		Port:     port,
		Username: Sensitive(username),
		Password: Sensitive(password),
		From:     from,
		UseTLS:   useTLS,
	}
}

// String returns a string representation of the Mailer instance, redacting sensitive information.
func (m *Mailer) String() string {
	return fmt.Sprintf(`Mailer{
		Host: %s,
		Port: %d,
		Username: %s, 
		Password: %s,
		From: %s,
		UseTLS: %s
	}`,
		m.Host, m.Port,
		m.Username.String(), m.Password.String(),
		m.From, strconv.FormatBool(m.UseTLS),
	)
}

// Send prepares the email message for sending and sends it using the configured SMTP server.
func (m *Mailer) Send(_ context.Context, msg *domain.Notification) error {
	msg.Trim()

	if valErr := msg.Validate(); valErr != nil {
		return fmt.Errorf("invalid notification message: %w", valErr)
	}

	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", m.From)
	message.SetHeader("To", msg.To)
	message.SetHeader("Subject", msg.Title)

	// Set email HTML body
	message.SetBody("text/html", msg.Message)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(
		m.Host,
		int(m.Port),
		string(m.Username), // Convert Sensitive to string for the dialer
		string(m.Password), // Convert Sensitive to string for the dialer
	)

	if dialer == nil {
		return errors.New("failed to create SMTP dialer")
	}

	if m.UseTLS {
		dialer.StartTLSPolicy = gomail.MandatoryStartTLS
	} else {
		dialer.TLSConfig = nil // Disable TLS if not using it
	}

	// Send the email
	return dialer.DialAndSend(message)
}

func RenderTemplate(
	_ context.Context,
	templateName string,
	data map[string]interface{},
) (string, error) {
	// templates.MailerTemplates is embed.FS with *.html files
	content, err := template.ParseFS(
		templates.MailerTemplates,
		templateName,
	)

	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var rendered bytes.Buffer

	if err = content.ExecuteTemplate(
		&rendered,
		templateName,
		data,
	); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return rendered.String(), nil
}
