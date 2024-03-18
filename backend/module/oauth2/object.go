package oauth2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Object struct {
	conf   Config
	driver Driver
}

func NewObject(cfg Config, d Driver) *Object {
	return &Object{
		conf:   cfg,
		driver: d,
	}
}

func (o *Object) GetAuthorizeUri(redirect string) string {
	q := url.Values{}
	q.Add("client_id", o.conf.ClientID)
	q.Add("scope", "user:email")
	if redirect != "" {
		q.Add("redirect_uri", redirect)
	}
	return fmt.Sprintf("%s?%s", o.driver.GetEndpoints().Authorize, q.Encode())
}

func (o *Object) GetUserByState(ctx context.Context, state string) (*AuthUser, error) {
	token, err := o.getAuthorizationToken(ctx, state)
	if err != nil {
		return nil, err
	}
	return o.getUser(ctx, token)
}

func (o *Object) getUser(ctx context.Context, token *Token) (*AuthUser, error) {
	return o.driver.GetUser(ctx, token)
}

func (o *Object) getAuthorizationToken(ctx context.Context, state string) (*Token, error) {
	q := url.Values{}
	q.Add("client_id", o.conf.ClientID)
	q.Add("client_secret", o.conf.ClientSecret)
	q.Add("code", state)
	uri := fmt.Sprintf("%s?%s", o.driver.GetEndpoints().Token, q.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	retparams, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}
	if retparams.Has("error") {
		return nil, errors.New(retparams.Get("error_description"))
	}
	var token Token
	token.Type = retparams.Get("token_type")
	token.Value = retparams.Get("access_token")
	return &token, nil
}
