package request

import (
	protobase "apicat-cloud/backend/route/proto/base"
)

type ProjectShareSwitchOption struct {
	protobase.ProjectIdOption
	Status bool `json:"status" binding:"boolean"`
}

type CheckProjectShareSecretKeyOption struct {
	protobase.ProjectIdOption
	protobase.SecretKeyOption
}
