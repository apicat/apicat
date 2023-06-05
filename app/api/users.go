package api

import (
	"math"
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/auth"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type SetUserInfoData struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Username string `json:"username" binding:"required,lte=255"`
}

type ChangePasswordData struct {
	Password           string `json:"password" binding:"required,gte=6,lte=255"`
	NewPassword        string `json:"new_password" binding:"required,gte=6,lte=255"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,gte=6,lte=255,eqfield=NewPassword"`
}

type GetUsersData struct {
	Page     int `form:"page" binding:"omitempty,gte=1"`
	PageSize int `form:"page_size" binding:"omitempty,gte=1,lte=100"`
}

func GetUserInfo(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	user, _ := CurrentUser.(*models.Users)

	ctx.JSON(200, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func SetUserInfo(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	currentUser, _ := CurrentUser.(*models.Users)

	var data SetUserInfoData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, _ := models.NewUsers()
	if err := user.GetByEmail(data.Email); err == nil && user.ID != currentUser.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.MailboxAlreadyExists"}),
		})
		return
	}

	currentUser.Email = data.Email
	currentUser.Username = data.Username
	if err := currentUser.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	token, err := auth.GenerateToken(currentUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": map[string]interface{}{
			"id":         currentUser.ID,
			"username":   currentUser.Username,
			"email":      currentUser.Email,
			"role":       currentUser.Role,
			"created_at": currentUser.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": currentUser.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func ChangePassword(ctx *gin.Context) {
	CurrentUser, _ := ctx.Get("CurrentUser")
	currentUser, _ := CurrentUser.(*models.Users)

	var data ChangePasswordData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !auth.CheckPasswordHash(data.Password, currentUser.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.WrongPassword"}),
		})
		return
	}

	hashedPassword, err := auth.HashPassword(data.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	currentUser.Password = hashedPassword
	if err := currentUser.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	token, err := auth.GenerateToken(currentUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.UpdateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

func GetUsers(ctx *gin.Context) {
	var data GetUsersData
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

	user, _ := models.NewUsers()
	totalUsers, err := user.Count()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.QueryFailed"}),
		})
		return
	}

	users, err := user.List(data.Page, data.PageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.QueryFailed"}),
		})
		return
	}

	userList := []any{}
	if len(users) > 0 {
		for _, v := range users {
			email := v.Email
			parts := strings.Split(email, "@")

			userList = append(userList, map[string]any{
				"id":         v.ID,
				"username":   v.Username,
				"email":      parts[0][0:1] + "***" + parts[0][len(parts[0])-1:] + "@" + parts[len(parts)-1],
				"role":       v.Role,
				"created_at": v.CreatedAt.Format("2006-01-02 15:04:05"),
				"updated_at": v.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"current_page": data.Page,
		"total_page":   int(math.Ceil(float64(totalUsers) / float64(data.PageSize))),
		"total_users":  totalUsers,
		"users":        userList,
	})
}
