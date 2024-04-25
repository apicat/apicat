package storage

import (
	"errors"
	"log/slog"

	"github.com/apicat/apicat/v2/backend/module/storage/cloudflare"
	"github.com/apicat/apicat/v2/backend/module/storage/common"
	"github.com/apicat/apicat/v2/backend/module/storage/local"
	"github.com/apicat/apicat/v2/backend/module/storage/qiniu"
)

const (
	LOCAL      = "disk"
	CLOUDFLARE = "cloudflare"
	QINIU      = "qiniu"
)

type Storage struct {
	Driver     string
	LocalDisk  local.Disk
	Qiniu      qiniu.QiniuOpt
	Cloudflare cloudflare.R2Opt
}

func NewStorage(cfg Storage) common.Storage {
	slog.Debug("storage.NewStorage", "cfg", cfg)

	switch cfg.Driver {
	case CLOUDFLARE:
		return cloudflare.NewR2(cfg.Cloudflare)
	case QINIU:
		return qiniu.NewQiniu(cfg.Qiniu)
	case LOCAL:
		return local.NewDisk(cfg.LocalDisk)
	default:
		return nil
	}
}

func Init(cfg Storage) error {
	switch cfg.Driver {
	case CLOUDFLARE:
		s := cloudflare.NewR2(cfg.Cloudflare)
		return s.Check()
	case QINIU:
		s := qiniu.NewQiniu(cfg.Qiniu)
		return s.Check()
	case LOCAL:
		s := local.NewDisk(cfg.LocalDisk)
		return s.Check()
	default:
		return errors.New("storage driver not found")
	}
}
