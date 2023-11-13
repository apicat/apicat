package user

import (
	"github.com/apicat/apicat/backend/model/user"
	"github.com/apicat/apicat/backend/module/auth"
	"github.com/apicat/apicat/backend/module/translator"
	"github.com/apicat/apicat/backend/route/proto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func EmailLogin(ctx *gin.Context) {
	var data proto.LoginEmail
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	u, _ := user.NewUsers()
	if err := u.GetByEmail(data.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.AccountDoesNotExist"}),
		})
		return
	}

	if !auth.CheckPasswordHash(data.Password, u.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.WrongPassword"}),
		})
		return
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.LoginFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"expires_in":   auth.TokenExpireDuration,
		"user": map[string]interface{}{
			"id":         u.ID,
			"username":   u.Username,
			"email":      u.Email,
			"role":       u.Role,
			"created_at": u.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": u.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func EmailRegister(ctx *gin.Context) {
	var data proto.RegisterEmail
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	u, _ := user.NewUsers()
	if err := u.GetByEmail(data.Email); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.MailboxAlreadyExists"}),
		})
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.RegistrationFailed"}),
		})
		return
	}

	u.Username = strings.Split(data.Email, "@")[0]
	u.Email = data.Email
	u.Password = hashedPassword
	// 第一个注册的用户权限为superadmin
	userCount, err := u.Count()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.RegistrationFailed"}),
		})
		return
	}
	if userCount == 0 {
		u.Role = "superadmin"
	} else {
		u.Role = "admin"
	}

	if err := u.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.RegistrationFailed"}),
		})
		return
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.RegistrationFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token": token,
		"expires_in":   auth.TokenExpireDuration,
		"user": map[string]interface{}{
			"id":         u.ID,
			"username":   u.Username,
			"email":      u.Email,
			"role":       u.Role,
			"created_at": u.CreatedAt.Format("2006-01-02 15:04:05"),
			"updated_at": u.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
