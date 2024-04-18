package mailer

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/model/user"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/module/onetime_token"

	"github.com/gin-gonic/gin"
)

type UserToken struct {
	Email  string
	UserID uint
}

func scheme(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme
}

// SendActiveAccountMail 发送激活账户邮件
func SendActiveAccountMail(ctx *gin.Context, usr *user.User) {
	tokenData := UserToken{
		Email:  usr.Email,
		UserID: usr.ID,
	}

	tokenKey := fmt.Sprintf(
		"SendActiveAccountMail-%d-%s-%d",
		usr.ID,
		usr.Email,
		time.Now().Unix(),
	)

	c, err := cache.NewCache(config.Get().Cache.ToCfg())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return
	}

	token, err := onetime_token.NewTokenHelper(c).GenerateToken(tokenKey, tokenData, time.Hour*2)
	if err != nil {
		slog.ErrorContext(ctx, "onetime_token.GenerateToken", "err", err)
		return
	}

	content := createContent("email_verification.tmpl", contentData{
		Link: fmt.Sprintf(
			"%s/email_verification/%s",
			config.Get().App.AppUrl,
			token,
		),
		Data: usr,
	})
	AsyncSend("Welcome to ApiCat! Please Verify Your Email", content, usr.Email)
}

// SendResetPasswordMail 发送重置密码邮件
func SendResetPasswordMail(ctx *gin.Context, usr *user.User) {
	tokenData := UserToken{
		Email:  usr.Email,
		UserID: usr.ID,
	}

	tokenKey := fmt.Sprintf(
		"SendResetPasswordMail-%d-%s-%d",
		usr.ID,
		usr.Email,
		time.Now().Unix(),
	)

	c, err := cache.NewCache(config.Get().Cache.ToCfg())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return
	}

	token, err := onetime_token.NewTokenHelper(c).GenerateToken(tokenKey, tokenData, time.Hour*2)
	if err != nil {
		slog.ErrorContext(ctx, "onetime_token.GenerateToken", "err", err)
		return
	}

	content := createContent("reset_password.tmpl", contentData{
		Link: fmt.Sprintf(
			"%s/reset_password/%s",
			config.Get().App.AppUrl,
			token,
		),
		Data: usr,
	})
	AsyncSend("ApiCat Password Reset Request", content, usr.Email)
}

// SendModifyEmailMail 发送修改邮箱邮件
func SendModifyEmailMail(ctx *gin.Context, usr *user.User, newEmail string) {
	tokenData := UserToken{
		Email:  newEmail,
		UserID: usr.ID,
	}

	tokenKey := fmt.Sprintf(
		"SendModifyEmailMail-%d-%s-%d",
		usr.ID,
		newEmail,
		time.Now().Unix(),
	)

	c, err := cache.NewCache(config.Get().Cache.ToCfg())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return
	}

	token, err := onetime_token.NewTokenHelper(c).GenerateToken(tokenKey, tokenData, time.Hour*2)
	if err != nil {
		slog.ErrorContext(ctx, "onetime_token.GenerateToken", "err", err)
		return
	}

	content := createContent("change_email.tmpl", contentData{
		Link: fmt.Sprintf(
			"%s/change_email/%s",
			config.Get().App.AppUrl,
			token,
		),
		Data: gin.H{
			"Email": newEmail,
			"usr":   usr,
		},
	})
	AsyncSend("Verify Your New Email Address for ApiCat", content, newEmail)
}

// SendTeamInviteMail 发送团队邀请邮件
func SendTeamInviteMail() {}
