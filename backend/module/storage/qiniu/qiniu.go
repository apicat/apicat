package qiniu

import (
	"context"
	"fmt"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuOpt struct {
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	Url             string
}

type qiniu struct {
	cfg      QiniuOpt
	qiniu    *auth.Credentials
	qiniuCfg *storage.Config
}

func NewQiniu(cfg QiniuOpt) *qiniu {
	cfg.Url = strings.TrimRight(cfg.Url, "/")
	return &qiniu{
		cfg: cfg,
	}
}

func (q *qiniu) init() {
	q.qiniu = qbox.NewMac(q.cfg.AccessKeyID, q.cfg.AccessKeySecret)
	q.qiniuCfg = &storage.Config{}
	if i := strings.Index(q.cfg.Url, "https"); i >= 0 {
		q.qiniuCfg.UseHTTPS = true
	} else {
		q.qiniuCfg.UseHTTPS = false
	}
}

func (q *qiniu) Check() error {
	q.init()
	bucketManager := storage.NewBucketManager(q.qiniu, q.qiniuCfg)
	if _, err := bucketManager.GetBucketInfo(q.cfg.BucketName); err != nil {
		return err
	}
	return nil
}

func (q *qiniu) PutObject(key string, data []byte, contentType string) (string, error) {
	q.init()

	putPolicy := storage.PutPolicy{
		Scope: q.cfg.BucketName,
	}
	upToken := putPolicy.UploadToken(q.qiniu)
	formUploader := storage.NewFormUploader(q.qiniuCfg)
	ret := storage.PutRet{}

	dataLen := int64(len(data))
	extra := &storage.PutExtra{
		MimeType: contentType,
	}
	err := formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), dataLen, extra)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", q.cfg.Url, key), nil
}

func (q *qiniu) GetObject(key string) (storage.FileInfo, error) {
	q.init()
	bucketManager := storage.NewBucketManager(q.qiniu, q.qiniuCfg)

	fileInfo, err := bucketManager.Stat(q.cfg.BucketName, key)
	if err != nil {
		return storage.FileInfo{}, err
	}
	return fileInfo, nil
}

func (q *qiniu) DelObject(key string) error {
	q.init()
	bucketManager := storage.NewBucketManager(q.qiniu, q.qiniuCfg)

	err := bucketManager.Delete(q.cfg.BucketName, key)
	if err != nil {
		return err
	}
	return nil
}
