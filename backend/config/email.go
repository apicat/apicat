package config

import (
	"net/mail"

	mailsender "github.com/apicat/apicat/v2/backend/module/mail"
	"github.com/apicat/apicat/v2/backend/module/mail/sendcloud"
	"github.com/apicat/apicat/v2/backend/module/mail/smtp"
)

type Email struct {
	Driver    string          `yaml:"Driver"`
	SendCloud *EmailSendCloud `yaml:"SendCloud"`
	Smtp      *EmailSmtp      `yaml:"Smtp"`
}

// sendCloud
// https://www.sendcloud.net
type EmailSendCloud struct {
	ApiUser  string `yaml:"ApiUser" json:"apiUser"`
	ApiKey   string `yaml:"ApiKey" json:"apiKey"`
	From     string `yaml:"From" json:"fromEmail"`
	FromName string `yaml:"FromName" json:"fromName"`
}

// Smtp
type EmailSmtp struct {
	Host     string       `yaml:"Host" json:"Host"`
	From     mail.Address `yaml:"From" json:"From"`
	Password string       `yaml:"Password" json:"Password"`
}

func SetEmail(emailConfig *Email) {
	globalConf.Email = emailConfig
}

func (e *Email) ToCfg() mailsender.Sender {
	switch e.Driver {
	case mailsender.SMTP:
		return mailsender.Sender{
			Driver: e.Driver,
			Smtp: smtp.SmtpSender{
				Host:     e.Smtp.Host,
				From:     e.Smtp.From,
				Password: e.Smtp.Password,
			},
		}
	case mailsender.SENDCLOUD:
		return mailsender.Sender{
			Driver: e.Driver,
			SendCloud: sendcloud.SendCloud{
				ApiUser:  e.SendCloud.ApiUser,
				ApiKey:   e.SendCloud.ApiKey,
				From:     e.SendCloud.From,
				FromName: e.SendCloud.FromName,
			},
		}
	default:
		return mailsender.Sender{
			Driver: "",
		}
	}
}
