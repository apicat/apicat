package iteration

import (
	"github.com/apicat/apicat/backend/model/collection"
	"github.com/apicat/apicat/backend/model/iteration"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/translator"
	"github.com/apicat/apicat/backend/route/proto"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v4"
)

func IterationsList(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		data proto.IterationListData
		res  proto.IterationListResData
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
	res.Iterations = []proto.IterationSchemaData{}

	pmDict := map[uint]project.ProjectMembers{}
	if data.ProjectID != "" {
		targetProject, err := project.NewProjects(data.ProjectID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.QueryFailed"}),
			})
			return
		}

		pm, _ := project.NewProjectMembers()
		pm.ProjectID = targetProject.ID
		pm.UserID = currentUser.(*user.Users).ID
		if err := pm.GetByUserIDAndProjectID(); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    proto.ProjectMemberInsufficientPermissionsCode,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
			})
			return
		}
		pIDs = append(pIDs, targetProject.ID)
		pmDict[targetProject.ID] = *pm
	} else {
		pms, err := project.GetUserInvolvedProject(currentUser.(*user.Users).ID)
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

	p, _ := project.NewProjects()
	projects, err := p.List(pIDs...)
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

	pDict := map[uint]project.Projects{}
	for _, v := range projects {
		pDict[v.ID] = v
	}

	iterationTotal, err := iteration.IterationsCount(pIDs...)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.QueryFailed"}),
		})
		return
	}
	res.TotalPage = int64(math.Ceil(float64(iterationTotal) / float64(data.PageSize)))
	res.Total = iterationTotal

	i, _ := iteration.NewIterations()
	iterations, err := i.List(int(data.Page), int(data.PageSize), pIDs...)
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

	iterationApi, _ := iteration.NewIterationApis()
	iterationApis, err := iterationApi.List(iterationIDs...)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.QueryFailed"}),
		})
		return
	}

	for _, i := range iterations {
		apiNum := 0
		for _, v := range iterationApis {
			if i.ID == v.IterationID && v.CollectionType != "category" {
				apiNum++
			}
		}

		res.Iterations = append(res.Iterations, proto.IterationSchemaData{
			ID:           i.PublicID,
			Title:        i.Title,
			Description:  i.Description,
			ProjectID:    pDict[i.ProjectID].PublicId,
			ProjectTitle: pDict[i.ProjectID].Title,
			ApiNum:       int64(apiNum),
			Authority:    pmDict[i.ProjectID].Authority,
			CreatedAt:    i.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, res)
}

func IterationsDetails(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		uriData proto.IterationUriData
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	i, err := iteration.NewIterations(uriData.IterationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Redirect404Page,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.NotFound"}),
		})
		return
	}

	p, err := project.NewProjects(i.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Redirect404Page,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	pm, _ := project.NewProjectMembers()
	pm.ProjectID = p.ID
	pm.UserID = currentUser.(*user.Users).ID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	apiNum, err := iteration.IterationApiCount(i.ID, "api")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.CreateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, proto.IterationSchemaData{
		ID:           i.PublicID,
		Title:        i.Title,
		Description:  i.Description,
		ProjectID:    p.PublicId,
		ProjectTitle: p.Title,
		ApiNum:       apiNum,
		Authority:    pm.Authority,
		CreatedAt:    i.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

func IterationsCreate(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		data proto.IterationCreateData
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	p, err := project.NewProjects(data.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Projects.NotFound"}),
		})
		return
	}

	pm, _ := project.NewProjectMembers()
	pm.ProjectID = p.ID
	pm.UserID = currentUser.(*user.Users).ID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}
	if !pm.MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	i, _ := iteration.NewIterations()
	i.PublicID = shortuuid.New()
	i.ProjectID = p.ID
	i.Title = data.Title
	i.Description = data.Description
	i.CreatedBy = currentUser.(*user.Users).ID
	if err := i.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.CreateFailed"}),
		})
		return
	}

	if err := collection.PlanningIterationApi(data.CollectionIDs, i); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.CreateFailed"}),
		})
		return
	}

	apiNum, err := iteration.IterationApiCount(i.ID, "api")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.CreateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusOK, proto.IterationSchemaData{
		ID:           i.PublicID,
		Title:        i.Title,
		Description:  i.Description,
		ProjectID:    p.PublicId,
		ProjectTitle: p.Title,
		ApiNum:       apiNum,
		Authority:    pm.Authority,
		CreatedAt:    i.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

func IterationsUpdate(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")

	var (
		uriData proto.IterationUriData
		data    proto.IterationUpdateData
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

	i, err := iteration.NewIterations(uriData.IterationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.NotFound"}),
		})
		return
	}

	pm, _ := project.NewProjectMembers()
	pm.ProjectID = i.ProjectID
	pm.UserID = currentUser.(*user.Users).ID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}
	if !pm.MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	i.Title = data.Title
	i.Description = data.Description
	i.UpdatedBy = currentUser.(*user.Users).ID
	if err := i.Update(); err == nil {
		if err := collection.PlanningIterationApi(data.CollectionIDs, i); err != nil {
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
		uriData proto.IterationUriData
	)

	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&uriData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	i, err := iteration.NewIterations(uriData.IterationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.NotFound"}),
		})
		return
	}

	pm, _ := project.NewProjectMembers()
	pm.ProjectID = i.ProjectID
	pm.UserID = currentUser.(*user.Users).ID
	if err := pm.GetByUserIDAndProjectID(); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}
	if !pm.MemberHasWritePermission() {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.ProjectMemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	i.DeletedBy = currentUser.(*user.Users).ID
	if err := i.Update(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.DeleteFailed"}),
		})
		return
	}

	if err := i.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Iteration.DeleteFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
