package request

import (
	protobase "apicat-cloud/backend/route/proto/base"
	teambase "apicat-cloud/backend/route/proto/team/base"
)

type SettingOption struct {
	protobase.TeamIdOption
	teambase.TeamDataOption
}
