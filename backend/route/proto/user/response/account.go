package response

import (
	userbase "apicat-cloud/backend/route/proto/user/base"
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
