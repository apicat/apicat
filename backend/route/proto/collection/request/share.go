package request

import (
	protobase "apicat-cloud/backend/route/proto/base"
	"apicat-cloud/backend/route/proto/collection/base"
)

type SwitchCollectionShareOption struct {
	base.ProjectCollectionIDOption
	Status bool `json:"status" binding:"boolean"`
}

type CheckCollectionShareSecretKeyOpt struct {
	base.ProjectCollectionIDOption
	protobase.SecretKeyOption
}
