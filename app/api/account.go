package api

import (
	"net/http"
	"strings"

	"github.com/apicat/apicat/commom/auth"
	"github.com/apicat/apicat/commom/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginEmail struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password" binding:"required,lte=255"`
}

type RegisterEmail struct {
	Email           string `json:"email" binding:"required,email,lte=255"`
	Password        string `json:"password" binding:"required,lte=255"`
	ConfirmPassword string `json:"confirm_password" binding:"required,lte=255,eqfield=Password"`
	Username        string `json:"username" binding:"lte=255"`
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

	if checkPasswordHash(data.Password, user.Password) {
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

	hashedPassword, err := hashPassword(data.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.PasswordEncryptionFailed"}),
		})
		return
	}

	if len(data.Username) == 0 {
		user.Username = strings.Split(data.Email, "@")[0]
	} else {
		user.Username = data.Username
	}
	user.Email = data.Email
	user.Password = hashedPassword
	user.Role = "user"
	if err := user.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.RegistrationFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}
