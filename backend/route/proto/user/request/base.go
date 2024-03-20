package request

import "github.com/apicat/apicat/backend/route/proto/user/base"

type PasswordOption struct {
	Password string `json:"password" binding:"required,gte=6,lte=64"`
}

type OauthOption struct {
	base.OauthTypeOption
	// oauth授权码
	Code string `query:"code" json:"code" binding:"required"`
}

type CodeOption struct {
	Code string `uri:"code" binding:"required,len=32"`
}
