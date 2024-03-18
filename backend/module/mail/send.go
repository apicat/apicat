package mail

import (
	"apicat-cloud/backend/module/mail/common"
	"apicat-cloud/backend/module/mail/sendcloud"
	"apicat-cloud/backend/module/mail/smtp"
	"errors"
	"fmt"
	"log/slog"
)

const (
	SMTP      = "smtp"
	SENDCLOUD = "sendcloud"
)

func NewSender(cfg map[string]interface{}) (common.Provider, error) {
	slog.Debug("mail.NewSender", "cfg", cfg)
	if cfg == nil {
		return nil, errors.New("mail config is nil")
	}

	if cfg["Driver"] == SMTP {
		return smtp.NewSMTP(cfg["Smtp"].(map[string]interface{}))
	} else if cfg["Driver"] == SENDCLOUD {
		return sendcloud.NewSendCloud(cfg["SendCloud"].(map[string]interface{}))
	}
	return nil, errors.New("mail driver not found")
}

func NewMessage(subject string, body fmt.Stringer) *common.Message {
	return &common.Message{
		Subject:  subject,
		Body:     body.String(),
		BodyType: "text/html",
	}
}
