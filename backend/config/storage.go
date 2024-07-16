package config

import (
	"errors"
	"os"

	"github.com/apicat/apicat/v2/backend/module/storage"
	"github.com/apicat/apicat/v2/backend/module/storage/cloudflare"
	"github.com/apicat/apicat/v2/backend/module/storage/local"
	"github.com/apicat/apicat/v2/backend/module/storage/qiniu"
)

type Storage struct {
	Driver     string
	LocalDisk  *LocalDisk
	Cloudflare *Cloudflare
	Qiniu      *Qiniu
}

type LocalDisk struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

type Cloudflare struct {
	AccountID       string `json:"accountID"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	BucketName      string `json:"bucketName"`
	Url             string `json:"bucketUrl"`
}

type Qiniu struct {
	AccessKeyID     string `json:"accessKey"`
	AccessKeySecret string `json:"secretKey"`
	BucketName      string `json:"bucketName"`
	Url             string `json:"bucketUrl"`
}

func LoadStorageConfig() {
	globalConf.Storage = &Storage{}
	if v, exists := os.LookupEnv("STORAGE_DRIVER"); exists {
		switch v {
		case storage.CLOUDFLARE:
			globalConf.Storage.Driver = storage.CLOUDFLARE
			loadCloudflareConfig()
		case storage.QINIU:
			globalConf.Storage.Driver = storage.QINIU
			loadQiniuConfig()
		case storage.LOCAL:
			globalConf.Storage.Driver = storage.LOCAL
			loadLocalDiskConfig()
		}
	}
}

func loadCloudflareConfig() {
	globalConf.Storage.Cloudflare = &Cloudflare{}
	if v, exists := os.LookupEnv("CLOUDFLARE_ACCOUNT_ID"); exists {
		globalConf.Storage.Cloudflare.AccountID = v
	}
	if v, exists := os.LookupEnv("CLOUDFLARE_ACCESS_KEY_ID"); exists {
		globalConf.Storage.Cloudflare.AccessKeyID = v
	}
	if v, exists := os.LookupEnv("CLOUDFLARE_ACCESS_KEY_SECRET"); exists {
		globalConf.Storage.Cloudflare.AccessKeySecret = v
	}
	if v, exists := os.LookupEnv("CLOUDFLARE_BUCKET_NAME"); exists {
		globalConf.Storage.Cloudflare.BucketName = v
	}
	if v, exists := os.LookupEnv("CLOUDFLARE_URL"); exists {
		globalConf.Storage.Cloudflare.Url = v
	}
}

func loadQiniuConfig() {
	globalConf.Storage.Qiniu = &Qiniu{}
	if v, exists := os.LookupEnv("QINIU_ACCESS_KEY_ID"); exists {
		globalConf.Storage.Qiniu.AccessKeyID = v
	}
	if v, exists := os.LookupEnv("QINIU_ACCESS_KEY_SECRET"); exists {
		globalConf.Storage.Qiniu.AccessKeySecret = v
	}
	if v, exists := os.LookupEnv("QINIU_BUCKET_NAME"); exists {
		globalConf.Storage.Qiniu.BucketName = v
	}
	if v, exists := os.LookupEnv("QINIU_URL"); exists {
		globalConf.Storage.Qiniu.Url = v
	}
}

func loadLocalDiskConfig() {
	globalConf.Storage.LocalDisk = &LocalDisk{}
	if v, exists := os.LookupEnv("LOCAL_DISK_PATH"); exists {
		globalConf.Storage.LocalDisk.Path = v
	}
}

func CheckStorageConfig() error {
	if globalConf.Storage.Driver == "" {
		return errors.New("storage driver is empty")
	}
	switch globalConf.Storage.Driver {
	case storage.CLOUDFLARE:
		if globalConf.Storage.Cloudflare == nil {
			return errors.New("cloudflare config is empty")
		}
		if globalConf.Storage.Cloudflare.AccountID == "" {
			return errors.New("cloudflare account id is empty")
		}
		if globalConf.Storage.Cloudflare.AccessKeyID == "" {
			return errors.New("cloudflare access key id is empty")
		}
		if globalConf.Storage.Cloudflare.AccessKeySecret == "" {
			return errors.New("cloudflare access key secret is empty")
		}
		if globalConf.Storage.Cloudflare.BucketName == "" {
			return errors.New("cloudflare bucket name is empty")
		}
		if globalConf.Storage.Cloudflare.Url == "" {
			return errors.New("cloudflare url is empty")
		}
	case storage.QINIU:
		if globalConf.Storage.Qiniu == nil {
			return errors.New("qiniu config is empty")
		}
		if globalConf.Storage.Qiniu.AccessKeyID == "" {
			return errors.New("qiniu access key id is empty")
		}
		if globalConf.Storage.Qiniu.AccessKeySecret == "" {
			return errors.New("qiniu access key secret is empty")
		}
		if globalConf.Storage.Qiniu.BucketName == "" {
			return errors.New("qiniu bucket name is empty")
		}
		if globalConf.Storage.Qiniu.Url == "" {
			return errors.New("qiniu url is empty")
		}
	case storage.LOCAL:
		if globalConf.Storage.LocalDisk == nil {
			return errors.New("local disk config is empty")
		}
		if globalConf.Storage.LocalDisk.Path == "" {
			return errors.New("local disk path is empty")
		}
	}
	return nil
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
		return storage.Storage{}
	}
}
