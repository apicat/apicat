package middleware

import (
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/auth"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

// 基于JWT认证中间件
func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Auth.TokenParsingFailed"}),
			})
			//阻止调用后续的函数
			ctx.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Auth.TokenParsingFailed"}),
			})
			ctx.Abort()
			return
		}

		mc, err := auth.ParseToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Auth.TokenParsingFailed"}),
			})
			ctx.Abort()
			return
		}

		if mc.User == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Auth.TokenParsingFailed"}),
			})
			ctx.Abort()
			return
		}

		user, err := models.NewUsers(mc.User.ID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Auth.TokenParsingFailed"}),
			})
			ctx.Abort()
			return
		}

		if user.IsEnabled == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
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
