package request

import (
	protobase "github.com/apicat/apicat/v2/backend/route/proto/base"
	teambase "github.com/apicat/apicat/v2/backend/route/proto/team/base"
)

type SettingOption struct {
	protobase.TeamIdOption
	teambase.TeamDataOption
}
