package config

import (
	"errors"
	"net/mail"
	"os"

	mailsender "github.com/apicat/apicat/v2/backend/module/mail"
	"github.com/apicat/apicat/v2/backend/module/mail/sendcloud"
	"github.com/apicat/apicat/v2/backend/module/mail/smtp"
)

type Email struct {
	Driver    string
	SendCloud *EmailSendCloud
	Smtp      *EmailSmtp
}

// sendCloud
// https://www.sendcloud.net
type EmailSendCloud struct {
	ApiUser     string `json:"apiUser"`
	ApiKey      string `json:"apiKey"`
	FromAddress string `json:"fromAddress"`
	FromName    string `json:"fromName"`
}

// Smtp
type EmailSmtp struct {
	Host     string `json:"host"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func LoadEmailConfig() {
	globalConf.Email = &Email{}
	if v, exists := os.LookupEnv("EMAIL_DRIVER"); exists {
		switch v {
		case mailsender.SMTP:
			globalConf.Email.Driver = mailsender.SMTP
			loadSmtpConfig()
		case mailsender.SENDCLOUD:
			globalConf.Email.Driver = mailsender.SENDCLOUD
			loadSendCloudConfig()
		}
	}
}

func loadSendCloudConfig() {
	globalConf.Email.SendCloud = &EmailSendCloud{}
	if v, exists := os.LookupEnv("SENDCLOUD_API_USER"); exists {
		globalConf.Email.SendCloud.ApiUser = v
	} else {
		return
	}
	if v, exists := os.LookupEnv("SENDCLOUD_API_KEY"); exists {
		globalConf.Email.SendCloud.ApiKey = v
	} else {
		return
	}
	if v, exists := os.LookupEnv("SENDCLOUD_FROM_ADDRESS"); exists {
		globalConf.Email.SendCloud.FromAddress = v
	} else {
		return
	}
	if v, exists := os.LookupEnv("SENDCLOUD_FROM_NAME"); exists {
		globalConf.Email.SendCloud.FromName = v
	} else {
		return
	}
}

func loadSmtpConfig() {
	globalConf.Email.Smtp = &EmailSmtp{}
	if v, exists := os.LookupEnv("SMTP_HOST"); exists {
		globalConf.Email.Smtp.Host = v
	} else {
		return
	}
	if v, exists := os.LookupEnv("SMTP_ADDRESS"); exists {
		globalConf.Email.Smtp.Address = v
	} else {
		return
	}
	if v, exists := os.LookupEnv("SMTP_NAME"); exists {
		globalConf.Email.Smtp.Name = v
	}
	if v, exists := os.LookupEnv("SMTP_PASSWORD"); exists {
		globalConf.Email.Smtp.Password = v
	} else {
		return
	}
}

func CheckEmailConfig() error {
	if globalConf.Email.Driver == "" {
		return nil
	}

	switch globalConf.Email.Driver {
	case mailsender.SMTP:
		if globalConf.Email.Smtp == nil {
			return errors.New("smtp config is empty")
		}
		if globalConf.Email.Smtp.Host == "" {
			return errors.New("smtp host is empty")
		}
		if globalConf.Email.Smtp.Address == "" {
			return errors.New("smtp address is empty")
		}
		if globalConf.Email.Smtp.Password == "" {
			return errors.New("smtp password is empty")
		}
	case mailsender.SENDCLOUD:
		if globalConf.Email.SendCloud == nil {
			return errors.New("sendcloud config is empty")
		}
		if globalConf.Email.SendCloud.ApiUser == "" {
			return errors.New("sendcloud api user is empty")
		}
		if globalConf.Email.SendCloud.ApiKey == "" {
			return errors.New("sendcloud api key is empty")
		}
		if globalConf.Email.SendCloud.FromAddress == "" {
			return errors.New("sendcloud from address is empty")
		}
		if globalConf.Email.SendCloud.FromName == "" {
			return errors.New("sendcloud from name is empty")
		}
	}
	return nil
}

func SetEmail(emailConfig *Email) {
	globalConf.Email = emailConfig
}

func (e *Email) ToCfg() mailsender.Sender {
	if e == nil {
		return mailsender.Sender{}
	}

	switch e.Driver {
	case mailsender.SMTP:
		return mailsender.Sender{
			Driver: e.Driver,
			Smtp: smtp.SmtpSender{
				Host: e.Smtp.Host,
				From: mail.Address{
					Address: e.Smtp.Address,
					Name:    e.Smtp.Name,
				},
				Password: e.Smtp.Password,
			},
		}
	case mailsender.SENDCLOUD:
		return mailsender.Sender{
			Driver: e.Driver,
			SendCloud: sendcloud.SendCloud{
				ApiUser:  e.SendCloud.ApiUser,
				ApiKey:   e.SendCloud.ApiKey,
				From:     e.SendCloud.FromAddress,
				FromName: e.SendCloud.FromName,
			},
		}
	default:
		return mailsender.Sender{}
	}
}
