package app

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"mini-sns-ws/internal/domain"
	"mini-sns-ws/internal/templates"
	"net/smtp"
)

type MailTransportConfig struct {
	Host     string
	Port     string
	Password string
}

func (cfg MailTransportConfig) Address() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

type MailTransport struct {
	config MailTransportConfig
	body   string
}

func (transport *MailTransport) Send(mail *Mail) error {
	buffer := new(bytes.Buffer)

	if err := mail.Template.Execute(buffer, mail.Data); err != nil {
		return err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + mail.Subject + "\n"
	message := []byte(subject + mime + buffer.String())

	if err := smtp.SendMail(transport.config.Address(), nil, mail.From, mail.To, message); err != nil {
		return err
	}

	return nil
}

func NewMailTransport(cfg MailTransportConfig) *MailTransport {
	return &MailTransport{config: cfg, body: ""}
}

type Mail struct {
	Template *template.Template
	To       []string
	From     string
	Subject  string
	Data     interface{}
}

func NewMail(template *template.Template, to []string, from, subject string, data interface{}) (*Mail, error) {
	return &Mail{
		Template: template,
		To:       to,
		From:     from,
		Subject:  subject,
		Data:     data,
	}, nil
}

func SendVerificationEmail(transport MailTransport, user domain.User, verificationToken string) error {
	data := struct {
		Name string
		URL  string
	}{
		Name: user.FullName(),
		URL:  "http://localhost:6943/api/v1/verify?token=" + verificationToken,
	}

	mailTo := []string{user.Email}
	subject := "Verify your account"

	template := template.Must(template.New("layout").Parse(templates.VerifyAccount))
	mail, err := NewMail(template, mailTo, "noreply@mini-sns.com", subject, data)

	if err != nil {
		log.Println(err)
		return err
	}

	if err := transport.Send(mail); err != nil {
		return err
	}

	return nil
}
