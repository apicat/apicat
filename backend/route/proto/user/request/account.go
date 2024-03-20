package request

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	userbase "github.com/apicat/apicat/backend/route/proto/user/base"
)

type LoginOption struct {
	userbase.EmailOption
	PasswordOption
	protobase.InvitationTokenOption
}

type RegisterUserOption struct {
	userbase.EmailOption
	userbase.NameOption
	userbase.AvatarOption
	userbase.LanguageOption
	PasswordOption
	Bind *userbase.UserOauthBindOption `json:"bind"`
	protobase.InvitationTokenOption
}

type Oauth2StateOption struct {
	OauthOption
	protobase.InvitationTokenOption
	userbase.LanguageOption
}

type ResetPasswordOption struct {
	CodeOption
	PasswordOption
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}
