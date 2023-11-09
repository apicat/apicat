package check

import (
	"net/http"
	"strings"

	"github.com/apicat/apicat/backend/common/auth"
	"github.com/apicat/apicat/backend/common/translator"
	"github.com/apicat/apicat/backend/enum"
	"github.com/apicat/apicat/backend/models"
	"github.com/gin-gonic/gin"
)

func checkMemberStatus(authorization string) *models.Users {
	if authorization == "" {
		return nil
	}

	parts := strings.SplitN(authorization, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil
	}

	mc, err := auth.ParseToken(parts[1])
	if err != nil {
		return nil
	}

	if mc.UserID == 0 {
		return nil
	}

	user, err := models.NewUsers(mc.UserID)
	if err != nil {
		return nil
	}

	return user
}

func CheckMember() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		user := checkMemberStatus(authorization)

		if user == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    enum.InvalidOrIncorrectLoginToken,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Auth.TokenParsingFailed"}),
			})
			ctx.Abort()
			return
		}

		if user == nil || user.IsEnabled == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    enum.InvalidOrIncorrectLoginToken,
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Auth.AccountDisabled"}),
			})
			ctx.Abort()
			return
		}

		//将当前请求的username信息保存到请求的上下文c上
		ctx.Set("CurrentUser", user)
		//后续的处理函数可以通过c.Get("CurrentUser")来获取请求的用户信息
		ctx.Next()
	}
}
