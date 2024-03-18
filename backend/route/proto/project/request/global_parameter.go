package request

import (
	protobase "apicat-cloud/backend/route/proto/base"
	projectbase "apicat-cloud/backend/route/proto/project/base"
)

type GetGlobalParameterOption struct {
	protobase.ProjectIdOption
	ParameterID uint `uri:"parameterID" json:"parameterID" query:"parameterID" binding:"required,numeric,gt=0"`
}

type CreateGlobalParameterOption struct {
	protobase.ProjectIdOption
	projectbase.GlobalParameterDataOption
}

type UpdateGlobalParameterOption struct {
	GetGlobalParameterOption
	projectbase.GlobalParameterDataOption
}

type DeleteGlobalParameterOption struct {
	GetGlobalParameterOption
	projectbase.DerefOption
}

type SortGlobalParameterOption struct {
	protobase.ProjectIdOption
	In           string `json:"in" binding:"required,oneof=header cookie query path"`
	ParameterIDs []uint `json:"parameterIDs" binding:"omitempty,dive,gte=0"`
}
