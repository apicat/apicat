package sysconfig

import (
	sysconfigbase "apicat-cloud/backend/route/proto/sysconfig/base"
	sysconfigrequest "apicat-cloud/backend/route/proto/sysconfig/request"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type ServiceApi interface {
	// Get Get service config
	// @route GET /sysconfigs/service
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigbase.ServiceOption, error)

	// Update Update service config
	// @route PUT /sysconfigs/service
	Update(*gin.Context, *sysconfigbase.ServiceOption) (*ginrpc.Empty, error)

	// GetDB Get MySQL config
	// @route GET /sysconfigs/db
	GetDB(*gin.Context, *ginrpc.Empty) (*sysconfigbase.MySQLDetail, error)
}

type OauthApi interface {
	// GetGithubClientID Get Github client id
	// @route GET /sysconfigs/github
	GetGithubClientID(*gin.Context, *ginrpc.Empty) (*sysconfigbase.GitHubClientID, error)

	// Get Get GitHub oauth config
	// @route GET /sysconfigs
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigbase.GitHubOption, error)

	// Update Update GitHub oauth config
	// @route PUT /sysconfigs
	Update(*gin.Context, *sysconfigbase.GitHubOption) (*ginrpc.Empty, error)
}

type StorageApi interface {
	// Get Get storage config list
	// @route GET /sysconfigs/storages
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigbase.ConfigList, error)

	// Update Update disk storage config
	// @route PUT /sysconfigs/storages/disk
	UpdateDisk(*gin.Context, *sysconfigrequest.DiskOption) (*ginrpc.Empty, error)

	// Update Update cloudflare storage config
	// @route PUT /sysconfigs/storages/cloudflare
	UpdateCloudflare(*gin.Context, *sysconfigrequest.CloudflareOption) (*ginrpc.Empty, error)

	// Update Update qiniu storage config
	// @route PUT /sysconfigs/storages/qiniu
	UpdateQiniu(*gin.Context, *sysconfigrequest.QiniuOption) (*ginrpc.Empty, error)
}

type CacheApi interface {
	// Get Get cache config list
	// @route GET /sysconfigs/caches
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigbase.ConfigList, error)

	// Update Update memory cache config
	// @route PUT /sysconfigs/caches/memory
	UpdateMemory(*gin.Context, *ginrpc.Empty) (*ginrpc.Empty, error)

	// Update Update redis cache config
	// @route PUT /sysconfigs/caches/redis
	UpdateRedis(*gin.Context, *sysconfigrequest.RedisOption) (*ginrpc.Empty, error)
}

type EmailApi interface {
	// Get Get email config list
	// @route GET /sysconfigs/emails
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigbase.ConfigList, error)

	// Update Update smtp email config
	// @route PUT /sysconfigs/emails/smtp
	UpdateSMTP(*gin.Context, *sysconfigrequest.SMTPOption) (*ginrpc.Empty, error)

	// Update Update sendcloud email config
	// @route PUT /sysconfigs/emails/sendcloud
	UpdateSendCloud(*gin.Context, *sysconfigrequest.SendCloudOption) (*ginrpc.Empty, error)
}

type ModelApi interface {
	// Get Get model config list
	// @route GET /sysconfigs/models
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigbase.ConfigList, error)

	// Update Update OpenAI model config
	// @route PUT /sysconfigs/models/openai
	UpdateOpenAI(*gin.Context, *sysconfigrequest.OpenAIOption) (*ginrpc.Empty, error)

	// Update Update Azure OpenAI model config
	// @route PUT /sysconfigs/models/azure-openai
	UpdateAzureOpenAI(*gin.Context, *sysconfigrequest.AzureOpenAIOption) (*ginrpc.Empty, error)
}
