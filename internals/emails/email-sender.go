package emails

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type EmailRequest struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewEmailRequest(to []string, subject string) *EmailRequest {
	return &EmailRequest{
		to:      to,
		subject: subject,
	}
}

func (r *EmailRequest) SendEmail() (bool, error) {
	config := getEmailConfig()
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := fmt.Sprintf("%s:%v", config.host, config.port)
	err := smtp.SendMail(addr, config.GetAuth(), config.senderEmail, r.to, msg)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *EmailRequest) ParseTemplate(data interface{}, templateFileName ...string) error {
	t, err := template.ParseFiles(templateFileName...)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		fmt.Println("Parsing Template Failed ", err)
		return err
	}
	r.body = buf.String()
	return nil
}
