package common

type Provider interface {
	Send(msg *Message, to ...string) error
}
