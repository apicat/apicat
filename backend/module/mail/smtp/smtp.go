package smtp

import (
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/mail/common"
)

type SmtpSender struct {
	From     mail.Address
	Password string
	Host     string
}

func NewSMTP(cfg SmtpSender) *SmtpSender {
	return &cfg
}

func (s *SmtpSender) Send(msg *common.Message, to ...string) error {
	msg.From = s.From
	msg.To = strings.Join(to, ",")
	bys, err := msg.Bytes()
	if err != nil {
		return err
	}
	// ps := strings.Split(c.From.Address, "@")
	auth := smtp.PlainAuth(
		"",
		s.From.Address,
		s.Password,
		strings.Split(s.Host, ":")[0],
	)
	return smtp.SendMail(s.Host, auth, s.From.Address, to, bys)
}
