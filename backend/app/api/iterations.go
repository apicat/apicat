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
	ID           string `json:"id,omitempty" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	ProjectID    string `json:"project_id" binding:"required"`
	ProjectTitle string `json:"project_title" binding:"required"`
	ApiNum       int64  `json:"api_num"`
	Authority    string `json:"authority"`
	CreatedAt    string `json:"created_at" binding:"required"`
}

type IterationListData struct {
	ProjectID string `form:"project_id"`
	Page      int64  `form:"page"`
	PageSize  int64  `form:"page_size"`
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
	ProjectID     string `json:"project_id" binding:"required"`
	CollectionIDs []uint `json:"collection_ids"`
}

type IterationUriData struct {
	IterationID string `uri:"iteration-id" binding:"required"`
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

	res.CurrentPage = data.Page
	res.Iterations = []IterationSchemaData{}

	pmDict := map[uint]models.ProjectMembers{}
	if data.ProjectID != "" {
		targetProject, err := models.NewProjects(data.ProjectID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.QueryFailed"}),
			})
			return
		}

		pm, _ := models.NewProjectMembers()
		pm.ProjectID = targetProject.ID
		pm.UserID = currentUser.(*models.Users).ID
		if err := pm.GetByUserIDAndProjectID(); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    enum.ProjectMemberInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
			})
			return
		}
		pIDs = append(pIDs, targetProject.ID)
		pmDict[targetProject.ID] = *pm
	} else {
		pms, err := models.GetUserInvolvedProject(currentUser.(*models.Users).ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.QueryFailed"}),
			})
			return
		}
		if len(pms) == 0 {
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.QueryFailed"}),
		})
		return
	}
	if len(projects) == 0 {
		ctx.JSON(http.StatusOK, res)
		return
	}

	pDict := map[uint]models.Projects{}
	for _, v := range projects {
		pDict[v.ID] = v
	}

	iteration, _ := models.NewIterations()
	iterations, err := iteration.List(int(data.Page), int(data.PageSize), pIDs...)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.QueryFailed"}),
		})
		return
	}
	if len(iterations) == 0 {
		ctx.JSON(http.StatusOK, res)
		return
	}

	var iterationIDs []uint
	for _, v := range iterations {
		iterationIDs = append(iterationIDs, v.ID)
	}

	iterationTotal, err := iteration.IterationsCount(pIDs...)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.QueryFailed"}),
		})
		return
	}

	iterationApi, _ := models.NewIterationApis()
	iterationApis, err := iterationApi.List(iterationIDs...)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.QueryFailed"}),
		})
		return
	}

	res.TotalPage = int64(math.Ceil(float64(iterationTotal) / float64(data.PageSize)))
	res.Total = iterationTotal

	for _, i := range iterations {
		apiNum := 0
		for _, v := range iterationApis {
			if i.ID == v.IterationID {
				apiNum++
			}
		}

		res.Iterations = append(res.Iterations, IterationSchemaData{
			ID:           i.PublicID,
			Title:        i.Title,
			Description:  i.Description,
			ProjectID:    pDict[i.ProjectID].PublicId,
			ProjectTitle: pDict[i.ProjectID].Title,
			ApiNum:       int64(apiNum),
			Authority:    pmDict[i.ProjectID].Authority,
			CreatedAt:    i.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	ctx.JSON(http.StatusOK, res)
}

func IterationsDetails(ctx *gin.Context) {
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

	project, err := models.NewProjects(iteration.ProjectID)
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

	apiNum, err := models.IterationApiCount(iteration.ID, "api")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.CreateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, IterationSchemaData{
		ID:           iteration.PublicID,
		Title:        iteration.Title,
		Description:  iteration.Description,
		ProjectID:    project.PublicId,
		ProjectTitle: project.Title,
		ApiNum:       apiNum,
		Authority:    pm.Authority,
		CreatedAt:    iteration.CreatedAt.Format("2006-01-02 15:04"),
	})
}

func IterationsCreate(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		data IterationCreateData
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

	ctx.JSON(http.StatusOK, IterationSchemaData{
		ID:           iteration.PublicID,
		Title:        iteration.Title,
		Description:  iteration.Description,
		ProjectID:    project.PublicId,
		ProjectTitle: project.Title,
		ApiNum:       apiNum,
		Authority:    pm.Authority,
		CreatedAt:    iteration.CreatedAt.Format("2006-01-02 15:04"),
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
	iteration.UpdatedBy = currentUser.(*models.Users).ID
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

	iteration.DeletedBy = currentUser.(*models.Users).ID
	if err := iteration.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.DeleteFailed"}),
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
