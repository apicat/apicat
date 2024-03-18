package response

import (
	protobase "apicat-cloud/backend/route/proto/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
)

type ProjectServer struct {
	protobase.OnlyIdInfo
	projectbase.ProjectServerDataOption
}

type ProjectServerList []*ProjectServer
