package collection

import (
	"crypto/md5"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/project"
	"github.com/apicat/apicat/v2/backend/model/share"
	"github.com/apicat/apicat/v2/backend/module/cache"
	"github.com/apicat/apicat/v2/backend/route/middleware/access"
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	protocollection "github.com/apicat/apicat/v2/backend/route/proto/collection"
	collectionbase "github.com/apicat/apicat/v2/backend/route/proto/collection/base"
	collectionrequest "github.com/apicat/apicat/v2/backend/route/proto/collection/request"
	collectionresponse "github.com/apicat/apicat/v2/backend/route/proto/collection/response"
	"github.com/apicat/apicat/v2/backend/utils/password"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type collectionShareApiImpl struct{}

func NewCollectionShareApi() protocollection.CollectionShareApi {
	return &collectionShareApiImpl{}
}

func (csai *collectionShareApiImpl) Status(ctx *gin.Context, opt *collectionbase.CollectionPublicIDOption) (*collectionbase.ProjectCollectionIDOption, error) {
	c := &collection.Collection{PublicID: opt.CollectionPublicID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.FailedToGetStatus"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	p := &project.Project{ID: c.ProjectID}
	exist, err = p.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "p.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.FailedToGetStatus"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	if p.Visibility == project.VisibilityPrivate && c.ShareKey == "" {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	return &collectionbase.ProjectCollectionIDOption{
		ProjectIdOption: protobase.ProjectIdOption{
			ProjectID: p.ID,
		},
		CollectionIDOption: collectionbase.CollectionIDOption{
			CollectionID: c.ID,
		},
	}, nil
}

func (csai *collectionShareApiImpl) Detail(ctx *gin.Context, opt *collectionbase.ProjectCollectionIDOption) (*collectionresponse.CollectionShareDetail, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfP.Visibility == project.VisibilityPrivate && selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfP.ID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.FailedToGetStatus"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	return &collectionresponse.CollectionShareDetail{
		ProjectVisibilityOption: protobase.ProjectVisibilityOption{
			Visibility: selfP.Visibility,
		},
		CollectionShareData: collectionresponse.CollectionShareData{
			CollectionPublicIDOption: collectionbase.CollectionPublicIDOption{
				CollectionPublicID: c.PublicID,
			},
			SecretKeyOption: protobase.SecretKeyOption{
				SecretKey: c.ShareKey,
			},
		},
	}, nil
}

func (csai *collectionShareApiImpl) Switch(ctx *gin.Context, opt *collectionrequest.SwitchCollectionShareOption) (*collectionresponse.CollectionShareData, error) {
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfPM.ProjectID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharingFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	if opt.Status {
		if c.ShareKey == "" {
			c.ShareKey = password.RandomPassword(4)
			if err := c.UpdateShareKey(ctx); err != nil {
				slog.ErrorContext(ctx, "c.UpdateShareKey", "err", err)
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharingFailed"))
			}
		}
	} else {
		if c.ShareKey != "" {
			if err := share.DeleteCollectionShareTmpTokens(ctx, c.ID); err != nil {
				slog.ErrorContext(ctx, "share.DeleteCollectionShareTmpTokens", "err", err)
				return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.FailedToDisable"))
			}
		}

		c.ShareKey = ""
		if err := c.UpdateShareKey(ctx); err != nil {
			slog.ErrorContext(ctx, "c.UpdateShareKey", "err", err)
			return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.FailedToDisable"))
		}
	}

	return &collectionresponse.CollectionShareData{
		CollectionPublicIDOption: collectionbase.CollectionPublicIDOption{
			CollectionPublicID: c.PublicID,
		},
		SecretKeyOption: protobase.SecretKeyOption{
			SecretKey: c.ShareKey,
		},
	}, nil
}

func (csai *collectionShareApiImpl) Reset(ctx *gin.Context, opt *collectionbase.ProjectCollectionIDOption) (*protobase.SecretKeyOption, error) {
	selfP := access.GetSelfProject(ctx)
	selfPM := access.GetSelfProjectMember(ctx)
	if selfPM.Permission.Lower(project.ProjectMemberWrite) {
		return nil, ginrpc.NewError(http.StatusForbidden, i18n.NewErr("common.PermissionDenied"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: selfP.ID}
	exist, err := c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyResetFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	if selfP.Visibility != project.VisibilityPrivate {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.PublicProjectShare"))
	}
	if c.ShareKey == "" {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.SharedKeyResetFailed"))
	}

	if err := share.DeleteCollectionShareTmpTokens(ctx, c.ID); err != nil {
		slog.ErrorContext(ctx, "share.DeleteCollectionShareTmpTokens", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyResetFailed"))
	}

	c.ShareKey = password.RandomPassword(4)
	if err := c.UpdateShareKey(ctx); err != nil {
		slog.ErrorContext(ctx, "c.UpdateShareKey", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyResetFailed"))
	}

	return &protobase.SecretKeyOption{
		SecretKey: c.ShareKey,
	}, nil
}

func (csai *collectionShareApiImpl) Check(ctx *gin.Context, opt *collectionrequest.CheckCollectionShareSecretKeyOpt) (*collectionbase.ShareCode, error) {
	pcache, err := cache.NewCache(config.Get().Cache.ToModuleStruct())
	if err != nil {
		slog.ErrorContext(ctx, "cache.NewCache", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyVerificationFailed"))
	}

	// 按照ip最大重试次数
	checkDocShareCodeTimesKey := fmt.Sprintf("checkDocShareCode-%s", ctx.ClientIP())
	ts, ok, _ := pcache.Get(checkDocShareCodeTimesKey)
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

	_ = pcache.Set(checkDocShareCodeTimesKey, strconv.Itoa(number+1), time.Hour)

	p := &project.Project{ID: opt.ProjectID}
	exist, err := p.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "p.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyVerificationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	if p.Visibility != project.VisibilityPrivate {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.PublicProjectShare"))
	}

	c := &collection.Collection{ID: opt.CollectionID, ProjectID: p.ID}
	exist, err = c.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "c.Get", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyVerificationFailed"))
	}
	if !exist {
		return nil, ginrpc.NewError(http.StatusNotFound, i18n.NewErr("collection.DoesNotExist"))
	}

	if c.ShareKey == "" {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.SharedKeyError"))
	}
	if c.ShareKey != opt.SecretKey {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("share.SharedKeyError"))
	}

	shareCode := "d" + fmt.Sprintf("%x", md5.Sum([]byte(c.ShareKey+fmt.Sprint(time.Now().UnixNano()))))
	stt := &share.ShareTmpToken{
		ShareToken:   fmt.Sprintf("%x", md5.Sum([]byte(shareCode))),
		Expiration:   time.Now().Add(time.Hour * 24),
		ProjectID:    p.ID,
		CollectionID: c.ID,
	}
	if err := stt.Create(ctx); err != nil {
		slog.ErrorContext(ctx, "stt.Create", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("share.SharedKeyVerificationFailed"))
	}

	_ = pcache.Del(checkDocShareCodeTimesKey)

	return &collectionbase.ShareCode{
		ShareCode:  shareCode,
		Expiration: stt.Expiration.Unix(),
	}, nil
}
