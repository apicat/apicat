package response

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	projectbase "github.com/apicat/apicat/v2/backend/route/proto/project/base"
)

type GlobalParameter struct {
	protobase.OnlyIdInfo
	projectbase.GlobalParameterDataOption
}

type GlobalParameterList struct {
	Header []*GlobalParameter `json:"header"`
	Cookie []*GlobalParameter `json:"cookie"`
	Query  []*GlobalParameter `json:"query"`
	Path   []*GlobalParameter `json:"path"`
}
