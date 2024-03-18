package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	"apicat-cloud/backend/route/proto/collection/base"
)

type CollectionShareDetail struct {
	protobase.ProjectVisibilityOption
	CollectionShareData
}

type CollectionShareData struct {
	base.CollectionPublicIDOption
	protobase.SecretKeyOption
}
