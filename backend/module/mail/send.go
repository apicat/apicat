package mail

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/apicat/apicat/backend/module/mail/common"
	"github.com/apicat/apicat/backend/module/mail/sendcloud"
	"github.com/apicat/apicat/backend/module/mail/smtp"
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
