package cloudflare

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type r2 struct {
	accountId       string
	accessKeyId     string
	accessKeySecret string
	bucketName      string
	url             string
	client          *s3.Client
}

func NewR2(cfg map[string]interface{}) (*r2, error) {
	for _, v := range []string{"AccountID", "AccessKeyID", "AccessKeySecret", "BucketName", "Url"} {
		if _, ok := cfg[v]; !ok {
			return nil, fmt.Errorf("cloudflare config %s is required", v)
		}
	}

	return &r2{
		accountId:       cfg["AccountID"].(string),
		accessKeyId:     cfg["AccessKeyID"].(string),
		accessKeySecret: cfg["AccessKeySecret"].(string),
		bucketName:      cfg["BucketName"].(string),
		url:             strings.TrimRight(cfg["Url"].(string), "/"),
	}, nil
}

func (r *r2) init() error {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r.accountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r.accessKeyId, r.accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return err
	}

	r.client = s3.NewFromConfig(cfg)
	return nil
}

func (r *r2) Check() error {
	if err := r.init(); err != nil {
		return err
	}

	_, err := r.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  aws.String(r.bucketName),
		MaxKeys: aws.Int32(1),
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *r2) PutObject(key string, data []byte, contentType string) (string, error) {
	if err := r.init(); err != nil {
		return "", err
	}

	key = strings.TrimLeft(key, "/")
	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("couldn't upload file %v to %v. Here's why: %v", key, r.bucketName, err)
	}
	return fmt.Sprintf("%s/%s", r.url, key), nil
}
