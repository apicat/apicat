package config

import (
	"encoding/json"
	"net/mail"
)

type Email struct {
	Driver    string          `yaml:"Driver"`
	SendCloud *EmailSendCloud `yaml:"SendCloud"`
	Smtp      *EmailSmtp      `yaml:"Smtp"`
}

// sendCloud
// https://www.sendcloud.net
type EmailSendCloud struct {
	ApiUser  string `yaml:"ApiUser"`
	ApiKey   string `yaml:"ApiKey"`
	From     string `yaml:"From"`
	FromName string `yaml:"FromName"`
}

// Smtp
type EmailSmtp struct {
	Host     string       `yaml:"Host"`
	From     mail.Address `yaml:"From"`
	Password string       `yaml:"Password"`
}

func SetEmail(emailConfig *Email) {
	globalConf.Email = emailConfig
}

func (e *Email) ToMapInterface() map[string]interface{} {
	var (
		res      map[string]interface{}
		jsonByte []byte
	)
	jsonByte, _ = json.Marshal(e)
	_ = json.Unmarshal(jsonByte, &res)
	return res
}
