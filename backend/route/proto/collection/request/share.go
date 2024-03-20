package request

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	"github.com/apicat/apicat/backend/route/proto/collection/base"
)

type SwitchCollectionShareOption struct {
	base.ProjectCollectionIDOption
	Status bool `json:"status" binding:"boolean"`
}

type CheckCollectionShareSecretKeyOpt struct {
	base.ProjectCollectionIDOption
	protobase.SecretKeyOption
}
