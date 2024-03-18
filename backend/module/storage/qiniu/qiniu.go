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

type qiniu struct {
	accessKeyId     string
	accessKeySecret string
	bucketName      string
	url             string
	qiniu           *auth.Credentials
	qiniuCfg        *storage.Config
}

func NewQiniu(cfg map[string]interface{}) (*qiniu, error) {
	for _, v := range []string{"AccessKeyID", "AccessKeySecret", "BucketName", "Url"} {
		if _, ok := cfg[v]; !ok {
			return nil, fmt.Errorf("sendcloud config %s is required", v)
		}
	}
	return &qiniu{
		accessKeyId:     cfg["AccessKeyID"].(string),
		accessKeySecret: cfg["AccessKeySecret"].(string),
		bucketName:      cfg["BucketName"].(string),
		url:             strings.TrimRight(cfg["Url"].(string), "/"),
	}, nil
}

func (q *qiniu) init() {
	q.qiniu = qbox.NewMac(q.accessKeyId, q.accessKeySecret)
	q.qiniuCfg = &storage.Config{}
	if i := strings.Index(q.url, "https"); i >= 0 {
		q.qiniuCfg.UseHTTPS = true
	} else {
		q.qiniuCfg.UseHTTPS = false
	}
}

func (q *qiniu) Check() error {
	q.init()
	bucketManager := storage.NewBucketManager(q.qiniu, q.qiniuCfg)
	if _, err := bucketManager.GetBucketInfo(q.bucketName); err != nil {
		return err
	}
	return nil
}

func (q *qiniu) PutObject(key string, data []byte, contentType string) (string, error) {
	q.init()

	putPolicy := storage.PutPolicy{
		Scope: q.bucketName,
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

	return fmt.Sprintf("%s/%s", q.url, key), nil
}

func (q *qiniu) GetObject(key string) (storage.FileInfo, error) {
	q.init()
	bucketManager := storage.NewBucketManager(q.qiniu, q.qiniuCfg)

	fileInfo, err := bucketManager.Stat(q.bucketName, key)
	if err != nil {
		return storage.FileInfo{}, err
	}
	return fileInfo, nil
}

func (q *qiniu) DelObject(key string) error {
	q.init()
	bucketManager := storage.NewBucketManager(q.qiniu, q.qiniuCfg)

	err := bucketManager.Delete(q.bucketName, key)
	if err != nil {
		return err
	}
	return nil
}
