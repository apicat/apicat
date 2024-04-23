package config

import (
	"github.com/apicat/apicat/v2/backend/module/storage"
	"github.com/apicat/apicat/v2/backend/module/storage/cloudflare"
	"github.com/apicat/apicat/v2/backend/module/storage/local"
	"github.com/apicat/apicat/v2/backend/module/storage/qiniu"
)

type Storage struct {
	Driver     string      `yaml:"Driver"`
	LocalDisk  *LocalDisk  `yaml:"LocalDisk"`
	Cloudflare *Cloudflare `yaml:"Cloudflare"`
	Qiniu      *Qiniu      `yaml:"Qiniu"`
}

type LocalDisk struct {
	Path string `yaml:"Path" json:"path"`
	Url  string `yaml:"Url" json:"url"`
}

type Cloudflare struct {
	AccountID       string `yaml:"AccountID" json:"accountID"`
	AccessKeyID     string `yaml:"AccessKeyID" json:"accessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret" json:"accessKeySecret"`
	BucketName      string `yaml:"BucketName" json:"bucketName"`
	Url             string `yaml:"Url" json:"bucketUrl"`
}

type Qiniu struct {
	AccessKeyID     string `yaml:"AccessKeyID" json:"accessKey"`
	AccessKeySecret string `yaml:"AccessKeySecret" json:"secretKey"`
	BucketName      string `yaml:"BucketName" json:"bucketName"`
	Url             string `yaml:"Url" json:"bucketUrl"`
}

func GetStorageDefault() *Storage {
	return &Storage{
		Driver: storage.LOCAL,
		LocalDisk: &LocalDisk{
			Path: "./uploads",
		},
	}
}

func SetStorage(storageConfig *Storage) {
	globalConf.Storage = storageConfig
}

func SetLocalDiskUrl(url string) {
	globalConf.Storage.LocalDisk.Url = url + "/uploads"
}

func (s *Storage) ToCfg() storage.Storage {
	switch s.Driver {
	case storage.CLOUDFLARE:
		return storage.Storage{
			Driver: s.Driver,
			Cloudflare: cloudflare.R2Opt{
				AccountID:       s.Cloudflare.AccountID,
				AccessKeyID:     s.Cloudflare.AccessKeyID,
				AccessKeySecret: s.Cloudflare.AccessKeySecret,
				BucketName:      s.Cloudflare.BucketName,
				Url:             s.Cloudflare.Url,
			},
		}
	case storage.QINIU:
		return storage.Storage{
			Driver: s.Driver,
			Qiniu: qiniu.QiniuOpt{
				AccessKeyID:     s.Qiniu.AccessKeyID,
				AccessKeySecret: s.Qiniu.AccessKeySecret,
				BucketName:      s.Qiniu.BucketName,
				Url:             s.Qiniu.Url,
			},
		}
	case storage.LOCAL:
		return storage.Storage{
			Driver: s.Driver,
			LocalDisk: local.Disk{
				Path: s.LocalDisk.Path,
				Url:  s.LocalDisk.Url,
			},
		}
	default:
		return storage.Storage{
			Driver: "",
		}
	}
}
