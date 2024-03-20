package request

import protobase "github.com/apicat/apicat/v2/backend/route/proto/base"

type GetMockOption struct {
	protobase.ProjectIdOption
	Path string `uri:"path" binding:"required"`
}
