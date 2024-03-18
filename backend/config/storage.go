package config

import (
	"apicat-cloud/backend/module/storage"
	"encoding/json"
)

type Storage struct {
	Driver     string      `yaml:"Driver"`
	LocalDisk  *LocalDisk  `yaml:"LocalDisk"`
	Cloudflare *Cloudflare `yaml:"Cloudflare"`
	Qiniu      *Qiniu      `yaml:"Qiniu"`
}

type LocalDisk struct {
	Path string `yaml:"Path"`
	Url  string `yaml:"Url"`
}

type Cloudflare struct {
	AccountID       string `yaml:"AccountID"`
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	BucketName      string `yaml:"BucketName"`
	Url             string `yaml:"Url"`
}

type Qiniu struct {
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	BucketName      string `yaml:"BucketName"`
	Url             string `yaml:"Url"`
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

func (s *Storage) ToMapInterface() map[string]interface{} {
	var (
		res      map[string]interface{}
		jsonByte []byte
	)
	jsonByte, _ = json.Marshal(s)
	_ = json.Unmarshal(jsonByte, &res)
	return res
}
