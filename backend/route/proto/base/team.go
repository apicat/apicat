package base

type TeamIdOption struct {
	TeamID string `uri:"teamID" json:"teamID" query:"teamID" binding:"required,len=24"`
}

type InvitationTokenOption struct {
	InvitationToken string `uri:"invitationToken" query:"invitationToken" json:"invitationToken" binding:"omitempty,len=32"`
}
