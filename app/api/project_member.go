package api

import (
	"net/http"

	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type GetPathUserID struct {
	UserID uint `uri:"user-id" binding:"required"`
}

type CreateProjectMemberData struct {
	UserID    uint   `json:"user_id" binding:"required"`
	Authority string `json:"authority" binding:"required,oneof=manage write read"`
}

type UpdateProjectMemberAuthData struct {
	Authority string `json:"authority" binding:"required,oneof=manage write read"`
}

type ProjectMemberData struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Authority string `json:"authority"`
	CreatedAt string `json:"created_at"`
}

// MembersList handles GET requests to retrieve a list of members in the current project.
func ProjectMembersList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	member, _ := models.NewProjectMembers()
	member.ProjectID = currentProject.(*models.Projects).ID
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

// MemberGet retrieves the project member data for a given user and project ID.
func MemberGetByUserID(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	data := GetPathUserID{}
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
	pm.ProjectID = currentProject.(*models.Projects).ID
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

// ProjectMembersCreate projects the creation of a new member.
func ProjectMembersCreate(ctx *gin.Context) {
	// 项目管理员才添加成员
	currentMember, _ := ctx.Get("CurrentMember")
	if currentMember.(*models.ProjectMembers).Authority != models.ProjectMembersManage {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := CreateProjectMemberData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := models.NewUsers(data.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.CreateFailed"}),
		})
		return
	}

	checkMember, _ := models.NewProjectMembers()
	checkMember.UserID = user.ID
	checkMember.ProjectID = currentMember.(*models.ProjectMembers).ProjectID
	if err := checkMember.Get(); err == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.AlreadyExists"}),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.UserID = user.ID
	pm.ProjectID = currentMember.(*models.ProjectMembers).ProjectID
	pm.Authority = data.Authority
	if err := pm.Create(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.CreateFailed"}),
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

// DeleteMember deletes a project member by checking if the given member exists in the project.
func ProjectMembersDelete(ctx *gin.Context) {
	currentMember, _ := ctx.Get("CurrentMember")
	if currentMember.(*models.ProjectMembers).Authority != models.ProjectMembersManage {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := GetPathUserID{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.UserID = data.UserID
	pm.ProjectID = currentMember.(*models.ProjectMembers).ProjectID
	if err := pm.Get(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.NotFound"}),
		})
		return
	}

	if err := pm.Delete(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.DeleteFailed"}),
		})
	}

	ctx.Status(http.StatusNoContent)
}

// UpdateMember updates the authority of a project member in the database.
func ProjectMembersAuthUpdate(ctx *gin.Context) {
	currentMember, _ := ctx.Get("CurrentMember")
	if currentMember.(*models.ProjectMembers).Authority != models.ProjectMembersManage {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := GetPathUserID{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.UserID = data.UserID
	pm.ProjectID = currentMember.(*models.ProjectMembers).ProjectID
	if err := pm.Get(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.NotFound"}),
		})
		return
	}

	bodyData := UpdateProjectMemberAuthData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&bodyData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm.Authority = bodyData.Authority
	if err := pm.Update(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.UpdateFailed"}),
		})
	}

	ctx.Status(http.StatusCreated)
}
