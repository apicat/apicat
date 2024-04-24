package project

import (
	"crypto/md5"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/model/share"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	protoproject "github.com/apicat/apicat/v2/backend/route/proto/project"
	projectbase "github.com/apicat/apicat/v2/backend/route/proto/project/base"
	projectrequest "github.com/apicat/apicat/v2/backend/route/proto/project/request"
	projectresponse "github.com/apicat/apicat/v2/backend/route/proto/project/response"
	"github.com/apicat/apicat/v2/backend/utils/password"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type projectShareApiImpl struct{}

func NewProjectShareApi() protoproject.ProjectShareApi {
	return &projectShareApiImpl{}
}

// Status 获取项目分享状态
func (psai *projectShareApiImpl) Status(ctx *gin.Context, opt *protobase.ProjectIdOption) (*projectresponse.ProjectShareStatus, error) {
	p := access.GetSelfProject(ctx)
	if p == nil {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("project.DoesNotExist"))
	}

	permission := project.ProjectMemberNone
	pm := access.GetSelfProjectMember(ctx)
	if pm != nil {
		permission = pm.Permission
	}

	return &projectresponse.ProjectShareStatus{
		ProjectMemberPermission: protobase.ProjectMemberPermission{
			Permission: permission,
		},
		ProjectVisibilityOption: protobase.ProjectVisibilityOption{
			Visibility: p.Visibility,
		},
		HasShare: p.ShareKey != "",
	}, nil
}

// Detail 项目分享详情
func (psai *projectShareApiImpl) Detail(ctx *gin.Context, opt *protobase.ProjectIdOption) (*projectresponse.ProjectShareDetail, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfP.Visibility == project.VisibilityPrivate && selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	return &projectresponse.ProjectShareDetail{
		ProjectMemberPermission: protobase.ProjectMemberPermission{
			Permission: selfPM.Permission,
		},
		ProjectVisibilityOption: protobase.ProjectVisibilityOption{
			Visibility: selfP.Visibility,
		},
		SecretKeyOption: protobase.SecretKeyOption{
			SecretKey: selfP.ShareKey,
		},
	}, nil
}

// Switch 切换项目分享状态
func (psai *projectShareApiImpl) Switch(ctx *gin.Context, opt *projectrequest.ProjectShareSwitchOption) (*protobase.SecretKeyOption, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if selfP.Visibility != project.VisibilityPrivate {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.PublicProjectShare"))
	}

	if opt.Status {
		if selfP.ShareKey == "" {
			selfP.ShareKey = password.RandomPassword(4)
			if err := selfP.UpdateShareKey(ctx); err != nil {
				slog.ErrorContext(ctx, "selfP.UpdateShareKey", "err", err)
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharingFailed"))
			}
		}
	} else {
		if selfP.ShareKey != "" {
			if err := share.DeleteProjectShareTmpTokens(ctx, selfP); err != nil {
				slog.ErrorContext(ctx, "share.DeleteProjectShareTmpTokens", "err", err)
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.FailedToDisable"))
			}

			selfP.ShareKey = ""
			if err := selfP.UpdateShareKey(ctx); err != nil {
				slog.ErrorContext(ctx, "selfP.UpdateShareKey", "err", err)
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.FailedToDisable"))
			}
		}
	}

	return &protobase.SecretKeyOption{
		SecretKey: selfP.ShareKey,
	}, nil
}

// Reset 重置项目分享密钥
func (psai *projectShareApiImpl) Reset(ctx *gin.Context, opt *protobase.ProjectIdOption) (*protobase.SecretKeyOption, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	if selfP.Visibility != project.VisibilityPrivate {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.PublicProjectShare"))
	}
	if selfP.ShareKey == "" {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.SharedKeyResetFailed"))
	}

	if err := share.DeleteProjectShareTmpTokens(ctx, selfP); err != nil {
		slog.ErrorContext(ctx, "share.DeleteProjectShareTmpTokens", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyResetFailed"))
	}

	selfP.ShareKey = password.RandomPassword(4)
	if err := selfP.UpdateShareKey(ctx); err != nil {
		slog.ErrorContext(ctx, "selfP.UpdateShareKey", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyResetFailed"))
	}

	return &protobase.SecretKeyOption{
		SecretKey: selfP.ShareKey,
	}, nil
}

// Check 检查分享密钥
func (psai *projectShareApiImpl) Check(ctx *gin.Context, opt *projectrequest.CheckProjectShareSecretKeyOption) (*projectbase.ShareCode, error) {
	pcache, err := cache.NewCache(config.Get().Cache.ToCfg())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyVerificationFailed"))
	}
	// 按照ip最大重试次数
	checkProjectShareCodeTimesKey := fmt.Sprintf("checkProjectShareCode-%s", ctx.ClientIP())
	ts, ok, _ := pcache.Get(checkProjectShareCodeTimesKey)
	var number int
	if ok {
		var err error
		number, err = strconv.Atoi(ts)
		if err != nil {
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyVerificationFailed"))
		}
		if number > 10 {
			return nil, ginrpc.NewError(http.StatusTooManyRequests, i18n.NewErr("common.TooManyOperations"))
		}
	}

	_ = pcache.Set(checkProjectShareCodeTimesKey, strconv.Itoa(number+1), time.Hour)

	p := &project.Project{ID: opt.ProjectID}
	exist, err := p.Get(ctx)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyVerificationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("project.DoesNotExist"))
	}

	t := &team.Team{ID: p.TeamID}
	exist, err = t.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "t.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyVerificationFailed"))
	}
	if !exist {
		// 团队都不存在了，项目也就不存在了
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("project.DoesNotExist"))
	}

	if p.Visibility != project.VisibilityPrivate {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.PublicProjectShare"))
	}
	if p.ShareKey == "" {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.SharedKeyError"))
	}
	if p.ShareKey != opt.SecretKey {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.SharedKeyError"))
	}

	shareCode := "p" + fmt.Sprintf("%x", md5.Sum([]byte(p.ShareKey+fmt.Sprint(time.Now().UnixNano()))))

	stt := &share.ShareTmpToken{
		ShareToken: fmt.Sprintf("%x", md5.Sum([]byte(shareCode))),
		Expiration: time.Now().Add(time.Hour * 24),
		ProjectID:  p.ID,
	}
	if err := stt.Create(ctx); err != nil {
		slog.ErrorContext(ctx, "stt.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyVerificationFailed"))
	}

	_ = pcache.Del(checkProjectShareCodeTimesKey)

	return &projectbase.ShareCode{
		ShareCode:  shareCode,
		Expiration: stt.Expiration.Unix(),
	}, nil
}
