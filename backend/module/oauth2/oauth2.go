package oauth2

import (
	"context"
)

type Driver interface {
	GetEndpoints() Endpoints
	GetUser(context.Context, *Token) (*AuthUser, error)
}

type Config struct {
	ClientID     string `yaml:"ClientID"`
	ClientSecret string `yaml:"ClientSecret"`
}

type Token struct {
	Value string
	Type  string
}

func (t Token) String() string {
	if t.Type == "" {
		return "bearer " + t.Value
	}
	return t.Type + " " + t.Value
}

type Endpoints struct {
	Authorize string
	Token     string
}

type AuthUser struct {
	ID     string
	Email  string
	Name   string
	Avatar string
}
