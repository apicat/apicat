package collection

import (
	protobase "apicat-cloud/backend/route/proto/base"
	"apicat-cloud/backend/route/proto/collection/base"
	"apicat-cloud/backend/route/proto/collection/request"
	"apicat-cloud/backend/route/proto/collection/response"

	"apicat-cloud/backend/module/spec"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type CollectionApi interface {
	// Create 创建集合
	// @route POST /projects/{projectID}/collections
	Create(*gin.Context, *request.CreateCollectionOption) (*response.Collection, error)

	// List 获取集合树
	// @route GET /projects/{projectID}/collections
	List(*gin.Context, *request.GetCollectionListOption) (*response.CollectionTree, error)

	// Get 集合详情
	// @route GET /projects/{projectID}/collections/{collectionID}
	Get(*gin.Context, *base.ProjectCollectionIDOption) (*response.Collection, error)

	// Update 编辑集合
	// @route PUT /projects/{projectID}/collections/{collectionID}
	Update(*gin.Context, *request.UpdateCollectionOption) (*ginrpc.Empty, error)

	// Delete 删除集合
	// @route DELETE /projects/{projectID}/collections/{collectionID}
	Delete(*gin.Context, *request.DeleteCollectionOption) (*ginrpc.Empty, error)

	// Move 移动集合
	// @route PUT /projects/{projectID}/collections/move
	Move(*gin.Context, *request.MoveCollectionOption) (*ginrpc.Empty, error)

	// Copy 复制集合
	// @route POST /projects/{projectID}/collections/{collectionID}/copy
	Copy(*gin.Context, *request.CopyCollectionOption) (*response.Collection, error)

	// Trashes 已删除的集合列表
	// @route GET /projects/{projectID}/collections/trashes
	Trashes(*gin.Context, *protobase.ProjectIdOption) (*response.TrashList, error)

	// Restore 恢复已删除的集合
	// @route PUT /projects/{projectID}/collections/restore
	Restore(*gin.Context, *request.RestoreOption) (*response.RestoreNum, error)

	// GetExportPath 获取导出path
	// @route GET /projects/{projectID}/collections/{collectionID}/export
	GetExportPath(*gin.Context, *request.GetExportPathOption) (*response.ExportCollection, error)

	// AIGenerate AI 生成集合
	// @route POST /projects/{projectID}/ai/collections
	AIGenerate(*gin.Context, *request.AIGenerateCollectionOption) (*response.Collection, error)
}

type CollectionShareApi interface {
	// Status 获取集合分享状态
	// @route GET /collections/{collectionPublicID}/share/status
	Status(*gin.Context, *base.CollectionPublicIDOption) (*base.ProjectCollectionIDOption, error)

	// Detail 集合分享详情
	// @route GET /projects/{projectID}/collections/{collectionID}/share
	Detail(*gin.Context, *base.ProjectCollectionIDOption) (*response.CollectionShareDetail, error)

	// Switch 切换集合分享状态
	// @route PUT /projects/{projectID}/collections/{collectionID}/share
	Switch(*gin.Context, *request.SwitchCollectionShareOption) (*response.CollectionShareData, error)

	// Reset 重置集合分享密钥
	// @route PUT /projects/{projectID}/collections/{collectionID}/share/reset
	Reset(*gin.Context, *base.ProjectCollectionIDOption) (*protobase.SecretKeyOption, error)

	// Check 检查分享密钥
	// @route POST /projects/{projectID}/collections/{collectionID}/share/check
	Check(*gin.Context, *request.CheckCollectionShareSecretKeyOpt) (*base.ShareCode, error)
}

type CollectionHistoryApi interface {
	// List 获取集合历史列表
	// @route GET /projectss/{projectID}/collectionss/{collectionID}/histories
	List(*gin.Context, *request.GetCollectionHistoryListOption) (*response.CollectionHistoryList, error)

	// Get 获取集合历史详情
	// @route GET /projectss/{projectID}/collectionss/{collectionID}/histories/{historyID}
	Get(*gin.Context, *request.CollectionHistoryIDOption) (*response.CollectionHistory, error)

	// Restore 恢复集合历史
	// @route PUT /projectss/{projectID}/collectionss/{collectionID}/histories/{historyID}/restore
	Restore(*gin.Context, *request.CollectionHistoryIDOption) (*ginrpc.Empty, error)

	// Diff 集合历史对比
	// @route GET /projectss/{projectID}/collectionss/{collectionID}/histories/diff
	Diff(*gin.Context, *request.DiffCollectionHistoriesOption) (*response.DiffCollectionHistories, error)
}

type CollectionMockApi interface {
	// Mock 集合mock
	// @route POST /mock/{projectID}/{path}
	Mock(*gin.Context, *request.GetMockOption) (*spec.HTTPResponses, error)
}

type TestCaseApi interface {
	// Generate 生成测试用例
	// @route POST /projects/{projectID}/collections/{collectionID}/testcases
	Generate(*gin.Context, *request.GenerateTestCaseOption) (*ginrpc.Empty, error)

	// List 获取测试用例列表
	// @route GET /projects/{projectID}/collections/{collectionID}/testcases
	List(*gin.Context, *base.ProjectCollectionIDOption) (*response.TestCaseList, error)

	// Get 测试用例详情
	// @route GET /projects/{projectID}/collections/{collectionID}/testcases/{testCaseID}
	Get(*gin.Context, *request.GetTestCaseOption) (*response.TestCaseDetail, error)

	// Regenerate 重新生成测试用例
	// @route PUT /projects/{projectID}/collections/{collectionID}/testcases/{testCaseID}
	Regenerate(*gin.Context, *request.RegenerateTestCaseOption) (*response.TestCaseDetail, error)

	// Delete 删除测试用例
	// @route DELETE /projects/{projectID}/collections/{collectionID}/testcases/{testCaseID}
	Delete(*gin.Context, *request.DeleteTestCaseOption) (*ginrpc.Empty, error)
}
