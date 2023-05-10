package api

type ProjectMember struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Authority string `json:"authority"`
	CreatedAt string `json:"created_at"`
}

type ProjectMembersData struct {
	ProjectID string `json:"project_id" binding:"required"`
}

// func MembersList(ctx *gin.Context) {
// 	currentProject, _ := ctx.Get("CurrentProject")
// 	project, _ := currentProject.(*models.Projects)
// }
