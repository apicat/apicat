package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/apicat/apicat/backend/module/oauth2"
)

type Github struct{}

func (*Github) GetEndpoints() oauth2.Endpoints {
	return oauth2.Endpoints{
		Authorize: "https://github.com/login/oauth/authorize",
		Token:     "https://github.com/login/oauth/access_token",
	}
}

func (*Github) GetUser(ctx context.Context, token *oauth2.Token) (*oauth2.AuthUser, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	req.Header.Set("Authorization", token.String())
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(res.StatusCode))
	}
	var ret struct {
		ID        int64   `json:"id"`
		Email     *string `json:"email"`
		Name      string  `json:"name"`
		AvatarUrl string  `json:"avatar_url"`
	}
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}
	o := &oauth2.AuthUser{
		ID:     fmt.Sprintf("%d", ret.ID),
		Name:   ret.Name,
		Avatar: ret.AvatarUrl,
	}
	if ret.Email != nil {
		o.Email = *ret.Email
	}
	return o, nil
}
