package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
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
