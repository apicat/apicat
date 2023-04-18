package middleware

import (
	"net/http"
	"strings"

	"github.com/apicat/apicat/commom/auth"
	"github.com/gin-gonic/gin"
)

// 基于JWT认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO 补充中间件中所有message的i18n

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Status(http.StatusUnauthorized)
			//阻止调用后续的函数
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		mc, err := auth.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err,
			})
			c.Abort()
			return
		}

		if mc.User == nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		//将当前请求的username信息保存到请求的上下文c上
		c.Set("CurrentUser", mc.User)
		//后续的处理函数可以通过c.Get("CurrentUser")来获取请求的用户信息
		c.Next()
	}
}
