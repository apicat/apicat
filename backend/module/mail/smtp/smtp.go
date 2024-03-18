package smtp

import (
	"apicat-cloud/backend/module/mail/common"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

type smtpSender struct {
	from     mail.Address
	password string
	host     string
}

func NewSMTP(cfg map[string]interface{}) (*smtpSender, error) {
	for _, v := range []string{"From", "Password", "Host"} {
		if _, ok := cfg[v]; !ok {
			return nil, fmt.Errorf("smtp config %s is required", v)
		}
	}
	return &smtpSender{
		from: mail.Address{
			Name:    cfg["From"].(map[string]interface{})["Name"].(string),
			Address: cfg["From"].(map[string]interface{})["Address"].(string),
		},
		password: cfg["Password"].(string),
		host:     cfg["Host"].(string),
	}, nil
}

func (s *smtpSender) Send(msg *common.Message, to ...string) error {
	msg.From = s.from
	msg.To = strings.Join(to, ",")
	bys, err := msg.Bytes()
	if err != nil {
		return err
	}
	// ps := strings.Split(c.From.Address, "@")
	auth := smtp.PlainAuth(
		"",
		s.from.Address,
		s.password,
		strings.Split(s.host, ":")[0],
	)
	return smtp.SendMail(s.host, auth, s.from.Address, to, bys)
}
