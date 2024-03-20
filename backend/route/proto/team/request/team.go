package request

import (
	protobase "github.com/apicat/apicat/backend/route/proto/base"
	teambase "github.com/apicat/apicat/backend/route/proto/team/base"
)

type SettingOption struct {
	protobase.TeamIdOption
	teambase.TeamDataOption
}
