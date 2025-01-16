package config

import (
	"net/smtp"
)

type SendOtpEmailData struct {
	OtpCode string
	Email   string
}

type Mailer struct {
	Auth smtp.Auth
	Cnf  *Config
}

func NewMailer(cnf *Config) *Mailer {
	username := cnf.Env.GetString("SMTP_USERNAME")
	password := cnf.Env.GetString("SMTP_PASSWORD")
	smtpHost := cnf.Env.GetString("SMTP_HOST")

	return &Mailer{
		Auth: smtp.PlainAuth("", username, password, smtpHost),
		Cnf:  cnf,
	}
}

func (m Mailer) SendEmail(to string, subject string, body string) error {
	from := m.Cnf.Env.GetString("SMTP_FROM")
	smtpPort := m.Cnf.Env.GetString("SMTP_PORT")
	smtpHost := m.Cnf.Env.GetString("SMTP_HOST")
	auth := m.Auth

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	return err
}
