package response

import (
	userbase "github.com/apicat/apicat/backend/route/proto/user/base"
)

type Oauth2User struct {
	UserData
	userbase.TokenResponse
	Bind *userbase.UserOauthBindOption `json:"bind"`
}

type RegisterFireRes struct {
	userbase.MessageTemplate
	userbase.TokenResponse
}
