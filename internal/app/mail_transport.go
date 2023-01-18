package app

import (
	"fmt"
	"net/smtp"
)

type MailTransportConfig struct {
	Host     string
	Port     string
	From     string
	Password string
}

func (cfg MailTransportConfig) Address() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

type MailTransport struct {
	config MailTransportConfig
}

func (mt *MailTransport) Send(to, message string) error {
	err := smtp.SendMail(mt.config.Address(), nil, mt.config.From, []string{to}, []byte(message))

	if err != nil {
		return err
	}

	return nil
}

func NewMailTransport(cfg MailTransportConfig) *MailTransport {
	return &MailTransport{cfg}
}
