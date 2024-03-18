package request

type DiskOption struct {
	Path string `json:"path" binding:"required,gt=1"`
}

type CloudflareOption struct {
	AccountID       string `json:"accountID" binding:"required,gt=1"`
	AccessKeyID     string `json:"accessKeyID" binding:"required,gt=1"`
	AccessKeySecret string `json:"accessKeySecret" binding:"required,gt=1"`
	BucketName      string `json:"bucketName" binding:"required,gt=1"`
	BucketUrl       string `json:"bucketUrl" binding:"required,url"`
}

type QiniuOption struct {
	AccessKey  string `json:"accessKey" binding:"required,gt=1"`
	SecretKey  string `json:"secretKey" binding:"required,gt=1"`
	BucketName string `json:"bucketName" binding:"required,gt=1"`
	BucketUrl  string `json:"bucketUrl" binding:"required,url"`
}
