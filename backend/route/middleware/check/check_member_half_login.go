package check

import (
	"github.com/gin-gonic/gin"
)

func CheckMemberHalfLogin() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		user := checkMemberStatus(authorization)

		if user != nil && user.IsEnabled != 0 {
			ctx.Set("CurrentUser", user)
		}
	}
}
