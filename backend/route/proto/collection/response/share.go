package response

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	"github.com/apicat/apicat/v2/backend/route/proto/collection/base"
)

type CollectionShareDetail struct {
	protobase.ProjectVisibilityOption
	CollectionShareData
}

type CollectionShareData struct {
	base.CollectionPublicIDOption
	protobase.SecretKeyOption
}
