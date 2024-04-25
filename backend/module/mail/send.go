package mail

import (
	"fmt"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/module/mail/common"
	"github.com/apicat/apicat/v2/backend/module/mail/sendcloud"
	"github.com/apicat/apicat/v2/backend/module/mail/smtp"
)

const (
	SMTP      = "smtp"
	SENDCLOUD = "sendcloud"
)

type Sender struct {
	Driver    string
	Smtp      smtp.SmtpSender
	SendCloud sendcloud.SendCloud
}

func NewSender(cfg Sender) common.Provider {
	slog.Debug("mail.NewSender", "cfg", cfg)

	if cfg.Driver == SMTP {
		return smtp.NewSMTP(cfg.Smtp)
	} else if cfg.Driver == SENDCLOUD {
		return sendcloud.NewSendCloud(cfg.SendCloud)
	}
	return nil
}

func NewMessage(subject string, body fmt.Stringer) *common.Message {
	return &common.Message{
		Subject:  subject,
		Body:     body.String(),
		BodyType: "text/html",
	}
}
