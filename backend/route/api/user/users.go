package user

import (
	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/auth"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	u, _ := CurrentUser.(*user.Users)

	ctx.JSON(200, gin.H{
		"id":         u.ID,
		"username":   u.Username,
		"email":      u.Email,
		"role":       u.Role,
		"is_enabled": u.IsEnabled,
		"created_at": u.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": u.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func SetUserInfo(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	currentUser, _ := CurrentUser.(*user.Users)

	var data proto.SetUserInfoData
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	u, _ := user.NewUsers()
	if err := u.GetByEmail(data.Email); err == nil && u.ID != currentUser.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "User.MailboxAlreadyExists"}),
		})
		return
	}

	currentUser.Email = data.Email
	currentUser.Username = data.Username
	if err := currentUser.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         currentUser.ID,
		"username":   currentUser.Username,
		"email":      currentUser.Email,
		"role":       currentUser.Role,
		"is_enabled": currentUser.IsEnabled,
		"created_at": currentUser.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": currentUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func ChangePassword(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	currentUser, _ := CurrentUser.(*user.Users)

	var data proto.ChangePasswordData
	if err := i18n.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !auth.CheckPasswordHash(data.Password, currentUser.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "User.WrongPassword"}),
		})
		return
	}

	hashedPassword, err := auth.HashPassword(data.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	currentUser.Password = hashedPassword
	if err := currentUser.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": i18n.Trasnlate(ctx, &i18n.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
