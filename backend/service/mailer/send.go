package mailer

import (
	"fmt"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/module/mail"

	"log/slog"
)

func Send(subject string, content fmt.Stringer, to ...string) error {
	sender := mail.NewSender(config.Get().Email.ToCfg())
	if sender == nil {
		return nil
	}
	msg := mail.NewMessage(subject, content)
	return sender.Send(msg, to...)
}

func AsyncSend(subject string, content fmt.Stringer, to ...string) {
	go func() {
		if err := Send(subject, content, to...); err != nil {
			slog.Error("async send mail faild", "err", err)
		}
	}()
}
