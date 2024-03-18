package user

import (
	"apicat-cloud/backend/config"
	"apicat-cloud/backend/i18n"
	"apicat-cloud/backend/model/user"
	"apicat-cloud/backend/module/cache"
	"apicat-cloud/backend/module/oauth2"
	"apicat-cloud/backend/module/oauth2/github"
	"apicat-cloud/backend/module/onetime_token"
	"apicat-cloud/backend/module/password"
	"apicat-cloud/backend/route/middleware/jwt"
	protouser "apicat-cloud/backend/route/proto/user"
	protouserbase "apicat-cloud/backend/route/proto/user/base"
	protouserrequest "apicat-cloud/backend/route/proto/user/request"
	protouserresponse "apicat-cloud/backend/route/proto/user/response"
	"apicat-cloud/backend/service/mailer"
	"apicat-cloud/backend/service/team_relations"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type accountApiImpl struct {
	oauths map[string]*oauth2.Object
}

func NewAccountApi() protouser.AccountApi {
	objs := make(map[string]*oauth2.Object)
	if oauthcfg := config.Get().Oauth2; oauthcfg != nil {
		for k, cfg := range oauthcfg {
			var dr oauth2.Driver
			switch k {
			case "github":
				dr = &github.Github{}
			// case "google":
			default:
				continue
			}
			objs[k] = oauth2.NewObject(cfg, dr)
		}
	}
	return &accountApiImpl{
		oauths: objs,
	}
}

// buildToken 生成token
func (s *accountApiImpl) buildToken(ctx *gin.Context, usr *user.User) protouserbase.TokenResponse {
	_ = usr.UpdateLastLogin(ctx, ctx.ClientIP())
	token, _ := jwt.Generate(usr.ID)
	return protouserbase.TokenResponse{
		AccessToken: token,
	}
}

// Login 登录
func (s *accountApiImpl) Login(ctx *gin.Context, opt *protouserrequest.LoginOption) (*protouserbase.TokenResponse, error) {
	// 按照ip和email组合最大重试次数
	var number int

	ucache, err := cache.NewCache(config.Get().Cache.ToMapInterface())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.LoginFailed"))
	}

	loginTimeKey := fmt.Sprintf("login-%s/%s", opt.Email, ctx.ClientIP())
	ts, ok, _ := ucache.Get(loginTimeKey)
	if ok {
		var err error
		number, err = strconv.Atoi(ts)
		if err != nil {
			slog.ErrorContext(ctx, "strconv.Atoi", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.LoginFailed"))
		}
		if number > 10 {
			return nil, ginrpc.NewError(http.StatusTooManyRequests, i18n.NewErr("common.TooManyOperations"))
		}
	}

	_ = ucache.Set(loginTimeKey, strconv.Itoa(number+1), time.Hour)

	usr := &user.User{Email: opt.Email}
	exist, err := usr.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "usr.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.LoginFailed"))
	}
	if !exist || !usr.CheckPassword(opt.Password) {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("user.IncorrectEmailOrPassword"))
	}

	_ = ucache.Del(loginTimeKey)

	if !usr.IsActive {
		// 还未激活
		return nil, &ginrpc.Error{
			Code: http.StatusBadRequest,
			Err:  i18n.NewErr("user.IncorrectEmailOrPassword"),
			Attrs: map[string]any{
				"errcode": 1,
			},
		}
	}

	// 如果有邀请码则加入团队
	if opt.InvitationToken != "" {
		if err = team_relations.JoinTeam(ctx, opt.InvitationToken, usr); err != nil {
			slog.ErrorContext(ctx, "team_relations.JoinTeam", "err", err)
		}
	}

	token := s.buildToken(ctx, usr)
	return &token, nil
}

// Register 注册
func (s *accountApiImpl) Register(ctx *gin.Context, opt *protouserrequest.RegisterUserOption) (*protouserbase.TokenResponse, error) {
	ucache, err := cache.NewCache(config.Get().Cache.ToMapInterface())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
	}

	// 按照ip最大重试次数
	registerTimeKey := fmt.Sprintf("register-%s", ctx.ClientIP())
	ts, ok, _ := ucache.Get(registerTimeKey)
	var number int
	if ok {
		var err error
		number, err = strconv.Atoi(ts)
		if err != nil {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.RegisterFailed"))
		}
		if number > 10 {
			return nil, ginrpc.NewError(
				http.StatusTooManyRequests,
				i18n.NewErr("common.TooManyOperations"),
			)
		}
	}
	_ = ucache.Set(registerTimeKey, strconv.Itoa(number+1), time.Hour)

	if _, exist := user.SupportedLanguages[opt.Language]; !exist {
		opt.Language = user.LanguageEnUS
	}

	var usr *user.User
	if opt.Bind != nil {
		// oauth注册
		usr, err = OauthRegister(ctx, opt)
	} else {
		// 邮箱注册
		usr, err = EmailRegister(ctx, opt)
	}
	if err != nil {
		return nil, err
	}

	_ = ucache.Del(registerTimeKey)

	// 如果有邀请码则加入团队
	if opt.InvitationToken != "" {
		if err := team_relations.JoinTeam(ctx, opt.InvitationToken, usr); err != nil {
			slog.ErrorContext(ctx, "team_relations.JoinTeam", "err", err)
		}
	}

	mailer.SendActiveAccountMail(ctx, usr)
	token := s.buildToken(ctx, usr)
	return &token, nil
}

// RegisterFire 激活账户
func (s *accountApiImpl) RegisterFire(ctx *gin.Context, opt *protouserrequest.CodeOption) (*protouserresponse.RegisterFireRes, error) {
	var (
		tData mailer.UserToken
		err   error
	)

	errResp := ginrpc.NewError(
		http.StatusBadRequest,
		i18n.NewErr("user.EmailVerificationFailed"),
	)

	c, err := cache.NewCache(config.Get().Cache.ToMapInterface())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, errResp
	}
	tokenHelper := onetime_token.NewTokenHelper(c)

	if !tokenHelper.CheckToken(opt.Code, &tData) {
		errResp.Err = i18n.NewErr("common.LinkExpired")
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("common.LinkExpiredTitle").Translate(ctx),
			"description": i18n.NewTran("user.ResendEmail").Translate(ctx),
		}
		return nil, errResp
	}

	usr := &user.User{ID: tData.UserID}
	exist, err := usr.Get(ctx)
	if err != nil {
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("user.EmailVerificationFailedTitle").Translate(ctx),
			"description": i18n.NewTran("user.ResendEmail").Translate(ctx),
		}
		return nil, errResp
	}
	if !exist || usr.Email != tData.Email {
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("user.EmailVerificationFailedTitle").Translate(ctx),
			"description": i18n.NewTran("user.ResendEmail").Translate(ctx),
		}
		return nil, errResp
	}

	if usr.IsActive {
		errResp.Err = i18n.NewErr("user.EmailHasBeenVerified")
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("user.EmailHasVerifiedTitle").Translate(ctx),
			"description": i18n.NewTran("user.EmailHasVerifiedDesc").Translate(ctx),
		}
		return nil, errResp
	}

	err = usr.SetActive(ctx)
	if err != nil {
		errResp.Attrs = map[string]any{
			"emoji":       "😳",
			"title":       i18n.NewTran("user.EmailVerificationFailedTitle").Translate(ctx),
			"description": i18n.NewTran("user.ResendEmail").Translate(ctx),
		}
		return nil, errResp
	}

	tokenHelper.DelToken(opt.Code)
	registerTimeKey := fmt.Sprintf("register-%s", ctx.ClientIP())
	_ = c.Del(registerTimeKey)

	return &protouserresponse.RegisterFireRes{
		MessageTemplate: protouserbase.MessageTemplate{
			Emoji:       "🎉",
			Title:       i18n.NewTran("user.EmailVerificationSuccessfulTitle").Translate(ctx),
			Description: i18n.NewTran("common.SuccessfulDesc").Translate(ctx),
		},
		TokenResponse: protouserbase.TokenResponse{
			AccessToken: s.buildToken(ctx, usr).AccessToken,
		},
	}, nil
}

// LoginWithOauthCode oauth2平台回调
func (s *accountApiImpl) LoginWithOauthCode(ctx *gin.Context, opt *protouserrequest.Oauth2StateOption) (*protouserresponse.Oauth2User, error) {
	o, ok := s.oauths[opt.Type]
	if !ok {
		return nil, ginrpc.NewError(
			http.StatusNotFound,
			i18n.NewErr("user.NotSupportOauth", opt.Type),
		)
	}

	// 根据code请求oauth平台获取用户信息
	oauthUser, err := o.GetUserByState(ctx, opt.Code)
	if err != nil {
		slog.ErrorContext(ctx, "o.GetUserByState", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.OauthLoginFailed"))
	}

	var usr *user.User
	defer func() {
		if usr != nil && opt.InvitationToken != "" {
			if err = team_relations.JoinTeam(ctx, opt.InvitationToken, usr); err != nil {
				slog.ErrorContext(ctx, "team_relations.JoinTeam", "err", err)
			}
		}
	}()

	// oauth已绑定过（已绑定过但解绑的账号恢复原绑定）
	usr, err = user.GetAndRecoverUserByOauth(ctx, oauthUser.ID, opt.Type)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.OauthLoginFailed"))
	}

	toAddInfo := &protouserresponse.Oauth2User{
		UserData: protouserresponse.UserData{
			NameOption:   protouserbase.NameOption{Name: oauthUser.Name},
			AvatarOption: protouserbase.AvatarOption{Avatar: oauthUser.Avatar},
		},
		Bind: &protouserbase.UserOauthBindOption{
			OauthTypeOption: protouserbase.OauthTypeOption{Type: opt.Type},
			OauthUserID:     oauthUser.ID,
		},
	}

	// oauth未绑定过
	if usr == nil {
		// 如果email为空则返回给前端，补充信息后继续调用注册接口
		if oauthUser.Email == "" {
			return toAddInfo, nil
		}

		// oauth成功获取到了邮箱
		// 判断邮箱是否已注册
		usr = &user.User{Email: oauthUser.Email}
		exist, err := usr.Get(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "usr.Get", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.OauthLoginFailed"))
		}

		// oauth邮箱已注册，直接绑定
		if exist {
			if err := usr.BindOrRecoverOauth(ctx, opt.Type, oauthUser.ID); err != nil {
				slog.ErrorContext(ctx, "usr.BindOrRecoverOauth", "err", err)
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.OauthLoginFailed"))
			}
			return &protouserresponse.Oauth2User{
				TokenResponse: s.buildToken(ctx, usr),
			}, nil
		}

		if _, exist := user.SupportedLanguages[opt.Language]; !exist {
			opt.Language = user.LanguageEnUS
		}

		// oauth邮箱未注册，自动注册并绑定，自动注册的账号密码随机生成（基本上都是使用oauth登录, 如果非要使用密码登录,可以通过忘记密码重置为新密码）
		usr = &user.User{
			Name:        oauthUser.Name,
			Email:       oauthUser.Email,
			Avatar:      oauthUser.Avatar,
			Language:    opt.Language,
			IsActive:    true,
			Password:    password.RandomPassword(8),
			LastLoginAt: time.Now(),
		}
		if err := usr.CreateAndBindOauth(ctx, opt.Type, oauthUser.ID); err != nil {
			slog.ErrorContext(ctx, "usr.CreateAndBindOauth", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.OauthLoginFailed"))
		}
	}

	if !usr.IsActive {
		return toAddInfo, nil
	}

	return &protouserresponse.Oauth2User{
		TokenResponse: s.buildToken(ctx, usr),
	}, nil
}

// SendResetPasswordMail 发送重置密码邮件
func (s *accountApiImpl) SendResetPasswordMail(ctx *gin.Context, opt *protouserbase.EmailOption) (*ginrpc.Empty, error) {
	ucache, err := cache.NewCache(config.Get().Cache.ToMapInterface())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.EmailSendFailed"))
	}

	// 按照ip最大重试次数
	resetPasswordTimeKey := fmt.Sprintf("resetPassword-%s", ctx.ClientIP())
	ts, ok, _ := ucache.Get(resetPasswordTimeKey)
	var number int
	if ok {
		var err error
		number, err = strconv.Atoi(ts)
		if err != nil {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.EmailSendFailed"))
		}
		if number > 10 {
			return nil, ginrpc.NewError(http.StatusTooManyRequests, i18n.NewErr("common.TooManyOperations"))
		}
	}

	_ = ucache.Set(resetPasswordTimeKey, strconv.Itoa(number+1), time.Hour)

	u := &user.User{Email: opt.Email}
	if exist, err := u.Get(ctx); err != nil {
		slog.ErrorContext(ctx, "u.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("common.EmailSendFailed"))
	} else if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("user.EmailDoesNotExist"))
	}
	if !u.IsActive {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("user.EmailDoesNotExist"))
	}

	mailer.SendResetPasswordMail(ctx, u)
	return &ginrpc.Empty{}, nil
}

// ResetPasswordCheck 检查重置密码令牌
func (s *accountApiImpl) ResetPasswordCheck(ctx *gin.Context, opt *protouserrequest.CodeOption) (*ginrpc.Empty, error) {
	var v mailer.UserToken

	errResp := ginrpc.NewError(
		http.StatusBadRequest,
		i18n.NewErr("common.LinkExpired"),
	)

	c, err := cache.NewCache(config.Get().Cache.ToMapInterface())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, errResp
	}

	if !onetime_token.NewTokenHelper(c).CheckToken(opt.Code, &v) {
		errResp.Attrs = map[string]any{
			"emoji":       "😭",
			"title":       i18n.NewTran("common.LinkExpiredTitle").Translate(ctx),
			"description": i18n.NewTran("user.ResendEmail").Translate(ctx),
		}
		return nil, errResp
	}

	return &ginrpc.Empty{}, nil
}

// ResetPassword 重置密码
func (s *accountApiImpl) ResetPassword(ctx *gin.Context, opt *protouserrequest.ResetPasswordOption) (*protouserbase.MessageTemplate, error) {
	if opt.Code == "" {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.PasswordResetFailed"))
	}
	// 通过code获取连接内容
	// 提取是要改的目标邮箱
	// 邮箱不能暴露否则用户可以随意修改任意邮箱的密码
	ucache, err := cache.NewCache(config.Get().Cache.ToMapInterface())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.PasswordResetFailed"))
	}

	tokenHelper := onetime_token.NewTokenHelper(ucache)
	var v mailer.UserToken
	if !tokenHelper.CheckToken(opt.Code, &v) {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.PasswordResetFailed"))
	}
	usr := &user.User{Email: v.Email}
	if exist, err := usr.Get(ctx); err != nil {
		slog.ErrorContext(ctx, "usr.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.PasswordResetFailed"))
	} else if !exist {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("user.PasswordResetFailed"))
	}
	usr.Password = opt.Password
	if err := usr.UpdatePassword(ctx); err != nil {
		slog.ErrorContext(ctx, "usr.UpdatePassword", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("user.PasswordResetFailed"))
	}

	// 重置完密码这个邮箱连接就失效了
	tokenHelper.DelToken(opt.Code)
	loginTimeKey := fmt.Sprintf("login-%s/%s", v.Email, ctx.ClientIP())
	_ = ucache.Del(loginTimeKey)
	resetPasswordTimeKey := fmt.Sprintf("resetPassword-%s", ctx.ClientIP())
	_ = ucache.Del(resetPasswordTimeKey)

	return &protouserbase.MessageTemplate{
		Emoji:       "🎉",
		Title:       i18n.NewTran("user.PasswordResetSuccessfulTitle").Translate(ctx),
		Description: i18n.NewTran("common.SuccessfulDesc").Translate(ctx),
	}, nil
}
