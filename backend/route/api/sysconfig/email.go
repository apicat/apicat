package sysconfig

import (
	"encoding/json"
	"log/slog"
	"net/mail"
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

func NewEmailApi() protosysconfig.EmailApi {
	return &emailApiImpl{}
}

func (e *emailApiImpl) Get(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigbase.ConfigList, error) {
	list, err := sysconfig.GetList(ctx, "email")
	if err != nil {
		slog.ErrorContext(ctx, "sysconfig.GetList", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.FailedToGetEmailList"))
	}
	slist := make(sysconfigbase.ConfigList, 0, len(list))
	for _, v := range list {
		if v.Driver == mailmodule.SMTP {
			d := &sysconfigbase.ConfigDetail{
				Driver: v.Driver,
				Use:    v.BeingUsed,
				Config: map[string]interface{}{
					"host":     "",
					"user":     "",
					"address":  "",
					"password": "",
				},
			}

			var smtpCfg config.EmailSmtp
			if err := json.Unmarshal([]byte(v.Config), &smtpCfg); err == nil {
				d.Config["host"] = smtpCfg.Host
				d.Config["user"] = smtpCfg.From.Name
				d.Config["address"] = smtpCfg.From.Address
				d.Config["password"] = smtpCfg.Password
			}
			slist = append(slist, d)
		} else {
			slist = append(slist, &sysconfigbase.ConfigDetail{
				Driver: v.Driver,
				Use:    v.BeingUsed,
				Config: cfgFormat(v),
			})
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
			From:     mail.Address{Name: opt.User, Address: opt.Address},
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
			ApiUser:  opt.ApiUser,
			ApiKey:   opt.ApiKey,
			From:     opt.FromEmail,
			FromName: opt.FromName,
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
