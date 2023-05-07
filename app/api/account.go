package api

import (
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/auth"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type LoginEmail struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=255"`
}

type RegisterEmail struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=255"`
}

func EmailLogin(ctx *gin.Context) {
	var data LoginEmail
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, _ := models.NewUsers()
	if err := user.GetByEmail(data.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !auth.CheckPasswordHash(data.Password, user.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.WrongPassword"}),
		})
		return
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.FailedToGenerateToken"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"expires_in":   auth.TokenExpireDuration,
		"user": map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func EmailRegister(ctx *gin.Context) {
	var data RegisterEmail
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, _ := models.NewUsers()
	if err := user.GetByEmail(data.Email); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.MailboxAlreadyExists"}),
		})
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.PasswordEncryptionFailed"}),
		})
		return
	}

	user.Username = strings.Split(data.Email, "@")[0]
	user.Email = data.Email
	user.Password = hashedPassword
	// 第一个注册的用户权限为superadmin
	userCount, err := user.Count()
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if userCount == 0 {
		user.Role = "superadmin"
	} else {
		user.Role = "admin"
	}

	if err := user.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.RegistrationFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
