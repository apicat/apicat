package api

import (
	"net/http"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type ProjectMemberIDData struct {
	MemberID uint `uri:"member-id" binding:"required"`
}

// CheckMember checks if a project member exists given their member ID.
// Return: A pointer to ProjectMembers and an error, if any.
func (pmd *ProjectMemberIDData) CheckMember(ctx *gin.Context) (*models.ProjectMembers, error) {
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&pmd)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.NotFound"}),
		})
		return nil, err
	}

	member, err := models.NewProjectMembers(pmd.MemberID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.NotFound"}),
		})
		return nil, err
	}

	return member, nil
}

type ProjectMemberData struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Authority string `json:"authority"`
	CreatedAt string `json:"created_at"`
}

type GetProjectMemberData struct {
	UserID uint `uri:"user_id" binding:"required"`
}

// MembersList handles GET requests to retrieve a list of members in the current project.
func MembersList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	data := ProjectMemberIDData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	member, _ := models.NewProjectMembers()
	member.ProjectID = project.ID
	members, err := member.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	userIDs := []uint{}
	for _, v := range members {
		userIDs = append(userIDs, v.UserID)
	}

	users, err := models.UserListByIDs(userIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	userIDToNameMap := map[uint]string{}
	for _, v := range users {
		userIDToNameMap[v.ID] = v.Username
	}

	membersList := []*ProjectMemberData{}
	for _, v := range members {
		membersList = append(membersList, &ProjectMemberData{
			ID:        v.ID,
			UserID:    v.UserID,
			Username:  userIDToNameMap[v.UserID],
			Authority: v.Authority,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, membersList)
}

// MemberGet retrieves the details of a project member from the database and returns it as JSON.
func MemberGet(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

	data := GetProjectMemberData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := models.NewUsers(data.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.UserID = user.ID
	pm.ProjectID = project.ID
	if err := pm.Get(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, ProjectMemberData{
		ID:        pm.ID,
		UserID:    user.ID,
		Username:  user.Username,
		Authority: pm.Authority,
		CreatedAt: pm.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}
