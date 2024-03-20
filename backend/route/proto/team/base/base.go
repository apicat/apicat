package base

import (
	"github.com/apicat/apicat/backend/model/team"
)

type TeamMemberStatusOption struct {
	Status string `json:"status" binding:"omitempty,oneof=active deactive"`
}

type TeamDataOption struct {
	Name   string `json:"name" binding:"required,lte=255"`
	Avatar string `json:"avatar" binding:"omitempty,url"`
}

type TeamMemberData struct {
	Role team.Role `json:"role" binding:"omitempty,oneof=admin member"`
	TeamMemberStatusOption
}
