package api

import (
	"math"
	"net/http"

	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/apicat/apicat/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

type IterationSchemaData struct {
	ID              uint   `json:"id,omitempty" binding:"required"`
	PublicID        string `json:"public_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	ProjectPublicID string `json:"project_public_id" binding:"required"`
	ProjectTitle    string `json:"project_title" binding:"required"`
	IsFollowed      bool   `json:"is_followed"`
	ApiNum          int64  `json:"api_num"`
	Authority       string `json:"authority"`
	CreatedAt       string `json:"created_at" binding:"required"`
	CreatedBy       string `json:"created_by" binding:"required"`
}

type IterationListData struct {
	ProjectID uint  `form:"project_id"`
	Page      int64 `form:"page"`
	PageSize  int64 `form:"page_size"`
}

type IterationListResData struct {
	CurrentPage int64                 `json:"current_page"`
	TotalPage   int64                 `json:"total_page"`
	Total       int64                 `json:"total"`
	Iterations  []IterationSchemaData `json:"iterations"`
}

type IterationCreateData struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	ProjectID     uint   `json:"project_id" binding:"required"`
	CollectionIDs []uint `json:"collection_ids"`
}

type IterationUriData struct {
	IterationID uint `uri:"iteration_id" binding:"required"`
}

type IterationUpdateData struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	CollectionIDs []uint `json:"collection_ids"`
}

func IterationsList(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		data IterationListData
		res  IterationListResData
		pIDs []uint
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.Page == 0 {
		data.Page = 1
	}
	if data.PageSize == 0 {
		data.PageSize = 15
	}

	pmDict := map[uint]models.ProjectMembers{}
	if data.ProjectID != 0 {
		pm, _ := models.NewProjectMembers()
		pm.ProjectID = data.ProjectID
		pm.UserID = currentUser.(*models.Users).ID
		if err := pm.GetByUserIDAndProjectID(); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ProjectMemberInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
			})
			return
		}
		pIDs = append(pIDs, data.ProjectID)
		pmDict[data.ProjectID] = *pm
	} else {
		pms, err := models.GetUserInvolvedProject(currentUser.(*models.Users).ID)
		if err != nil {
			ctx.JSON(http.StatusOK, res)
			return
		}
		for _, pm := range pms {
			pIDs = append(pIDs, pm.ProjectID)
			pmDict[pm.ProjectID] = pm
		}
	}

	project, _ := models.NewProjects()
	projects, err := project.List(pIDs...)
	if err != nil {
		ctx.JSON(http.StatusOK, res)
		return
	}
	pDict := map[uint]models.Projects{}
	for _, v := range projects {
		pDict[v.ID] = v
	}

	user, _ := models.NewUsers()
	users, err := user.List(0, 0)
	if err != nil {
		ctx.JSON(http.StatusOK, res)
		return
	}
	uDict := map[uint]models.Users{}
	for _, v := range users {
		uDict[v.ID] = v
	}

	iteration, _ := models.NewIterations()
	iterations, err := iteration.List(int(data.Page), int(data.PageSize), pIDs...)
	if err != nil {
		ctx.JSON(http.StatusOK, res)
		return
	}

	pf, _ := models.NewProjectFollows()
	pfs, err := pf.List(currentUser.(*models.Users).ID)
	if err != nil {
		ctx.JSON(http.StatusOK, res)
		return
	}

	iterationTotal, err := iteration.IterationsCount(pIDs...)
	if err != nil {
		ctx.JSON(http.StatusOK, res)
		return
	}

	iterationApi, _ := models.NewIterationApis()
	iterationApis, err := iterationApi.List()
	if err != nil {
		ctx.JSON(http.StatusOK, res)
		return
	}

	res.CurrentPage = data.Page
	res.TotalPage = int64(math.Ceil(float64(iterationTotal) / float64(data.PageSize)))
	res.Total = iterationTotal

	for _, i := range iterations {
		isFollowed := false
		for _, v := range pfs {
			if i.ProjectID == v.ProjectID {
				isFollowed = true
				break
			}
		}
		apiNum := 0
		for _, v := range iterationApis {
			if i.ID == v.IterationID {
				apiNum++
			}
		}

		res.Iterations = append(res.Iterations, IterationSchemaData{
			ID:              i.ID,
			PublicID:        i.PublicID,
			Title:           i.Title,
			Description:     i.Description,
			ProjectPublicID: pDict[i.ProjectID].PublicId,
			ProjectTitle:    pDict[i.ProjectID].Title,
			IsFollowed:      isFollowed,
			ApiNum:          int64(apiNum),
			Authority:       pmDict[i.ProjectID].Authority,
			CreatedAt:       i.CreatedAt.Format("2006-01-02 15:04"),
			CreatedBy:       uDict[i.CreatedBy].Username,
		})
	}

	ctx.JSON(http.StatusOK, res)
}

func IterationsCreate(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		data       IterationCreateData
		isFollowed bool
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	project, err := models.NewProjects(data.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.ProjectID = project.ID
	pm.UserID = currentUser.(*models.Users).ID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}
	if !pm.MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	iteration, _ := models.NewIterations()
	iteration.PublicID = shortuuid.New()
	iteration.ProjectID = project.ID
	iteration.Title = data.Title
	iteration.Description = data.Description
	iteration.CreatedBy = currentUser.(*models.Users).ID
	if err := iteration.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.CreateFailed"}),
		})
		return
	}

	if err := iteration.PlanningIterationApi(data.CollectionIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.CreateFailed"}),
		})
		return
	}

	apiNum, err := models.IterationApiCount(iteration.ID, "api")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.CreateFailed"}),
		})
		return
	}

	pf, _ := models.NewProjectFollows()
	pf.UserID = currentUser.(*models.Users).ID
	pf.ProjectID = project.ID
	if err := pf.GetByUserIDAndProjectID(); err == nil {
		isFollowed = true
	}

	ctx.JSON(http.StatusOK, IterationSchemaData{
		ID:              iteration.ID,
		PublicID:        iteration.PublicID,
		Title:           iteration.Title,
		Description:     iteration.Description,
		ProjectPublicID: project.PublicId,
		ProjectTitle:    project.Title,
		IsFollowed:      isFollowed,
		ApiNum:          apiNum,
		Authority:       pm.Authority,
		CreatedAt:       iteration.CreatedAt.Format("2006-01-02 15:04"),
		CreatedBy:       currentUser.(*models.Users).Username,
	})
}

func IterationsUpdate(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		uriData IterationUriData
		data    IterationUpdateData
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	iteration, err := models.NewIterations(uriData.IterationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iterations.NotFound"}),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.ProjectID = iteration.ProjectID
	pm.UserID = currentUser.(*models.Users).ID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}
	if !pm.MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	iteration.Title = data.Title
	iteration.Description = data.Description
	if err := iteration.Update(); err == nil {
		if err := iteration.PlanningIterationApi(data.CollectionIDs); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.UpdateFailed"}),
			})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func IterationsDelete(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		uriData IterationUriData
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	iteration, err := models.NewIterations(uriData.IterationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iterations.NotFound"}),
		})
		return
	}

	pm, _ := models.NewProjectMembers()
	pm.ProjectID = iteration.ProjectID
	pm.UserID = currentUser.(*models.Users).ID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}
	if !pm.MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    enum.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	if err := iteration.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.DeleteFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
