package request

import protobase "github.com/apicat/apicat/backend/route/proto/base"

type GetMockOption struct {
	protobase.ProjectIdOption
	Path string `uri:"path" binding:"required"`
}
