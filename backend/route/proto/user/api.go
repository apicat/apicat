package user

import (
	"github.com/apicat/apicat/backend/route/proto/user/base"
	"github.com/apicat/apicat/backend/route/proto/user/request"
	"github.com/apicat/apicat/backend/route/proto/user/response"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

// AccountApi 账户不需要登录
type AccountApi interface {
	// Login 登录
	// @route POST /account/login
	Login(*gin.Context, *request.LoginOption) (*base.TokenResponse, error)

	// Register 注册
	// @route POST /account/register
	Register(*gin.Context, *request.RegisterUserOption) (*base.TokenResponse, error)

	// LoginWithOauthCode oauth登录 需要先执行上一步 GotoOauth2 如果已绑定返回token 否则需要去注册
	// @route POST /account/oauth/{type}/login
	LoginWithOauthCode(*gin.Context, *request.Oauth2StateOption) (*response.Oauth2User, error)

	// RegisterFire 邮箱验证
	// @route PUT /account/email_verification/{code}
	RegisterFire(*gin.Context, *request.CodeOption) (*response.RegisterFireRes, error)

	// SendResetPasswordMail 发送找回密码邮件
	// @route POST /account/retrieve_password
	SendResetPasswordMail(*gin.Context, *base.EmailOption) (*ginrpc.Empty, error)

	// ResetPasswordCheck 检查重置密码令牌
	// @route GET /account/reset_password/check/{code}
	ResetPasswordCheck(*gin.Context, *request.CodeOption) (*ginrpc.Empty, error)

	// ResetPassword 重置密码
	// @route PUT /account/reset_password/{code}
	ResetPassword(*gin.Context, *request.ResetPasswordOption) (*base.MessageTemplate, error)
}

// UserApi 用户需要登录
type UserApi interface {
	// GetList 用户列表
	// @route GET /users
	GetList(*gin.Context, *request.UserListOption) (*response.UserList, error)

	// ChangePasswordByAdmin 管理员修改用户密码
	// @route PUT /users/{userID}
	ChangePasswordByAdmin(*gin.Context, *request.ChangePasswordOption) (*ginrpc.Empty, error)

	// DelUser 删除用户
	// @route DELETE /users/{userID}
	DelUser(*gin.Context, *request.UserIDOption) (*ginrpc.Empty, error)

	// GetSelf 当前登录的用户
	// @route GET /user
	GetSelf(*gin.Context, *ginrpc.Empty) (*response.User, error)

	// SetSelf 设置当前用户自身的信息
	// @route PUT /user
	SetSelf(*gin.Context, *request.SetUserSelfOption) (*ginrpc.Empty, error)

	// ChangePassword 修改密码
	// @route PUT /user/password
	ChangePassword(*gin.Context, *request.ChangePwdOption) (*ginrpc.Empty, error)

	// SendChangeEmail 发送修改邮箱邮件
	// @route POST /user/email
	SendChangeEmail(*gin.Context, *base.EmailOption) (*ginrpc.Empty, error)

	// ChangeEmailFire 修改邮箱
	// @route PUT /user/email/{code}
	ChangeEmailFire(*gin.Context, *request.CodeOption) (*base.MessageTemplate, error)

	// UploadAvatar 上传头像
	// @route POST /user/avatar
	UploadAvatar(*gin.Context, *request.UploadAvatarOption) (*base.AvatarOption, error)

	// OauthConnect 绑定oauth
	// @route POST /user/oauth/{type}/connect
	OauthConnect(*gin.Context, *request.OauthOption) (*ginrpc.Empty, error)

	// OauthDisConnect 解绑oauth
	// @route DELETE /user/oauth/{type}/disconnect
	OauthDisconnect(*gin.Context, *base.OauthTypeOption) (*ginrpc.Empty, error)
}
