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

// MembersList handles GET requests to retrieve a list of members in the current project.
func ProjectMembersList(ctx *gin.Context) {
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

type GetProjectMemberData struct {
	UserID uint `uri:"user_id" binding:"required"`
}

// MemberGet retrieves the project member data for a given user and project ID.
func MemberGetByUserID(ctx *gin.Context) {
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

type CreateProjectMemberData struct {
	UserID    uint   `json:"user_id" binding:"required"`
	Authority string `json:"authority" binding:"required,oneof=manage write read"`
}

func ProjectMembersCreate(ctx *gin.Context) {
	currentProject, _ := ctx.Get("CurrentProject")
	project, _ := currentProject.(*models.Projects)

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

	pm, _ := models.NewProjectMembers()
	pm.UserID = user.ID
	pm.ProjectID = project.ID
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
	currentProject, _ := ctx.Get("CurrentProject")
	currentUser, _ := ctx.Get("CurrentUser")

	checkMember, _ := models.NewProjectMembers()
	checkMember.UserID = currentUser.(*models.Users).ID
	checkMember.ProjectID = currentProject.(*models.Projects).ID
	if err := checkMember.Get(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.NotFound"}),
		})
		return
	}

	pmd := ProjectMemberIDData{}
	pm, err := pmd.CheckMember(ctx)
	if err != nil {
		return
	}

	if err := pm.Delete(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.DeleteFailed"}),
		})
	}

	ctx.Status(http.StatusNoContent)
}

type UpdateProjectMemberAuthData struct {
	Authority string `json:"authority" binding:"required,oneof=manage write read"`
}

// UpdateMember updates the authority of a project member in the database.
func ProjectMembersAuthUpdate(ctx *gin.Context) {
	// currentProject, _ := ctx.Get("CurrentProject")
	// currentUser, _ := ctx.Get("CurrentUser")

	// checkMember, _ := models.NewProjectMembers()
	// checkMember.UserID = currentUser.(*models.Users).ID
	// checkMember.ProjectID = currentProject.(*models.Projects).ID
	// if err := checkMember.Get(); err != nil {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{
	// 		"message": translator.Trasnlate(ctx, &translator.TT{ID: "Auth.TokenParsingFailed"}),
	// 	})
	// 	return
	// }

	// if !slices.Contains([]string{"manage", "write", "read"}, v) {
	// 	return
	// }

	pmd := ProjectMemberIDData{}
	pm, err := pmd.CheckMember(ctx)
	if err != nil {
		return
	}

	data := UpdateProjectMemberAuthData{}
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	pm.Authority = data.Authority
	if err := pm.Update(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "ProjectMember.UpdateFailed"}),
		})
	}

	ctx.Status(http.StatusCreated)
}
