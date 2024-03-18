package base

type MessageTemplate struct {
	Emoji       string `json:"emoji"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TokenResponse
}

type OauthTypeOption struct {
	// oauth平台类型 如github
	Type string `uri:"type" json:"type" query:"type" binding:"required,oneof=github"`
}

type UserOauthBindOption struct {
	// oauth平台类型 如github
	OauthTypeOption
	// oauth对应平台的用户id
	OauthUserID string `json:"oauthUserID" binding:"required"`
}

type EmailOption struct {
	Email string `json:"email" binding:"required,email,lte=255"`
}

type NameOption struct {
	Name string `json:"name" binding:"required,gte=2,lte=64"`
}

type AvatarOption struct {
	Avatar string `json:"avatar" binding:"omitempty,url"`
}

type LanguageOption struct {
	Language string `json:"language" binding:"required,oneof=zh-CN en-US"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
}
