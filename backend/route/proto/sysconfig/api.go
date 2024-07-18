package sysconfig

import (
	sysconfigbase "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/base"
	sysconfigrequest "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/request"
	sysconfigresponse "github.com/apicat/apicat/v2/backend/route/proto/sysconfig/response"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type ServiceApi interface {
	// Get service config
	// @route GET /sysconfigs/service
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigbase.ServiceOption, error)

	// Update service config
	// @route PUT /sysconfigs/service
	Update(*gin.Context, *sysconfigbase.ServiceOption) (*ginrpc.Empty, error)
}

type OauthApi interface {
	// Get Github client id
	// @route GET /sysconfigs/github
	GetGithubClientID(*gin.Context, *ginrpc.Empty) (*sysconfigbase.GitHubClientID, error)

	// Get GitHub oauth config
	// @route GET /sysconfigs
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigbase.GitHubOption, error)

	// Update GitHub oauth config
	// @route PUT /sysconfigs
	Update(*gin.Context, *sysconfigbase.GitHubOption) (*ginrpc.Empty, error)
}

type EmailApi interface {
	// Get email config list
	// @route GET /sysconfigs/emails
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigbase.ConfigList, error)

	// Update smtp email config
	// @route PUT /sysconfigs/emails/smtp
	UpdateSMTP(*gin.Context, *sysconfigrequest.SMTPOption) (*ginrpc.Empty, error)

	// Update sendcloud email config
	// @route PUT /sysconfigs/emails/sendcloud
	UpdateSendCloud(*gin.Context, *sysconfigrequest.SendCloudOption) (*ginrpc.Empty, error)
}

type ModelApi interface {
	// Get model config list
	// @route GET /sysconfigs/models
	Get(*gin.Context, *ginrpc.Empty) (*sysconfigresponse.ModelConfigList, error)

	// Update OpenAI model config
	// @route PUT /sysconfigs/models/openai
	UpdateOpenAI(*gin.Context, *sysconfigrequest.OpenAIOption) (*ginrpc.Empty, error)

	// Update Azure OpenAI model config
	// @route PUT /sysconfigs/models/azure-openai
	UpdateAzureOpenAI(*gin.Context, *sysconfigrequest.AzureOpenAIOption) (*ginrpc.Empty, error)

	// Get default model
	// @route GET /sysconfigs/models/default
	GetDefault(*gin.Context, *ginrpc.Empty) (*sysconfigresponse.DefaultModelMap, error)

	// Update default model
	// @route PUT /sysconfigs/models/default
	UpdateDefault(*gin.Context, *sysconfigrequest.DefaultModelMapOption) (*ginrpc.Empty, error)
}
