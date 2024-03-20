package user

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/user"
	protouserrequest "github.com/apicat/apicat/v2/backend/route/proto/user/request"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

func EmailRegister(ctx *gin.Context, opt *protouserrequest.RegisterUserOption) (*user.User, error) {
	usr := &user.User{
		Email: opt.Email,
	}
	exist, err := usr.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "EmailRegister.usr.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
	}

	if exist {
		if usr.IsActive {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.EmailHasRegistered"))
		}

		// 邮箱已被注册但未激活，将邮箱对应账号信息改为本次注册信息
		usr.Name = opt.Name
		usr.Language = opt.Language
		if err := usr.Update(ctx); err != nil {
			slog.ErrorContext(ctx, "EmailRegister.usr.Update", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
		}
		usr.Password = opt.Password
		if err := usr.UpdatePassword(ctx); err != nil {
			slog.ErrorContext(ctx, "EmailRegister.usr.UpdatePassword", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
		}
	} else {
		usr = &user.User{
			Name:        opt.Name,
			Password:    opt.Password,
			Email:       opt.Email,
			Language:    opt.Language,
			Role:        user.RoleUser,
			LastLoginAt: time.Now(),
			IsActive:    true,
		}
		if err := usr.Create(ctx); err != nil {
			slog.ErrorContext(ctx, "EmailRegister.usr.Create", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
		}
		if usr.ID == 1 {
			usr.SetSysAdmin(ctx)
		}
	}

	return usr, nil
}

func OauthRegister(ctx *gin.Context, opt *protouserrequest.RegisterUserOption) (*user.User, error) {
	usr, err := user.GetUserByOauth(ctx, opt.Bind.OauthUserID, opt.Bind.Type)
	if err != nil {
		slog.ErrorContext(ctx, "OauthRegister.user.GetUserByOauth", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
	}
	if usr != nil {
		if usr.IsActive {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.OauthConnectRepeat", opt.Bind.Type))
		}

		// oauth账号已被注册但未激活，将oauth账号对应账号信息改为本次注册信息
		usr.Name = opt.Name
		usr.Avatar = opt.Avatar
		usr.Language = opt.Language
		if err := usr.Update(ctx); err != nil {
			slog.ErrorContext(ctx, "OauthRegister.usr.Update", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
		}
		usr.Email = opt.Email
		if err := usr.UpdateEmail(ctx); err != nil {
			slog.ErrorContext(ctx, "OauthRegister.usr.UpdateEmail", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
		}
		usr.Password = opt.Password
		if err := usr.UpdatePassword(ctx); err != nil {
			slog.ErrorContext(ctx, "OauthRegister.usr.UpdatePassword", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
		}
	} else {
		usr = &user.User{
			Email: opt.Email,
		}
		exist, _ := usr.Get(ctx)
		if exist && usr.IsActive {
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.EmailHasRegistered"))
		}

		usr = &user.User{
			Name:        opt.Name,
			Password:    opt.Password,
			Email:       opt.Email,
			Avatar:      opt.Avatar,
			Language:    opt.Language,
			Role:        user.RoleUser,
			LastLoginAt: time.Now(),
			IsActive:    true,
		}
		if err := usr.CreateAndBindOauth(ctx, opt.Bind.Type, opt.Bind.OauthUserID); err != nil {
			slog.ErrorContext(ctx, "OauthRegister.usr.CreateAndBindOauth", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
		}
		if usr.ID == 1 {
			usr.SetSysAdmin(ctx)
		}
	}

	return usr, nil
}
