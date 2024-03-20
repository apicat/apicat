package sysconfig

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/apicat/apicat/v2/backend/config"
	"github.com/apicat/apicat/v2/backend/i18n"
	"github.com/apicat/apicat/v2/backend/model/sysconfig"
	"github.com/apicat/apicat/v2/backend/module/storage"
	"github.com/apicat/apicat/v2/backend/module/storage/cloudflare"
	"github.com/apicat/apicat/v2/backend/module/storage/local"
	"github.com/apicat/apicat/v2/backend/module/storage/qiniu"
	protosysconfig "github.com/apicat/apicat/v2/backend/route/proto/sysconfig"
	sysconfigbase "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/base"
	sysconfigrequest "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/request"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type storageApiImpl struct{}

func NewStorageApi() protosysconfig.StorageApi {
	return &storageApiImpl{}
}

func (s *storageApiImpl) Get(ctx *gin.Context, _ *ginrpc.Empty) (*sysconfigbase.ConfigList, error) {
	list, err := sysconfig.GetList(ctx, "storage")
	if err != nil {
		slog.ErrorContext(ctx, "sysconfig.GetList", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.FailedToGetStorageList"))
	}
	slist := make(sysconfigbase.ConfigList, 0, len(list))
	for _, v := range list {
		slist = append(slist, &sysconfigbase.ConfigDetail{
			Driver: v.Driver,
			Use:    v.BeingUsed,
			Config: cfgFormat(v),
		})
	}
	return &slist, nil
}

func (s *storageApiImpl) UpdateDisk(ctx *gin.Context, opt *sysconfigrequest.DiskOption) (*ginrpc.Empty, error) {
	storageConfig := &config.Storage{
		Driver: "localdisk",
		LocalDisk: &config.LocalDisk{
			Path: opt.Path,
			Url:  config.Get().App.AppUrl + "/uploads",
		},
	}

	diskc := storageConfig.ToMapInterface()
	if disk, err := local.NewDisk(diskc["LocalDisk"].(map[string]interface{})); err != nil {
		slog.ErrorContext(ctx, "local.NewDisk", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.LocalPathInvalid"))
	} else {
		if err := disk.Check(); err != nil {
			slog.ErrorContext(ctx, "disk.Check", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.LocalPathInvalid"))
		}
	}

	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.StorageUpdateFailed"))
	}

	storage := &sysconfig.Sysconfig{
		Type:      "storage",
		Driver:    "disk",
		BeingUsed: true,
		Config:    string(jsonData),
	}

	if err := sysconfig.UpdateOrCreate(ctx, storage); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.StorageUpdateFailed"))
	}
	config.SetStorage(storageConfig)
	return nil, nil
}

func (s *storageApiImpl) UpdateCloudflare(ctx *gin.Context, opt *sysconfigrequest.CloudflareOption) (*ginrpc.Empty, error) {
	storageConfig := &config.Storage{
		Driver: storage.CLOUDFLARE,
		Cloudflare: &config.Cloudflare{
			AccountID:       opt.AccountID,
			AccessKeyID:     opt.AccessKeyID,
			AccessKeySecret: opt.AccessKeySecret,
			BucketName:      opt.BucketName,
			Url:             opt.BucketUrl,
		},
	}

	r2c := storageConfig.ToMapInterface()
	if r2, err := cloudflare.NewR2(r2c["Cloudflare"].(map[string]interface{})); err != nil {
		slog.ErrorContext(ctx, "cloudflare.NewR2", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.CloudflareConfigInvalid"))
	} else {
		if err := r2.Check(); err != nil {
			slog.ErrorContext(ctx, "r2.Check", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.CloudflareConfigInvalid"))
		}
	}

	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.StorageUpdateFailed"))
	}

	storage := &sysconfig.Sysconfig{
		Type:      "storage",
		Driver:    storage.CLOUDFLARE,
		BeingUsed: true,
		Config:    string(jsonData),
	}

	if err := sysconfig.UpdateOrCreate(ctx, storage); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.StorageUpdateFailed"))
	}
	config.SetStorage(storageConfig)
	return nil, nil
}

func (s *storageApiImpl) UpdateQiniu(ctx *gin.Context, opt *sysconfigrequest.QiniuOption) (*ginrpc.Empty, error) {
	storageConfig := &config.Storage{
		Driver: storage.QINIU,
		Qiniu: &config.Qiniu{
			AccessKeyID:     opt.AccessKey,
			AccessKeySecret: opt.SecretKey,
			BucketName:      opt.BucketName,
			Url:             opt.BucketUrl,
		},
	}

	qc := storageConfig.ToMapInterface()
	if q, err := qiniu.NewQiniu(qc["Qiniu"].(map[string]interface{})); err != nil {
		slog.ErrorContext(ctx, "qiniu.NewQiniu", "err", err)
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.QiniuConfigInvalid"))
	} else {
		if err := q.Check(); err != nil {
			slog.ErrorContext(ctx, "q.Check", "err", err)
			return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("sysConfig.QiniuConfigInvalid"))
		}
	}

	jsonData, err := json.Marshal(opt)
	if err != nil {
		slog.ErrorContext(ctx, "json.Marshal", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.StorageUpdateFailed"))
	}

	storage := &sysconfig.Sysconfig{
		Type:      "storage",
		Driver:    storage.QINIU,
		BeingUsed: true,
		Config:    string(jsonData),
	}

	if err := sysconfig.UpdateOrCreate(ctx, storage); err != nil {
		slog.ErrorContext(ctx, "sysconfig.UpdateOrCreate", "err", err)
		return nil, ginrpc.NewError(http.StatusInternalServerError, i18n.NewErr("sysConfig.StorageUpdateFailed"))
	}
	config.SetStorage(storageConfig)
	return nil, nil
}
