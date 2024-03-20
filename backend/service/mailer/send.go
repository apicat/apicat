package mailer

import (
	"fmt"

	"github.com/apicat/apicat/backend/config"
	"github.com/apicat/apicat/backend/module/mail"

	"log/slog"
)

func Send(subject string, content fmt.Stringer, to ...string) error {
	sender, err := mail.NewSender(config.Get().Email.ToMapInterface())
	if err != nil {
		return err
	}
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
