package emails

import (
	"net/smtp"
	"os"
)

var (
	config     *smtpConfig
	smptClient *SMTP
)

type SMTP struct{}

func createNewSMTP() *SMTP {
	return &SMTP{}
}

func GetSMTPClient() *SMTP {
	if smptClient == nil {
		smptClient = createNewSMTP()
	}

	return smptClient
}

var auth smtp.Auth

type smtpConfig struct {
	username    string
	password    string
	host        string
	port        string
	senderEmail string
}

func (s *smtpConfig) GetAuth() smtp.Auth {
	if auth == nil {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}
	return auth
}

func getEmailConfig() *smtpConfig {
	if config == nil {
		config = createSmtpConfig()
	}
	return config
}

func createSmtpConfig() *smtpConfig {
	return &smtpConfig{
		username:    os.Getenv("EMAIL_USERNAME"),
		password:    os.Getenv("EMAIL_PASSWORD"),
		host:        os.Getenv("EMAIL_HOST"),
		port:        os.Getenv("EMAIL_PORT"),
		senderEmail: os.Getenv("EMAIL_SENDER_EMAIL"),
	}
}
