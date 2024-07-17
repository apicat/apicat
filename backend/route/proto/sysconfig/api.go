package sysconfig

import (
	sysconfigbase "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/base"
	sysconfigrequest "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/request"

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
