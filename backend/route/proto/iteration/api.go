package iteration

import (
	"github.com/apicat/apicat/v2/backend/route/proto/iteration/base"
	"github.com/apicat/apicat/v2/backend/route/proto/iteration/request"
	"github.com/apicat/apicat/v2/backend/route/proto/iteration/response"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type IterationApi interface {
	// Create 创建迭代
	// @route POST /teams/{teamID}/iterations
	Create(*gin.Context, *request.CreateIterationOption) (*ginrpc.Empty, error)

	// List 获取迭代列表
	// @route GET /teams/{teamID}/iterations
	List(*gin.Context, *request.GetIterationListOption) (*response.IterationList, error)

	// Get 迭代详情
	// @route GET /iterations/{iterationID}
	Get(*gin.Context, *base.IterationIDOption) (*response.Iteration, error)

	// Update 编辑迭代
	// @route PUT /iterations/{iterationID}
	Update(*gin.Context, *request.UpdateIterationOption) (*ginrpc.Empty, error)

	// Delete 删除迭代
	// @route DELETE /iterations/{iterationID}
	Delete(*gin.Context, *base.IterationIDOption) (*ginrpc.Empty, error)
}
