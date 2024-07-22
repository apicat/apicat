package sysconfig

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/sysconfig"
	mailmodule "github.com/apicat/apicat/v2/backend/module/mail"
	protosysconfig "github.com/apicat/apicat/v2/backend/route/proto/sysconfig"
	sysconfigbase "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/base"
	sysconfigrequest "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/request"

	"net/http"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type emailApiImpl struct{}

var emailDrivers = []string{mailmodule.SMTP, mailmodule.SENDCLOUD}

func NewEmailApi() protosysconfig.EmailApi {
	return &emailApiImpl{}
}

func (e *emailApiImpl) Get(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigbase.ConfigList, error) {
	configMap := make(map[string]*sysconfigbase.ConfigDetail)

	emailCfg := config.GetEmail()
	switch emailCfg.Driver {
	case mailmodule.SMTP:
		if js, err := json.Marshal(emailCfg.Smtp); err == nil {
			configMap[mailmodule.SMTP] = &sysconfigbase.ConfigDetail{
				Driver: mailmodule.SMTP,
				Use:    true,
				Config: cfgFormat(&sysconfig.Sysconfig{Config: string(js)}),
			}
		}
	case mailmodule.SENDCLOUD:
		if js, err := json.Marshal(emailCfg.SendCloud); err == nil {
			configMap[mailmodule.SENDCLOUD] = &sysconfigbase.ConfigDetail{
				Driver: mailmodule.SENDCLOUD,
				Use:    true,
				Config: cfgFormat(&sysconfig.Sysconfig{Config: string(js)}),
			}
		}
	}

	if list, err := sysconfig.GetList(ctx, "email"); err == nil {
		for _, v := range list {
			if _, ok := configMap[v.Driver]; !ok {
				configMap[v.Driver] = &sysconfigbase.ConfigDetail{
					Driver: v.Driver,
					Use:    v.BeingUsed,
					Config: cfgFormat(v),
				}
			}
		}
	}

	slist := make(sysconfigbase.ConfigList, 0)
	for _, v := range emailDrivers {
		if _, ok := configMap[v]; ok {
			slist = append(slist, configMap[v])
		}
	}
	return &slist, nil
}

func (e *emailApiImpl) UpdateSMTP(ctx *gin.Context, opt *sysconfigrequest.SMTPOption) (*ginrpc.Empty, error) {
	if i := strings.Index(opt.Host, ":"); i < 7 {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.SMTPConfigInvalid"))
	}

	emailConfig := &config.Email{
		Driver: mailmodule.SMTP,
		Smtp: &config.EmailSmtp{
			Host:     opt.Host,
			Address:  opt.Address,
			Name:     opt.User,
			Password: opt.Password,
		},
	}

	jsonData, err := json.Marshal(emailConfig.Smtp)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.EmailUpdateFailed"))
	}

	email := &sysconfig.Sysconfig{
		Type:      "email",
		Driver:    mailmodule.SMTP,
		BeingUsed: true,
		Config:    string(jsonData),
	}

	if err := sysconfig.UpdateOrCreate(ctx, email); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.EmailUpdateFailed"))
	}
	config.SetEmail(emailConfig)
	return nil, nil
}

func (e *emailApiImpl) UpdateSendCloud(ctx *gin.Context, opt *sysconfigrequest.SendCloudOption) (*ginrpc.Empty, error) {
	emailConfig := &config.Email{
		Driver: mailmodule.SENDCLOUD,
		SendCloud: &config.EmailSendCloud{
			ApiUser:     opt.ApiUser,
			ApiKey:      opt.ApiKey,
			FromAddress: opt.FromEmail,
			FromName:    opt.FromName,
		},
	}

	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.EmailUpdateFailed"))
	}

	email := &sysconfig.Sysconfig{
		Type:      "email",
		Driver:    mailmodule.SENDCLOUD,
		BeingUsed: true,
		Config:    string(jsonData),
	}

	if err := sysconfig.UpdateOrCreate(ctx, email); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.EmailUpdateFailed"))
	}
	config.SetEmail(emailConfig)
	return nil, nil
}
