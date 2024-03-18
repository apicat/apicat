package request

import protobase "apicat-cloud/backend/route/proto/base"

type GetMockOption struct {
	protobase.ProjectIdOption
	Path string `uri:"path" binding:"required"`
}
