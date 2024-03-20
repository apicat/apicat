package project

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/v2/backend/route/proto/project/base"
	"github.com/apicat/apicat/v2/backend/route/proto/project/request"
	"github.com/apicat/apicat/v2/backend/route/proto/project/response"
	teamresponse "github.com/apicat/apicat/v2/backend/route/proto/team/response"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

// ProjectApi 项目相关接口定义
type ProjectApi interface {
	// Create 创建项目
	// @route POST /teams/{teamID}/projects
	Create(*gin.Context, *request.CreateProjectOption) (*response.ProjectListItem, error)

	// List 获取项目列表
	// @route GET /teams/{teamID}/projects
	List(*gin.Context, *request.GetProjectListOption) (*response.GetProjectsResponse, error)

	// Get 获取项目
	// @route GET /projects/{projectID}
	Get(*gin.Context, *request.GetProjectDetailOption) (*response.ProjectDetail, error)

	// ChangeGroup 切换项目分组
	// @route PUT /projects/{projectID}/group
	ChangeGroup(*gin.Context, *request.SwitchProjectGroupOption) (*ginrpc.Empty, error)

	// Follow 关注项目
	// @route POST /projects/{projectID}/follow
	Follow(*gin.Context, *protobase.ProjectIdOption) (*ginrpc.Empty, error)

	// UnFollow 取消关注项目
	// @route DELETE /projects/{projectID}/follow
	UnFollow(*gin.Context, *protobase.ProjectIdOption) (*ginrpc.Empty, error)

	// Setting 项目设置
	// @route PUT /projects/{projectID}
	Setting(*gin.Context, *request.UpdateProjectOption) (*ginrpc.Empty, error)

	// Delete 删除项目
	// @route DELETE /projects/{projectID}
	Delete(*gin.Context, *protobase.ProjectIdOption) (*ginrpc.Empty, error)

	// Transfer 移交项目
	// @route PUT /projects/{projectID}/transfer
	Transfer(*gin.Context, *request.ProjectMemberIDOption) (*ginrpc.Empty, error)

	// Exit 退出项目
	// @route DELETE /projects/{projectID}/exit
	Exit(*gin.Context, *protobase.ProjectIdOption) (*ginrpc.Empty, error)

	// GetExportPath 获取导出path
	// @route GET /projects/{projectID}/export
	GetExportPath(*gin.Context, *request.GetExportPathOption) (*response.ExportProject, error)
}

type ProjectMemberApi interface {
	// Create 创建项目成员
	// @route POST /projects/{projectID}/member
	Create(*gin.Context, *request.CreateProjectMemberOption) (*ginrpc.Empty, error)

	// Members 获取项目成员列表
	// @route GET /projects/{projectID}/member
	Members(*gin.Context, *request.GetProjectMemberListOption) (*response.GetProjectMemberListResponse, error)

	// NotInProjectMembers 不在此项目的成员列表
	// @route GET /projects/{projectID}/member/without
	NotInProjectMembers(*gin.Context, *protobase.ProjectIdOption) (*teamresponse.TeamMembers, error)

	// Update 更新项目成员
	// @route PUT /projects/{projectID}/member/{memberID}
	Update(*gin.Context, *request.UpdateProjectMemberOption) (*ginrpc.Empty, error)

	// Delete 删除项目成员
	// @route DELETE /projects/{projectID}/member/{memberID}
	Delete(*gin.Context, *request.ProjectMemberIDOption) (*ginrpc.Empty, error)
}

type ProjectGroupApi interface {
	// Create 创建项目分组
	// @route POST /teams/{teamID}/projects/group
	Create(*gin.Context, *request.CreateProjectGroupOption) (*response.ProjectGroup, error)

	// List 获取项目分组列表
	// @route GET /teams/{teamID}/projects/group
	List(*gin.Context, *protobase.TeamIdOption) (*response.ProjectGroups, error)

	// Delete 删除项目分组
	// @route DELETE /projects/group/{groupID}
	Delete(*gin.Context, *request.ProjectGroupIdOption) (*ginrpc.Empty, error)

	// Rename 重命名项目分组
	// @route PUT /projects/group/{groupID}
	Rename(*gin.Context, *request.RenameProjectGroupOption) (*ginrpc.Empty, error)

	// Sort 项目分组排序
	// @route PUT /teams/{teamID}/projects/group/sort
	Sort(*gin.Context, *request.SortProjectGroupOption) (*ginrpc.Empty, error)
}

type DefinitionResponseApi interface {
	// Create 创建定义响应
	// @route POST /projects/{projectID}/definition/responses
	Create(*gin.Context, *request.CreateDefinitionResponseOption) (*response.DefinitionResponse, error)

	// List 获取定义响应树
	// @route GET /projects/{projectID}/definition/responses
	List(*gin.Context, *protobase.ProjectIdOption) (*response.DefinitionResponseTree, error)

	// Get 定义响应详情
	// @route GET /projects/{projectID}/definition/responses/{responseID}
	Get(*gin.Context, *request.GetDefinitionResponseOption) (*response.DefinitionResponse, error)

	// Update 编辑定义响应
	// @route PUT /projects/{projectID}/definition/responses/{responseID}
	Update(*gin.Context, *request.UpdateDefinitionResponseOption) (*ginrpc.Empty, error)

	// Delete 删除定义响应
	// @route DELETE /projects/{projectID}/definition/responses/{responseID}
	Delete(*gin.Context, *request.DeleteDefinitionResponseOption) (*ginrpc.Empty, error)

	// Move 移动定义响应
	// @route PUT /projects/{projectID}/definition/responses/move
	Move(*gin.Context, *request.SortDefinitionResponseOption) (*ginrpc.Empty, error)

	// Copy 复制定义响应
	// @route POST /projects/{projectID}/definition/responses/{responseID}/copy
	Copy(*gin.Context, *request.GetDefinitionResponseOption) (*response.DefinitionResponse, error)
}

type DefinitionSchemaApi interface {
	// Create 创建定义模型
	// @route POST /projects/{projectID}/definition/schemas
	Create(*gin.Context, *request.CreateDefinitionSchemaOption) (*response.DefinitionSchema, error)

	// List 获取定义模型树
	// @route GET /projects/{projectID}/definition/schemas
	List(*gin.Context, *protobase.ProjectIdOption) (*response.DefinitionSchemaTree, error)

	// Get 定义模型详情
	// @route GET /projects/{projectID}/definition/schemas/{schemaID}
	Get(*gin.Context, *request.GetDefinitionSchemaOption) (*response.DefinitionSchema, error)

	// Update 编辑定义模型
	// @route PUT /projects/{projectID}/definition/schemas/{schemaID}
	Update(*gin.Context, *request.UpdateDefinitionSchemaOption) (*ginrpc.Empty, error)

	// Delete 删除定义模型
	// @route DELETE /projects/{projectID}/definition/schemas/{schemaID}
	Delete(*gin.Context, *request.DeleteDefinitionSchemaOption) (*ginrpc.Empty, error)

	// Move 移动定义模型
	// @route PUT /projects/{projectID}/definition/schemas/move
	Move(*gin.Context, *request.SortDefinitionSchemaOption) (*ginrpc.Empty, error)

	// Copy 复制定义模型
	// @route POST /projects/{projectID}/definition/schemas/{schemaID}/copy
	Copy(*gin.Context, *request.GetDefinitionSchemaOption) (*response.DefinitionSchema, error)

	// AIGenerate AI 生成模型
	// @route POST /projects/{projectID}/definition/ai/schemas
	AIGenerate(*gin.Context, *request.AIGenerateSchemaOption) (*response.DefinitionSchema, error)
}

// GlobalParameterApi 全局参数相关
type GlobalParameterApi interface {
	// Create 创建全局参数
	// @route POST /projects/{projectID}/global/parameters
	Create(*gin.Context, *request.CreateGlobalParameterOption) (*response.GlobalParameter, error)

	// List 获取全局参数列表
	// @route GET /projects/{projectID}/global/parameters
	List(*gin.Context, *protobase.ProjectIdOption) (*response.GlobalParameterList, error)

	// Update 更新全局参数
	// @route PUT /projects/{projectID}/global/parameters/{parameterID}
	Update(*gin.Context, *request.UpdateGlobalParameterOption) (*ginrpc.Empty, error)

	// Delete 删除全局参数
	// @route DELETE /projects/{projectID}/global/parameters/{parameterID}
	Delete(*gin.Context, *request.DeleteGlobalParameterOption) (*ginrpc.Empty, error)

	// Sort 全局参数排序
	// @route PUT /projects/{projectID}/global/parameters/sort
	Sort(*gin.Context, *request.SortGlobalParameterOption) (*ginrpc.Empty, error)
}

type ProjectServerApi interface {
	// Create 创建项目URL
	// @route POST /projects/{projectID}/servers
	Create(*gin.Context, *request.CreateProjectServerOption) (*response.ProjectServer, error)

	// List 获取项目URL列表
	// @route GET /projects/{projectID}/servers
	List(*gin.Context, *protobase.ProjectIdOption) (*response.ProjectServerList, error)

	// Update 修改项目URL
	// @route PUT /projects/{projectID}/servers/{serverID}
	Update(*gin.Context, *request.UpdateProjectServerOption) (*ginrpc.Empty, error)

	// Delete 删除项目URL
	// @route Delete /projects/{projectID}/servers/{serverID}
	Delete(*gin.Context, *request.GetProjectServerOption) (*ginrpc.Empty, error)

	// Sort 项目URL排序
	// @route PUT /projects/{projectID}/servers/sort
	Sort(*gin.Context, *request.SortProjectServerOpt) (*ginrpc.Empty, error)
}

type ProjectShareApi interface {
	// Status 获取项目分享状态
	// @route GET /projects/{projectID}/share/status
	Status(*gin.Context, *protobase.ProjectIdOption) (*response.ProjectShareStatus, error)

	// Detail 项目分享详情
	// @route GET /projects/{projectID}/share
	Detail(*gin.Context, *protobase.ProjectIdOption) (*response.ProjectShareDetail, error)

	// Switch 切换项目分享状态
	// @route PUT /projects/{projectID}/share
	Switch(*gin.Context, *request.ProjectShareSwitchOption) (*protobase.SecretKeyOption, error)

	// Reset 重置项目分享密钥
	// @route PUT /projects/{projectID}/share/reset
	Reset(*gin.Context, *protobase.ProjectIdOption) (*protobase.SecretKeyOption, error)

	// Check 检查分享密钥
	// @route POST /projects/{projectID}/share/check
	Check(*gin.Context, *request.CheckProjectShareSecretKeyOption) (*projectbase.ShareCode, error)
}

type DefinitionSchemaHistoryAPi interface {
	// List 获取定义模型历史列表
	// @route GET /projects/{projectID}/definition/schemas/{schemaID}/histories
	List(*gin.Context, *request.GetDefinitionSchemaHistoryListOption) (*response.DefinitionSchemaHistoryList, error)

	// Get 获取定义模型历史详情
	// @route GET /projects/{projectID}/definition/schemas/{schemaID}/histories/{historyID}
	Get(*gin.Context, *request.DefinitionSchemaHistoryIDOption) (*response.DefinitionSchemaHistory, error)

	// Restore 恢复定义模型历史
	// @route PUT /projects/{projectID}/definition/schemas/{schemaID}/histories/{historyID}/restore
	Restore(*gin.Context, *request.DefinitionSchemaHistoryIDOption) (*ginrpc.Empty, error)

	// Diff 定义模型历史对比
	// @route GET /projects/{projectID}/definition/schemas/{schemaID}/histories/diff
	Diff(*gin.Context, *request.DiffDefinitionSchemaHistoriesOption) (*response.DiffDefinitionSchemaHistories, error)
}
