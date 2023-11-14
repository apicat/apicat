package project

import (
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/route/proto"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MembersList handles GET requests to retrieve a list of members in the current project.
func ProjectMembersList(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	var data proto.GetMembersData
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.Page <= 0 {
		data.Page = 1
	}
	if data.PageSize <= 0 {
		data.PageSize = 15
	}

	member, _ := project.NewProjectMembers()
	member.ProjectID = currentProject.(*project.Projects).ID
	totalMember, err := member.Count()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	member.ProjectID = currentProject.(*project.Projects).ID
	members, err := member.List(data.Page, data.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	userIDs := []uint{}
	for _, v := range members {
		userIDs = append(userIDs, v.UserID)
	}

	users, err := user.UserListByIDs(userIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	userIDToNameMap := map[uint]*user.Users{}
	for _, v := range users {
		userIDToNameMap[v.ID] = v
	}

	membersList := []any{}
	for _, v := range members {
		membersList = append(membersList, map[string]any{
			"id":         v.ID,
			"user_id":    v.UserID,
			"username":   userIDToNameMap[v.UserID].Username,
			"authority":  v.Authority,
			"is_enabled": userIDToNameMap[v.UserID].IsEnabled,
			"email":      userIDToNameMap[v.UserID].Email,
			"created_at": v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"current_page": data.Page,
		"total_page":   int(math.Ceil(float64(totalMember) / float64(data.PageSize))),
		"total":        totalMember,
		"records":      membersList,
	})
}

// (Abandoned) MemberGet retrieves the project member data for a given user and project ID.
func MemberGetByUserID(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")

	data := proto.GetPathUserID{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	u, err := user.NewUsers(data.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	pm, _ := project.NewProjectMembers()
	pm.UserID = u.ID
	pm.ProjectID = currentProject.(*project.Projects).ID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, proto.ProjectMemberData{
		ID:        pm.ID,
		UserID:    u.ID,
		Username:  u.Username,
		Authority: pm.Authority,
		CreatedAt: pm.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

// ProjectMembersCreate projects the creation of a new member.
func ProjectMembersCreate(ctx *gin.Context) {
	// 项目管理员才添加成员
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := proto.CreateProjectMemberData{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := []gin.H{}
	for _, v := range data.UserIDs {
		u, err := user.NewUsers(v)
		if err != nil {
			continue
		}

		pm, _ := project.NewProjectMembers()
		pm.UserID = u.ID
		pm.ProjectID = currentProjectMember.(*project.ProjectMembers).ProjectID
		if err := pm.GetByUserIDAndProjectID(); err == nil {
			continue
		}

		pm.Authority = data.Authority
		if err := pm.Create(); err != nil {
			continue
		}

		result = append(result, gin.H{
			"id":         pm.ID,
			"user_id":    u.ID,
			"username":   u.Username,
			"email":      u.Email,
			"is_enabled": u.IsEnabled,
			"authority":  pm.Authority,
			"created_at": pm.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, result)
}

// DeleteMember deletes a project member by checking if the given member exists in the project.
func ProjectMembersDelete(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := proto.GetPathUserID{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.UserID == currentProjectMember.(*project.ProjectMembers).UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.DeleteFailed"}),
		})
		return
	}

	pm, _ := project.NewProjectMembers()
	pm.UserID = data.UserID
	pm.ProjectID = currentProjectMember.(*project.ProjectMembers).ProjectID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.NotFound"}),
		})
		return
	}

	if err := pm.Delete(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.DeleteFailed"}),
		})
	}

	ctx.Status(http.StatusNoContent)
}

// UpdateMember updates the authority of a project member in the database.
func ProjectMembersAuthUpdate(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	data := proto.GetPathUserID{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindUri(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm, _ := project.NewProjectMembers()
	pm.UserID = data.UserID
	pm.ProjectID = currentProjectMember.(*project.ProjectMembers).ProjectID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.NotFound"}),
		})
		return
	}

	bodyData := proto.UpdateProjectMemberAuthData{}
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&bodyData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm.Authority = bodyData.Authority
	if err := pm.Update(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.UpdateFailed"}),
		})
	}

	ctx.Status(http.StatusCreated)
}

func ProjectMembersWithout(ctx *gin.Context) {
	currentProjectMember, _ := ctx.Get("CurrentProjectMember")
	if !currentProjectMember.(*project.ProjectMembers).MemberIsManage() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	u, _ := user.NewUsers()
	users, err := u.List(0, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	projectMember, _ := project.NewProjectMembers()
	projectMember.ProjectID = currentProjectMember.(*project.ProjectMembers).ProjectID
	projectMembers, err := projectMember.List(0, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "ProjectMember.QueryFailed"}),
		})
		return
	}

	pmMap := map[uint]project.ProjectMembers{}
	for _, v := range projectMembers {
		pmMap[v.UserID] = v
	}

	result := []map[string]any{}
	for _, u := range users {
		if _, ok := pmMap[u.ID]; !ok {
			result = append(result, map[string]any{
				"user_id":  u.ID,
				"username": u.Username,
				"email":    u.Email,
			})
		}
	}

	ctx.JSON(http.StatusOK, result)
}
