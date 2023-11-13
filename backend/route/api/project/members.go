package project

import (
	"fmt"
	"github.com/apicat/apicat/backend/model/project"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/auth"
	"github.com/apicat/apicat/backend/module/translator"
	"github.com/apicat/apicat/backend/route/proto"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetMembers(ctx *gin.Context) {
	var data proto.GetMembersData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
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

	u, _ := user.NewUsers()
	totalUsers, err := u.Count()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.QueryFailed"}),
		})
		return
	}

	users, err := u.List(data.Page, data.PageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.QueryFailed"}),
		})
		return
	}

	userList := []any{}
	if len(users) > 0 {
		for _, v := range users {
			userList = append(userList, map[string]any{
				"id":         v.ID,
				"username":   v.Username,
				"email":      v.Email,
				"role":       v.Role,
				"is_enabled": v.IsEnabled,
				"created_at": v.CreatedAt.Format("2006-01-02 15:04:05"),
				"updated_at": v.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"current_page": data.Page,
		"total_page":   int(math.Ceil(float64(totalUsers) / float64(data.PageSize))),
		"total":        totalUsers,
		"records":      userList,
	})
}

func AddMember(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")
	if currentUser.(*user.Users).Role != "superadmin" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.MemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var data proto.AddMemberData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	u, _ := user.NewUsers()
	u.Email = data.Email
	if err := u.GetByEmail(data.Email); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.MailboxAlreadyExists"}),
		})
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.CreateFailed"}),
		})
		return
	}

	u.Username = strings.Split(data.Email, "@")[0]
	u.Role = data.Role
	u.Password = hashedPassword

	if err := u.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.CreateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         u.ID,
		"email":      u.Email,
		"username":   u.Username,
		"role":       u.Role,
		"is_enabled": u.IsEnabled,
		"created_at": u.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": u.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func SetMember(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")
	if currentUser.(*user.Users).Role != "superadmin" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.MemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var userIDData proto.UserIDData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&userIDData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if currentUser.(*user.Users).ID == userIDData.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.UpdateFailed"}),
		})
		return
	}

	var data proto.SetMemberData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.Email != "" {
		checkUser, _ := user.NewUsers()
		if err := checkUser.GetByEmail(data.Email); err == nil && checkUser.ID != userIDData.UserID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.MailboxAlreadyExists"}),
			})
			return
		}
	}

	u, err := user.NewUsers(userIDData.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.UpdateFailed"}),
		})
		return
	}

	if data.Email != "" {
		u.Email = data.Email
	}
	if data.Role != "" {
		u.Role = data.Role
	}
	u.IsEnabled = data.IsEnabled
	if data.Password != "" {
		hashedPassword, err := auth.HashPassword(data.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.UpdateFailed"}),
			})
			return
		}
		u.Password = hashedPassword
	}

	if err := u.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func DeleteMember(ctx *gin.Context) {
	currentUser, _ := ctx.Get("CurrentUser")
	if currentUser.(*user.Users).Role != "superadmin" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    proto.MemberInsufficientPermissionsCode,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Common.InsufficientPermissions"}),
		})
		return
	}

	var userIDData proto.UserIDData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&userIDData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if currentUser.(*user.Users).ID == userIDData.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.DeleteFailed"}),
		})
		return
	}

	u, err := user.NewUsers(userIDData.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    proto.Display404ErrorMessage,
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.DeleteFailed"}),
		})
		return
	}

	// Determine whether the member to be deleted is an administrator of a project
	pm, err := project.GetUserInvolvedProject(u.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.DeleteFailed"}),
		})
		return
	}

	ps := []project.ProjectMembers{}
	for _, v := range pm {
		if v.Authority == project.ProjectMembersManage {
			ps = append(ps, v)
		}
	}
	if len(ps) > 0 {
		tm := translator.Trasnlate(ctx, &translator.TT{ID: "Member.IsProjectManage"})
		pns := []string{}
		for _, v := range ps {
			if p, err := project.NewProjects(v.ProjectID); err == nil {
				pns = append(pns, p.Title)
			}
		}

		pn := strings.Join(pns, "、")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf(tm, pn),
		})
		return
	}

	if err := u.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.DeleteFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
